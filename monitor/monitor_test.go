package monitor

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewDefaultSchedulerMonitor tests the creation of a new default scheduler monitor.
// TestNewDefaultSchedulerMonitor 测试创建新的默认调度器监控。
//
// Test scenarios:
// - Happy path: create monitor with default settings
// - Boundary conditions: create monitor with custom options
// - Exception cases: invalid option values
func TestNewDefaultSchedulerMonitor(t *testing.T) {
	t.Parallel()

	t.Run("creates monitor with default settings", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()

		require.NotNil(t, monitor, "monitor should not be nil")
		assert.NotNil(t, monitor.counter, "counter map should be initialized")
		assert.NotNil(t, monitor.time, "time map should be initialized")
		assert.NotNil(t, monitor.jobChan, "jobs channel should be initialized")
		assert.NotNil(t, monitor.jobRecord, "jobs record map should be initialized")
		assert.NotNil(t, monitor.eventIDCli, "event ID generator should be initialized")
		assert.Equal(t, defaultMaxRecords, monitor.maxRecords, "maxRecords should be default value")
	})

	t.Run("creates monitor with custom max records", func(t *testing.T) {
		t.Parallel()

		customMaxRecords := 10
		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(customMaxRecords))

		require.NotNil(t, monitor, "monitor should not be nil")
		assert.Equal(t, customMaxRecords, monitor.maxRecords, "maxRecords should be custom value")
	})

	t.Run("creates monitor with zero max records uses default", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(0))

		assert.Equal(t, defaultMaxRecords, monitor.maxRecords, "maxRecords should be default when zero")
	})

	t.Run("creates monitor with negative max records uses default", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(-5))

		assert.Equal(t, defaultMaxRecords, monitor.maxRecords, "maxRecords should be default when negative")
	})

	t.Run("creates monitor with custom event ID generator", func(t *testing.T) {
		t.Parallel()

		customGenerator := &TimeEventIDGenerator{timeFormat: "20060102150405"}
		monitor := NewDefaultSchedulerMonitor(WithEventIDGenerator(customGenerator))

		require.NotNil(t, monitor, "monitor should not be nil")
		assert.Equal(t, customGenerator, monitor.eventIDCli, "event ID generator should be custom")
	})

	t.Run("creates monitor with nil event ID generator uses default", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithEventIDGenerator(nil))

		assert.Equal(t, defaultEventIDGenerator, monitor.eventIDCli, "event ID generator should be default when nil")
	})

	t.Run("creates monitor with multiple options", func(t *testing.T) {
		t.Parallel()

		customMaxRecords := 20
		customGenerator := &TimeEventIDGenerator{}
		monitor := NewDefaultSchedulerMonitor(
			WithMaxRecords(customMaxRecords),
			WithEventIDGenerator(customGenerator),
		)

		assert.Equal(t, customMaxRecords, monitor.maxRecords, "maxRecords should be custom")
		assert.Equal(t, customGenerator, monitor.eventIDCli, "event ID generator should be custom")
	})
}

// TestIncrementJob tests the IncrementJob method.
// TestIncrementJob 测试 IncrementJob 方法。
//
// Test scenarios:
// - Happy path: increment jobs count
// - Boundary conditions: first increment, multiple increments
func TestIncrementJob(t *testing.T) {
	t.Parallel()

	t.Run("increments jobs count for new jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "test-jobs"
		tags := []string{"tag1", "tag2"}

		monitor.IncrementJob(jobID, jobName, tags, gocron.Success)

		assert.Equal(t, 1, monitor.counter[jobName], "jobs count should be 1 after first increment")
	})

	t.Run("increments jobs count for existing jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "existing-jobs"
		tags := []string{"tag1"}

		monitor.IncrementJob(jobID, jobName, tags, gocron.Success)
		monitor.IncrementJob(jobID, jobName, tags, gocron.Success)
		monitor.IncrementJob(jobID, jobName, tags, gocron.Success)

		assert.Equal(t, 3, monitor.counter[jobName], "jobs count should be 3 after three increments")
	})

	t.Run("thread-safe concurrent increments", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "concurrent-jobs"
		tags := []string{"concurrent"}

		const goroutines = 100
		const incrementsPerGoroutine = 10

		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < incrementsPerGoroutine; j++ {
					monitor.IncrementJob(jobID, jobName, tags, gocron.Success)
				}
			}()
		}

		wg.Wait()

		expectedCount := goroutines * incrementsPerGoroutine
		assert.Equal(t, expectedCount, monitor.counter[jobName], "concurrent increments should be accurate")
	})
}

