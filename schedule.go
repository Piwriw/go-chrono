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

var defaultWatch = EmptyWatchFunc

var DefaultScheduler *Scheduler

// InitScheduler initializes the default scheduler.
// InitScheduler 初始化默认调度器。
//
// Parameters:
//
//	ctx     - The context for scheduler lifecycle / 调度器生命周期的上下文
//	monitor - The scheduler monitor / 调度器监控器
//	options - Variable number of scheduler options / 可变数量的调度器选项
//
// Returns:
//
//	error - Error if initialization fails / 如果初始化失败的错误
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
	// Context for scheduler lifecycle
	// 调度器生命周期的上下文
	ctx context.Context
	// Underlying gocron scheduler
	// 底层 gocron 调度器
	scheduler gocron.Scheduler
	// Scheduler monitor
	// 调度器监控器
	monitor SchedulerMonitor
	// Alias to jobID mapping
	// 别名到 jobID 的映射
	aliasMap map[string]string
	// JobID to watch function mapping
	// jobID 到监听函数的映射
	watchFuncMap map[string]func(event JobWatchInterface)
	// Mutex to protect watchFuncMap
	// 用于保护 watchFuncMap 的互斥锁
	mu sync.Mutex
	// Job type map
	// 任务类型映射
	jobTypeMap map[string]JobType
	// Mutex to protect jobTypeMap
	// 用于保护 jobTypeMap 的互斥锁
	jobTypeMu sync.Mutex
	// Scheduler options
	// 调度器选项
	schOptions *SchedulerOptions
}

// SchedulerOptions holds options for the scheduler.
// SchedulerOptions 保存调度器的选项。
type SchedulerOptions struct {
	// 别名选项
	// Alias option
	alias *AliasOption
	// Watch option
	// 监听选项
	watch *WatchOption
	// Timeout option
	timeout *TimeoutOption
	// WebMonitor option
	webMonitor *WebMonitorOption
	// Limit option
	limit *LimitOption
	// Prometheus option
	prometheus *PrometheusOption
}

// Enable checks if a specific option is enabled.
// Enable 用于查询某个选项是否启用。
func (s *Scheduler) Enable(option string) bool {
	switch option {
	case AliasOptionName:
		if s.schOptions.alias != nil {
			return s.schOptions.alias.Enable()
		}
	case WatchOptionName:
		if s.schOptions.watch != nil {
			return s.schOptions.watch.Enable()
		}
	case WebMonitorOptionName:
		if s.schOptions.webMonitor != nil {
			return s.schOptions.webMonitor.Enable()
		}
	case LimitOptionName:
		if s.schOptions.limit != nil {
			return s.schOptions.limit.Enable()
		}
	case PrometheusOptionName:
		if s.schOptions.prometheus != nil {
			return s.schOptions.prometheus.Enable()
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
		s.alias = &AliasOption{enabled: true}
	}
}

// WithWatch sets the watch option.
// WithWatch 设置监听选项。
func WithWatch(watchFunc func(event JobWatchInterface)) SchedulerOption {
	return func(s *SchedulerOptions) {
		if watchFunc != nil {
			s.watch = &WatchOption{enabled: true, watchFunc: watchFunc}
			return
		}
		s.watch = &WatchOption{enabled: true, watchFunc: defaultWatch}
	}
}

// WithWebMonitor sets the web monitor option.
// WithWebMonitor 设置 Web 监控器选项。
//
// Parameters:
//
//	address - The web monitor address / Web监控器地址
//
// Returns:
//
//	SchedulerOption - The scheduler option / 调度器选项
func WithWebMonitor(address string) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.webMonitor = &WebMonitorOption{enabled: true, address: address}
	}
}

// WithLimit sets the limit option.
// WithLimit 设置限制选项。
//
// Parameters:
//
//	limit - The maximum number of jobs / 最大任务数量
//
// Returns:
//
//	SchedulerOption - The scheduler option / 调度器选项
func WithLimit(limit int) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.limit = &LimitOption{enabled: true, number: limit}
	}
}

// WithPrometheus sets the prometheus option.
// need WithWatch first
// WithPrometheus 设置 Prometheus 监控指标选项（需要先启用 WithWatch）。
//
// Parameters:
//
//	address - The prometheus address / Prometheus地址
//
// Returns:
//
//	SchedulerOption - The scheduler option / 调度器选项
func WithPrometheus(address string) SchedulerOption {
	return func(s *SchedulerOptions) {
		s.prometheus = &PrometheusOption{enabled: true, address: address}
	}
}

