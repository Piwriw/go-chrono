package chrono

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

const (
	defaultMaxRecords = 3 // 默认最大记录数
)

var defaultEventIDGenerator = &UUIDEventIDGenerator{}

// JobWatchInterface defines the interface for job event watching.
// JobWatchInterface 定义了任务事件监听的接口。
type JobWatchInterface interface {
	// GetJobID Gets the job ID
	// 获取任务 ID
	GetJobID() string
	// GetJobName Gets the job name
	// 获取任务名称
	GetJobName() string
	// GetStartTime Gets the job start time
	// 获取任务开始时间
	//GetStartTime() time.Time
	// GetEndTime Gets the job end time
	// 获取任务结束时间
	//GetEndTime() time.Time
	// GetStatus Gets the job status
	// 获取任务状态
	//GetStatus() gocron.JobStatus
	// GetTags Gets the job tags
	// 获取任务标签
	GetTags() []string
	// Error Gets the job error
	// 获取任务错误
	//Error() error
	GetCurrentEvent() *JobEvent
}

// SchedulerMonitor defines the interface for a scheduler monitor.
// SchedulerMonitor 定义了调度器监控的接口。
type SchedulerMonitor interface {
	gocron.MonitorStatus
	// Watch Watches job events
	// 监听任务事件
	Watch() chan JobWatchInterface
	// UpdateJobEvents Updates job events
	// 更新任务事件
	UpdateJobEvents(jobID uuid.UUID, jobName string, jobnewEvent *JobEvent, jobTags ...string)
	// GetJobEvents Gets job events
	// 获取任务事件
	GetJobEvents(jobID string) []*JobEvent
}

type defaultSchedulerMonitor struct {
	mu         sync.Mutex
	counter    map[string]int
	time       map[string][]time.Duration
	jobChan    chan JobWatchInterface
	maxRecords int
	jobRecord  map[string]MonitorJobSpec
	eventIDCli EventIDGenerator
}

var _ SchedulerMonitor = (*defaultSchedulerMonitor)(nil)

type SchedulerMonitorOption func(*defaultSchedulerMonitor)

func WithMaxRecords(maxRecords int) func(*defaultSchedulerMonitor) {
	return func(s *defaultSchedulerMonitor) {
		if s.maxRecords <= 0 {
			s.maxRecords = defaultMaxRecords
			return
		}
		s.maxRecords = maxRecords
	}
}

func WithEventIDGenerator(eventIDGenerator EventIDGenerator) func(*defaultSchedulerMonitor) {
	return func(s *defaultSchedulerMonitor) {
		if eventIDGenerator == nil {
			s.eventIDCli = defaultEventIDGenerator
			return
		}
		s.eventIDCli = eventIDGenerator
	}
}