// TestRecordJobTiming tests the RecordJobTiming method.
// TestRecordJobTiming 测试 RecordJobTiming 方法。
//
// Test scenarios:
// - Happy path: record jobs timing
// - Boundary conditions: first recording, multiple recordings
func TestRecordJobTiming(t *testing.T) {
	t.Parallel()

	t.Run("records timing for new jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "timing-jobs"
		tags := []string{"timing"}
		startTime := time.Now().Add(-1 * time.Second)
		endTime := time.Now()

		monitor.RecordJobTiming(startTime, endTime, jobID, jobName, tags)

		require.Contains(t, monitor.time, jobName, "jobs should have timing records")
		assert.Len(t, monitor.time[jobName], 1, "should have one timing record")
	})

	t.Run("records multiple timings for same jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "multi-timing-jobs"
		tags := []string{"multi"}

		for i := 0; i < 5; i++ {
			startTime := time.Now().Add(-time.Duration(i+1) * time.Second)
			endTime := time.Now()
			monitor.RecordJobTiming(startTime, endTime, jobID, jobName, tags)
		}

		assert.Len(t, monitor.time[jobName], 5, "should have five timing records")
	})

	t.Run("calculates correct duration", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "duration-jobs"
		tags := []string{"duration"}
		startTime := time.Now().Add(-100 * time.Millisecond)
		endTime := time.Now()

		monitor.RecordJobTiming(startTime, endTime, jobID, jobName, tags)

		require.Len(t, monitor.time[jobName], 1, "should have one timing record")
		duration := monitor.time[jobName][0]
		assert.Greater(t, duration, 90*time.Millisecond, "duration should be at least 90ms")
		assert.Less(t, duration, 200*time.Millisecond, "duration should be less than 200ms")
	})

	t.Run("thread-safe concurrent recordings", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "concurrent-timing-jobs"
		tags := []string{"concurrent"}

		const goroutines = 50
		const recordingsPerGoroutine = 10

		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < recordingsPerGoroutine; j++ {
					startTime := time.Now().Add(-10 * time.Millisecond)
					endTime := time.Now()
					monitor.RecordJobTiming(startTime, endTime, jobID, jobName, tags)
				}
			}()
		}

		wg.Wait()

		expectedCount := goroutines * recordingsPerGoroutine
		assert.Len(t, monitor.time[jobName], expectedCount, "concurrent recordings should be accurate")
	})
}