// Event represents a job event.
// Event 表示一个任务事件。
type Event struct {
	// Job ID
	// 任务 ID
	JobID string
	// Job name
	// 任务名称
	JobName string
	// Next run time
	// 下次运行时间
	NextRunTime time.Time
	// Last run time
	// 上次运行时间
	LastTime time.Time
	// Error
	// 错误
	Err error
}

// GetJobID returns the job ID.
// GetJobID 返回任务ID。
//
// Returns:
//
//	string - The job ID / 任务ID
func (e Event) GetJobID() string {
	return e.JobID
}

// GetJobName returns the job name.
// GetJobName 返回任务名称。
//
// Returns:
//
//	string - The job name / 任务名称
func (e Event) GetJobName() string {
	return e.JobName
}

// GetNextRunTime returns the next run time.
// GetNextRunTime 返回下次运行时间。
//
// Returns:
//
//	time.Time - The next run time / 下次运行时间
func (e Event) GetNextRunTime() time.Time {
	return e.NextRunTime
}

// GetLastTime returns the last run time.
// GetLastTime 返回上次运行时间。
//
// Returns:
//
//	time.Time - The last run time / 上次运行时间
func (e Event) GetLastTime() time.Time {
	return e.LastTime
}

// GetError returns the error.
// GetError 返回错误。
//
// Returns:
//
//	error - The error / 错误
func (e Event) GetError() error {
	return e.Err
}

// Watch starts watching job events.
// Watch 开始监听任务事件。
func (s *Scheduler) Watch() {
	// Check if watch option is enabled
	// 检查是否启用监听选项
	if !s.Enable(WatchOptionName) {
		slog.Error("need watch option")
		return
	}
	// Get event channel from monitor
	// 从监控器获取事件通道
	event := s.monitor.Watch()
	for {
		select {
		// Exit when context is cancelled
		// 当上下文取消时退出
		case <-s.ctx.Done():
			return
		// Process job event
		// 处理任务事件
		case e := <-event:
			jobID := e.GetJobID()
			// Get watch function for this job
			// 获取此任务的监听函数
			fn, ok := s.watchFuncMap[jobID]
			if !ok {
				slog.Error("chrono:job not found", slog.Any("jobID", e.GetJobID()))
				continue
			}
			// Call watch function
			// 调用监听函数
			fn(e)
			jobName := e.GetJobName()
			currentEvent := e.GetCurrentEvent()
			// Update Prometheus metrics if enabled
			// 如果启用Prometheus，则更新指标
			if s.Enable(PrometheusOptionName) {
				s.updatePrometheusMetrics(jobID, jobName, currentEvent)
			}
		}
	}
}

// updatePromJobStatus updates Prometheus job status metrics.
// updatePromJobStatus 更新 Prometheus 任务状态指标。
//
// Parameters:
//
//	jobType - The job type / 任务类型
//	jobName - The job name / 任务名称
//	status  - The job status / 任务状态
func (s *Scheduler) updatePromJobStatus(jobType JobType, jobName string, status gocron.JobStatus) {
	switch status {
	// Increment running count on success
	// 成功时增加运行计数
	case gocron.Success:
		IncJobRunning(jobType, jobName, jobName)
	// Decrement running count on failure
	// 失败时减少运行计数
	case gocron.Fail:
		DecJobRunning(jobType, jobName, jobName)
	}
}

// updatePrometheusJobTime updates Prometheus job execution time metrics.
// updatePrometheusJobTime 更新 Prometheus 任务执行时间指标。
//
// Parameters:
//
//	jobType - The job type / 任务类型
//	jobID   - The job ID / 任务ID
//	jobName - The job name / 任务名称
//	status  - The job status / 任务状态
//	event   - The job event / 任务事件
func (s *Scheduler) updatePrometheusJobTime(jobType JobType, jobID, jobName string, status gocron.JobStatus, event *JobEvent) {
	// Record successful job execution
	// 记录成功的任务执行
	if status == gocron.Success {
		RecordJobExecution(jobType, jobName, jobName, float64(event.GetSpendTime()), true, nil)
		return
	}
	// Record failed job execution
	// 记录失败的任务执行
	RecordJobExecution(jobType, jobID, jobName, float64(event.GetSpendTime()), false, event.GetError())
}

// updatePrometheusMetrics updates Prometheus metrics for a job.
// updatePrometheusMetrics 更新任务的 Prometheus 指标。
//
// Parameters:
//
//	jobID   - The job ID / 任务ID
//	jobName - The job name / 任务名称
//	event   - The job event / 任务事件
func (s *Scheduler) updatePrometheusMetrics(jobID, jobName string, event *JobEvent) {
	// Get job type
	// 获取任务类型
	jobType := s.getJobType(jobID)
	// Get job status
	// 获取任务状态
	status := event.GetStatus()
	// Update job status metrics
	// 更新任务状态指标
	s.updatePromJobStatus(jobType, jobName, status)
	// Update job execution time metrics
	// 更新任务执行时间指标
	s.updatePrometheusJobTime(jobType, jobID, jobName, status, event)
}

