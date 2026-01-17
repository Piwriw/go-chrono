// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RetryConfig 重试配置
// RetryConfig is the retry configuration
type RetryConfig struct {
	// MaxRetries 最大重试次数（0 表示不重试）
	// MaxRetries is the maximum number of retries (0 means no retry)
	MaxRetries int

	// Policy 重试策略
	// Policy is the retry strategy
	Policy RetryPolicy

	// RetryableErrors 可重试的错误类型（为空时所有错误都重试）
	// RetryableErrors are the error types that can be retried
	RetryableErrors []error

	// OnFinalFailure 最终失败回调
	// OnFinalFailure is the callback when all retries fail
	OnFinalFailure func(ctx context.Context, jobID uuid.UUID, jobName string, err error)
}

// DefaultRetryConfig 返回默认的重试配置
// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      0,
		Policy:          NewFixedIntervalPolicy(5 * time.Second),
		RetryableErrors: nil,
	}
}