// TestRecordJobTimingWithStatus tests the RecordJobTimingWithStatus method.
// TestRecordJobTimingWithStatus 测试 RecordJobTimingWithStatus 方法。
//
// Test scenarios:
// - Happy path: record jobs timing with status
// - Boundary conditions: success and failure statuses
// - Exception cases: with and without errors
func TestRecordJobTimingWithStatus(t *testing.T) {
	t.Parallel()

	t.Run("records successful jobs execution", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "success-jobs"
		tags := []string{"success"}
		startTime := time.Now().Add(-50 * time.Millisecond)
		endTime := time.Now()

		monitor.RecordJobTimingWithStatus(startTime, endTime, jobID, jobName, tags, gocron.Success, nil)

		// Verify event was created
		// 验证事件已创建
		events := monitor.GetJobEvents(jobID.String())
		require.Len(t, events, 1, "should have one event")

		event := events[0]
		assert.Equal(t, gocron.Success, event.Status, "status should be Success")
		assert.NoError(t, event.Err, "error should be nil")
		assert.WithinDuration(t, startTime, event.StartTime, 10*time.Millisecond, "start time should match")
		assert.WithinDuration(t, endTime, event.EndTime, 10*time.Millisecond, "end time should match")
	})

	t.Run("records failed jobs execution with error", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "failed-jobs"
		tags := []string{"failed"}
		startTime := time.Now().Add(-50 * time.Millisecond)
		endTime := time.Now()
		jobError := assert.AnError

		monitor.RecordJobTimingWithStatus(startTime, endTime, jobID, jobName, tags, gocron.Fail, jobError)

		events := monitor.GetJobEvents(jobID.String())
		require.Len(t, events, 1, "should have one event")

		event := events[0]
		assert.Equal(t, gocron.Fail, event.Status, "status should be Fail")
		assert.Equal(t, jobError, event.Err, "error should match")
	})

	t.Run("sends event to watch channel", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "watch-jobs"
		tags := []string{"watch"}
		startTime := time.Now()
		endTime := time.Now().Add(100 * time.Millisecond)

		// Receive from channel in goroutine
		// 在goroutine中从通道接收
		type result struct {
			event JobWatchInterface
			ok    bool
		}
		resultChan := make(chan result, 1)
		go func() {
			event, ok := <-monitor.jobChan
			resultChan <- result{event, ok}
		}()

		monitor.RecordJobTimingWithStatus(startTime, endTime, jobID, jobName, tags, gocron.Success, nil)

		// Wait for event or timeout
		// 等待事件或超时
		select {
		case res := <-resultChan:
			assert.True(t, res.ok, "should receive event from channel")
			assert.NotNil(t, res.event, "event should not be nil")
		case <-time.After(1 * time.Second):
			t.Fatal("should have received event within 1 second")
		}
	})

	t.Run("generates unique event IDs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "unique-id-jobs"
		tags := []string{"unique"}

		// Record multiple events
		// 记录多个事件
		eventIDs := make(map[string]bool)
		for i := 0; i < 10; i++ {
			startTime := time.Now()
			endTime := time.Now().Add(10 * time.Millisecond)
			monitor.RecordJobTimingWithStatus(startTime, endTime, jobID, jobName, tags, gocron.Success, nil)

			events := monitor.GetJobEvents(jobID.String())
			eventIDs[events[len(events)-1].EventID] = true
			time.Sleep(5 * time.Millisecond) // Ensure different timestamps / 确保不同的时间戳
		}

		assert.Equal(t, 10, len(eventIDs), "all event IDs should be unique")
	})
}

// TestUpdateJobEvents tests the UpdateJobEvents method.
// TestUpdateJobEvents 测试 UpdateJobEvents 方法。
//
// Test scenarios:
// - Happy path: update jobs events
// - Boundary conditions: first update, max records limit
// - Exception cases: circular buffer behavior
func TestUpdateJobEvents(t *testing.T) {
	t.Parallel()

	t.Run("creates new jobs record", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(5))
		jobID := uuid.New()
		jobName := "new-jobs"
		tags := []string{"new"}
		newEvent := &JobEvent{
			EventID:   "event-1",
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Status:    gocron.Success,
		}

		monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)

		events := monitor.GetJobEvents(jobID.String())
		require.Len(t, events, 1, "should have one event")
		assert.Equal(t, newEvent, events[0], "event should match")

		spec, ok := monitor.jobRecord[jobID.String()]
		assert.True(t, ok, "jobs record should exist")
		assert.Equal(t, jobID.String(), spec.JobSpec.JobID, "jobs ID should match")
		assert.Equal(t, jobName, spec.JobSpec.JobName, "jobs name should match")
		assert.Equal(t, tags, spec.JobSpec.Tags, "tags should match")
	})

	t.Run("appends events until max records", func(t *testing.T) {
		t.Parallel()

		maxRecords := 3
		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(maxRecords))
		jobID := uuid.New()
		jobName := "append-jobs"
		tags := []string{"append"}

		for i := 0; i < maxRecords; i++ {
			newEvent := &JobEvent{
				EventID:   "event-" + string(rune('A'+i)),
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Status:    gocron.Success,
			}
			monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)
		}

		events := monitor.GetJobEvents(jobID.String())
		assert.Len(t, events, maxRecords, "should have maxRecords events")
	})

	t.Run("implements circular buffer after max records", func(t *testing.T) {
		t.Parallel()

		maxRecords := 3
		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(maxRecords))
		jobID := uuid.New()
		jobName := "circular-jobs"
		tags := []string{"circular"}

		// Add maxRecords events
		// 添加maxRecords个事件
		eventIDs := make([]string, maxRecords+2)
		for i := 0; i < maxRecords+2; i++ {
			eventID := "event-" + string(rune('0'+i))
			eventIDs[i] = eventID
			newEvent := &JobEvent{
				EventID:   eventID,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Status:    gocron.Success,
			}
			monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)
			time.Sleep(5 * time.Millisecond) // Ensure different timestamps / 确保不同的时间戳
		}

		events := monitor.GetJobEvents(jobID.String())
		assert.Len(t, events, maxRecords, "should only keep maxRecords events")

		// First event should be removed (circular buffer)
		// 第一个事件应该被移除（环形缓冲区）
		assert.NotEqual(t, eventIDs[0], events[0].EventID, "first event should be removed")

		// Last event should be present
		// 最后一个事件应该存在
		assert.Equal(t, eventIDs[len(eventIDs)-1], events[len(events)-1].EventID, "last event should be present")
	})

	t.Run("handles empty tags", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New()
		jobName := "no-tags-jobs"
		newEvent := &JobEvent{
			EventID:   "event-no-tags",
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Status:    gocron.Success,
		}

		monitor.UpdateJobEvents(jobID, jobName, newEvent)

		events := monitor.GetJobEvents(jobID.String())
		require.Len(t, events, 1, "should have one event")
	})
}

