package chrono

import (
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

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
	GetStartTime() time.Time
	// GetEndTime Gets the job end time
	// 获取任务结束时间
	GetEndTime() time.Time
	// GetStatus Gets the job status
	// 获取任务状态
	GetStatus() gocron.JobStatus
	// GetTags Gets the job tags
	// 获取任务标签
	GetTags() []string
	// Error Gets the job error
	// 获取任务错误
	Error() error
}

// SchedulerMonitor defines the interface for a scheduler monitor.
// SchedulerMonitor 定义了调度器监控的接口。
type SchedulerMonitor interface {
	gocron.MonitorStatus
	Watch() chan JobWatchInterface // Watches job events
	// 监听任务事件
}

type defaultSchedulerMonitor struct {
	mu      sync.Mutex
	counter map[string]int
	time    map[string][]time.Duration
	jobChan chan JobWatchInterface
}

// MonitorJobSpec represents the specification of a monitored job.
// MonitorJobSpec 表示被监控任务的规范。
type MonitorJobSpec struct {
	JobID uuid.UUID // Job ID
	// 任务 ID
	JobName string // Job name
	// 任务名称
	StartTime time.Time // Start time
	// 开始时间
	EndTime time.Time // End time
	// 结束时间
	Status gocron.JobStatus // Job status
	// 任务状态
	Tags []string // Job tags
	// 任务标签
	Err error // Job error
	// 任务错误
}

// GetJobID gets the job ID.
// GetJobID 获取任务 ID。
func (m MonitorJobSpec) GetJobID() string {
	return m.JobID.String()
}

// GetJobName gets the job name.
// GetJobName 获取任务名称。
func (m MonitorJobSpec) GetJobName() string {
	return m.JobName
}

// GetStartTime gets the job start time.
// GetStartTime 获取任务开始时间。
func (m MonitorJobSpec) GetStartTime() time.Time {
	return m.StartTime
}

// GetEndTime gets the job end time.
// GetEndTime 获取任务结束时间。
func (m MonitorJobSpec) GetEndTime() time.Time {
	return m.EndTime
}

// GetStatus gets the job status.
// GetStatus 获取任务状态。
func (m MonitorJobSpec) GetStatus() gocron.JobStatus {
	return m.Status
}

// GetTags gets the job tags.
// GetTags 获取任务标签。
func (m MonitorJobSpec) GetTags() []string {
	return m.Tags
}

// Error gets the job error.
// Error 获取任务错误。
func (m MonitorJobSpec) Error() error {
	return m.Err
}

// NewMonitorJobSpec creates a new MonitorJobSpec.
// NewMonitorJobSpec 创建一个新的 MonitorJobSpec。
func NewMonitorJobSpec(id uuid.UUID, name string, startTime, endTime time.Time, tags []string, status gocron.JobStatus, err error) MonitorJobSpec {
	return MonitorJobSpec{
		JobID:     id,
		JobName:   name,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
		Tags:      tags,
		Err:       err,
	}
}

// newDefaultSchedulerMonitor creates a new default scheduler monitor.
// newDefaultSchedulerMonitor 创建一个默认的调度器监控。
func newDefaultSchedulerMonitor() *defaultSchedulerMonitor {
	return &defaultSchedulerMonitor{
		counter: make(map[string]int),
		time:    make(map[string][]time.Duration),
		jobChan: make(chan JobWatchInterface, 100),
	}
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
	slog.Debug("chrono:RecordJobTiming", "JobID", id, "JobName", name, "startTime", startTime.Format("2006-01-02 15:04:05"),
		"endTime", endTime.Format("2006-01-02 15:04:05"), "duration", endTime.Sub(startTime), "tags", tags)
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
	slog.Debug("chrono:RecordJobTimingWithStatus", "JobID", id, "JobName", name, "startTime", startTime.Format("2006-01-02 15:04:05"),
		"endTime", endTime.Format("2006-01-02 15:04:05"), "duration", endTime.Sub(startTime), "status", status, "err", err)
	jobSpec := NewMonitorJobSpec(id, name, startTime, endTime, tags, status, err)
	s.jobChan <- jobSpec
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