// NewScheduler creates a new scheduler.
// NewScheduler 创建一个新的调度器。
//
// Parameters:
//
//	ctx     - The context for scheduler lifecycle / 调度器生命周期的上下文
//	monitor - The scheduler monitor / 调度器监控器
//	options - Variable number of scheduler options / 可变数量的调度器选项
//
// Returns:
//
//	*Scheduler - The created scheduler / 创建的调度器
//	error      - Error if creation fails / 如果创建失败的错误
func NewScheduler(ctx context.Context, monitor SchedulerMonitor, options ...SchedulerOption) (*Scheduler, error) {
	// Use background context if not provided
	// 如果未提供上下文，则使用后台上下文
	if ctx == nil {
		ctx = context.Background()
	}
	// Create default monitor if not provided
	// 如果未提供监控器，则创建默认监控器
	if monitor == nil {
		monitor = newDefaultSchedulerMonitor()
	}
	// Create gocron scheduler with monitor
	// 使用监控器创建 gocron 调度器
	s, err := gocron.NewScheduler(gocron.WithMonitorStatus(monitor), gocron.WithMonitor(monitor))
	// Error handling
	// 错误处理
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to create scheduler: %w", err)
	}
	// Initialize scheduler options
	// 初始化调度器选项
	schOptions := &SchedulerOptions{}
	// Apply options
	// 应用选项
	for _, option := range options {
		option(schOptions)
	}
	return &Scheduler{
		scheduler:    s,
		monitor:      monitor,
		ctx:          ctx,
		watchFuncMap: make(map[string]func(event JobWatchInterface)),
		aliasMap:     make(map[string]string),
		jobTypeMap:   make(map[string]JobType),
		schOptions:   schOptions,
	}, nil
}

// Start starts the scheduler.
// Start 启动调度器。
func (s *Scheduler) Start() {
	// Start web monitor if enabled
	// 如果启用，则启动Web监控器
	if s.Enable(WebMonitorOptionName) {
		if err := NewWebMonitor(s, s.schOptions.webMonitor.Address()).Start(); err != nil {
			panic("chrono:failed to start web monitor")
		}
	}
	// Start Prometheus endpoint if enabled
	// 如果启用，则启动Prometheus端点
	if s.Enable(PrometheusOptionName) {
		StartPrometheusEndpoint(s.schOptions.prometheus.Address())
	}
	// Start the scheduler
	// 启动调度器
	s.scheduler.Start()
}

// Stop stops the scheduler.
// Stop 停止调度器。
//
// Returns:
//
//	error - Error if shutdown fails / 如果关闭失败的错误
func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}

// RemoveJob removes a job by jobID.
// RemoveJob 通过 jobID 移除任务。
//
// Parameters:
//
//	jobID - The job ID to remove / 要移除的任务ID
//
// Returns:
//
//	error - Error if removal fails / 如果移除失败的错误
func (s *Scheduler) RemoveJob(jobID string) error {
	// Check if jobID is empty
	// 检查任务ID是否为空
	if jobID == "" {
		return ErrJobIDNil
	}
	// Parse jobID as UUID
	// 将任务ID解析为UUID
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("chrono:invalid job ID %s: %w", jobID, err)
	}
	// Increment limit if limit option is enabled
	// 如果启用限制选项，则增加限制
	if s.Enable(LimitOptionName) {
		if err := s.incLimit(); err != nil {
			return err
		}
	}
	// Remove alias if alias option is enabled
	// 如果启用别名选项，则移除别名
	if s.Enable(AliasOptionName) {
		s.removeAlias(jobID)
	}
	// Remove watch function if watch option is enabled
	// 如果启用监听选项，则移除监听函数
	if s.Enable(WatchOptionName) {
		s.removeWatchFunc(jobID)
	}
	// Remove job type
	// 移除任务类型
	s.removeJobType(jobID)
	return s.scheduler.RemoveJob(jobUUID)
}