// TestGetJobEvents tests the GetJobEvents method.
// TestGetJobEvents 测试 GetJobEvents 方法。
//
// Test scenarios:
// - Happy path: get jobs events
// - Boundary conditions: no events, empty jobs record
// - Exception cases: non-existent jobs
func TestGetJobEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns nil for non-existent jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		nonExistentJobID := uuid.New().String()

		events := monitor.GetJobEvents(nonExistentJobID)
		assert.Nil(t, events, "should return nil for non-existent jobs")
	})

	t.Run("returns nil for empty monitor", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		// Don't add any events
		// 不添加任何事件

		events := monitor.GetJobEvents(uuid.New().String())
		assert.Nil(t, events, "should return nil when no jobs records exist")
	})

	t.Run("returns all events for existing jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(10))
		jobID := uuid.New()
		jobName := "get-events-jobs"
		tags := []string{"get"}

		// Add multiple events
		// 添加多个事件
		expectedEvents := []*JobEvent{
			{EventID: "event-1", StartTime: time.Now(), EndTime: time.Now(), Status: gocron.Success},
			{EventID: "event-2", StartTime: time.Now(), EndTime: time.Now(), Status: gocron.Success},
			{EventID: "event-3", StartTime: time.Now(), EndTime: time.Now(), Status: gocron.Success},
		}

		for _, event := range expectedEvents {
			monitor.UpdateJobEvents(jobID, jobName, event, tags...)
		}

		events := monitor.GetJobEvents(jobID.String())
		require.Len(t, events, len(expectedEvents), "should return all events")
	})

	t.Run("returns events in correct order", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(5))
		jobID := uuid.New()
		jobName := "ordered-jobs"
		tags := []string{"ordered"}

		// Add events in specific order
		// 按特定顺序添加事件
		eventIDs := []string{"first", "second", "third"}
		for _, id := range eventIDs {
			newEvent := &JobEvent{
				EventID:   id,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Status:    gocron.Success,
			}
			monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)
			time.Sleep(5 * time.Millisecond)
		}

		events := monitor.GetJobEvents(jobID.String())
		require.Len(t, events, len(eventIDs), "should have all events")

		// Verify order
		// 验证顺序
		for i, event := range events {
			assert.Equal(t, eventIDs[i], event.EventID, "event %d should be in correct order", i)
		}
	})
}

// TestWatch tests the Watch method.
// TestWatch 测试 Watch 方法。
//
// Test scenarios:
// - Happy path: get watch channel
// - Boundary conditions: channel capacity
func TestWatch(t *testing.T) {
	t.Parallel()

	t.Run("returns jobs channel", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		channel := monitor.Watch()

		require.NotNil(t, channel, "channel should not be nil")
	})

	t.Run("channel has correct capacity", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		channel := monitor.Watch()

		// Default capacity is 100
		// 默认容量是100
		for i := 0; i < 100; i++ {
			select {
			case channel <- &MonitorJobSpec{}:
				// OK
			default:
				t.Fatalf("channel should accept at least 100 items, failed at %d", i)
			}
		}
	})
}

