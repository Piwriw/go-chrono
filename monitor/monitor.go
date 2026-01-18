package monitor

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/retry"
)

const (
	// defaultMaxRecords is the default maximum number of jobs event records to keep.
	// defaultMaxRecords 是默认保留的任务事件记录的最大数量。
	defaultMaxRecords = 3
	// defaultMaxRetryHistory is the default maximum number of retry records to keep.
	// defaultMaxRetryHistory 是默认保留的重试记录的最大数量。
	defaultMaxRetryHistory = 10
)

// defaultEventIDGenerator is the default event ID generator.
// defaultEventIDGenerator 是默认的事件 ID 生成器。
var defaultEventIDGenerator = &UUIDEventIDGenerator{}

// JobWatchInterface defines the interface for jobs event watching.
// JobWatchInterface 定义了任务事件监听的接口。
type JobWatchInterface interface {
	// GetJobID gets the jobs ID.
	// 获取任务 ID。
	//
	// Returns:
	//	string - The jobs ID / 任务 ID
	GetJobID() string
	// GetJobName gets the jobs name.
	// 获取任务名称。
	//
	// Returns:
	//	string - The jobs name / 任务名称
	GetJobName() string
	// GetStartTime gets the jobs start time.
	// 获取任务开始时间。
	//
	// Returns:
	//	time.Time - The start time / 开始时间
	GetStartTime() time.Time
	// GetEndTime gets the jobs end time.
	// 获取任务结束时间。
	//
	// Returns:
	//	time.Time - The end time / 结束时间
	GetEndTime() time.Time
	// GetStatus gets the jobs status.
	// 获取任务状态。
	//
	// Returns:
	//	gocron.JobStatus - The jobs status / 任务状态
	GetStatus() gocron.JobStatus
	// GetTags gets the jobs tags.
	// 获取任务标签。
	//
	// Returns:
	//	[]string - The jobs tags / 任务标签
	GetTags() []string
	// Error gets the jobs error.
	// 获取任务错误。
	//
	// Returns:
	//	error - The jobs error / 任务错误
	Error() error
	// GetCurrentEvent gets the current jobs event.
	// 获取当前任务事件。
	//
	// Returns:
	//	*JobEvent - The current jobs event / 当前任务事件
	GetCurrentEvent() *JobEvent
}

// SchedulerMonitor defines the interface for a scheduler monitor.
// SchedulerMonitor 定义了调度器监控的接口。
type SchedulerMonitor interface {
	gocron.MonitorStatus
	// Watch watches jobs events.
	// 监听任务事件。
	//
	// Returns:
	//	chan JobWatchInterface - The channel for jobs events / 任务事件的通道
	Watch() chan JobWatchInterface
	// UpdateJobEvents updates jobs events.
	// 更新任务事件。
	//
	// Parameters:
	//	jobID      - The jobs ID / 任务 ID
	//	jobName    - The jobs name / 任务名称
	//	jobnewEvent - The new jobs event / 新的任务事件
	//	jobTags    - Variable number of jobs tags / 可变数量的任务标签
	UpdateJobEvents(jobID uuid.UUID, jobName string, jobnewEvent *JobEvent, jobTags ...string)
	// GetJobEvents gets jobs events.
	// 获取任务事件。
	//
	// Parameters:
	//	jobID - The jobs ID / 任务 ID
	// Returns:
	//	[]*JobEvent - The list of jobs events / 任务事件列表
	GetJobEvents(jobID string) []*JobEvent

	// GetRetryHistory gets the retry history for a jobs.
	// 获取任务的重试历史。
	//
	// Parameters:
	//	jobID - The jobs ID / 任务 ID
	//
	// Returns:
	//	[]*retry.RetryEvent - The list of retry events / 重试事件列表
	GetRetryHistory(jobID string) []*retry.RetryEvent

	// RecordRetryEvent records a retry event.
	// 记录重试事件。
	//
	// Parameters:
	//	event - The retry event to record / 要记录的重试事件
	RecordRetryEvent(event *retry.RetryEvent)
}