// RemoveJobByName removes a job by name.
// RemoveJobByName 通过名称移除任务。
//
// Parameters:
//
//	name - The job name to remove / 要移除的任务名称
//
// Returns:
//
//	error - Error if removal fails / 如果移除失败的错误
func (s *Scheduler) RemoveJobByName(name string) error {
	// Get all jobs
	// 获取所有任务
	jobs, err := s.GetJobs()
	if err != nil {
		return err
	}
	// Find and remove job by name
	// 按名称查找并移除任务
	for _, job := range jobs {
		if job.Name() == name {
			if err := s.RemoveJob(job.ID().String()); err != nil {
				return err
			}
			s.removeJobType(job.ID().String())
		}
	}
	return fmt.Errorf("job with name %s not found", name)
}

// RemoveJobByAlias removes a job by alias.
// RemoveJobByAlias 通过别名移除任务。
//
// Parameters:
//
//	alias - The alias to remove / 要移除的别名
//
// Returns:
//
//	error - Error if removal fails / 如果移除失败的错误
func (s *Scheduler) RemoveJobByAlias(alias string) error {
	// Check if alias option is enabled
	// 检查是否启用别名选项
	if !s.Enable(AliasOptionName) {
		return ErrDisEnableAlias
	}
	// Increment limit if limit option is enabled
	// 如果启用限制选项，则增加限制
	if s.Enable(LimitOptionName) {
		if err := s.incLimit(); err != nil {
			return err
		}
	}
	// Get jobID by alias
	// 通过别名获取任务ID
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return fmt.Errorf("chrono:alias %s not found", alias)
	}
	// Parse jobID as UUID
	// 将任务ID解析为UUID
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("chrono:invalid job ID %s: %w", jobID, err)
	}
	// Remove alias
	// 移除别名
	s.removeAlias(jobID)
	// Remove job type
	// 移除任务类型
	s.removeJobType(jobID)
	// Remove watch function if watch option is enabled
	// 如果启用监听选项，则移除监听函数
	if s.Enable(WatchOptionName) {
		s.removeWatchFunc(jobID)
	}
	return s.scheduler.RemoveJob(jobUUID)
}

// GetAlias gets alias by jobID.
// GetAlias 通过 jobID 获取别名。
//
// Parameters:
//
//	jobID - The job ID / 任务ID
//
// Returns:
//
//	string - The alias / 别名
//	error  - Error if alias not found / 如果找不到别名的错误
func (s *Scheduler) GetAlias(jobID string) (string, error) {
	// Check if alias option is enabled
	// 检查是否启用别名选项
	if !s.Enable(AliasOptionName) {
		return "", ErrDisEnableAlias
	}
	// Find alias by jobID
	// 通过任务ID查找别名
	for alias, realJobID := range s.aliasMap {
		if jobID == realJobID {
			return alias, nil
		}
	}
	return "", ErrFoundAlias
}

// RunJobNow runs a job immediately by jobID.
// RunJobNow 通过 jobID 立即运行任务。
//
// Parameters:
//
//	jobID - The job ID to run / 要运行的任务ID
//
// Returns:
//
//	error - Error if run fails / 如果运行失败的错误
func (s *Scheduler) RunJobNow(jobID string) error {
	// Get job by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return err
	}
	// Run job immediately
	// 立即运行任务
	return job.RunNow()
}

// RunJobNowByAlias runs a job immediately by alias.
// RunJobNowByAlias 通过别名立即运行任务。
//
// Parameters:
//
//	alias - The alias to run / 要运行的别名
//
// Returns:
//
//	error - Error if run fails / 如果运行失败的错误
func (s *Scheduler) RunJobNowByAlias(alias string) error {
	// Check if alias option is enabled
	// 检查是否启用别名选项
	if !s.Enable(AliasOptionName) {
		return ErrDisEnableAlias
	}
	// Get job by alias
	// 通过别名获取任务
	job, err := s.GetJobByAlias(alias)
	if err != nil {
		return err
	}
	// Run job immediately
	// 立即运行任务
	return job.RunNow()
}

// addJobType adds a job type mapping.
// addJobType 添加任务类型映射。
//
// Parameters:
//
//	jobID   - The job ID / 任务ID
//	jobType - The job type / 任务类型
func (s *Scheduler) addJobType(jobID string, jobType JobType) {
	s.jobTypeMu.Lock()
	defer s.jobTypeMu.Unlock()
	// Check if jobID is empty
	// 检查任务ID是否为空
	if jobID == "" {
		slog.Warn("chrono:jobID is empty", "jobType", jobType)
		return
	}
	s.jobTypeMap[jobID] = jobType
}

// removeJobType removes a job type mapping.
// removeJobType 移除任务类型映射。
//
// Parameters:
//
//	jobID - The job ID to remove / 要移除的任务ID
func (s *Scheduler) removeJobType(jobID string) {
	s.jobTypeMu.Lock()
	defer s.jobTypeMu.Unlock()
	// Check if jobID exists in map
	// 检查任务ID是否存在于映射中
	if _, exists := s.jobTypeMap[jobID]; exists {
		delete(s.jobTypeMap, jobID)
	} else {
		slog.Warn("chrono:jobTypeMap not found in jobTypeMap", "jobID", jobID)
	}
}