// TestMonitorJobSpec tests the MonitorJobSpec struct and methods.
// TestMonitorJobSpec 测试 MonitorJobSpec 结构体和方法。
func TestMonitorJobSpec(t *testing.T) {
	t.Parallel()

	t.Run("GetJobID returns correct jobs ID", func(t *testing.T) {
		t.Parallel()

		spec := MonitorJobSpec{
			JobSpec: JobSpec{
				JobID: "test-jobs-id-123",
			},
		}

		assert.Equal(t, "test-jobs-id-123", spec.GetJobID(), "should return correct jobs ID")
	})

	t.Run("GetJobName returns correct jobs name", func(t *testing.T) {
		t.Parallel()

		spec := MonitorJobSpec{
			JobSpec: JobSpec{
				JobName: "test-jobs-name",
			},
		}

		assert.Equal(t, "test-jobs-name", spec.GetJobName(), "should return correct jobs name")
	})

	t.Run("GetTags returns correct tags", func(t *testing.T) {
		t.Parallel()

		tags := []string{"tag1", "tag2", "tag3"}
		spec := MonitorJobSpec{
			JobSpec: JobSpec{
				Tags: tags,
			},
		}

		assert.Equal(t, tags, spec.GetTags(), "should return correct tags")
	})

	t.Run("GetCurrentEvent returns latest event", func(t *testing.T) {
		t.Parallel()

		events := []*JobEvent{
			{EventID: "event-1", StartTime: time.Now(), EndTime: time.Now(), Status: gocron.Success},
			{EventID: "event-2", StartTime: time.Now(), EndTime: time.Now(), Status: gocron.Success},
			{EventID: "event-3", StartTime: time.Now(), EndTime: time.Now(), Status: gocron.Success},
		}
		spec := MonitorJobSpec{
			JobEvents: events,
		}

		currentEvent := spec.GetCurrentEvent()
		require.NotNil(t, currentEvent, "should return current event")
		assert.Equal(t, "event-3", currentEvent.EventID, "should return latest event")
	})

	t.Run("GetCurrentEvent returns nil when no events", func(t *testing.T) {
		t.Parallel()

		spec := MonitorJobSpec{
			JobEvents: []*JobEvent{},
		}

		assert.Nil(t, spec.GetCurrentEvent(), "should return nil when no events")
	})

	t.Run("implements JobWatchInterface", func(t *testing.T) {
		t.Parallel()

		var _ JobWatchInterface = &MonitorJobSpec{}
	})
}

// TestJobEvent tests the JobEvent struct and methods.
// TestJobEvent 测试 JobEvent 结构体和方法。
func TestJobEvent(t *testing.T) {
	t.Parallel()

	t.Run("GetStartTime returns correct start time", func(t *testing.T) {
		t.Parallel()

		expectedTime := time.Now()
		event := JobEvent{
			StartTime: expectedTime,
		}

		assert.WithinDuration(t, expectedTime, event.GetStartTime(), time.Microsecond, "should return correct start time")
	})

	t.Run("GetEndTime returns correct end time", func(t *testing.T) {
		t.Parallel()

		expectedTime := time.Now()
		event := JobEvent{
			EndTime: expectedTime,
		}

		assert.WithinDuration(t, expectedTime, event.GetEndTime(), time.Microsecond, "should return correct end time")
	})

	t.Run("GetStatus returns correct status", func(t *testing.T) {
		t.Parallel()

		event := JobEvent{
			Status: gocron.Success,
		}

		assert.Equal(t, gocron.Success, event.GetStatus(), "should return correct status")
	})

	t.Run("GetSpendTime returns correct duration", func(t *testing.T) {
		t.Parallel()

		startTime := time.Now()
		endTime := startTime.Add(100 * time.Millisecond)
		event := JobEvent{
			StartTime: startTime,
			EndTime:   endTime,
		}

		spendTime := event.GetSpendTime()
		assert.GreaterOrEqual(t, spendTime, int64(99), "spend time should be at least 99ms")
		assert.LessOrEqual(t, spendTime, int64(150), "spend time should be at most 150ms")
	})

	t.Run("GetError returns correct error", func(t *testing.T) {
		t.Parallel()

		expectedError := assert.AnError
		event := JobEvent{
			Err: expectedError,
		}

		assert.Equal(t, expectedError, event.GetError(), "should return correct error")
	})

	t.Run("MarshalJSON produces valid JSON", func(t *testing.T) {
		t.Parallel()

		startTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 15, 14, 31, 0, 0, time.UTC)
		event := JobEvent{
			EventID:   "test-event-123",
			StartTime: startTime,
			EndTime:   endTime,
			Status:    gocron.Success,
			Err:       nil,
		}

		jsonData, err := json.Marshal(event)
		require.NoError(t, err, "should marshal without error")

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		require.NoError(t, err, "should unmarshal without error")

		assert.Equal(t, "test-event-123", result["event_id"], "JSON should contain event_id")
		assert.Contains(t, result["start_time"], "2024-01-15", "JSON should contain formatted start_time")
		assert.Contains(t, result["end_time"], "2024-01-15", "JSON should contain formatted end_time")
	})

	t.Run("MarshalJSON includes error message", func(t *testing.T) {
		t.Parallel()

		event := JobEvent{
			EventID:   "error-event",
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Status:    gocron.Fail,
			Err:       assert.AnError,
		}

		jsonData, err := json.Marshal(event)
		require.NoError(t, err, "should marshal without error")

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		require.NoError(t, err, "should unmarshal without error")

		assert.NotEmpty(t, result["error"], "JSON should contain error message")
	})
}

