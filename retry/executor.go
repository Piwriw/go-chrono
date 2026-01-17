// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// RetryExecutor 重试执行器
// RetryExecutor handles task execution with retry logic
type RetryExecutor struct {
	config  *RetryConfig
	jobID   uuid.UUID
	jobName string
}

// NewRetryExecutor 创建重试执行器
// NewRetryExecutor creates a retry executor
func NewRetryExecutor(config *RetryConfig, jobID uuid.UUID, jobName string) *RetryExecutor {
	return &RetryExecutor{
		config:  config,
		jobID:   jobID,
		jobName: jobName,
	}
}

// ExecuteWithRetry 执行任务，支持重试
// ExecuteWithRetry executes a task with retry support
func (e *RetryExecutor) ExecuteWithRetry(ctx context.Context, task func() error) error {
	var lastErr error
	attempt := 0

	for attempt <= e.config.MaxRetries {
		startTime := time.Now()

		// 执行任务
		// Execute task
		err := task()
		duration := time.Since(startTime)

		if err == nil {
			// 成功，记录并返回
			// Success, log and return
			slog.Info("chrono: task succeeded",
				"job_id", e.jobID,
				"job_name", e.jobName,
				"attempt", attempt,
				"duration", duration)
			return nil
		}

		lastErr = err
		slog.Warn("chrono: task failed",
			"job_id", e.jobID,
			"job_name", e.jobName,
			"attempt", attempt,
			"error", err,
			"duration", duration)

		// 检查是否应该重试
		// Check if should retry
		if attempt >= e.config.MaxRetries {
			// 达到最大重试次数
			// Reached max retries
			break
		}

		if !e.config.Policy.ShouldRetry(attempt, err) {
			// 策略认为不应重试
			// Policy says don't retry
			break
		}

		// 检查错误是否可重试
		// Check if error is retryable
		if len(e.config.RetryableErrors) > 0 && !e.isRetryableError(err) {
			break
		}

		// 计算下次重试时间
		// Calculate next retry time
		nextRetryIn := e.config.Policy.NextRetry(attempt, err)

		slog.Info("chrono: scheduling retry",
			"job_id", e.jobID,
			"job_name", e.jobName,
			"attempt", attempt,
			"next_retry_in", nextRetryIn)

		// 等待重试或上下文取消
		// Wait for retry or context cancellation
		select {
		case <-time.After(nextRetryIn):
			// 继续重试
			// Continue retry
		case <-ctx.Done():
			// 上下文被取消
			// Context canceled
			return fmt.Errorf("retry canceled: %w", ctx.Err())
		}

		attempt++
	}

	// 所有重试都失败
	// All retries failed
	if e.config.OnFinalFailure != nil {
		e.config.OnFinalFailure(ctx, e.jobID, e.jobName, lastErr)
	}

	return fmt.Errorf("task failed after %d attempts: %w", attempt, lastErr)
}

// isRetryableError 检查错误是否可重试
// isRetryableError checks if error is retryable
func (e *RetryExecutor) isRetryableError(err error) bool {
	for _, retryableErr := range e.config.RetryableErrors {
		if errors.Is(err, retryableErr) {
			return true
		}
	}
	return false
}