// getJobType gets a job type by jobID.
// getJobType 通过 jobID 获取任务类型。
//
// Parameters:
//
//	jobID - The job ID / 任务ID
//
// Returns:
//
//	JobType - The job type / 任务类型
func (s *Scheduler) getJobType(jobID string) JobType {
	s.jobTypeMu.Lock()
	defer s.jobTypeMu.Unlock()
	// Return job type if exists, otherwise return unknown
	// 如果存在则返回任务类型，否则返回未知类型
	if jobType, exists := s.jobTypeMap[jobID]; exists {
		return jobType
	}
	return JobTypeUnknown
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

// CheckLimit checks if the limit is reached.
// CheckLimit 检查是否达到限制。
//
// Returns:
//
//	bool - True if limit is not reached, false otherwise / 如果未达到限制返回true，否则返回false
func (s *Scheduler) CheckLimit() bool {
	// Decrease limit
	// 减少限制
	if err := s.decLimit(); err != nil {
		return false
	}
	return s.schOptions.limit.number >= 0
}

// incLimit increases the limit.
// incLimit 增加限制。
//
// Returns:
//
//	error - Error if limit option is disabled / 如果限制选项未启用的错误
func (s *Scheduler) incLimit() error {
	// Check if limit option is enabled
	// 检查是否启用限制选项
	if !s.Enable(LimitOptionName) {
		return ErrDisEnableLimit
	}
	s.schOptions.limit.number++
	return nil
}

// decLimit decreases the limit.
// decLimit 减少限制。
//
// Returns:
//
//	error - Error if limit option is disabled / 如果限制选项未启用的错误
func (s *Scheduler) decLimit() error {
	// Check if limit option is enabled
	// 检查是否启用限制选项
	if !s.Enable(LimitOptionName) {
		return ErrDisEnableLimit
	}
	s.schOptions.limit.number--
	return nil
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
//
// Returns:
//
//	[]gocron.Job - The list of jobs / 任务列表
//	error         - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobs() ([]gocron.Job, error) {
	return s.scheduler.Jobs(), nil
}

// GetJobLastTimeByAlias gets the last run time of a job by alias.
// GetJobLastTimeByAlias 通过别名获取任务的最后运行时间。
//
// Parameters:
//
//	alias - The job alias / 任务别名
//
// Returns:
//
//	time.Time - The last run time / 最后运行时间
//	error     - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobLastTimeByAlias(alias string) (time.Time, error) {
	// Check if alias option is enabled
	// 检查是否启用别名选项
	if !s.Enable(AliasOptionName) {
		return time.Time{}, ErrDisEnableAlias
	}
	// Get jobID by alias
	// 通过别名获取任务ID
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return time.Time{}, ErrFoundAlias
	}
	// Get job by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	// Get last run time
	// 获取最后运行时间
	lastRun, err := job.LastRun()
	if err != nil {
		return time.Time{}, err
	}
	return lastRun, nil
}

// GetJobLastTime gets the last run time of a job by jobID.
// GetJobLastTime 通过 jobID 获取任务的最后运行时间。
//
// Parameters:
//
//	jobID - The job ID / 任务ID
//
// Returns:
//
//	time.Time - The last run time / 最后运行时间
//	error     - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobLastTime(jobID string) (time.Time, error) {
	// Get job by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	// Get last run time
	// 获取最后运行时间
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
//
// Parameters:
//
//	jobID - The job ID / 任务ID
//
// Returns:
//
//	time.Time - The next run time / 下次运行时间
//	error     - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobNextTime(jobID string) (time.Time, error) {
	// Get job by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return time.Time{}, err
	}
	// Get next run time
	// 获取下次运行时间
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
//
// Parameters:
//
//	jobName - The job name / 任务名称
//
// Returns:
//
//	gocron.Job - The found job / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByName(jobName string) (gocron.Job, error) {
	// Search through all jobs
	// 遍历所有任务
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == jobName {
			return job, nil
		}
	}
	return nil, fmt.Errorf("chrono:job %s not found", jobName)
}

// GetJobByID gets a job by its jobID.
// GetJobByID 通过 jobID 获取任务。
//
// Parameters:
//
//	jobID - The job ID / 任务ID
//
// Returns:
//
//	gocron.Job - The found job / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByID(jobID string) (gocron.Job, error) {
	// Search through all jobs
	// 遍历所有任务
	for _, job := range s.scheduler.Jobs() {
		if job.ID().String() == jobID {
			return job, nil
		}
	}
	return nil, fmt.Errorf("chrono:job %s not found", jobID)
}