// BenchmarkIncrementJob benchmarks the IncrementJob method.
// BenchmarkIncrementJob 对 IncrementJob 方法进行基准测试。
func BenchmarkIncrementJob(b *testing.B) {
	monitor := NewDefaultSchedulerMonitor()
	jobID := uuid.New()
	jobName := "benchmark-jobs"
	tags := []string{"benchmark"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		monitor.IncrementJob(jobID, jobName, tags, gocron.Success)
	}
}

// BenchmarkRecordJobTiming benchmarks the RecordJobTiming method.
// BenchmarkRecordJobTiming 对 RecordJobTiming 方法进行基准测试。
func BenchmarkRecordJobTiming(b *testing.B) {
	monitor := NewDefaultSchedulerMonitor()
	jobID := uuid.New()
	jobName := "benchmark-jobs"
	tags := []string{"benchmark"}
	startTime := time.Now()
	endTime := time.Now()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		monitor.RecordJobTiming(startTime, endTime, jobID, jobName, tags)
	}
}

// BenchmarkUpdateJobEvents benchmarks the UpdateJobEvents method.
// BenchmarkUpdateJobEvents 对 UpdateJobEvents 方法进行基准测试。
func BenchmarkUpdateJobEvents(b *testing.B) {
	monitor := NewDefaultSchedulerMonitor(WithMaxRecords(100))
	jobID := uuid.New()
	jobName := "benchmark-jobs"
	tags := []string{"benchmark"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		newEvent := &JobEvent{
			EventID:   "event-" + string(rune('0'+i%10)),
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Status:    gocron.Success,
		}
		monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)
	}
}

// TestGetRetryHistory tests the GetRetryHistory method.
// TestGetRetryHistory 测试 GetRetryHistory 方法。
//
// Test scenarios:
// - Happy path: get retry history for a jobs
// - Boundary conditions: no retry history, empty history
// - Exception cases: non-existent jobs
func TestGetRetryHistory(t *testing.T) {
	t.Parallel()

	t.Run("returns nil for non-existent jobs", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		nonExistentJobID := "non-existent-jobs"

		history := monitor.GetRetryHistory(nonExistentJobID)
		assert.Nil(t, history, "should return nil for non-existent jobs")
	})

	t.Run("returns empty history for jobs with no retries", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New().String()

		history := monitor.GetRetryHistory(jobID)
		assert.Nil(t, history, "should return nil for jobs with no retry history")
	})

	t.Run("returns retry history for jobs with retries", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New().String()

		// Add retry events
		// 添加重试事件
		now := time.Now()
		events := []*retry.RetryEvent{
			{
				EventID:         "retry-1",
				OriginalEventID: jobID,
				Attempt:         0,
				StartTime:       now.Add(-3 * time.Second),
				EndTime:         now.Add(-2 * time.Second),
				NextRetryIn:     1 * time.Second,
			},
			{
				EventID:         "retry-2",
				OriginalEventID: jobID,
				Attempt:         1,
				StartTime:       now.Add(-1 * time.Second),
				EndTime:         now,
				NextRetryIn:     2 * time.Second,
			},
		}

		for _, event := range events {
			monitor.RecordRetryEvent(event)
		}

		history := monitor.GetRetryHistory(jobID)
		require.Len(t, history, len(events), "should return all retry events")
	})

	t.Run("returns a copy not the original slice", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New().String()

		event := &retry.RetryEvent{
			EventID:         "retry-1",
			OriginalEventID: jobID,
			Attempt:         0,
			StartTime:       time.Now(),
			EndTime:         time.Now(),
			NextRetryIn:     1 * time.Second,
		}
		monitor.RecordRetryEvent(event)

		history1 := monitor.GetRetryHistory(jobID)
		history2 := monitor.GetRetryHistory(jobID)

		// Modify the first copy
		// 修改第一个副本
		history1[0] = &retry.RetryEvent{EventID: "modified"}

		// Second copy should not be affected
		// 第二个副本不应受影响
		assert.NotEqual(t, history1[0].EventID, history2[0].EventID, "should return independent copies")
	})
}

