package chrono

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

const (
	defaultTimeout = 15 * time.Second
)

var DefaultScheduler *Scheduler

// InitScheduler initializes the default scheduler.
func InitScheduler(ctx context.Context, monitor SchedulerMonitor, options ...SchedulerOption) error {
	var err error
	DefaultScheduler, err = NewScheduler(ctx, monitor, options...)
	if err != nil {
		return err
	}
	return nil
}

// Scheduler is the base gocron scheduler.
// Scheduler 是基础的 gocron 调度器。
type Scheduler struct {
	ctx context.Context // Context for scheduler lifecycle
	// 调度器生命周期的上下文
	scheduler gocron.Scheduler // Underlying gocron scheduler
	// 底层 gocron 调度器
	monitor SchedulerMonitor // Scheduler monitor
	// 调度器监控器
	aliasMap map[string]string // Alias to jobID mapping
	// 别名到 jobID 的映射
	watchFuncMap map[string]func(event JobWatchInterface) // JobID to watch function mapping
	// jobID 到监听函数的映射
	mu sync.Mutex // Mutex to protect watchFuncMap
	// 用于保护 watchFuncMap 的互斥锁
	schOptions *SchedulerOptions // Scheduler options
	// 调度器选项
}

// SchedulerOptions holds options for the scheduler.
// SchedulerOptions 保存调度器的选项。
type SchedulerOptions struct {
	// 别名选项
	// Alias option
	aliasOption ChronoOption
	// Watch option
	// 监听选项
	watchOption *WatchOption
	// Timeout option
	timeoutOption *TimeoutOption
	// WebMonitor option
	webMonitorOption *WebMonitorOption
}

// Enable checks if a specific option is enabled.
// Enable 用于查询某个选项是否启用。
func (s *Scheduler) Enable(option string) bool {
	switch option {
	case AliasOptionName:
		if s.schOptions.aliasOption != nil {
			return s.schOptions.aliasOption.Enable()
		}
	case WatchOptionName:
		if s.schOptions.watchOption != nil {
			return s.schOptions.watchOption.Enable()
		}
	case TimoutOptionName:
		if s.schOptions.timeoutOption != nil {
			return s.schOptions.timeoutOption.Enable()
		}
	case WebMonitorOptionName:
		if s.schOptions.webMonitorOption != nil {
			return s.schOptions.webMonitorOption.Enable()
		}
	}
	return false
}

// SchedulerOption is a function that sets options for SchedulerOptions.
// SchedulerOption 是用于设置 SchedulerOptions 的函数类型。
type SchedulerOption func(*SchedulerOptions)

// WithAliasMode sets the alias mode option.
// WithAliasMode 设置别名模式选项。
func WithAliasMode() SchedulerOption {
	return func(s *SchedulerOptions) {
		s.aliasOption = &AliasOption{enabled: true}
	}
}

// WithWatch sets the watch option.
// WithWatch 设置监听选项。
func WithWatch(watchFunc func(event JobWatchInterface)) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.aliasOption = &WatchOption{enabled: true, watchFunc: watchFunc}
	}
}

// WithTimeout sets the watch option.
// WithTimeout 设置监听选项。
func WithTimeout(timeout time.Duration) SchedulerOption {
	return func(s *SchedulerOptions) {
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		s.timeoutOption = &TimeoutOption{enabled: true, timeout: timeout}
	}
}

func WithWebMonitor(address string) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.webMonitorOption = &WebMonitorOption{enabled: true, address: address}
	}
}

// Event represents a job event.
// Event 表示一个任务事件。
type Event struct {
	JobID string // Job ID
	// 任务 ID
	JobName string // Job name
	// 任务名称
	NextRunTime time.Time // Next run time
	// 下次运行时间
	LastTime time.Time // Last run time
	// 上次运行时间
	Err error // Error
	// 错误
}

// Watch starts watching job events.
// Watch 开始监听任务事件。
func (s *Scheduler) Watch() {
	if !s.Enable(WatchOptionName) {
		slog.Error("need watch option")
		return
	}
	event := s.monitor.Watch()
	for {
		select {
		case <-s.ctx.Done():
			return
		case e := <-event:
			fn, ok := s.watchFuncMap[e.GetJobID()]
			if !ok {
				slog.Error("chrono:job not found", "jobID", e.GetJobID())
				continue
			}
			fn(e)
		}
	}
}

