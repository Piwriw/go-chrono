package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/piwriw/go-chrono/monitor"
	"github.com/piwriw/go-chrono/retry"
)

var defaultWatch = EmptyWatchFuncMonitor

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
func InitScheduler(ctx context.Context, schedMonitor monitor.SchedulerMonitor, options ...SchedulerOption) error {
	var err error
	DefaultScheduler, err = NewScheduler(ctx, schedMonitor, options...)
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
	schedMonitor monitor.SchedulerMonitor
	// Alias to jobID mapping
	// 别名到 jobID 的映射
	aliasMap map[string]string
	// JobID to watch function mapping
	// jobID 到监听函数的映射
	watchFuncMap map[string]func(event monitor.JobWatchInterface)
	// Mutex to protect watchFuncMap
	// 用于保护 watchFuncMap 的互斥锁
	mu sync.Mutex
	// Job type map
	// 任务类型映射
	jobTypeMap map[string]jobs.JobType
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
	alias *common.AliasOption
	// Watch option
	// 监听选项
	watch *common.WatchOption
	// Timeout option
	timeout *common.TimeoutOption
	// WebMonitor option
	webMonitor *common.WebMonitorOption
	// Limit option
	limit *common.LimitOption
	// Prometheus option
	prometheus *common.PrometheusOption
	// Retry option
	// 重试选项
	retry *common.RetryOption
}

// Enable checks if a specific option is Enabled.
// Enable 用于查询某个选项是否启用。
func (s *Scheduler) Enable(option string) bool {
	switch option {
	case common.AliasOptionName:
		if s.schOptions.alias != nil {
			return s.schOptions.alias.Enable()
		}
	case common.WatchOptionName:
		if s.schOptions.watch != nil {
			return s.schOptions.watch.Enable()
		}
	case common.WebMonitorOptionName:
		if s.schOptions.webMonitor != nil {
			return s.schOptions.webMonitor.Enable()
		}
	case common.LimitOptionName:
		if s.schOptions.limit != nil {
			return s.schOptions.limit.Enable()
		}
	case common.PrometheusOptionName:
		if s.schOptions.prometheus != nil {
			return s.schOptions.prometheus.Enable()
		}
	case common.RetryOptionName:
		if s.schOptions.retry != nil {
			return s.schOptions.retry.Enable()
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
		s.alias = &common.AliasOption{Enabled: true}
	}
}

// WithWatch sets the watch option.
// WithWatch 设置监听选项。
func WithWatch(watchFunc func(event monitor.JobWatchInterface)) SchedulerOption {
	return func(s *SchedulerOptions) {
		if watchFunc != nil {
			s.watch = &common.WatchOption{Enabled: true, WatchFunc: watchFunc}
			return
		}
		s.watch = &common.WatchOption{Enabled: true, WatchFunc: defaultWatch}
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
		s.webMonitor = &common.WebMonitorOption{Enabled: true, Address: address}
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
		s.limit = &common.LimitOption{Enabled: true, Limit: limit}
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
		s.prometheus = &common.PrometheusOption{Enabled: true, Address: address}
	}
}

// WithRetry sets the retry configuration for all jobs.
// WithRetry 为所有任务设置重试配置。
//
// Parameters:
//
//	config - The retry configuration / 重试配置
//
// Returns:
//
//	SchedulerOption - The scheduler option / 调度器选项
func WithRetry(config *retry.RetryConfig) SchedulerOption {
	return func(s *SchedulerOptions) {
		retryOpt := &common.RetryOption{
			Enabled: config != nil && config.MaxRetries > 0,
			Config:  config,
		}
		s.retry = retryOpt
	}
}

// Event represents a jobs event.
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
	NextRunTime *time.Time
	// Last run time
	// 上次运行时间
	LastTime *time.Time
	// Error
	// 错误
	Err error
}

// GetJobID returns the jobs ID.
// GetJobID 返回任务ID。
//
// Returns:
//
//	string - The jobs ID / 任务ID
func (e Event) GetJobID() string {
	return e.JobID
}

// GetJobName returns the jobs name.
// GetJobName 返回任务名称。
//
// Returns:
//
//	string - The jobs name / 任务名称
func (e Event) GetJobName() string {
	return e.JobName
}

// GetNextRunTime returns the next run time.
// GetNextRunTime 返回下次运行时间。
//
// Returns:
//
//	time.Time - The next run time / 下次运行时间
func (e Event) GetNextRunTime() *time.Time {
	return e.NextRunTime
}