// defaultSchedulerMonitor is the default implementation of SchedulerMonitor.
// defaultSchedulerMonitor 是 SchedulerMonitor 的默认实现。
type defaultSchedulerMonitor struct {
	// mu is the mutex for thread-safe operations.
	// mu 是用于线程安全操作的互斥锁。
	mu sync.Mutex
	// counter tracks the execution count of jobs.
	// counter 跟踪任务的执行次数。
	counter map[string]int
	// time tracks the execution time of jobs.
	// time 跟踪任务的执行时间。
	time map[string][]time.Duration
	// jobChan is the channel for jobs events.
	// jobChan 是任务事件的通道。
	jobChan chan JobWatchInterface
	// maxRecords is the maximum number of jobs event records to keep.
	// maxRecords 是保留的任务事件记录的最大数量。
	maxRecords int
	// jobRecord stores the jobs specifications and events.
	// jobRecord 存储任务规范和事件。
	jobRecord map[string]MonitorJobSpec
	// eventIDCli is the event ID generator.
	// eventIDCli 是事件 ID 生成器。
	eventIDCli EventIDGenerator
	// maxRetryHistory is the maximum number of retry records to keep per jobs.
	// maxRetryHistory 是每个任务保留的重试记录的最大数量。
	maxRetryHistory int
	// retryHistory stores the retry history for jobs.
	// retryHistory 存储任务的重试历史。
	retryHistory map[string][]*retry.RetryEvent
}

var _ SchedulerMonitor = (*defaultSchedulerMonitor)(nil)

// SchedulerMonitorOption is the option type for configuring the scheduler monitor.
// SchedulerMonitorOption 是用于配置调度器监控的选项类型。
type SchedulerMonitorOption func(*defaultSchedulerMonitor)

// WithMaxRecords sets the maximum number of jobs event records to keep.
// WithMaxRecords 设置保留的任务事件记录的最大数量。
//
// Parameters:
//
//	maxRecords - The maximum number of records / 最大记录数量
//
// Returns:
//
//	func(*defaultSchedulerMonitor) - The scheduler monitor option / 调度器监控选项
func WithMaxRecords(maxRecords int) func(*defaultSchedulerMonitor) {
	return func(s *defaultSchedulerMonitor) {
		if s.maxRecords <= 0 {
			s.maxRecords = defaultMaxRecords
			return
		}
		s.maxRecords = maxRecords
	}
}

// WithEventIDGenerator sets the event ID generator.
// WithEventIDGenerator 设置事件 ID 生成器。
//
// Parameters:
//
//	eventIDGenerator - The event ID generator / 事件 ID 生成器
//
// Returns:
//
//	func(*defaultSchedulerMonitor) - The scheduler monitor option / 调度器监控选项
func WithEventIDGenerator(eventIDGenerator EventIDGenerator) func(*defaultSchedulerMonitor) {
	return func(s *defaultSchedulerMonitor) {
		if eventIDGenerator == nil {
			s.eventIDCli = defaultEventIDGenerator
			return
		}
		s.eventIDCli = eventIDGenerator
	}
}

// WithMaxRetryHistory sets the maximum number of retry records to keep per jobs.
// WithMaxRetryHistory 设置每个任务保留的重试记录的最大数量。
//
// Parameters:
//
//	maxRetryHistory - The maximum number of retry records / 最大重试记录数量
//
// Returns:
//
//	func(*defaultSchedulerMonitor) - The scheduler monitor option / 调度器监控选项
func WithMaxRetryHistory(maxRetryHistory int) func(*defaultSchedulerMonitor) {
	return func(s *defaultSchedulerMonitor) {
		if maxRetryHistory <= 0 {
			s.maxRetryHistory = defaultMaxRetryHistory
			return
		}
		s.maxRetryHistory = maxRetryHistory
	}
}