// TestRecordRetryEvent tests the RecordRetryEvent method.
// TestRecordRetryEvent 测试 RecordRetryEvent 方法。
//
// Test scenarios:
// - Happy path: record retry events
// - Boundary conditions: max retry history limit
// - Exception cases: circular buffer behavior
func TestRecordRetryEvent(t *testing.T) {
	t.Parallel()

	t.Run("records single retry event", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New().String()

		event := &retry.RetryEvent{
			EventID:         "retry-1",
			OriginalEventID: jobID,
			Attempt:         0,
			StartTime:       time.Now(),
			EndTime:         time.Now().Add(1 * time.Second),
			NextRetryIn:     5 * time.Second,
		}

		monitor.RecordRetryEvent(event)

		history := monitor.GetRetryHistory(jobID)
		require.Len(t, history, 1, "should have one retry event")
		assert.Equal(t, event.EventID, history[0].EventID, "event ID should match")
	})

	t.Run("records multiple retry events", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New().String()

		// Add multiple retry events
		// 添加多个重试事件
		for i := 0; i < 5; i++ {
			event := &retry.RetryEvent{
				EventID:         fmt.Sprintf("retry-%d", i),
				OriginalEventID: jobID,
				Attempt:         i,
				StartTime:       time.Now(),
				EndTime:         time.Now(),
				NextRetryIn:     time.Duration(i) * time.Second,
			}
			monitor.RecordRetryEvent(event)
		}

		history := monitor.GetRetryHistory(jobID)
		assert.Len(t, history, 5, "should have 5 retry events")
	})

	t.Run("implements circular buffer at max retry history", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor()
		jobID := uuid.New().String()

		// Add more than maxRetryHistory (100) events
		// 添加超过 maxRetryHistory (100) 个事件
		const maxRetryHistory = 100
		for i := 0; i < maxRetryHistory+5; i++ {
			event := &retry.RetryEvent{
				EventID:         fmt.Sprintf("retry-%d", i),
				OriginalEventID: jobID,
				Attempt:         i,
				StartTime:       time.Now(),
				EndTime:         time.Now(),
				NextRetryIn:     1 * time.Second,
			}
			monitor.RecordRetryEvent(event)
		}

		history := monitor.GetRetryHistory(jobID)
		assert.Len(t, history, maxRetryHistory, "should only keep maxRetryHistory events")

		// First event should be removed
		// 第一个事件应该被移除
		assert.NotEqual(t, "retry-0", history[0].EventID, "first event should be removed")

		// Last event should be present
		// 最后一个事件应该存在
		assert.Equal(t, fmt.Sprintf("retry-%d", maxRetryHistory+4), history[len(history)-1].EventID, "last event should be present")
	})
}