// GetLastTime returns the last run time.
// GetLastTime 返回上次运行时间。
//
// Returns:
//
//	time.Time - The last run time / 上次运行时间
func (e Event) GetLastTime() *time.Time {
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

// Watch starts watching jobs events.
// Watch 开始监听任务事件。
func (s *Scheduler) Watch() {
	// Check if watch option is Enabled
	// 检查是否启用监听选项
	if !s.Enable(common.WatchOptionName) {
		slog.Error("need watch option")
		return
	}
	// Get event channel from monitor
	// 从监控器获取事件通道
	event := s.schedMonitor.Watch()
	for {
		select {
		// Exit when context is cancelled
		// 当上下文取消时退出
		case <-s.ctx.Done():
			return
		// Process jobs event
		// 处理任务事件
		case e := <-event:
			jobID := e.GetJobID()
			// Get watch function for this jobs
			// 获取此任务的监听函数
			fn, ok := s.watchFuncMap[jobID]
			if !ok {
				slog.Error("chrono:jobs not found", slog.Any("jobID", e.GetJobID()))
				continue
			}
			// Call watch function
			// 调用监听函数
			fn(e)
			jobName := e.GetJobName()
			currentEvent := e.GetCurrentEvent()
			// Update Prometheus metrics if Enabled
			// 如果启用Prometheus，则更新指标
			if s.Enable(common.PrometheusOptionName) {
				s.updatePrometheusMetrics(jobID, jobName, currentEvent)
			}
		}
	}
}

// WatchFuncAdapter converts a pkg.JobWatchInterface function to monitor.JobWatchInterface
// WatchFuncAdapter 将 pkg.JobWatchInterface 函数转换为 monitor.JobWatchInterface
func WatchFuncAdapter(jobFunc func(event common.JobWatchInterface)) func(event monitor.JobWatchInterface) {
	if jobFunc == nil {
		return nil
	}
	return func(event monitor.JobWatchInterface) {
		// Convert monitor.JobWatchInterface to pkg.JobWatchInterface
		// This is a shallow adapter - in a real scenario, you might need a more robust conversion
		jobFunc(&jobWatchAdapter{event})
	}
}

// jobWatchAdapter adapts monitor.JobWatchInterface to pkg.JobWatchInterface
// jobWatchAdapter 将 monitor.JobWatchInterface 适配为 pkg.JobWatchInterface
type jobWatchAdapter struct {
	event monitor.JobWatchInterface
}

// GetJobID gets the jobs ID.
// 获取任务 ID。
//
// Returns:
//
//	string - The jobs ID / 任务 ID
func (a *jobWatchAdapter) GetJobID() string {
	return a.event.GetJobID()
}

// GetJobName gets the jobs name.
// 获取任务名称。
//
// Returns:
//
//	string - The jobs name / 任务名称
func (a *jobWatchAdapter) GetJobName() string {
	return a.event.GetJobName()
}

// GetStartTime gets the jobs start time.
// 获取任务开始时间。
//
// Returns:
//
//	*time.Time - The start time / 开始时间
func (a *jobWatchAdapter) GetStartTime() *time.Time {
	return a.event.GetStartTime()
}

// GetEndTime gets the jobs end time.
// 获取任务结束时间。
//
// Returns:
//
//	*time.Time - The end time / 结束时间
func (a *jobWatchAdapter) GetEndTime() *time.Time {
	return a.event.GetEndTime()
}

// GetStatus gets the jobs status.
// 获取任务状态。
//
// Returns:
//
//	int - The jobs status / 任务状态
func (a *jobWatchAdapter) GetStatus() int {
	// Convert string status to int - this is a simple mapping
	status := a.event.GetStatus()
	switch status {
	case gocron.Success:
		return 0
	case gocron.Fail:
		return 1
	default:
		return 2 // Unknown/running status
	}
}

// GetTags gets the jobs tags.
// 获取任务标签。
//
// Returns:
//
//	[]string - The jobs tags / 任务标签
func (a *jobWatchAdapter) GetTags() []string {
	return a.event.GetTags()
}

// Error gets the jobs error.
// 获取任务错误。
//
// Returns:
//
//	error - The jobs error / 任务错误
func (a *jobWatchAdapter) Error() error {
	return a.event.Error()
}

// GetCurrentEvent gets the current jobs event.
// 获取当前任务事件。
//
// Returns:
//
//	*pkg.JobEvent - The current jobs event / 当前任务事件
func (a *jobWatchAdapter) GetCurrentEvent() *common.JobEvent {
	// This is a limitation - we can't perfectly convert without a full implementation
	// For now, return nil or create a basic conversion
	monitorEvent := a.event.GetCurrentEvent()
	if monitorEvent == nil {
		return nil
	}
	return &common.JobEvent{
		EventID:         monitorEvent.EventID,
		StartTime:       monitorEvent.StartTime,
		EndTime:         monitorEvent.EndTime,
		Status:          a.GetStatus(),
		Err:             monitorEvent.Err,
		RetryCount:      monitorEvent.RetryCount,
		IsRetry:         monitorEvent.IsRetry,
		OriginalEventID: monitorEvent.OriginalEventID,
	}
}

// updatePromJobStatus updates Prometheus jobs status metrics.
// updatePromJobStatus 更新 Prometheus 任务状态指标。
//
// Parameters:
//
//	jobType - The jobs type / 任务类型
//	jobName - The jobs name / 任务名称
//	status  - The jobs status / 任务状态
func (s *Scheduler) updatePromJobStatus(jobType jobs.JobType, jobName string, status gocron.JobStatus) {
	switch status {
	// Increment running count on success
	// 成功时增加运行计数
	case gocron.Success:
		monitor.IncJobRunning(jobType, jobName, jobName)
	// Decrement running count on failure
	// 失败时减少运行计数
	case gocron.Fail:
		monitor.DecJobRunning(jobType, jobName, jobName)
	}
}

// updatePrometheusJobTime updates Prometheus jobs execution time metrics.
// updatePrometheusJobTime 更新 Prometheus 任务执行时间指标。
//
// Parameters:
//
//	jobType - The jobs type / 任务类型
//	jobID   - The jobs ID / 任务ID
//	jobName - The jobs name / 任务名称
//	status  - The jobs status / 任务状态
//	event   - The jobs event / 任务事件
func (s *Scheduler) updatePrometheusJobTime(jobType jobs.JobType, jobID, jobName string, status gocron.JobStatus, event *monitor.JobEvent) {
	// Record successful jobs execution
	// 记录成功的任务执行
	if status == gocron.Success {
		monitor.RecordJobExecution(jobType, jobName, jobName, float64(event.GetSpendTime()), true, nil)
		return
	}
	// Record failed jobs execution
	// 记录失败的任务执行
	monitor.RecordJobExecution(jobType, jobID, jobName, float64(event.GetSpendTime()), false, event.GetError())
}

// updatePrometheusMetrics updates Prometheus metrics for a jobs.
// updatePrometheusMetrics 更新任务的 Prometheus 指标。
//
// Parameters:
//
//	jobID   - The jobs ID / 任务ID
//	jobName - The jobs name / 任务名称
//	event   - The jobs event / 任务事件
func (s *Scheduler) updatePrometheusMetrics(jobID, jobName string, event *monitor.JobEvent) {
	// Get jobs type
	// 获取任务类型
	jobType := s.getJobType(jobID)
	// Get jobs status
	// 获取任务状态
	status := event.GetStatus()
	// Update jobs status metrics
	// 更新任务状态指标
	s.updatePromJobStatus(jobType, jobName, status)
	// Update jobs execution time metrics
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
func NewScheduler(ctx context.Context, schedMonitor monitor.SchedulerMonitor, options ...SchedulerOption) (*Scheduler, error) {
	// Use background context if not provided
	// 如果未提供上下文，则使用后台上下文
	if ctx == nil {
		ctx = context.Background()
	}
	// Create default monitor if not provided
	// 如果未提供监控器，则创建默认监控器
	if schedMonitor == nil {
		schedMonitor = monitor.NewDefaultSchedulerMonitor()
	}
	// Create gocron scheduler with monitor
	// 使用监控器创建 gocron 调度器
	s, err := gocron.NewScheduler(gocron.WithMonitorStatus(schedMonitor), gocron.WithMonitor(schedMonitor))
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
		schedMonitor: schedMonitor,
		ctx:          ctx,
		watchFuncMap: make(map[string]func(event monitor.JobWatchInterface)),
		aliasMap:     make(map[string]string),
		jobTypeMap:   make(map[string]jobs.JobType),
		schOptions:   schOptions,
	}, nil
}

// Start starts the scheduler.
// Start 启动调度器。
func (s *Scheduler) Start() {
	// Start web monitor if Enabled
	// 如果启用，则启动Web监控器
	if s.Enable(common.WebMonitorOptionName) {
		if err := monitor.NewWebMonitor(s, s.schedMonitor, s.schOptions.webMonitor.Address).Start(); err != nil {
			panic("chrono:failed to start web monitor")
		}
	}
	// Start Prometheus endpoint if Enabled
	// 如果启用，则启动Prometheus端点
	if s.Enable(common.PrometheusOptionName) {
		monitor.StartPrometheusEndpoint(s.schOptions.prometheus.Address)
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

// RemoveJob removes a jobs by jobID.
// RemoveJob 通过 jobID 移除任务。
//
// Parameters:
//
//	jobID - The jobs ID to remove / 要移除的任务ID
//
// Returns:
//
//	error - Error if removal fails / 如果移除失败的错误
func (s *Scheduler) RemoveJob(jobID string) error {
	// Check if jobID is empty
	// 检查任务ID是否为空
	if jobID == "" {
		return common.ErrJobIDNil
	}
	// Parse jobID as UUID
	// 将任务ID解析为UUID
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		return fmt.Errorf("chrono:invalid jobs ID %s: %w", jobID, err)
	}
	// Increment limit if limit option is Enabled
	// 如果启用限制选项，则增加限制
	if s.Enable(common.LimitOptionName) {
		if err := s.incLimit(); err != nil {
			return err
		}
	}
	// Remove alias if alias option is Enabled
	// 如果启用别名选项，则移除别名
	if s.Enable(common.AliasOptionName) {
		s.removeAliasByJobID(jobID)
	}
	// Remove watch function if watch option is Enabled
	// 如果启用监听选项，则移除监听函数
	if s.Enable(common.WatchOptionName) {
		s.removeWatchFunc(jobID)
	}
	// Remove jobs type
	// 移除任务类型
	s.removeJobType(jobID)
	return s.scheduler.RemoveJob(jobUUID)
}

// RemoveJobByName removes a jobs by name.
// RemoveJobByName 通过名称移除任务。
//
// Parameters:
//
//	name - The jobs name to remove / 要移除的任务名称
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
	// Find and remove jobs by name
	// 按名称查找并移除任务
	for _, job := range jobs {
		if job.Name() == name {
			if err := s.RemoveJob(job.ID().String()); err != nil {
				return err
			}
			s.removeJobType(job.ID().String())
		}
	}
	return fmt.Errorf("jobs with name %s not found", name)
}

// RemoveJobByAlias removes a jobs by alias.
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
	// Check if alias option is Enabled
	// 检查是否启用别名选项
	if !s.Enable(common.AliasOptionName) {
		return common.ErrDisEnableAlias
	}
	// Increment limit if limit option is Enabled
	// 如果启用限制选项，则增加限制
	if s.Enable(common.LimitOptionName) {
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
		return fmt.Errorf("chrono:invalid jobs ID %s: %w", jobID, err)
	}
	// Remove alias
	// 移除别名
	s.removeAlias(jobID)
	// Remove jobs type
	// 移除任务类型
	s.removeJobType(jobID)
	// Remove watch function if watch option is Enabled
	// 如果启用监听选项，则移除监听函数
	if s.Enable(common.WatchOptionName) {
		s.removeWatchFunc(jobID)
	}
	return s.scheduler.RemoveJob(jobUUID)
}