// NewScheduler creates a new scheduler.
// NewScheduler 创建一个新的调度器。
func NewScheduler(ctx context.Context, monitor SchedulerMonitor, options ...SchedulerOption) (*Scheduler, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 根据 monitor 是否为空来决定如何创建调度器
	if monitor == nil {
		monitor = newDefaultSchedulerMonitor()
	}
	s, err := gocron.NewScheduler(gocron.WithMonitorStatus(monitor), gocron.WithMonitor(monitor))
	// 错误处理
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to create scheduler: %w", err)
	}
	schOptions := &SchedulerOptions{}
	for _, option := range options {
		option(schOptions)
	}

	return &Scheduler{
		scheduler:    s,
		monitor:      monitor,
		ctx:          ctx,
		watchFuncMap: make(map[string]func(event JobWatchInterface)),
		aliasMap:     make(map[string]string),
		schOptions:   schOptions,
	}, nil
}

// Start starts the scheduler.
// Start 启动调度器。
func (s *Scheduler) Start() {
	if s.Enable(WebMonitorOptionName) {
		if err := NewWebMonitor(s, s.schOptions.webMonitorOption.Address()).Start(); err != nil {
			panic("chrono:failed to start web monitor")
		}
	}
	s.scheduler.Start()
}

// Stop stops the scheduler.
// Stop 停止调度器。
func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}

// RemoveJob removes a job by jobID.
// RemoveJob 通过 jobID 移除任务。
func (s *Scheduler) RemoveJob(jobID string) error {
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("chrono:invalid job ID %s: %w", jobID, err)
	}
	if s.Enable(AliasOptionName) {
		s.removeAlias(jobID)
	}
	if s.Enable(WatchOptionName) {
		s.removeWatchFunc(jobID)
	}
	slog.Info("chrono:job removed", "jobID", jobID)
	return s.scheduler.RemoveJob(jobUUID)
}

// RemoveJobByName removes a job by name.
func (s *Scheduler) RemoveJobByName(name string) error {
	jobs, err := s.GetJobs()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if job.Name() == name {
			return s.RemoveJob(job.ID().String())
		}
	}
	return fmt.Errorf("job with name %s not found", name)
}

// RemoveJobByAlias removes a job by alias.
// RemoveJobByAlias 通过别名移除任务。
func (s *Scheduler) RemoveJobByAlias(alias string) error {
	if !s.Enable(AliasOptionName) {
		return ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return fmt.Errorf("chrono:alias %s not found", alias)
	}
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("chrono:invalid job ID %s: %w", jobID, err)
	}
	s.removeAlias(jobID)
	if s.Enable(WatchOptionName) {
		s.removeWatchFunc(jobID)
	}
	return s.scheduler.RemoveJob(jobUUID)
}

// GetAlias gets alias by jobID.
// GetAlias 通过 jobID 获取别名。
func (s *Scheduler) GetAlias(jobID string) (string, error) {
	if !s.Enable(AliasOptionName) {
		return "", ErrDisEnableAlias
	}
	for alias, realJobID := range s.aliasMap {
		if jobID == realJobID {
			return alias, nil
		}
	}
	return "", ErrFoundAlias
}

// RunJobNow runs a job immediately by jobID.
// RunJobNow 通过 jobID 立即运行任务。
func (s *Scheduler) RunJobNow(jobID string) error {
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return err
	}
	return job.RunNow()
}

// RunJobNowByAlias runs a job immediately by alias.
// RunJobNowByAlias 通过别名立即运行任务。
func (s *Scheduler) RunJobNowByAlias(alias string) error {
	if !s.Enable(AliasOptionName) {
		return ErrDisEnableAlias
	}
	job, err := s.GetJobByAlias(alias)
	if err != nil {
		return err
	}
	return job.RunNow()
}

// addAlias adds an alias for a job.
// addAlias 为任务添加别名。
func (s *Scheduler) addAlias(alias string, jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if alias == "" {
		slog.Warn("chrono:alias is empty", "alias", alias)
		return
	}
	if jobID == "" {
		slog.Warn("chrono:jobID is empty", "jobID", jobID)
		return
	}
	s.aliasMap[alias] = jobID
}

