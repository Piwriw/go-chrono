// Package monitor provides integration tests for retry mechanism.
// 包 monitor 提供重试机制的集成测试。
package monitor

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/piwriw/go-chrono/retry"
	"github.com/stretchr/testify/assert"
)

// TestRetryPolicyCreation tests various retry policy creation methods.
// TestRetryPolicyCreation 测试各种重试策略创建方法。
func TestRetryPolicyCreation(t *testing.T) {
	t.Parallel()

	t.Run("FixedIntervalPolicy creates valid policy", func(t *testing.T) {
		t.Parallel()

		policy := retry.NewFixedIntervalPolicy(5 * time.Second)
		assert.NotNil(t, policy, "policy should not be nil")

		// Test NextRetry
		interval := policy.NextRetry(0, nil)
		assert.Equal(t, 5*time.Second, interval, "interval should match")

		// Test ShouldRetry
		assert.True(t, policy.ShouldRetry(0, nil), "should retry")
	})

	t.Run("ExponentialBackoffPolicy creates valid policy", func(t *testing.T) {
		t.Parallel()

		policy := retry.NewExponentialBackoffPolicy(1*time.Second, 60*time.Second)
		assert.NotNil(t, policy, "policy should not be nil")

		// Test exponential progression
		interval0 := policy.NextRetry(0, nil)
		interval1 := policy.NextRetry(1, nil)
		interval2 := policy.NextRetry(2, nil)

		assert.Equal(t, 1*time.Second, interval0, "first interval should be base")
		assert.Equal(t, 2*time.Second, interval1, "second interval should double")
		assert.Equal(t, 4*time.Second, interval2, "third interval should double again")
	})

	t.Run("JitterPolicy wraps base policy", func(t *testing.T) {
		t.Parallel()

		basePolicy := retry.NewFixedIntervalPolicy(10 * time.Second)
		jitterPolicy := retry.NewJitterPolicy(basePolicy, 0.2)
		assert.NotNil(t, jitterPolicy, "jitter policy should not be nil")

		// Test that jitter produces values near base
		baseInterval := basePolicy.NextRetry(0, nil)
		jitterInterval := jitterPolicy.NextRetry(0, nil)

		// Should be within ±20% of base
		minJitter := time.Duration(float64(baseInterval) * 0.8)
		maxJitter := time.Duration(float64(baseInterval) * 1.2)
		assert.GreaterOrEqual(t, jitterInterval, minJitter, "jitter should not be below minimum")
		assert.LessOrEqual(t, jitterInterval, maxJitter, "jitter should not exceed maximum")
	})
}

// TestRetryConfigCreation tests RetryConfig creation.
// TestRetryConfigCreation 测试 RetryConfig 创建。
func TestRetryConfigCreation(t *testing.T) {
	t.Parallel()

	t.Run("RetryConfig can be created with various policies", func(t *testing.T) {
		t.Parallel()

		// Fixed interval
		config1 := &retry.RetryConfig{
			MaxRetries: 5,
			Policy:     retry.NewFixedIntervalPolicy(2 * time.Second),
		}
		assert.NotNil(t, config1, "fixed interval config should not be nil")

		// Exponential backoff
		config2 := &retry.RetryConfig{
			MaxRetries: 3,
			Policy:     retry.NewExponentialBackoffPolicy(1*time.Second, 10*time.Second),
		}
		assert.NotNil(t, config2, "exponential backoff config should not be nil")

		// Jitter
		config3 := &retry.RetryConfig{
			MaxRetries: 2,
			Policy: retry.NewJitterPolicy(
				retry.NewFixedIntervalPolicy(1*time.Second),
				0.3,
			),
		}
		assert.NotNil(t, config3, "jitter config should not be nil")
	})
}

// TestRetryEventCreation tests RetryEvent creation and methods.
// TestRetryEventCreation 测试 RetryEvent 创建和方法。
func TestRetryEventCreation(t *testing.T) {
	t.Parallel()

	t.Run("NewRetryEvent creates valid event", func(t *testing.T) {
		t.Parallel()

		startTime := time.Now()
		endTime := startTime.Add(5 * time.Second)
		testErr := errors.New("test error")

		event := retry.NewRetryEvent(
			"event-123",
			"original-event-456",
			2,
			&startTime,
			&endTime,
			testErr,
			10*time.Second,
		)

		assert.NotNil(t, event, "event should not be nil")
		assert.Equal(t, "event-123", event.EventID, "event ID should match")
		assert.Equal(t, "original-event-456", event.OriginalEventID, "original event ID should match")
		assert.Equal(t, 2, event.Attempt, "attempt should match")
		assert.Equal(t, 10*time.Second, event.NextRetryIn, "next retry interval should match")

		// Test GetDuration
		duration := event.GetDuration()
		assert.Equal(t, 5*time.Second, duration, "duration should be 5 seconds")
	})
}