// GetAlias gets alias by jobID.
// GetAlias 通过 jobID 获取别名。
//
// Parameters:
//
//	jobID - The jobs ID / 任务ID
//
// Returns:
//
//	string - The alias / 别名
//	error  - Error if alias not found / 如果找不到别名的错误
func (s *Scheduler) GetAlias(jobID string) (string, error) {
	// Check if alias option is Enabled
	// 检查是否启用别名选项
	if !s.Enable(common.AliasOptionName) {
		return "", common.ErrDisEnableAlias
	}
	// Find alias by jobID
	// 通过任务ID查找别名
	for alias, realJobID := range s.aliasMap {
		if jobID == realJobID {
			return alias, nil
		}
	}
	return "", common.ErrFoundAlias
}

// RunJobNow runs a jobs immediately by jobID.
// RunJobNow 通过 jobID 立即运行任务。
//
// Parameters:
//
//	jobID - The jobs ID to run / 要运行的任务ID
//
// Returns:
//
//	error - Error if run fails / 如果运行失败的错误
func (s *Scheduler) RunJobNow(jobID string) error {
	// Get jobs by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return err
	}
	// Run jobs immediately
	// 立即运行任务
	return job.RunNow()
}

// RunJobNowByAlias runs a jobs immediately by alias.
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
	// Check if alias option is Enabled
	// 检查是否启用别名选项
	if !s.Enable(common.AliasOptionName) {
		return common.ErrDisEnableAlias
	}
	// Get jobs by alias
	// 通过别名获取任务
	job, err := s.GetJobByAlias(alias)
	if err != nil {
		return err
	}
	// Run jobs immediately
	// 立即运行任务
	return job.RunNow()
}