// UpdateJobEvents updates jobs events.
// UpdateJobEvents 更新任务事件。
//
// Parameters:
//
//	jobID      - The jobs ID / 任务 ID
//	jobName    - The jobs name / 任务名称
//	jobnewEvent - The new jobs event / 新的任务事件
//	jobTags    - Variable number of jobs tags / 可变数量的任务标签
func (s *defaultSchedulerMonitor) UpdateJobEvents(jobID uuid.UUID, jobName string, jobnewEvent *JobEvent, jobTags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the current event record
	// 获取当前的事件记录
	event, ok := s.jobRecord[jobID.String()]
	if !ok {
		// If not exists, create a new event record and return directly
		// 如果不存在，创建一个新的事件记录并直接返回
		s.jobRecord[jobID.String()] = MonitorJobSpec{
			JobSpec: JobSpec{
				JobID:   jobID.String(),
				JobName: jobName,
				Tags:    jobTags,
			},
			JobEvents: []*JobEvent{jobnewEvent},
		}
		return
	}

	// Add new event record
	// 添加新的事件记录
	if len(event.JobEvents) < s.maxRecords {
		event.JobEvents = append(event.JobEvents, jobnewEvent)
	} else {
		// Overwrite the earliest event, implementing a circular buffer
		// 覆盖最早的事件，实现环形缓冲区
		event.JobEvents = append(event.JobEvents[1:], jobnewEvent)
	}
	s.jobRecord[jobID.String()] = event
}

// GetJobEvents gets jobs events.
// GetJobEvents 获取任务事件。
//
// Parameters:
//
//	jobID - The jobs ID / 任务 ID
//
// Returns:
//
//	[]*JobEvent - The list of jobs events / 任务事件列表
func (s *defaultSchedulerMonitor) GetJobEvents(jobID string) []*JobEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.jobRecord) == 0 {
		return nil
	}
	events, ok := s.jobRecord[jobID]
	if !ok {
		return nil
	}
	// Return a copy to avoid race conditions
	// 返回副本以避免竞态条件
	result := make([]*JobEvent, len(events.JobEvents))
	copy(result, events.JobEvents)
	return result
}

// GetRetryHistory gets the retry history for a jobs.
// GetRetryHistory 获取任务的重试历史。
//
// Parameters:
//
//	jobID - The jobs ID / 任务 ID
//
// Returns:
//
//	[]*retry.RetryEvent - The list of retry events / 重试事件列表
func (s *defaultSchedulerMonitor) GetRetryHistory(jobID string) []*retry.RetryEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	history, ok := s.retryHistory[jobID]
	if !ok {
		return nil
	}
	// Return a copy
	// 返回副本
	result := make([]*retry.RetryEvent, len(history))
	copy(result, history)
	return result
}

// RecordRetryEvent records a retry event.
// RecordRetryEvent 记录重试事件。
//
// Parameters:
//
//	event - The retry event to record / 要记录的重试事件
func (s *defaultSchedulerMonitor) RecordRetryEvent(event *retry.RetryEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobID := event.OriginalEventID

	if _, ok := s.retryHistory[jobID]; !ok {
		s.retryHistory[jobID] = make([]*retry.RetryEvent, 0)
	}

	if len(s.retryHistory[jobID]) >= s.maxRetryHistory {
		s.retryHistory[jobID] = s.retryHistory[jobID][1:]
	}

	s.retryHistory[jobID] = append(s.retryHistory[jobID], event)

	slog.Info("chrono: retry event recorded",
		"job_id", jobID,
		"attempt", event.Attempt,
		"error", event.ErrorMessage)
}

// MonitorJobSpec represents the specification of a monitored jobs.
// MonitorJobSpec 表示被监控任务的规范。
type MonitorJobSpec struct {
	// JobSpec is the jobs specification.
	// JobSpec 是任务规范。
	JobSpec JobSpec
	// JobEvents is the list of jobs events.
	// JobEvents 是任务事件列表。
	JobEvents []*JobEvent
}