// GetJobByAlias gets a job by its alias.
// GetJobByAlias 通过别名获取任务。
//
// Parameters:
//
//	alias - The job alias / 任务别名
//
// Returns:
//
//	gocron.Job - The found job / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByAlias(alias string) (gocron.Job, error) {
	// Check if alias option is enabled
	// 检查是否启用别名选项
	if !s.Enable(AliasOptionName) {
		return nil, ErrDisEnableAlias
	}
	// Get jobID by alias
	// 通过别名获取任务ID
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, fmt.Errorf("chrono:alias %s not found", alias)
	}
	return s.GetJobByID(jobID)
}

// GetJobByIDOrAlias gets a job by jobID or alias, first by jobID, then by alias.
// GetJobByIDOrAlias 通过 jobID 或别名获取任务，优先通过 jobID。
//
// Parameters:
//
//	identifier - The job ID or alias / 任务ID或别名
//
// Returns:
//
//	gocron.Job - The found job / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByIDOrAlias(identifier string) (gocron.Job, error) {
	// Try ID lookup first
	// 优先尝试ID查找
	if jobID, err := s.GetJobByID(identifier); err == nil {
		return jobID, nil
	}

	// Fallback to alias lookup if alias feature is enabled
	// 如果启用别名功能，回退到别名查找
	if s.Enable(AliasOptionName) {
		if jobID, exists := s.aliasMap[identifier]; exists {
			return s.GetJobByAlias(jobID)
		}
	}

	return nil, fmt.Errorf("chrono:job with identifier %s not found", identifier)
}

// GetJobByIDS gets jobs by a list of jobIDs.
// GetJobByIDS 通过 jobID 列表获取任务。
//
// Parameters:
//
//	jobIDS - Variable number of job IDs / 可变数量的任务ID
//
// Returns:
//
//	[]gocron.Job - The list of found jobs / 找到的任务列表
//	error          - Error if any job is not found / 如果任何任务找不到的错误
func (s *Scheduler) GetJobByIDS(jobIDS ...string) ([]gocron.Job, error) {
	// Create a slice to store found jobs
	// 创建一个切片用于存储找到的任务
	jobs := make([]gocron.Job, 0, len(jobIDS))

	// Iterate through jobIDS and find each job
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
//
// Parameters:
//
//	job - The cron job to add / 要添加的cron任务
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddCronJob(job *CronJob) (gocron.Job, error) {
	// Check if job is nil
	// 检查任务是否为空
	if job == nil {
		return nil, ErrInvalidJob
	}
	// Check if job has error
	// 检查任务是否有错误
	if job.err != nil {
		return nil, job.err
	}
	// check if job has a task function
	// 检查任务是否有任务函数
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// Check if cron expression is set
	// 检查是否设置了cron表达式
	if job.Expr == "" {
		return nil, fmt.Errorf("chrono:job %s has nil expr", job.Name)
	}
	// Check limit if limit option is enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(LimitOptionName) {
		if !s.CheckLimit() {
			return nil, ErrMoreLimit
		}
	}
	// Prepare job options
	// 准备任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	// Set job ID if provided
	// 如果提供了任务ID，则设置
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("ichrono:nvalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create cron job
	// 创建cron任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false), // 使用 cron 表达式
		gocron.NewTask(job.TaskFunc),    // 任务函数
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add cron job: %w", err)
	}
	// Add job type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), job.Type)
	// Add alias if alias option is enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	// Add watch function if watch option is enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc())
	}
	return jobInstance, nil
}