// addJobType adds a jobs type mapping.
// addJobType 添加任务类型映射。
//
// Parameters:
//
//	jobID   - The jobs ID / 任务ID
//	jobType - The jobs type / 任务类型
func (s *Scheduler) addJobType(jobID string, jobType jobs.JobType) {
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

// removeJobType removes a jobs type mapping.
// removeJobType 移除任务类型映射。
//
// Parameters:
//
//	jobID - The jobs ID to remove / 要移除的任务ID
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

// getJobType gets a jobs type by jobID.
// getJobType 通过 jobID 获取任务类型。
//
// Parameters:
//
//	jobID - The jobs ID / 任务ID
//
// Returns:
//
//	jobs.JobType - The jobs type / 任务类型
func (s *Scheduler) getJobType(jobID string) jobs.JobType {
	s.jobTypeMu.Lock()
	defer s.jobTypeMu.Unlock()
	// Return jobs type if exists, otherwise return unknown
	// 如果存在则返回任务类型，否则返回未知类型
	if jobType, exists := s.jobTypeMap[jobID]; exists {
		return jobType
	}
	return jobs.JobTypeUnknown
}

// addAlias adds an alias for a jobs.
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

// removeAliasByJobID removes an alias by jobID.
// removeAliasByJobID 通过 jobID 移除别名。
func (s *Scheduler) removeAliasByJobID(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for alias, jid := range s.aliasMap {
		if jid == jobID {
			delete(s.aliasMap, alias)
			slog.Info("chrono:alias removed", "alias", alias, "jobID", jobID)
			return
		}
	}
	slog.Warn("chrono:alias not found for jobID", "jobID", jobID)
}

// addWatchFunc adds a watch function for a jobs.
// addWatchFunc 为任务添加监听函数。
func (s *Scheduler) addWatchFunc(jobID string, fn func(event monitor.JobWatchInterface)) {
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

// removeWatchFunc removes a watch function for a jobs.
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
	// If limit option is disabled, return true (unlimited)
	// 如果限制选项未启用，返回true（无限制）
	if !s.Enable(common.LimitOptionName) {
		return true
	}
	// Decrease limit
	// 减少限制
	if err := s.decLimit(); err != nil {
		return false
	}
	return s.schOptions.limit.Limit >= 0
}

// incLimit increases the limit.
// incLimit 增加限制。
//
// Returns:
//
//	error - Error if limit option is disabled / 如果限制选项未启用的错误
func (s *Scheduler) incLimit() error {
	// Check if limit option is Enabled
	// 检查是否启用限制选项
	if !s.Enable(common.LimitOptionName) {
		return common.ErrDisEnableLimit
	}
	s.schOptions.limit.Limit++
	return nil
}

// decLimit decreases the limit.
// decLimit 减少限制。
//
// Returns:
//
//	error - Error if limit option is disabled / 如果限制选项未启用的错误
func (s *Scheduler) decLimit() error {
	// Check if limit option is Enabled
	// 检查是否启用限制选项
	if !s.Enable(common.LimitOptionName) {
		return common.ErrDisEnableLimit
	}
	s.schOptions.limit.Limit--
	return nil
}

// TODO 批量移除任务
// RemoveJobs Removes jobs list with rollback support.
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
// 	for _, jobs := range jobs {
// 		if err := s.RemoveJob(jobs.ID().String()); err != nil {
// 			// 如果删除失败，回滚已删除的任务
// 			if rollbackErr := s.rollbackRemovedJobs(removedJobs); rollbackErr != nil {
// 				return fmt.Errorf("failed to remove jobs %s: %w; rollback failed: %v", jobs.ID(), err, rollbackErr)
// 			}
// 			return fmt.Errorf("failed to remove jobs %s: %w", jobs.ID(), err)
// 		}
// 		// 记录成功删除的任务
// 		removedJobs = append(removedJobs, jobs)
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
// 	for _, jobs := range jobs {
// 		if err := s.scheduler.AddJob(jobs); err != nil {
// 			rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to re-add jobs %s: %w", jobs.ID(), err))
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
//	[]pkg.Job - The list of jobs / 任务列表
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobs() ([]common.Job, error) {
	gocronJobs := s.scheduler.Jobs()
	jobs := make([]common.Job, 0, len(gocronJobs))
	for _, job := range gocronJobs {
		jobs = append(jobs, NewGocronJobWrapper(job))
	}
	return jobs, nil
}

// GetJobLastTimeByAlias gets the last run time of a jobs by alias.
// GetJobLastTimeByAlias 通过别名获取任务的最后运行时间。
//
// Parameters:
//
//	alias - The jobs alias / 任务别名
//
// Returns:
//
//	*time.Time - The last run time / 最后运行时间
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobLastTimeByAlias(alias string) (*time.Time, error) {
	// Check if alias option is Enabled
	// 检查是否启用别名选项
	if !s.Enable(common.AliasOptionName) {
		return nil, common.ErrDisEnableAlias
	}
	// Get jobID by alias
	// 通过别名获取任务ID
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, common.ErrFoundAlias
	}
	// Get jobs by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}
	// Get last run time
	// 获取最后运行时间
	lastRun, err := job.LastRun()
	if err != nil {
		return nil, err
	}
	return &lastRun, nil
}

// GetJobLastTime gets the last run time of a jobs by jobID.
// GetJobLastTime 通过 jobID 获取任务的最后运行时间。
//
// Parameters:
//
//	jobID - The jobs ID / 任务ID
//
// Returns:
//
//	*time.Time - The last run time / 最后运行时间
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobLastTime(jobID string) (*time.Time, error) {
	// Get jobs by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}
	// Get last run time
	// 获取最后运行时间
	lastRun, err := job.LastRun()
	if err != nil {
		return nil, err
	}
	return &lastRun, nil
}

