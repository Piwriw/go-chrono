// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
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
	multiplier   float64
}

// NewExponentialBackoffPolicy creates an exponential backoff retry policy
// NewExponentialBackoffPolicy 创建指数退避策略
func NewExponentialBackoffPolicy(base, max time.Duration) *ExponentialBackoffPolicy {
	return &ExponentialBackoffPolicy{
		baseInterval: base,
		maxInterval:  max,
		multiplier:   2.0,
	}
}

// NextRetry calculates exponentially increasing interval
// NextRetry 计算指数增长的间隔时间
func (p *ExponentialBackoffPolicy) NextRetry(attempt int, lastErr error) time.Duration {
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