// AddCronJobs adds a list of new cron jobs.
// AddCronJobs 添加一组新的 cron 任务。
//
// Parameters:
//
//	jobs - Variable number of cron jobs to add / 可变数量的要添加的cron任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error          - Error if any addition fails / 如果任何添加失败的错误
func (s *Scheduler) AddCronJobs(jobs ...*CronJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	// Add each cron job
	// 逐个添加cron任务
	for _, cronJob := range jobs {
		cronJobInstance, err := s.AddCronJob(cronJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	// Return error if any job failed to add
	// 如果任何任务添加失败，则返回错误
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
//
// Parameters:
//
//	job - The once job to add / 要添加的单次任务
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddOnceJob(job *OnceJob) (gocron.Job, error) {
	// Check if job is nil
	// 检查任务是否为空
	if job == nil {
		return nil, ErrInvalidJob
	}
	// Check if job has error
	// 检查任务是否有错误
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// Check limit if limit option is enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(LimitOptionName) {
		if !s.CheckLimit() {
			return nil, ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name), gocron.WithTags(job.Tag...))
	// Set job ID if provided
	// 如果提供了任务ID，则设置
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create once job
	// 创建单次任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(job.WorkTime...)),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add once job: %w", err)
	}
	// Add job type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), job.Type)
	// Add watch function if watch option is enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc())
	}
	// Add alias if alias option is enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddOnceJobs adds a list of new once jobs.
// AddOnceJobs 添加一组新的单次任务。
//
// Parameters:
//
//	jobs - Variable number of once jobs to add / 可变数量的要添加的单次任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error          - Error if any addition fails / 如果任何添加失败的错误
func (s *Scheduler) AddOnceJobs(jobs ...*OnceJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	// Add each once job
	// 逐个添加单次任务
	for _, onceJob := range jobs {
		cronJobInstance, err := s.AddOnceJob(onceJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	// Return error if any job failed to add
	// 如果任何任务添加失败，则返回错误
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
//
// Parameters:
//
//	job - The interval job to add / 要添加的间隔任务
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddIntervalJob(job *IntervalJob) (gocron.Job, error) {
	// Check if job is nil
	// 检查任务是否为空
	if job == nil {
		return nil, ErrInvalidJob
	}
	// Check if job has error
	// 检查任务是否有错误
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// Check limit if limit option is enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(LimitOptionName) {
		if !s.CheckLimit() {
			return nil, ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	// Set job ID if provided
	// 如果提供了任务ID，则设置
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}

	// Create interval job
	// 创建间隔任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.DurationJob(job.Interval),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add job: %w", err)
	}
	// Add job type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), job.Type)
	// Add watch function if watch option is enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc())
	}
	// Add alias if alias option is enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddIntervalJobs adds a list of new interval jobs.
// AddIntervalJobs 添加一组新的间隔任务。
//
// Parameters:
//
//	jobs - Variable number of interval jobs to add / 可变数量的要添加的间隔任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error          - Error if any addition fails / 如果任何添加失败的错误
func (s *Scheduler) AddIntervalJobs(jobs ...*IntervalJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	// Add each interval job
	// 逐个添加间隔任务
	for _, intervalJob := range jobs {
		intervalJobInstance, err := s.AddIntervalJob(intervalJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, intervalJobInstance)
	}
	// Return error if any job failed to add
	// 如果任何任务添加失败，则返回错误
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
//
// Parameters:
//
//	job - The daily job to add / 要添加的每日任务
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddDailyJob(job *DailyJob) (gocron.Job, error) {
	// Check if job is nil
	// 检查任务是否为空
	if job == nil {
		return nil, ErrTaskFuncNil
	}
	// Check if job has error
	// 检查任务是否有错误
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// Check limit if limit option is enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(LimitOptionName) {
		if !s.CheckLimit() {
			return nil, ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	// Set job ID if provided
	// 如果提供了任务ID，则设置
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create daily job
	// 创建每日任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.DailyJob(job.Interval, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add job: %w", err)
	}
	// Add job type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), job.Type)
	// Add watch function if watch option is enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc())
	}
	// Add alias if alias option is enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddDailyJobs adds a list of new daily jobs.
// AddDailyJobs 添加一组新的每日任务。
//
// Parameters:
//
//	jobs - Variable number of daily jobs to add / 可变数量的要添加的每日任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error          - Error if any addition fails / 如果任何添加失败的错误
func (s *Scheduler) AddDailyJobs(jobs ...*DailyJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	// Add each daily job
	// 逐个添加每日任务
	for _, dailyJob := range jobs {
		dailyJobInstance, err := s.AddDailyJob(dailyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	// Return error if any job failed to add
	// 如果任何任务添加失败，则返回错误
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
//
// Parameters:
//
//	job - The weekly job to add / 要添加的每周任务
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddWeeklyJob(job *WeeklyJob) (gocron.Job, error) {
	// Check if job is nil
	// 检查任务是否为空
	if job == nil {
		return nil, ErrInvalidJob
	}
	// Check if job has error
	// 检查任务是否有错误
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// Check limit if limit option is enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(LimitOptionName) {
		if !s.CheckLimit() {
			return nil, ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	// Set job ID if provided
	// 如果提供了任务ID，则设置
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create weekly job
	// 创建每周任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.WeeklyJob(job.Interval, job.DaysOfTheWeek, job.WorkTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add weekly job: %w", err)
	}
	// Add job type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), job.Type)
	// Add watch function if watch option is enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc())
	}
	// Add alias if alias option is enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(AliasOptionName) {
		s.addAlias(job.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddWeeklyJobs adds a list of new weekly jobs.
// AddWeeklyJobs 添加一组新的每周任务。
//
// Parameters:
//
//	jobs - Variable number of weekly jobs to add / 可变数量的要添加的每周任务
//
// Returns:
//
//	[]gocron.Job - The list of added jobs / 添加的任务列表
//	error          - Error if any addition fails / 如果任何添加失败的错误
func (s *Scheduler) AddWeeklyJobs(jobs ...*WeeklyJob) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(jobs))
	// Add each weekly job
	// 逐个添加每周任务
	for _, weeklyJob := range jobs {
		weeklyJobInstance, err := s.AddWeeklyJob(weeklyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, weeklyJobInstance)
	}
	// Return error if any job failed to add
	// 如果任何任务添加失败，则返回错误
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
//
// Parameters:
//
//	job - The monthly job to add / 要添加的每月任务
//
// Returns:
//
//	gocron.Job - The added job / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddMonthlyJob(job *MonthJob) (gocron.Job, error) {
	// Check if job is nil
	// 检查任务是否为空
	if job == nil {
		return nil, ErrInvalidJob
	}
	// Check if job has error
	// 检查任务是否有错误
	if job.err != nil {
		return nil, job.err
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:job %s has no task function", job.Name)
	}
	// Check limit if limit option is enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(LimitOptionName) {
		if !s.CheckLimit() {
			return nil, ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(job.Hooks...), gocron.WithName(job.Name))
	// Set job ID if provided
	// 如果提供了任务ID，则设置
	if job.ID != "" {
		jobID, err := uuid.Parse(job.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid job ID %s: %w", job.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create monthly job
	// 创建每月任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.MonthlyJob(job.Interval, job.DaysOfTheMonth, job.AtTimes),
		gocron.NewTask(job.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add monthly job: %w", err)
	}
	// Add job type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), job.Type)
	// Add watch function if watch option is enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(WatchOptionName) {
		if job.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), job.WatchFunc)
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc())
	}
	// Add alias if alias option is enabled
	// 如果启用别名选项，则添加别名
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

// OnceJob creates a new once job client.
// OnceJob 创建一个新的单次任务客户端。
//
// Returns:
//
//	OnceJobClientInterface - The once job client interface / 单次任务客户端接口
func (s *Scheduler) OnceJob() OnceJobClientInterface {
	return &OnceJobClient{
		scheduler: s,
		job: &OnceJob{
			Type: JobTypeOnce,
		},
	}
}

// CronJob creates a new cron job client.
// CronJob 创建一个新的定时任务客户端。
//
// Returns:
//
//	CronJobClientInterface - The cron job client interface / 定时任务客户端接口
func (s *Scheduler) CronJob() CronJobClientInterface {
	return &CronJobClient{
		scheduler: s,
		job: &CronJob{
			Type: JobTypeCron,
		},
	}
}

// DailyJob creates a new daily job client.
// DailyJob 创建一个新的每日任务客户端。
//
// Returns:
//
//	DailyJobClientInterface - The daily job client interface / 每日任务客户端接口
func (s *Scheduler) DailyJob() DailyJobClientInterface {
	return &DailyJobClient{
		scheduler: s,
		job: &DailyJob{
			Type: JobTypeDaily,
		},
	}
}

// IntervalJob creates a new interval job client.
// IntervalJob 创建一个新的间隔任务客户端。
//
// Returns:
//
//	IntervalJobClientInterface - The interval job client interface / 间隔任务客户端接口
func (s *Scheduler) IntervalJob() IntervalJobClientInterface {
	return &IntervalJobClient{
		scheduler: s,
		job: &IntervalJob{
			Type: JobInterval,
		},
	}
}

// WeeklyJob creates a new weekly job client.
// WeeklyJob 创建一个新的每周任务客户端。
//
// Returns:
//
//	WeeklyJobClientInterface - The weekly job client interface / 每周任务客户端接口
func (s *Scheduler) WeeklyJob() WeeklyJobClientInterface {
	return &WeeklyJobClient{
		scheduler: s,
		job: &WeeklyJob{
			Type: JobTypeWeekly,
		},
	}
}

// Monthly creates a new monthly job client.
// Monthly 创建一个新的每月任务客户端。
//
// Returns:
//
//	MonthJobClientInterface - The monthly job client interface / 每月任务客户端接口
func (s *Scheduler) Monthly() MonthJobClientInterface {
	return &MonthJobClient{
		scheduler: s,
		job: &MonthJob{
			Type: JobTypeMonthly,
		},
	}
}