// GetJobNextTimeByAlias gets the next run time of a jobs by alias.
// GetJobNextTimeByAlias 通过别名获取任务的下次运行时间。
//
// Parameters:
//
//	alias - The jobs alias / 任务别名
//
// Returns:
//
//	*time.Time - The next run time / 下次运行时间
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobNextTimeByAlias(alias string) (*time.Time, error) {
	if !s.Enable(common.AliasOptionName) {
		return nil, common.ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, common.ErrFoundAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return nil, err
	}
	return &nextRun, nil
}

// GetJobNextTime gets the next run time of a jobs by jobID.
// GetJobNextTime 通过 jobID 获取任务的下次运行时间。
//
// Parameters:
//
//	jobID - The jobs ID / 任务ID
//
// Returns:
//
//	*time.Time - The next run time / 下次运行时间
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobNextTime(jobID string) (*time.Time, error) {
	// Get jobs by ID
	// 通过ID获取任务
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}
	// Get next run time
	// 获取下次运行时间
	nextRun, err := job.NextRun()
	if err != nil {
		return nil, err
	}
	return &nextRun, nil
}

// GetJobLastAndNextByAlias gets the last and next run times of a jobs by alias.
// GetJobLastAndNextByAlias 通过别名获取任务的最后和下次运行时间。
//
// Parameters:
//
//	alias - The jobs alias / 任务别名
//
// Returns:
//
//	*time.Time - The last run time / 最后运行时间
//	*time.Time - The next run time / 下次运行时间
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobLastAndNextByAlias(alias string) (*time.Time, *time.Time, error) {
	if !s.Enable(common.AliasOptionName) {
		return nil, nil, common.ErrDisEnableAlias
	}
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, nil, common.ErrFoundAlias
	}
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, nil, err
	}
	lastRun, err := job.LastRun()
	if err != nil {
		return nil, nil, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return nil, nil, err
	}
	return &lastRun, &nextRun, nil
}

// GetJobLastAndNextByID gets the last and next run times of a jobs by jobID.
// GetJobLastAndNextByID 通过 jobID 获取任务的最后和下次运行时间。
//
// Parameters:
//
//	jobID - The jobs ID / 任务ID
//
// Returns:
//
//	*time.Time - The last run time / 最后运行时间
//	*time.Time - The next run time / 下次运行时间
//	error      - Error if retrieval fails / 如果获取失败的错误
func (s *Scheduler) GetJobLastAndNextByID(jobID string) (*time.Time, *time.Time, error) {
	job, err := s.GetJobByID(jobID)
	if err != nil {
		return nil, nil, err
	}
	lastRun, err := job.LastRun()
	if err != nil {
		return nil, nil, err
	}
	nextRun, err := job.NextRun()
	if err != nil {
		return nil, nil, err
	}
	return &lastRun, &nextRun, nil
}

// GetJobByName gets a jobs by its name.
// GetJobByName 通过名称获取任务。
//
// Parameters:
//
//	jobName - The jobs name / 任务名称
//
// Returns:
//
//	gocron.Job - The found jobs / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByName(jobName string) (gocron.Job, error) {
	// Search through all jobs
	// 遍历所有任务
	for _, job := range s.scheduler.Jobs() {
		if job.Name() == jobName {
			return job, nil
		}
	}
	return nil, fmt.Errorf("chrono:jobs %s not found", jobName)
}

// GetJobByID gets a jobs by its jobID.
// GetJobByID 通过 jobID 获取任务。
//
// Parameters:
//
//	jobID - The jobs ID / 任务ID
//
// Returns:
//
//	gocron.Job - The found jobs / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByID(jobID string) (gocron.Job, error) {
	// Search through all jobs
	// 遍历所有任务
	for _, job := range s.scheduler.Jobs() {
		if job.ID().String() == jobID {
			return job, nil
		}
	}
	return nil, fmt.Errorf("chrono:jobs %s not found", jobID)
}

// GetJobByAlias gets a jobs by its alias.
// GetJobByAlias 通过别名获取任务。
//
// Parameters:
//
//	alias - The jobs alias / 任务别名
//
// Returns:
//
//	gocron.Job - The found jobs / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByAlias(alias string) (gocron.Job, error) {
	// Check if alias option is Enabled
	// 检查是否启用别名选项
	if !s.Enable(common.AliasOptionName) {
		return nil, common.ErrDisEnableAlias
	}
	// Get jobID by alias
	// 通过别名获取任务ID
	jobID, ok := s.aliasMap[alias]
	if !ok {
		return nil, fmt.Errorf("chrono:alias %s not found", alias)
	}
	return s.GetJobByID(jobID)
}

// GetJobByIDOrAlias gets a jobs by jobID or alias, first by jobID, then by alias.
// GetJobByIDOrAlias 通过 jobID 或别名获取任务，优先通过 jobID。
//
// Parameters:
//
//	identifier - The jobs ID or alias / 任务ID或别名
//
// Returns:
//
//	gocron.Job - The found jobs / 找到的任务
//	error       - Error if not found / 如果找不到任务的错误
func (s *Scheduler) GetJobByIDOrAlias(identifier string) (gocron.Job, error) {
	// Try ID lookup first
	// 优先尝试ID查找
	if jobID, err := s.GetJobByID(identifier); err == nil {
		return jobID, nil
	}

	// Fallback to alias lookup if alias feature is Enabled
	// 如果启用别名功能，回退到别名查找
	if s.Enable(common.AliasOptionName) {
		if jobID, exists := s.aliasMap[identifier]; exists {
			return s.GetJobByAlias(jobID)
		}
	}

	return nil, fmt.Errorf("chrono:jobs with identifier %s not found", identifier)
}