// removeAlias removes an alias.
// removeAlias 移除别名。
func (s *Scheduler) removeAlias(alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.aliasMap[alias]; exists {
		delete(s.aliasMap, alias)
		slog.Info("chrono:alias  removed", "alias", alias)
	} else {
		slog.Warn("chrono:alias not found in aliasMap", "alias", alias)
	}
}

// addWatchFunc adds a watch function for a job.
// addWatchFunc 为任务添加监听函数。
func (s *Scheduler) addWatchFunc(jobID string, fn func(event JobWatchInterface)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if jobID == "" {
		slog.Warn("chrono:jobID is empty", "jobID", jobID)
		return
	}
	if fn == nil {
		slog.Warn("chrono:watchFunc is empty", "jobID", jobID)
		return
	}
	s.watchFuncMap[jobID] = fn
}

// removeWatchFunc removes a watch function for a job.
// removeWatchFunc 移除任务的监听函数。
func (s *Scheduler) removeWatchFunc(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.watchFuncMap[jobID]; exists {
		delete(s.watchFuncMap, jobID)
		slog.Info("chrono:Watch function removed", "jobID", jobID)
	} else {
		slog.Warn("chrono:Job not found in watchFuncMap", "jobID", jobID)
	}
}

// TODO 批量移除任务
// RemoveJobs Removes job list with rollback support.
// func (s *Scheduler) RemoveJobs(jobIDS ...string) error {
// 	// 获取所有需要删除的任务
// 	jobs, err := s.GetJobByIDS(jobIDS...)
// 	if err != nil {
// 		return fmt.Errorf("failed to get jobs: %w", err)
// 	}
//
// 	// 记录成功删除的任务
// 	removedJobs := make([]gocron.Job, 0, len(jobs))
//
// 	// 遍历任务列表，逐个删除任务
// 	for _, job := range jobs {
// 		if err := s.RemoveJob(job.ID().String()); err != nil {
// 			// 如果删除失败，回滚已删除的任务
// 			if rollbackErr := s.rollbackRemovedJobs(removedJobs); rollbackErr != nil {
// 				return fmt.Errorf("failed to remove job %s: %w; rollback failed: %v", job.ID(), err, rollbackErr)
// 			}
// 			return fmt.Errorf("failed to remove job %s: %w", job.ID(), err)
// 		}
// 		// 记录成功删除的任务
// 		removedJobs = append(removedJobs, job)
// 	}
//
// 	return nil
// }
//
// // rollbackRemovedJobs 回滚已删除的任务
// func (s *Scheduler) rollbackRemovedJobs(jobs []gocron.Job) error {
// 	var rollbackErrors []error
//
// 	// 遍历已删除的任务，逐个重新添加
// 	for _, job := range jobs {
// 		if err := s.scheduler.AddJob(job); err != nil {
// 			rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to re-add job %s: %w", job.ID(), err))
// 		}
// 	}
//
// 	// 如果有回滚错误，返回合并后的错误
// 	if len(rollbackErrors) > 0 {
// 		return fmt.Errorf("rollback errors: %v", rollbackErrors)
// 	}
//
// 	return nil
// }

// GetJobs gets all jobs.
// GetJobs 获取所有任务。
func (s *Scheduler) GetJobs() ([]gocron.Job, error) {
	return s.scheduler.Jobs(), nil
}

