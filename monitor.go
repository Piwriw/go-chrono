package chrono

import (
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type JobWatchInterface interface {
	GetJobID() string
	GetJobName() string
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetStatus() gocron.JobStatus
	GetTags() []string
	Error() error
}

type SchedulerMonitor interface {
	gocron.MonitorStatus
	Watch() chan JobWatchInterface
}

type defaultSchedulerMonitor struct {
	mu      sync.Mutex
	counter map[string]int
	time    map[string][]time.Duration
	jobChan chan JobWatchInterface
}

type MonitorJobSpec struct {
	JobID     uuid.UUID
	JobName   string
	StartTime time.Time
	EndTime   time.Time
	Status    gocron.JobStatus
	Tags      []string
	Err       error
}

func (m MonitorJobSpec) GetJobID() string {
	return m.JobID.String()
}

func (m MonitorJobSpec) GetJobName() string {
	return m.JobName
}

func (m MonitorJobSpec) GetStartTime() time.Time {
	return m.StartTime
}

func (m MonitorJobSpec) GetEndTime() time.Time {
	return m.EndTime
}

func (m MonitorJobSpec) GetStatus() gocron.JobStatus {
	return m.Status
}

func (m MonitorJobSpec) GetTags() []string {
	return m.Tags
}

func (m MonitorJobSpec) Error() error {
	return m.Err
}

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

func newDefaultSchedulerMonitor() *defaultSchedulerMonitor {
	return &defaultSchedulerMonitor{
		counter: make(map[string]int),
		time:    make(map[string][]time.Duration),
		jobChan: make(chan JobWatchInterface, 100),
	}
}

// IncrementJob 增加任务的执行次数
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

// RecordJobTiming 记录任务的执行时间
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

// RecordJobTimingWithStatus 记录任务的执行时间、状态
func (s *defaultSchedulerMonitor) RecordJobTimingWithStatus(startTime, endTime time.Time, id uuid.UUID, name string, tags []string, status gocron.JobStatus, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Debug("chrono:RecordJobTimingWithStatus", "JobID", id, "JobName", name, "startTime", startTime.Format("2006-01-02 15:04:05"),
		"endTime", endTime.Format("2006-01-02 15:04:05"), "duration", endTime.Sub(startTime), "status", status, "err", err)
	jobSpec := NewMonitorJobSpec(id, name, startTime, endTime, tags, status, err)
	s.jobChan <- jobSpec
}

// Watch 监听任务的执行情况
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