// GetJobByIDS gets jobs by a list of jobIDs.
// GetJobByIDS 通过 jobID 列表获取任务。
//
// Parameters:
//
//	jobIDS - Variable number of jobs IDs / 可变数量的任务ID
//
// Returns:
//
//	[]gocron.Job - The list of found jobs / 找到的任务列表
//	error          - Error if any jobs is not found / 如果任何任务找不到的错误
func (s *Scheduler) GetJobByIDS(jobIDS ...string) ([]gocron.Job, error) {
	// Create a slice to store found jobs
	// 创建一个切片用于存储找到的任务
	jobs := make([]gocron.Job, 0, len(jobIDS))

	// Iterate through jobIDS and find each jobs
	// 遍历 jobIDS，逐个查找任务
	for _, jobID := range jobIDS {
		job, err := s.GetJobByID(jobID)
		if err != nil {
			return nil, fmt.Errorf("chrono:failed to get jobs %s: %w", jobID, err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// AddCronJob adds a new cron jobs.
// AddCronJob 添加一个新的 cron 任务。
//
// Parameters:
//
//	jobs - The cron jobs to add / 要添加的cron任务
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddCronJob(job any) (gocron.Job, error) {
	// Type assert to *jobs.CronJob
	// 类型断言为 *jobs.CronJob
	cronJob, ok := job.(*jobs.CronJob)
	if !ok || cronJob == nil {
		return nil, common.ErrInvalidJob
	}
	// Check if jobs has error
	// 检查任务是否有错误
	if cronJob.GetError() != nil {
		return nil, cronJob.GetError()
	}
	// check if jobs has a task function
	// 检查任务是否有任务函数
	if cronJob.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", cronJob.Name)
	}
	// Check if cron expression is set
	// 检查是否设置了cron表达式
	if cronJob.Expr == "" {
		return nil, fmt.Errorf("chrono:jobs %s has nil expr", cronJob.Name)
	}
	// Check limit if limit option is Enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(common.LimitOptionName) {
		if !s.CheckLimit() {
			return nil, common.ErrMoreLimit
		}
	}
	// Prepare jobs options
	// 准备任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(cronJob.Hooks...), gocron.WithName(cronJob.Name))
	// Set jobs ID if provided
	// 如果提供了任务ID，则设置
	if cronJob.ID != "" {
		jobID, err := uuid.Parse(cronJob.ID)
		if err != nil {
			return nil, fmt.Errorf("ichrono:nvalid jobs ID %s: %w", cronJob.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create cron jobs
	// 创建cron任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(cronJob.Expr, false), // 使用 cron 表达式
		gocron.NewTask(cronJob.TaskFunc),    // 任务函数
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add cron jobs: %w", err)
	}
	// Add jobs type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), jobs.JobType(cronJob.Type))
	// Add alias if alias option is Enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(common.AliasOptionName) {
		s.addAlias(cronJob.Ali, jobInstance.ID().String())
	}
	// Add watch function if watch option is Enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(common.WatchOptionName) {
		if cronJob.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), WatchFuncAdapter(cronJob.WatchFunc))
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc.(func(event monitor.JobWatchInterface)))
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
func (s *Scheduler) AddCronJobs(cronJobs ...any) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(cronJobs))
	// Add each cron jobs
	// 逐个添加cron任务
	for _, j := range cronJobs {
		cronJob, ok := j.(*jobs.CronJob)
		if !ok || cronJob == nil {
			return nil, common.ErrInvalidJob
		}
		cronJobInstance, err := s.AddCronJob(cronJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	// Return error if any jobs failed to add
	// 如果任何任务添加失败，则返回错误
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add cron jobs: %v", errs)
	}
	return jobList, nil
}

// AddCronJobWithOptions adds a cron jobs with options.
// AddCronJobWithOptions 添加带选项的 cron 任务。
func (s *Scheduler) AddCronJobWithOptions(job *jobs.CronJob, options ...gocron.JobOption) (gocron.Job, error) {
	if job == nil {
		return nil, common.ErrInvalidJob
	}
	if job.GetError() != nil {
		return nil, job.GetError()
	}
	// 检查任务函数是否存在
	if job.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", job.Name)
	}
	jobInstance, err := s.scheduler.NewJob(
		gocron.CronJob(job.Expr, false),
		gocron.NewTask(job.TaskFunc),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add jobs: %w", err)
	}
	return jobInstance, nil
}

// AddOnceJob adds a new once jobs.
// AddOnceJob 添加一个新的单次任务。
//
// Parameters:
//
//	jobs - The once jobs to add / 要添加的单次任务
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddOnceJob(job any) (gocron.Job, error) {
	onceJob, ok := job.(*jobs.OnceJob)
	if !ok || onceJob == nil {
		return nil, common.ErrInvalidJob
	}
	if onceJob.GetError() != nil {
		return nil, onceJob.GetError()
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if onceJob.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", onceJob.Name)
	}
	// Check limit if limit option is Enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(common.LimitOptionName) {
		if !s.CheckLimit() {
			return nil, common.ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(onceJob.Hooks...), gocron.WithName(onceJob.Name), gocron.WithTags(onceJob.Tag...))
	// Set jobs ID if provided
	// 如果提供了任务ID，则设置
	if onceJob.ID != "" {
		jobID, err := uuid.Parse(onceJob.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid jobs ID %s: %w", onceJob.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create once jobs
	// 创建单次任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTimes(onceJob.WorkTime...)),
		gocron.NewTask(onceJob.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add once jobs: %w", err)
	}
	// Add jobs type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), jobs.JobType(onceJob.Type))
	// Add watch function if watch option is Enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(common.WatchOptionName) {
		if onceJob.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), WatchFuncAdapter(onceJob.WatchFunc))
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc.(func(event monitor.JobWatchInterface)))
	}
	// Add alias if alias option is Enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(common.AliasOptionName) {
		s.addAlias(onceJob.Ali, jobInstance.ID().String())
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
func (s *Scheduler) AddOnceJobs(onceJobs ...any) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(onceJobs))
	// Add each once jobs
	// 逐个添加单次任务
	for _, j := range onceJobs {
		onceJob, ok := j.(*jobs.OnceJob)
		if !ok || onceJob == nil {
			return nil, common.ErrInvalidJob
		}
		cronJobInstance, err := s.AddOnceJob(onceJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, cronJobInstance)
	}
	// Return error if any jobs failed to add
	// 如果任何任务添加失败，则返回错误
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add once jobs: %v", errs)
	}
	return jobList, nil
}