// TestUpdateJobEventsConcurrency tests concurrent access to UpdateJobEvents.
// TestUpdateJobEventsConcurrency 测试 UpdateJobEvents 的并发访问。
func TestUpdateJobEventsConcurrency(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping concurrency test in short mode")
	}

	t.Run("concurrent updates are safe", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(100))
		jobID := uuid.New()
		jobName := "concurrent-jobs"
		tags := []string{"concurrent"}

		const goroutines = 50
		const updatesPerGoroutine = 20

		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				for j := 0; j < updatesPerGoroutine; j++ {
					newEvent := &JobEvent{
						EventID:   fmt.Sprintf("event-%d-%d", index, j),
						StartTime: time.Now(),
						EndTime:   time.Now(),
						Status:    gocron.Success,
					}
					monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)
				}
			}(i)
		}

		wg.Wait()

		// Verify no data race occurred
		// 验证没有发生数据竞争
		events := monitor.GetJobEvents(jobID.String())
		assert.LessOrEqual(t, len(events), 100, "should not exceed maxRecords")
	})
}

// TestGetJobEventsConcurrency tests concurrent access to GetJobEvents.
// TestGetJobEventsConcurrency 测试 GetJobEvents 的并发访问。
func TestGetJobEventsConcurrency(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping concurrency test in short mode")
	}

	t.Run("concurrent reads and writes are safe", func(t *testing.T) {
		t.Parallel()

		monitor := NewDefaultSchedulerMonitor(WithMaxRecords(50))
		jobID := uuid.New()
		jobName := "race-jobs"
		tags := []string{"race"}

		// First add an initial event
		// 首先添加一个初始事件
		initialEvent := &JobEvent{
			EventID:   "initial",
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Status:    gocron.Success,
		}
		monitor.UpdateJobEvents(jobID, jobName, initialEvent, tags...)

		const goroutines = 50
		var wg sync.WaitGroup

		// Writer goroutines
		// 写入 goroutine
		for i := 0; i < goroutines/2; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					newEvent := &JobEvent{
						EventID:   fmt.Sprintf("write-%d-%d", index, j),
						StartTime: time.Now(),
						EndTime:   time.Now(),
						Status:    gocron.Success,
					}
					monitor.UpdateJobEvents(jobID, jobName, newEvent, tags...)
				}
			}(i)
		}

		// Reader goroutines
		// 读取 goroutine
		for i := 0; i < goroutines/2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					_ = monitor.GetJobEvents(jobID.String())
				}
			}()
		}

		wg.Wait()

		// Verify no panic or data race
		// 验证没有 panic 或数据竞争
		events := monitor.GetJobEvents(jobID.String())
		assert.LessOrEqual(t, len(events), 50, "should not exceed maxRecords")
	})
}

// TestJobEventRetryFields tests the retry-related fields in JobEvent.
// TestJobEventRetryFields 测试 JobEvent 中的重试相关字段。
func TestJobEventRetryFields(t *testing.T) {
	t.Parallel()

	t.Run("JobEvent includes retry fields in JSON", func(t *testing.T) {
		t.Parallel()

		startTime := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 15, 14, 31, 0, 0, time.UTC)
		event := JobEvent{
			EventID:         "retry-event-123",
			StartTime:       startTime,
			EndTime:         endTime,
			Status:          gocron.Fail,
			Err:             assert.AnError,
			RetryCount:      3,
			IsRetry:         true,
			OriginalEventID: "original-event-456",
		}

		jsonData, err := json.Marshal(event)
		require.NoError(t, err, "should marshal without error")

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		require.NoError(t, err, "should unmarshal without error")

		assert.Equal(t, float64(3), result["retry_count"], "JSON should contain retry_count")
		assert.Equal(t, true, result["is_retry"], "JSON should contain is_retry")
		assert.Equal(t, "original-event-456", result["original_event_id"], "JSON should contain original_event_id")
	})

	t.Run("JobEvent omits empty retry fields from JSON", func(t *testing.T) {
		t.Parallel()

		event := JobEvent{
			EventID:   "normal-event",
			StartTime: time.Now(),
			EndTime:   time.Now(),
			Status:    gocron.Success,
			Err:       nil,
			// Retry fields are zero values / 重试字段为零值
		}

		jsonData, err := json.Marshal(event)
		require.NoError(t, err, "should marshal without error")

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		require.NoError(t, err, "should unmarshal without error")

		_, hasRetryCount := result["retry_count"]
		_, hasIsRetry := result["is_retry"]
		_, hasOriginalEventID := result["original_event_id"]

		assert.False(t, hasRetryCount, "should omit retry_count when zero")
		assert.False(t, hasIsRetry, "should omit is_retry when false")
		assert.False(t, hasOriginalEventID, "should omit original_event_id when empty")
	})
}
