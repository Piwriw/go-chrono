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

	// Test ExponentialBackoffPolicy implements RetryPolicy interface
	// 测试 ExponentialBackoffPolicy 实现 RetryPolicy 接口
	var _ RetryPolicy = policy

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

// TestExponentialBackoffPolicy_ShouldRetry tests ShouldRetry method
// TestExponentialBackoffPolicy_ShouldRetry 测试 ShouldRetry 方法
func TestExponentialBackoffPolicy_ShouldRetry(t *testing.T) {
	policy := NewExponentialBackoffPolicy(1*time.Second, 60*time.Second)

	t.Run("ShouldRetry always returns true", func(t *testing.T) {
		tests := []int{0, 1, 10, 100}
		for _, attempt := range tests {
			if !policy.ShouldRetry(attempt, errors.New("test error")) {
				t.Errorf("ShouldRetry(%d) should return true", attempt)
			}
		}
	})
}

// TestExponentialBackoffPolicy_BoundaryConditions tests boundary conditions
// TestExponentialBackoffPolicy_BoundaryConditions 测试边界条件
func TestExponentialBackoffPolicy_BoundaryConditions(t *testing.T) {
	t.Run("attempt = 0 returns base interval", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(2*time.Second, 100*time.Second)
		result := policy.NextRetry(0, errors.New("test"))
		expected := 2 * time.Second
		if result != expected {
			t.Errorf("NextRetry(0) = %v, want %v", result, expected)
		}
	})

	t.Run("negative attempt returns base interval", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(2*time.Second, 100*time.Second)
		result := policy.NextRetry(-1, errors.New("test"))
		expected := 2 * time.Second
		if result != expected {
			t.Errorf("NextRetry(-1) = %v, want %v", result, expected)
		}
	})

	t.Run("very large attempt does not overflow", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(1*time.Second, 60*time.Second)
		// attempt = 100 would cause overflow without the fix
		// attempt = 100 会在没有修复的情况下导致溢出
		result := policy.NextRetry(100, errors.New("test"))
		expected := 60 * time.Second // Should be capped at max / 应该被限制为 max
		if result != expected {
			t.Errorf("NextRetry(100) = %v, want %v", result, expected)
		}
	})

	t.Run("attempt at boundary (62) works correctly", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(1*time.Second, 3600*time.Second)
		result := policy.NextRetry(62, errors.New("test"))
		// 1 << 62 should not overflow
		// 1 << 62 不应该溢出
		if result <= 0 {
			t.Error("NextRetry(62) returned invalid duration (possible overflow)")
		}
	})
}

// TestExponentialBackoffPolicy_ParameterValidation tests parameter validation
// TestExponentialBackoffPolicy_ParameterValidation 测试参数验证
func TestExponentialBackoffPolicy_ParameterValidation(t *testing.T) {
	t.Run("zero base uses default", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(0, 60*time.Second)
		result := policy.NextRetry(0, errors.New("test"))
		expected := 1 * time.Second // Default base / 默认 base
		if result != expected {
			t.Errorf("NextRetry() with zero base = %v, want %v", result, expected)
		}
	})

	t.Run("negative base uses default", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(-1*time.Second, 60*time.Second)
		result := policy.NextRetry(0, errors.New("test"))
		expected := 1 * time.Second // Default base / 默认 base
		if result != expected {
			t.Errorf("NextRetry() with negative base = %v, want %v", result, expected)
		}
	})

	t.Run("zero max uses default", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(1*time.Second, 0)
		result := policy.NextRetry(100, errors.New("test"))
		expected := 60 * time.Second // Default max / 默认 max
		if result != expected {
			t.Errorf("NextRetry() with zero max = %v, want %v", result, expected)
		}
	})

	t.Run("base > max swaps values", func(t *testing.T) {
		policy := NewExponentialBackoffPolicy(60*time.Second, 1*time.Second)
		result := policy.NextRetry(0, errors.New("test"))
		expected := 1 * time.Second // Should use swapped base / 应该使用交换后的 base
		if result != expected {
			t.Errorf("NextRetry() with base > max = %v, want %v", result, expected)
		}
	})
}

// TestJitterPolicy tests the jittered retry policy
// TestJitterPolicy 测试带抖动的重试策略
func TestJitterPolicy(t *testing.T) {
	base := NewFixedIntervalPolicy(10 * time.Second)
	jitterFactor := 0.2
	policy := NewJitterPolicy(base, jitterFactor)

	t.Run("NextRetry adds jitter to base interval", func(t *testing.T) {
		baseInterval := 10 * time.Second
		minJitter := time.Duration(float64(baseInterval) * (1 - jitterFactor))
		maxJitter := time.Duration(float64(baseInterval) * (1 + jitterFactor))

		// Multiple tests ensure randomness is within reasonable range
		// 多次测试确保随机性在合理范围内
		for range 10 {
			got := policy.NextRetry(0, errors.New("test"))
			if got < minJitter || got > maxJitter {
				t.Errorf("NextRetry() = %v, want between %v and %v",
					got, minJitter, maxJitter)
			}
		}
	})
}