// AddOnceJobWithOptions adds a once jobs with options.
// AddOnceJobWithOptions 添加带选项的单次任务。
func (s *Scheduler) AddOnceJobWithOptions(startAt gocron.OneTimeJobStartAtOption, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.OneTimeJob(startAt),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add once jobs: %w", err)
	}
	return job, nil
}

// AddIntervalJob adds a new interval jobs.
// AddIntervalJob 添加一个新的间隔任务。
//
// Parameters:
//
//	jobs - The interval jobs to add / 要添加的间隔任务
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddIntervalJob(job any) (gocron.Job, error) {
	intervalJob, ok := job.(*jobs.IntervalJob)
	if !ok || intervalJob == nil {
		return nil, common.ErrInvalidJob
	}
	if intervalJob.GetError() != nil {
		return nil, intervalJob.GetError()
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if intervalJob.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", intervalJob.Name)
	}
	// Check limit if limit option is Enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(common.LimitOptionName) {
		if !s.CheckLimit() {
			return nil, common.ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(intervalJob.Hooks...), gocron.WithName(intervalJob.Name))
	// Set jobs ID if provided
	// 如果提供了任务ID，则设置
	if intervalJob.ID != "" {
		jobID, err := uuid.Parse(intervalJob.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid jobs ID %s: %w", intervalJob.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}

	// Create interval jobs
	// 创建间隔任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.DurationJob(intervalJob.Interval),
		gocron.NewTask(intervalJob.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add jobs: %w", err)
	}
	// Add jobs type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), jobs.JobType(intervalJob.Type))
	// Add watch function if watch option is Enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(common.WatchOptionName) {
		if intervalJob.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), WatchFuncAdapter(intervalJob.WatchFunc))
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc.(func(event monitor.JobWatchInterface)))
	}
	// Add alias if alias option is Enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(common.AliasOptionName) {
		s.addAlias(intervalJob.Ali, jobInstance.ID().String())
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
func (s *Scheduler) AddIntervalJobs(intervalJobs ...any) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(intervalJobs))
	// Add each interval jobs
	// 逐个添加间隔任务
	for _, j := range intervalJobs {
		intervalJob, ok := j.(*jobs.IntervalJob)
		if !ok || intervalJob == nil {
			return nil, common.ErrInvalidJob
		}
		intervalJobInstance, err := s.AddIntervalJob(intervalJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, intervalJobInstance)
	}
	// Return error if any jobs failed to add
	// 如果任何任务添加失败，则返回错误
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add interval jobs: %v", errs)
	}
	return jobList, nil
}

// AddIntervalJobWithOptions adds an interval jobs with options.
// AddIntervalJobWithOptions 添加带选项的间隔任务。
func (s *Scheduler) AddIntervalJobWithOptions(interval time.Duration, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add interval jobs: %w", err)
	}
	return job, nil
}

// AddDailyJob adds a new daily jobs.
// AddDailyJob 添加一个新的每日任务。
//
// Parameters:
//
//	jobs - The daily jobs to add / 要添加的每日任务
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddDailyJob(job any) (gocron.Job, error) {
	dailyJob, ok := job.(*jobs.DailyJob)
	if !ok || dailyJob == nil {
		return nil, common.ErrInvalidJob
	}
	if dailyJob.GetError() != nil {
		return nil, dailyJob.GetError()
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if dailyJob.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", dailyJob.Name)
	}
	// Check limit if limit option is Enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(common.LimitOptionName) {
		if !s.CheckLimit() {
			return nil, common.ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(dailyJob.Hooks...), gocron.WithName(dailyJob.Name))
	// Set jobs ID if provided
	// 如果提供了任务ID，则设置
	if dailyJob.ID != "" {
		jobID, err := uuid.Parse(dailyJob.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid jobs ID %s: %w", dailyJob.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create daily jobs
	// 创建每日任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.DailyJob(dailyJob.Interval, dailyJob.AtTimes),
		gocron.NewTask(dailyJob.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add jobs: %w", err)
	}
	// Add jobs type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), jobs.JobType(dailyJob.Type))
	// Add watch function if watch option is Enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(common.WatchOptionName) {
		if dailyJob.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), WatchFuncAdapter(dailyJob.WatchFunc))
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc.(func(event monitor.JobWatchInterface)))
	}
	// Add alias if alias option is Enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(common.AliasOptionName) {
		s.addAlias(dailyJob.Ali, jobInstance.ID().String())
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
func (s *Scheduler) AddDailyJobs(dailyJobs ...any) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(dailyJobs))
	// Add each daily jobs
	// 逐个添加每日任务
	for _, j := range dailyJobs {
		dailyJob, ok := j.(*jobs.DailyJob)
		if !ok || dailyJob == nil {
			return nil, common.ErrInvalidJob
		}
		dailyJobInstance, err := s.AddDailyJob(dailyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, dailyJobInstance)
	}
	// Return error if any jobs failed to add
	// 如果任何任务添加失败，则返回错误
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add daily jobs: %v", errs)
	}
	return jobList, nil
}

// AddDailyJobWithOptions adds a daily jobs with options.
// AddDailyJobWithOptions 添加带选项的每日任务。
func (s *Scheduler) AddDailyJobWithOptions(interval uint, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.DailyJob(interval, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add daily jobs: %w", err)
	}
	return job, nil
}