// JobSpec represents the specification of a jobs.
// JobSpec 表示任务的规范。
type JobSpec struct {
	// JobID is the jobs ID.
	// JobID 是任务 ID。
	JobID string
	// JobName is the jobs name.
	// JobName 是任务名称。
	JobName string
	// Tags are the jobs tags.
	// Tags 是任务标签。
	Tags []string
}

// GetJobID gets the jobs ID.
// GetJobID 获取任务 ID。
//
// Returns:
//
//	string - The jobs ID / 任务 ID
func (m MonitorJobSpec) GetJobID() string {
	return m.JobSpec.JobID
}

// GetJobName gets the jobs name.
// GetJobName 获取任务名称。
//
// Returns:
//
//	string - The jobs name / 任务名称
func (m MonitorJobSpec) GetJobName() string {
	return m.JobSpec.JobName
}

// GetCurrentEvent gets the current jobs event.
// GetCurrentEvent 获取当前任务事件。
//
// Returns:
//
//	*JobEvent - The current jobs event / 当前任务事件
func (m MonitorJobSpec) GetCurrentEvent() *JobEvent {
	if len(m.JobEvents) == 0 {
		return nil
	}
	return m.JobEvents[len(m.JobEvents)-1]
}

// Error gets the jobs error.
// Error 获取任务错误。
//
// Returns:
//
//	error - The jobs error / 任务错误
func (m MonitorJobSpec) Error() error {
	if len(m.JobEvents) == 0 {
		return nil
	}
	return m.JobEvents[len(m.JobEvents)-1].Err
}

// GetStartTime gets the jobs start time.
// GetStartTime 获取任务开始时间。
//
// Returns:
//
//	time.Time - The start time / 开始时间
func (m MonitorJobSpec) GetStartTime() time.Time {
	if len(m.JobEvents) == 0 {
		return time.Time{}
	}
	return m.JobEvents[len(m.JobEvents)-1].StartTime
}

// GetEndTime gets the jobs end time.
// GetEndTime 获取任务结束时间。
//
// Returns:
//
//	time.Time - The end time / 结束时间
func (m MonitorJobSpec) GetEndTime() time.Time {
	if len(m.JobEvents) == 0 {
		return time.Time{}
	}
	return m.JobEvents[len(m.JobEvents)-1].EndTime
}

// GetStatus gets the jobs status.
// GetStatus 获取任务状态。
//
// Returns:
//
//	gocron.JobStatus - The jobs status / 任务状态
func (m MonitorJobSpec) GetStatus() gocron.JobStatus {
	if len(m.JobEvents) == 0 {
		return ""
	}
	return m.JobEvents[len(m.JobEvents)-1].Status
}

var _ JobWatchInterface = (*MonitorJobSpec)(nil)

// JobEvent represents a jobs event.
// JobEvent 表示一个任务事件。
type JobEvent struct {
	// EventID is the event ID.
	// EventID 是事件 ID。
	EventID string
	// StartTime is the start time.
	// StartTime 是开始时间。
	StartTime time.Time
	// EndTime is the end time.
	// EndTime 是结束时间。
	EndTime time.Time
	// Status is the jobs status.
	// Status 是任务状态。
	Status gocron.JobStatus
	// Err is the jobs error.
	// Err 是任务错误。
	Err error

	// New retry-related fields
	// 新增重试相关字段
	// RetryCount is the number of retry attempts.
	// RetryCount 重试次数。
	RetryCount int `json:"retry_count,omitempty"`
	// IsRetry indicates if this is a retry event.
	// IsRetry 表示是否为重试事件。
	IsRetry bool `json:"is_retry,omitempty"`
	// OriginalEventID is the original event ID for retries.
	// OriginalEventID 原始事件 ID（用于重试）。
	OriginalEventID string `json:"original_event_id,omitempty"`
}

