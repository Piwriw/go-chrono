// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"math/rand"
	"time"
)

// RetryPolicy defines the retry strategy interface
// RetryPolicy 定义重试策略接口
type RetryPolicy interface {
	// NextRetry calculates the interval for the next retry
	// NextRetry 计算下次重试的间隔时间
	NextRetry(attempt int, lastErr error) time.Duration

	// ShouldRetry determines whether a retry should be attempted
	// ShouldRetry 判断是否应该重试
	ShouldRetry(attempt int, err error) bool
}

// FixedIntervalPolicy retries at a fixed interval
// FixedIntervalPolicy 固定间隔重试策略
type FixedIntervalPolicy struct {
	interval time.Duration
}

// NewFixedIntervalPolicy creates a fixed interval retry policy
// NewFixedIntervalPolicy 创建固定间隔策略
func NewFixedIntervalPolicy(interval time.Duration) *FixedIntervalPolicy {
	return &FixedIntervalPolicy{
		interval: interval,
	}
}

// NextRetry returns a fixed interval
// NextRetry 返回固定的间隔时间
func (p *FixedIntervalPolicy) NextRetry(attempt int, lastErr error) time.Duration {
	return p.interval
}

// ShouldRetry always returns true (caller controls max attempts)
// ShouldRetry 始终返回 true（由调用方控制最大次数）
func (p *FixedIntervalPolicy) ShouldRetry(attempt int, err error) bool {
	return true
}

// ExponentialBackoffPolicy retries with exponential backoff
// ExponentialBackoffPolicy 指数退避重试策略
type ExponentialBackoffPolicy struct {
	baseInterval time.Duration
	maxInterval  time.Duration
}

// NewExponentialBackoffPolicy creates an exponential backoff retry policy
// NewExponentialBackoffPolicy 创建指数退避策略
func NewExponentialBackoffPolicy(base, max time.Duration) *ExponentialBackoffPolicy {
	// Validate parameters: base and max must be positive
	// 验证参数：base 和 max 必须为正数
	if base <= 0 {
		base = 1 * time.Second // Default to 1 second / 默认为 1 秒
	}
	if max <= 0 {
		max = 60 * time.Second // Default to 60 seconds / 默认为 60 秒
	}
	// Ensure base is not greater than max
	// 确保 base 不大于 max
	if base > max {
		base, max = max, base
	}
	return &ExponentialBackoffPolicy{
		baseInterval: base,
		maxInterval:  max,
	}
}

// NextRetry calculates exponentially increasing interval
// NextRetry 计算指数增长的间隔时间
func (p *ExponentialBackoffPolicy) NextRetry(attempt int, lastErr error) time.Duration {
	// Prevent integer overflow: limit attempt to prevent 1 << uint(attempt) overflow
	// 防止整数溢出：限制 attempt 值以防止 1 << uint(attempt) 溢出
	// For 64-bit systems, 1 << 62 is the maximum safe shift
	// 对于 64 位系统，1 << 62 是最大安全位移值
	maxSafeShift := 62
	if attempt > maxSafeShift {
		attempt = maxSafeShift
	}
	if attempt < 0 {
		attempt = 0
	}
	multiplier := 1 << uint(attempt)
	interval := time.Duration(float64(p.baseInterval) * float64(multiplier))
	if interval > p.maxInterval {
		return p.maxInterval
	}
	return interval
}

// ShouldRetry always returns true
// ShouldRetry 始终返回 true
func (p *ExponentialBackoffPolicy) ShouldRetry(attempt int, err error) bool {
	return true
}

// JitterPolicy wraps a retry policy with random jitter
// JitterPolicy 带随机抖动的重试策略包装器
type JitterPolicy struct {
	policy       RetryPolicy
	jitterFactor float64
}

// NewJitterPolicy creates a jittered retry policy
// NewJitterPolicy 创建带抖动的策略
func NewJitterPolicy(policy RetryPolicy, jitterFactor float64) *JitterPolicy {
	return &JitterPolicy{
		policy:       policy,
		jitterFactor: jitterFactor,
	}
}

// NextRetry returns interval with random jitter
// NextRetry 返回带随机抖动的间隔时间
func (p *JitterPolicy) NextRetry(attempt int, lastErr error) time.Duration {
	baseInterval := p.policy.NextRetry(attempt, lastErr)
	jitter := (rand.Float64()*2 - 1) * p.jitterFactor // [-jitterFactor, +jitterFactor]
	return time.Duration(float64(baseInterval) * (1 + jitter))
}

// ShouldRetry delegates to underlying policy
// ShouldRetry 委托给底层策略
func (p *JitterPolicy) ShouldRetry(attempt int, err error) bool {
	return p.policy.ShouldRetry(attempt, err)
}