// AddWeeklyJob adds a new weekly jobs.
// AddWeeklyJob 添加一个新的每周任务。
//
// Parameters:
//
//	jobs - The weekly jobs to add / 要添加的每周任务
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddWeeklyJob(job any) (gocron.Job, error) {
	weeklyJob, ok := job.(*jobs.WeeklyJob)
	if !ok || weeklyJob == nil {
		return nil, common.ErrInvalidJob
	}
	if weeklyJob.GetError() != nil {
		return nil, weeklyJob.GetError()
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if weeklyJob.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", weeklyJob.Name)
	}
	// Check limit if limit option is Enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(common.LimitOptionName) {
		if !s.CheckLimit() {
			return nil, common.ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(weeklyJob.Hooks...), gocron.WithName(weeklyJob.Name))
	// Set jobs ID if provided
	// 如果提供了任务ID，则设置
	if weeklyJob.ID != "" {
		jobID, err := uuid.Parse(weeklyJob.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid jobs ID %s: %w", weeklyJob.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create weekly jobs
	// 创建每周任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.WeeklyJob(weeklyJob.Interval, weeklyJob.DaysOfTheWeek, weeklyJob.WorkTimes),
		gocron.NewTask(weeklyJob.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add weekly jobs: %w", err)
	}
	// Add jobs type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), jobs.JobType(weeklyJob.Type))
	// Add watch function if watch option is Enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(common.WatchOptionName) {
		if weeklyJob.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), WatchFuncAdapter(weeklyJob.WatchFunc))
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc.(func(event monitor.JobWatchInterface)))
	}
	// Add alias if alias option is Enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(common.AliasOptionName) {
		s.addAlias(weeklyJob.Ali, jobInstance.ID().String())
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
func (s *Scheduler) AddWeeklyJobs(weeklyJobs ...any) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(weeklyJobs))
	// Add each weekly jobs
	// 逐个添加每周任务
	for _, j := range weeklyJobs {
		weeklyJob, ok := j.(*jobs.WeeklyJob)
		if !ok || weeklyJob == nil {
			return nil, common.ErrInvalidJob
		}
		weeklyJobInstance, err := s.AddWeeklyJob(weeklyJob)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		jobList = append(jobList, weeklyJobInstance)
	}
	// Return error if any jobs failed to add
	// 如果任何任务添加失败，则返回错误
	if len(errs) > 0 {
		return jobList, fmt.Errorf("chrono:failed to add weekly jobs: %v", errs)
	}
	return jobList, nil
}

// AddWeeklyJobWithOptions adds a weekly jobs with options.
// AddWeeklyJobWithOptions 添加带选项的每周任务。
func (s *Scheduler) AddWeeklyJobWithOptions(interval uint, daysOfTheWeek gocron.Weekdays, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.WeeklyJob(interval, daysOfTheWeek, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add weekly jobs: %w", err)
	}
	return job, nil
}

// AddMonthlyJob adds a new monthly jobs.
// AddMonthlyJob 添加一个新的每月任务。
//
// Parameters:
//
//	jobs - The monthly jobs to add / 要添加的每月任务
//
// Returns:
//
//	gocron.Job - The added jobs / 添加的任务
//	error      - Error if addition fails / 如果添加失败的错误
func (s *Scheduler) AddMonthlyJob(job any) (gocron.Job, error) {
	monthlyJob, ok := job.(*jobs.MonthJob)
	if !ok || monthlyJob == nil {
		return nil, common.ErrInvalidJob
	}
	if monthlyJob.GetError() != nil {
		return nil, monthlyJob.GetError()
	}
	// 检查任务函数是否存在
	// Check if task function exists
	if monthlyJob.TaskFunc == nil {
		return nil, fmt.Errorf("chrono:jobs %s has no task function", monthlyJob.Name)
	}
	// Check limit if limit option is Enabled
	// 如果启用限制选项，则检查限制
	if s.Enable(common.LimitOptionName) {
		if !s.CheckLimit() {
			return nil, common.ErrMoreLimit
		}
	}
	// Job options
	// 任务选项
	opts := make([]gocron.JobOption, 0)
	opts = append(opts, gocron.WithEventListeners(monthlyJob.Hooks...), gocron.WithName(monthlyJob.Name))
	// Set jobs ID if provided
	// 如果提供了任务ID，则设置
	if monthlyJob.ID != "" {
		jobID, err := uuid.Parse(monthlyJob.ID)
		if err != nil {
			return nil, fmt.Errorf("chrono:invalid jobs ID %s: %w", monthlyJob.ID, err)
		}
		opts = append(opts, gocron.WithIdentifier(jobID))
	}
	// Create monthly jobs
	// 创建每月任务
	jobInstance, err := s.scheduler.NewJob(
		gocron.MonthlyJob(monthlyJob.Interval, monthlyJob.DaysOfTheMonth, monthlyJob.AtTimes),
		gocron.NewTask(monthlyJob.TaskFunc),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add monthly jobs: %w", err)
	}
	// Add jobs type mapping
	// 添加任务类型映射
	s.addJobType(jobInstance.ID().String(), jobs.JobType(monthlyJob.Type))
	// Add watch function if watch option is Enabled
	// 如果启用监听选项，则添加监听函数
	if s.Enable(common.WatchOptionName) {
		if monthlyJob.WatchFunc != nil {
			s.addWatchFunc(jobInstance.ID().String(), WatchFuncAdapter(monthlyJob.WatchFunc))
		}
		s.addWatchFunc(jobInstance.ID().String(), s.schOptions.watch.WatchFunc.(func(event monitor.JobWatchInterface)))
	}
	// Add alias if alias option is Enabled
	// 如果启用别名选项，则添加别名
	if s.Enable(common.AliasOptionName) {
		s.addAlias(monthlyJob.Ali, jobInstance.ID().String())
	}
	return jobInstance, nil
}

// AddMonthlyJobs adds a list of new monthly jobs.
// AddMonthlyJobs 添加一组新的每月任务。
func (s *Scheduler) AddMonthlyJobs(monthlyJobs ...any) ([]gocron.Job, error) {
	var errs []error
	jobList := make([]gocron.Job, 0, len(monthlyJobs))
	for _, j := range monthlyJobs {
		monthlyJob, ok := j.(*jobs.MonthJob)
		if !ok || monthlyJob == nil {
			return nil, common.ErrInvalidJob
		}
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

// AddMonthlyJobWithOptions adds a monthly jobs with options.
// AddMonthlyJobWithOptions 添加带选项的每月任务。
func (s *Scheduler) AddMonthlyJobWithOptions(interval uint, daysOfTheMonth gocron.DaysOfTheMonth, atTimes gocron.AtTimes, task any, options ...gocron.JobOption) (gocron.Job, error) {
	job, err := s.scheduler.NewJob(
		gocron.MonthlyJob(interval, daysOfTheMonth, atTimes),
		gocron.NewTask(task),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("chrono:failed to add monthly jobs: %w", err)
	}
	return job, nil
}