// UpdateJobEvents 更新任务事件
func (s *defaultSchedulerMonitor) UpdateJobEvents(jobID uuid.UUID, jobName string, jobnewEvent *JobEvent, jobTags ...string) {
	// 获取当前的事件记录
	event, ok := s.jobRecord[jobID.String()]
	if !ok {
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

	// 添加新的事件记录
	if len(event.JobEvents) < s.maxRecords {
		event.JobEvents = append(event.JobEvents, jobnewEvent)
	} else {
		// 覆盖最早的事件，实现环形缓冲区
		event.JobEvents = append(event.JobEvents[1:], jobnewEvent)
	}
	s.jobRecord[jobID.String()] = event
}

func (s *defaultSchedulerMonitor) GetJobEvents(jobID string) []*JobEvent {
	if len(s.jobRecord) == 0 {
		return nil
	}
	events, ok := s.jobRecord[jobID]
	if !ok {
		return nil
	}
	return events.JobEvents
}

// MonitorJobSpec represents the specification of a monitored job.
// MonitorJobSpec 表示被监控任务的规范。
type MonitorJobSpec struct {
	JobSpec   JobSpec
	JobEvents []*JobEvent
}

type JobSpec struct {
	// Job ID
	// 任务 ID
	JobID string
	// Job name
	// 任务名称
	JobName string
	// Job tags
	// 任务标签
	Tags []string
}

// GetJobID gets the job ID.
// GetJobID 获取任务 ID。
func (m MonitorJobSpec) GetJobID() string {
	return m.JobSpec.JobID
}

// GetJobName gets the job name.
// GetJobName 获取任务名称。
func (m MonitorJobSpec) GetJobName() string {
	return m.JobSpec.JobName
}

// GetCurrentEvent gets the current job event.
// GetCurrentEvent 获取当前任务事件
func (m MonitorJobSpec) GetCurrentEvent() *JobEvent {
	if len(m.JobEvents) == 0 {
		return nil
	}
	return m.JobEvents[len(m.JobEvents)-1]
}

var _ JobWatchInterface = (*MonitorJobSpec)(nil)

type JobEvent struct {
	EventID string
	// Start time
	// 开始时间
	StartTime time.Time
	// End time
	// 结束时间
	EndTime time.Time
	// Job status
	// 任务状态
	Status gocron.JobStatus
	// Job error
	// 任务错误
	Err error
}

func (m JobEvent) MarshalJSON() ([]byte, error) {
	type Alias struct {
		EventID   string           `json:"event_id"`
		StartTime string           `json:"start_time"`
		EndTime   string           `json:"end_time"`
		Status    gocron.JobStatus `json:"status"`
		Err       string           `json:"error"`
	}

	var errStr string
	if m.Err != nil {
		errStr = m.Err.Error()
	}

	return json.Marshal(&Alias{
		EventID:   m.EventID,
		StartTime: m.StartTime.Format(time.DateTime),
		EndTime:   m.EndTime.Format(time.DateTime),
		Status:    m.Status,
		Err:       errStr,
	})
}

// GetStartTime gets the job start time.
// GetStartTime 获取任务开始时间。
func (m JobEvent) GetStartTime() time.Time {
	return m.StartTime
}

// GetEndTime gets the job end time.
// GetEndTime 获取任务结束时间。
func (m JobEvent) GetEndTime() time.Time {
	return m.EndTime
}

// GetStatus gets the job status.
// GetStatus 获取任务状态。
func (m JobEvent) GetStatus() gocron.JobStatus {
	return m.Status
}

// GetTags gets the job tags.
// GetTags 获取任务标签。
func (m MonitorJobSpec) GetTags() []string {
	return m.JobSpec.Tags
}

// Error gets the job error.
// Error 获取任务错误。
func (m JobEvent) Error() error {
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
func newDefaultSchedulerMonitor(opts ...SchedulerMonitorOption) *defaultSchedulerMonitor {
	defaultSchedulerMonitor := &defaultSchedulerMonitor{
		counter:    make(map[string]int),
		time:       make(map[string][]time.Duration),
		jobChan:    make(chan JobWatchInterface, 100),
		jobRecord:  make(map[string]MonitorJobSpec),
		eventIDCli: defaultEventIDGenerator,
	}
	for _, opt := range opts {
		opt(defaultSchedulerMonitor)
	}
	return defaultSchedulerMonitor
}

// IncrementJob increments the execution count of a job.
// IncrementJob 增加任务的执行次数。
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

// RecordJobTiming records the execution time of a job.
// RecordJobTiming 记录任务的执行时间。
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

// RecordJobTimingWithStatus records the execution time and status of a job.
// RecordJobTimingWithStatus 记录任务的执行时间和状态。
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
	jobMonitorSpec := MonitorJobSpec{
		JobSpec:   jobSpec,
		JobEvents: []*JobEvent{},
	}

	newEvent := &JobEvent{
		EventID:   s.eventIDCli.NextID(jobSpec),
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
		Err:       err,
	}

	s.UpdateJobEvents(id, name, newEvent)
	// 创建包含新事件的 MonitorJobSpec
	jobMonitorSpec = MonitorJobSpec{
		JobSpec:   jobSpec,
		JobEvents: []*JobEvent{newEvent}, // 包含当前事件
	}
	s.jobChan <- jobMonitorSpec
}

// Watch watches the execution of jobs.
// Watch 监听任务的执行情况。
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