// MarshalJSON marshals the JobEvent to JSON.
// MarshalJSON 将 JobEvent 序列化为 JSON。
//
// Returns:
//
//	[]byte - The JSON bytes / JSON 字节
//	error    - Error if marshaling fails / 如果序列化失败则返回错误
func (m JobEvent) MarshalJSON() ([]byte, error) {
	type Alias struct {
		EventID         string           `json:"event_id"`
		StartTime       string           `json:"start_time"`
		EndTime         string           `json:"end_time"`
		Status          gocron.JobStatus `json:"status"`
		Err             string           `json:"error"`
		RetryCount      int              `json:"retry_count,omitempty"`
		IsRetry         bool             `json:"is_retry,omitempty"`
		OriginalEventID string           `json:"original_event_id,omitempty"`
	}

	var errStr string
	if m.Err != nil {
		errStr = m.Err.Error()
	}

	return json.Marshal(&Alias{
		EventID:         m.EventID,
		StartTime:       m.StartTime.Format(time.DateTime),
		EndTime:         m.EndTime.Format(time.DateTime),
		Status:          m.Status,
		Err:             errStr,
		RetryCount:      m.RetryCount,
		IsRetry:         m.IsRetry,
		OriginalEventID: m.OriginalEventID,
	})
}

// GetStartTime gets the jobs start time.
// GetStartTime 获取任务开始时间。
//
// Returns:
//
//	time.Time - The start time / 开始时间
func (m JobEvent) GetStartTime() time.Time {
	return m.StartTime
}

// GetEndTime gets the jobs end time.
// GetEndTime 获取任务结束时间。
//
// Returns:
//
//	time.Time - The end time / 结束时间
func (m JobEvent) GetEndTime() time.Time {
	return m.EndTime
}

// GetStatus gets the jobs status.
// GetStatus 获取任务状态。
//
// Returns:
//
//	gocron.JobStatus - The jobs status / 任务状态
func (m JobEvent) GetStatus() gocron.JobStatus {
	return m.Status
}

// GetSpendTime gets the jobs spend time.
// GetSpendTime 获取任务花费时间。
//
// Returns:
//
//	int64 - The spend time in milliseconds / 花费的时间（毫秒）
func (m JobEvent) GetSpendTime() int64 {
	return m.EndTime.UnixMilli() - m.StartTime.UnixMilli()
}

// GetTags gets the jobs tags.
// GetTags 获取任务标签。
//
// Returns:
//
//	[]string - The jobs tags / 任务标签
func (m MonitorJobSpec) GetTags() []string {
	return m.JobSpec.Tags
}

// GetError gets the jobs error.
// GetError 获取任务错误。
//
// Returns:
//
//	error - The jobs error / 任务错误
func (m JobEvent) GetError() error {
	return m.Err
}

// NewMonitorJobSpec creates a new MonitorJobSpec.
// NewMonitorJobSpec 创建一个新的 MonitorJobSpec。
//func NewMonitorJobSpec(id uuid.UUID, name string, tags []string) MonitorJobSpec {
//	return MonitorJobSpec{
//		JobID:     id,
//		JobName:   name,
//		Tags:      tags,
//	}
//}

// newDefaultSchedulerMonitor creates a new default scheduler monitor.
// newDefaultSchedulerMonitor 创建一个默认的调度器监控。
//
// Parameters:
//
//	opts - Variable number of scheduler monitor options / 可变数量的调度器监控选项
//
// Returns:
//
//	*defaultSchedulerMonitor - The new scheduler monitor / 新的调度器监控
func NewDefaultSchedulerMonitor(opts ...SchedulerMonitorOption) *defaultSchedulerMonitor {
	defaultSchedulerMonitor := &defaultSchedulerMonitor{
		counter:         make(map[string]int),
		time:            make(map[string][]time.Duration),
		jobChan:         make(chan JobWatchInterface, 100),
		jobRecord:       make(map[string]MonitorJobSpec),
		eventIDCli:      defaultEventIDGenerator,
		maxRetryHistory: defaultMaxRetryHistory,
		retryHistory:    make(map[string][]*retry.RetryEvent),
	}
	for _, opt := range opts {
		opt(defaultSchedulerMonitor)
	}
	return defaultSchedulerMonitor
}

