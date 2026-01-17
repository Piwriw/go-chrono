// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"errors"
	"testing"
	"time"
)

// TestRetryPolicy_Interface tests the RetryPolicy interface implementation
// TestRetryPolicy_Interface 测试 RetryPolicy 接口实现
func TestRetryPolicy_Interface(t *testing.T) {
	// Test FixedIntervalPolicy implements RetryPolicy interface
	// 测试 FixedIntervalPolicy 实现 RetryPolicy 接口
	policy := NewFixedIntervalPolicy(5 * time.Second)

	var _ RetryPolicy = policy

	t.Run("NextRetry returns fixed interval", func(t *testing.T) {
		result := policy.NextRetry(0, errors.New("test error"))
		expected := 5 * time.Second
		if result != expected {
			t.Errorf("NextRetry() = %v, want %v", result, expected)
		}
	})

	t.Run("ShouldRetry returns true when under max attempts", func(t *testing.T) {
		if !policy.ShouldRetry(0, errors.New("test error")) {
			t.Error("ShouldRetry() should return true for attempt 0")
		}
	})
}

// TestExponentialBackoffPolicy tests the exponential backoff retry policy
// TestExponentialBackoffPolicy 测试指数退避重试策略
func TestExponentialBackoffPolicy(t *testing.T) {
	base := 1 * time.Second
	max := 60 * time.Second
	policy := NewExponentialBackoffPolicy(base, max)

	t.Run("NextRetry increases exponentially", func(t *testing.T) {
		tests := []struct {
			attempt int
			wantMin time.Duration
			wantMax time.Duration
		}{
			{0, 1 * time.Second, 1 * time.Second},
			{1, 2 * time.Second, 2 * time.Second},
			{2, 4 * time.Second, 4 * time.Second},
			{3, 8 * time.Second, 8 * time.Second},
			{6, 60 * time.Second, 60 * time.Second}, // 限制为 max
		}

		for _, tt := range tests {
			got := policy.NextRetry(tt.attempt, errors.New("test"))
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("NextRetry(%d) = %v, want between %v and %v",
					tt.attempt, got, tt.wantMin, tt.wantMax)
			}
		}
	})
}