// TestSchedulerMonitorRetryMethods tests retry-related monitor methods.
// TestSchedulerMonitorRetryMethods 测试重试相关的监控方法。
func TestSchedulerMonitorRetryMethods(t *testing.T) {
	t.Parallel()

	schedMonitor := NewDefaultSchedulerMonitor()

	t.Run("GetRetryHistory returns nil for non-existent jobs", func(t *testing.T) {
		t.Parallel()

		history := schedMonitor.GetRetryHistory("non-existent-jobs")
		assert.Nil(t, history, "should return nil for non-existent jobs")
	})

	t.Run("RecordRetryEvent adds event to history", func(t *testing.T) {
		t.Parallel()

		jobID := "test-jobs-id"
		event := &retry.RetryEvent{
			EventID:         "retry-1",
			OriginalEventID: jobID,
			Attempt:         0,
			StartTime:       func() *time.Time { t := time.Now(); return &t }(),
			EndTime:         func() *time.Time { t := time.Now(); return &t }(),
			NextRetryIn:     1 * time.Second,
		}

		schedMonitor.RecordRetryEvent(event)

		history := schedMonitor.GetRetryHistory(jobID)
		assert.NotNil(t, history, "history should not be nil")
		assert.Len(t, history, 1, "should have one retry event")
		assert.Equal(t, "retry-1", history[0].EventID, "event ID should match")
	})

	t.Run("GetRetryHistory returns copy not reference", func(t *testing.T) {
		t.Parallel()

		jobID := "test-jobs-copy"
		event := &retry.RetryEvent{
			EventID:         "retry-copy",
			OriginalEventID: jobID,
			Attempt:         0,
			StartTime:       func() *time.Time { t := time.Now(); return &t }(),
			EndTime:         func() *time.Time { t := time.Now(); return &t }(),
			NextRetryIn:     1 * time.Second,
		}

		schedMonitor.RecordRetryEvent(event)

		history1 := schedMonitor.GetRetryHistory(jobID)
		history2 := schedMonitor.GetRetryHistory(jobID)

		// Modify first copy
		history1[0] = &retry.RetryEvent{EventID: "modified"}

		// Second copy should not be affected
		assert.NotEqual(t, history1[0].EventID, history2[0].EventID, "should return independent copies")
	})

	t.Run("WithMaxRetryHistory configures retry history limit", func(t *testing.T) {
		t.Parallel()

		// Create monitor with custom max retry history
		customSchedMonitor := NewDefaultSchedulerMonitor(WithMaxRetryHistory(3))

		jobID := "test-limit-jobs"
		// Add 5 retry events
		for i := 0; i < 5; i++ {
			event := &retry.RetryEvent{
				EventID:         fmt.Sprintf("retry-%d", i),
				OriginalEventID: jobID,
				Attempt:         i,
				StartTime:       func() *time.Time { t := time.Now(); return &t }(),
				EndTime:         func() *time.Time { t := time.Now(); return &t }(),
				NextRetryIn:     1 * time.Second,
			}
			customSchedMonitor.RecordRetryEvent(event)
		}

		history := customSchedMonitor.GetRetryHistory(jobID)
		// Should only keep the last 3 events
		assert.Len(t, history, 3, "should keep only maxRetryHistory records")
		assert.Equal(t, "retry-2", history[0].EventID, "first event should be retry-2")
		assert.Equal(t, "retry-4", history[2].EventID, "last event should be retry-4")
	})

	t.Run("WithMaxRetryHistory zero uses default", func(t *testing.T) {
		t.Parallel()

		// Create monitor with zero max retry history (should use default)
		zeroSchedMonitor := NewDefaultSchedulerMonitor(WithMaxRetryHistory(0))

		jobID := "test-zero-jobs"
		// Add 15 retry events (more than default 10)
		for i := 0; i < 15; i++ {
			event := &retry.RetryEvent{
				EventID:         fmt.Sprintf("retry-%d", i),
				OriginalEventID: jobID,
				Attempt:         i,
				StartTime:       func() *time.Time { t := time.Now(); return &t }(),
				EndTime:         func() *time.Time { t := time.Now(); return &t }(),
				NextRetryIn:     1 * time.Second,
			}
			zeroSchedMonitor.RecordRetryEvent(event)
		}

		history := zeroSchedMonitor.GetRetryHistory(jobID)
		// Should keep default 10 records
		assert.Len(t, history, 10, "should use default maxRetryHistory when zero is specified")
	})
}