// GetJobLastTimeByAlias gets the last run time of a job by alias.
// GetJobLastTimeByAlias 通过别名获取任务的最后运行时间。
func (s *Scheduler) GetJobLastTimeByAlias(alias string) (time.Time, error) {
	if !s.Enable(AliasOptionName) {
		return time.Time{}, ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return time.Time{}, ErrFoundAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	lastRun, err := job.LastRun()
	if err != nil {
		return time.Time{}, err
	}
	return lastRun, nil
}

// GetJobLastTime gets the last run time of a job by jobID.
// GetJobLastTime 通过 jobID 获取任务的最后运行时间。
func (s *Scheduler) GetJobLastTime(jobID string) (time.Time, error) {
	if !s.Enable(AliasOptionName) {
		return time.Time{}, ErrDisEnableAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	lastRun, err := job.LastRun()
	if err != nil {
		return time.Time{}, err
	}
	return lastRun, nil
}

// GetJobNextTimeByAlias gets the next run time of a job by alias.
// GetJobNextTimeByAlias 通过别名获取任务的下次运行时间。
func (s *Scheduler) GetJobNextTimeByAlias(alias string) (time.Time, error) {
	if !s.Enable(AliasOptionName) {
		return time.Time{}, ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return time.Time{}, ErrFoundAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return time.Time{}, err
	}
	return nextRun, nil
}

// GetJobNextTime gets the next run time of a job by jobID.
// GetJobNextTime 通过 jobID 获取任务的下次运行时间。
func (s *Scheduler) GetJobNextTime(jobID string) (time.Time, error) {
	if !s.Enable(AliasOptionName) {
		return time.Time{}, ErrDisEnableAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return time.Time{}, err
	}
	return nextRun, nil
}

// GetJobLastAndNextByAlias gets the last and next run times of a job by alias.
// GetJobLastAndNextByAlias 通过别名获取任务的最后和下次运行时间。
func (s *Scheduler) GetJobLastAndNextByAlias(alias string) (time.Time, time.Time, error) {
	if !s.Enable(AliasOptionName) {
		return time.Time{}, time.Time{}, ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return time.Time{}, time.Time{}, ErrFoundAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	lastRun, err := job.LastRun()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return lastRun, nextRun, nil
}

// GetJobLastAndNextByID gets the last and next run times of a job by jobID.
// GetJobLastAndNextByID 通过 jobID 获取任务的最后和下次运行时间。
func (s *Scheduler) GetJobLastAndNextByID(jobID string) (time.Time, time.Time, error) {
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	lastRun, err := job.LastRun()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return lastRun, nextRun, nil
}

// GetJobByName gets a job by its name.
// GetJobByName 通过名称获取任务。
func (s *Scheduler) GetJobByName(jobName string) (gocron.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == jobName {
			return job, nil
		}
	}
	return nil, fmt.Errorf("chrono:job %s not found", jobName)
}

// GetJobByID gets a job by its jobID.
// GetJobByID 通过 jobID 获取任务。
func (s *Scheduler) GetJobByID(jobID string) (gocron.Job, error) {
	for _, job := range s.scheduler.Jobs() {
		if job.ID().String() == jobID {
			return job, nil
		}
	}
	return nil, fmt.Errorf("chrono:job %s not found", jobID)
}

// GetJobByAlias gets a job by its alias.
// GetJobByAlias 通过别名获取任务。
func (s *Scheduler) GetJobByAlias(alias string) (gocron.Job, error) {
	if !s.Enable(AliasOptionName) {
		return nil, ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, fmt.Errorf("chrono:alias %s not found", alias)
	}
	return s.GetJobByID(jobID)
}

// GetJobByIDOrAlias gets a job by jobID or alias, first by jobID, then by alias.
// GetJobByIDOrAlias 通过 jobID 或别名获取任务，优先通过 jobID。
func (s *Scheduler) GetJobByIDOrAlias(identifier string) (gocron.Job, error) {
	// 优先通过ID查找
	if jobID, err := s.GetJobByID(identifier); err == nil {
		return jobID, nil
	}

	// 如果没有找到ID，尝试通过别名查找
	if s.Enable(AliasOptionName) {
		if jobID, exists := s.aliasMap[identifier]; exists {
			return s.GetJobByAlias(jobID)
		}
	}

	return nil, fmt.Errorf("chrono:job with identifier %s not found", identifier)
}

// GetJobByIDS gets jobs by a list of jobIDs.
// GetJobByIDS 通过 jobID 列表获取任务。
func (s *Scheduler) GetJobByIDS(jobIDS ...string) ([]gocron.Job, error) {
	// 创建一个切片用于存储找到的任务
	jobs := make([]gocron.Job, 0, len(jobIDS))

	// 遍历 jobIDS，逐个查找任务
	for _, jobID := range jobIDS {
		job, err := s.GetJobByID(jobID)
		if err != nil {
			return nil, fmt.Errorf("chrono:failed to get job %s: %w", jobID, err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// AddCronJob adds a new cron job.
// AddCronJob 添加一个新的 cron 任务。
func (s *Scheduler) AddCronJob(job *CronJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// check if job has a task function
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	if job.Expr == "" {
		return nil, fmt.Errorf("chrono:job %s has nil expr", job.Name)
	}
	// 优先使用每个Job的Timeout
	if s.Enable(TimoutOptionName) {
		_ = job.Timeout(s.schOptions.timeoutOption.Timeout())
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("ichrono:nvalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false), // 使用 cron 表达式
		gocron.NewTask(job.TaskFunc),    // 任务函数
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add cron job: %w", err)
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watchOption.WatchFunc())
	}
	return jobInstance, nil
}

// AddCronJobs adds a list of new cron jobs.
// AddCronJobs 添加一组新的 cron 任务。
func (s *Scheduler) AddCronJobs(jobs ...*CronJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, cronJob := range jobs {
		cronJobInstance, err := s.AddCronJob(cronJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add cron jobs: %v", errs)
	}
	return jobList, nil
}

// AddCronJobWithOptions adds a cron job with options.
// AddCronJobWithOptions 添加带选项的 cron 任务。
func (s *Scheduler) AddCronJobWithOptions(job *CronJob, options ...gocron.JobOption) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false),
		gocron.NewTask(job.TaskFunc),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add job: %w", err)
	}
	return jobInstance, nil
}

// AddOnceJob adds a new once job.
// AddOnceJob 添加一个新的单次任务。
func (s *Scheduler) AddOnceJob(job *OnceJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// 优先使用每个Job的Timeout
	if s.Enable(TimoutOptionName) && job.timeout <= 0 {
		_ = job.Timeout(s.schOptions.timeoutOption.Timeout())
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(job.WorkTime...)),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add once job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watchOption.WatchFunc())
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddOnceJobs adds a list of new once jobs.
// AddOnceJobs 添加一组新的单次任务。
func (s *Scheduler) AddOnceJobs(jobs ...*OnceJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, onceJob := range jobs {
		cronJobInstance, err := s.AddOnceJob(onceJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add once jobs: %v", errs)
	}
	return jobList, nil
}

// AddOnceJobWithOptions adds a once job with options.
// AddOnceJobWithOptions 添加带选项的单次任务。
func (s *Scheduler) AddOnceJobWithOptions(startAt gocron.OneTimeJobStartAtOption, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.OneTimeJob(startAt),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add once job: %w", err)
	}
	return job, nil
}

// AddIntervalJob adds a new interval job.
// AddIntervalJob 添加一个新的间隔任务。
func (s *Scheduler) AddIntervalJob(job *IntervalJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// 优先使用每个Job的Timeout
	if s.Enable(TimoutOptionName) && job.timeout <= 0 {
		_ = job.Timeout(s.schOptions.timeoutOption.Timeout())
	}
	// Job options
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}

	jobInstance, err := s.scheduler.NewJob(
		gocron.DurationJob(job.Interval),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watchOption.WatchFunc())
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddIntervalJobs adds a list of new interval jobs.
// AddIntervalJobs 添加一组新的间隔任务。
func (s *Scheduler) AddIntervalJobs(jobs ...*IntervalJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, intervalJob := range jobs {
		intervalJobInstance, err := s.AddIntervalJob(intervalJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, intervalJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add interval jobs: %v", errs)
	}
	return jobList, nil
}

// AddIntervalJobWithOptions adds an interval job with options.
// AddIntervalJobWithOptions 添加带选项的间隔任务。
func (s *Scheduler) AddIntervalJobWithOptions(interval time.Duration, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add interval job: %w", err)
	}
	return job, nil
}

// AddDailyJob adds a new daily job.
// AddDailyJob 添加一个新的每日任务。
func (s *Scheduler) AddDailyJob(job *DailyJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrTaskFuncNil
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// 优先使用每个Job的Timeout
	if s.Enable(TimoutOptionName) && job.timeout <= 0 {
		_ = job.Timeout(s.schOptions.timeoutOption.Timeout())
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.DailyJob(job.Interval, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watchOption.WatchFunc())
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddDailyJobs adds a list of new daily jobs.
// AddDailyJobs 添加一组新的每日任务。
func (s *Scheduler) AddDailyJobs(jobs ...*DailyJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, dailyJob := range jobs {
		dailyJobInstance, err := s.AddDailyJob(dailyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add daily jobs: %v", errs)
	}
	return jobList, nil
}

// AddDailyJobWithOptions adds a daily job with options.
// AddDailyJobWithOptions 添加带选项的每日任务。
func (s *Scheduler) AddDailyJobWithOptions(interval uint, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DailyJob(interval, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add daily job: %w", err)
	}
	return job, nil
}

// AddWeeklyJob adds a new weekly job.
// AddWeeklyJob 添加一个新的每周任务。
func (s *Scheduler) AddWeeklyJob(job *WeeklyJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// 优先使用每个Job的Timeout
	if s.Enable(TimoutOptionName) && job.timeout <= 0 {
		_ = job.Timeout(s.schOptions.timeoutOption.Timeout())
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.WeeklyJob(job.Interval, job.DaysOfTheWeek, job.WorkTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add weekly job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watchOption.WatchFunc())
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddWeeklyJobs adds a list of new weekly jobs.
// AddWeeklyJobs 添加一组新的每周任务。
func (s *Scheduler) AddWeeklyJobs(jobs ...*WeeklyJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, weeklyJob := range jobs {
		dailyJobInstance, err := s.AddWeeklyJob(weeklyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add weekly jobs: %v", errs)
	}
	return jobList, nil
}

// AddWeeklyJobWithOptions adds a weekly job with options.
// AddWeeklyJobWithOptions 添加带选项的每周任务。
func (s *Scheduler) AddWeeklyJobWithOptions(interval uint, daysOfTheWeek gocron.Weekdays, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.WeeklyJob(interval, daysOfTheWeek, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add weekly job: %w", err)
	}
	return job, nil
}

// AddMonthlyJob adds a new monthly job.
// AddMonthlyJob 添加一个新的每月任务。
func (s *Scheduler) AddMonthlyJob(job *MonthJob) (gocron.Job, error) {
	if job == nil {
		return nil, ErrInvalidJob
	}
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// 优先使用每个Job的Timeout
	if s.Enable(TimoutOptionName) && job.timeout <= 0 {
		_ = job.Timeout(s.schOptions.timeoutOption.Timeout())
	}
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.MonthlyJob(job.Interval, job.DaysOfTheMonth, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add monthly job: %w", err)
	}
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watchOption.WatchFunc())
	}
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddMonthlyJobs adds a list of new monthly jobs.
// AddMonthlyJobs 添加一组新的每月任务。
func (s *Scheduler) AddMonthlyJobs(jobs ...*MonthJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	for _, monthlyJob := range jobs {
		dailyJobInstance, err := s.AddMonthlyJob(monthlyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add monthly jobs: %v", errs)
	}
	return jobList, nil
}

// AddMonthlyJobWithOptions adds a monthly job with options.
// AddMonthlyJobWithOptions 添加带选项的每月任务。
func (s *Scheduler) AddMonthlyJobWithOptions(interval uint, daysOfTheMonth gocron.DaysOfTheMonth, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.MonthlyJob(interval, daysOfTheMonth, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add monthly job: %w", err)
	}
	return job, nil
}

// OnceJob 表示一个一次性任务的客户端
func (s *Scheduler) OnceJob() OnceJobClientInterface {
	return &OnceJobClient{
		scheduler: s,
		job:       &OnceJob{},
	}
}

// CronJob adds a new cron job.
// CronJob 添加一个新的定时任务
func (s *Scheduler) CronJob() CronJobClientInterface {
	return &CronJobClient{
		scheduler: s,
		job:       &CronJob{},
	}
}

// DailyJob adds a new daily job.
// DailyJob 添加一个新的每日任务
func (s *Scheduler) DailyJob() DailyJobClientInterface {
	return &DailyJobClient{
		scheduler: s,
		job:       &DailyJob{},
	}
}

// IntervalJob 设置间隔Job
func (s *Scheduler) IntervalJob() IntervalJobClientInterface {
	return &IntervalJobClient{
		scheduler: s,
		job:       &IntervalJob{},
	}
}

func (s *Scheduler) WeeklyJob() WeeklyJobClientInterface {
	return &WeeklyJobClient{
		scheduler: s,
		job:       &WeeklyJob{},
	}
}

// Monthly 添加一个新的每月任务
func (s *Scheduler) Monthly() MonthJobClientInterface {
	return &MonthJobClient{
		scheduler: s,
		job:       &MonthJob{},
	}
}
