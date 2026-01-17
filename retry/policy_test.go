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