// IncrementJob increments the execution count of a jobs.
// IncrementJob 增加任务的执行次数。
//
// Parameters:
//
//	id     - The jobs ID / 任务 ID
//	name   - The jobs name / 任务名称
//	tags   - The jobs tags / 任务标签
//	status - The jobs status / 任务状态
func (s *defaultSchedulerMonitor) IncrementJob(id uuid.UUID, name string, tags []string, status gocron.JobStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Debug("chrono:IncrementJob", "JobID", id, "JobName", name, "tags", tags, "status", status)
	_, ok := s.counter[name]
	if !ok {
		s.counter[name] = 0
	}
	s.counter[name]++
}

// RecordJobTiming records the execution time of a jobs.
// RecordJobTiming 记录任务的执行时间。
//
// Parameters:
//
//	startTime - The start time / 开始时间
//	endTime   - The end time / 结束时间
//	id        - The jobs ID / 任务 ID
//	name      - The jobs name / 任务名称
//	tags      - The jobs tags / 任务标签
func (s *defaultSchedulerMonitor) RecordJobTiming(startTime, endTime time.Time, id uuid.UUID, name string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Debug("chrono:RecordJobTiming", "JobID", id, "JobName", name, "startTime", startTime.Format(time.DateTime),
		"endTime", endTime.Format(time.DateTime), "duration", endTime.Sub(startTime), "tags", tags)
	_, ok := s.time[name]
	if !ok {
		s.time[name] = make([]time.Duration, 0)
	}
	s.time[name] = append(s.time[name], endTime.Sub(startTime))
}

// RecordJobTimingWithStatus records the execution time and status of a jobs.
// RecordJobTimingWithStatus 记录任务的执行时间和状态。
//
// Parameters:
//
//	startTime - The start time / 开始时间
//	endTime   - The end time / 结束时间
//	id        - The jobs ID / 任务 ID
//	name      - The jobs name / 任务名称
//	tags      - The jobs tags / 任务标签
//	status    - The jobs status / 任务状态
//	err       - The jobs error / 任务错误
func (s *defaultSchedulerMonitor) RecordJobTimingWithStatus(startTime, endTime time.Time, id uuid.UUID, name string, tags []string, status gocron.JobStatus, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Debug("chrono:RecordJobTimingWithStatus", "JobID", id, "JobName", name, "startTime", startTime.Format(time.DateTime),
		"endTime", endTime.Format(time.DateTime), "duration", endTime.Sub(startTime), "status", status, "err", err)
	jobSpec := JobSpec{
		JobID:   id.String(),
		JobName: name,
		Tags:    tags,
	}

	newEvent := &JobEvent{
		EventID:   s.eventIDCli.NextID(jobSpec),
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
		Err:       err,
	}

	s.UpdateJobEvents(id, name, newEvent)
	// Create MonitorJobSpec containing the new event
	// 创建包含新事件的 MonitorJobSpec
	jobMonitorSpec := MonitorJobSpec{
		JobSpec:   jobSpec,
		JobEvents: []*JobEvent{newEvent},
	}
	s.jobChan <- jobMonitorSpec
}

// Watch watches the execution of jobs.
// Watch 监听任务的执行情况。
//
// Returns:
//
//	chan JobWatchInterface - The channel for jobs events / 任务事件的通道
func (s *defaultSchedulerMonitor) Watch() chan JobWatchInterface {
	return s.jobChan
	// for {
	// 	select {
	// 	case taskSpec := <-s.taskChan:
	// 		slog.Info("Watch", "taskSpec", taskSpec)
	// 		// 在这里可以添加更多的处理逻辑
	// 	case <-ctx.Done():
	// 		slog.Info("Watch stopped")
	// 		return
	// 	}
	// }
}
