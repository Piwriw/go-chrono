// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestRetryExecutor_ExecuteWithRetry_Success tests successful execution after retries
// TestRetryExecutor_ExecuteWithRetry_Success 测试重试后成功执行
func TestRetryExecutor_ExecuteWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries: 3,
		Policy:     NewFixedIntervalPolicy(10 * time.Millisecond),
	}

	executor := NewRetryExecutor(config, uuid.New(), "test-job")

	callCount := 0
	task := func() error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := executor.ExecuteWithRetry(ctx, task)

	if err != nil {
		t.Errorf("ExecuteWithRetry() error = %v, want nil", err)
	}
	if callCount != 3 {
		t.Errorf("task called %d times, want 3", callCount)
	}
}

// TestRetryExecutor_ExecuteWithRetry_AllRetriesFailed tests all retries failing
// TestRetryExecutor_ExecuteWithRetry_AllRetriesFailed 测试所有重试都失败
func TestRetryExecutor_ExecuteWithRetry_AllRetriesFailed(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries: 2,
		Policy:     NewFixedIntervalPolicy(10 * time.Millisecond),
	}

	executor := NewRetryExecutor(config, uuid.New(), "test-job")

	callCount := 0
	task := func() error {
		callCount++
		return errors.New("persistent error")
	}

	err := executor.ExecuteWithRetry(ctx, task)

	if err == nil {
		t.Error("ExecuteWithRetry() error = nil, want error")
	}
	if callCount != 3 { // initial + 2 retries
		t.Errorf("task called %d times, want 3", callCount)
	}
}

// TestRetryExecutor_ExecuteWithRetry_ContextCanceled tests context cancellation
// TestRetryExecutor_ExecuteWithRetry_ContextCanceled 测试上下文取消
func TestRetryExecutor_ExecuteWithRetry_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := &RetryConfig{
		MaxRetries: 5,
		Policy:     NewFixedIntervalPolicy(1 * time.Second),
	}

	executor := NewRetryExecutor(config, uuid.New(), "test-job")

	callCount := 0
	task := func() error {
		callCount++
		cancel() // Cancel context after first execution / 第一次执行后取消上下文
		return errors.New("error")
	}

	err := executor.ExecuteWithRetry(ctx, task)

	if err == nil {
		t.Error("ExecuteWithRetry() error = nil, want context canceled error")
	}
	if callCount > 2 { // Should stop quickly / 应该很快停止
		t.Errorf("task called %d times, should stop quickly", callCount)
	}
}

// TestRetryExecutor_ExecuteWithRetry_RetryableErrors tests retryable error filtering
// TestRetryExecutor_ExecuteWithRetry_RetryableErrors 测试可重试错误过滤
func TestRetryExecutor_ExecuteWithRetry_RetryableErrors(t *testing.T) {
	ctx := context.Background()

	// Define custom errors / 定义自定义错误
	temporaryErr := errors.New("temporary error")
	persistentErr := errors.New("persistent error")

	config := &RetryConfig{
		MaxRetries:       3,
		Policy:           NewFixedIntervalPolicy(10 * time.Millisecond),
		RetryableErrors:  []error{temporaryErr},
	}

	executor := NewRetryExecutor(config, uuid.New(), "test-job")

	// Test with persistent error (should not retry) / 测试持久性错误（不应重试）
	callCount := 0
	task := func() error {
		callCount++
		return persistentErr
	}

	err := executor.ExecuteWithRetry(ctx, task)

	if err == nil {
		t.Error("ExecuteWithRetry() error = nil, want error")
	}
	if callCount != 1 { // Should only call once / 应该只调用一次
		t.Errorf("task called %d times, want 1", callCount)
	}
}

// TestRetryExecutor_ExecuteWithRetry_OnFinalFailure tests final failure callback
// TestRetryExecutor_ExecuteWithRetry_OnFinalFailure 测试最终失败回调
func TestRetryExecutor_ExecuteWithRetry_OnFinalFailure(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New()
	jobName := "test-job"

	callbackCalled := false
	var callbackJobID uuid.UUID
	var callbackJobName string
	var callbackErr error

	config := &RetryConfig{
		MaxRetries: 1,
		Policy:     NewFixedIntervalPolicy(10 * time.Millisecond),
		OnFinalFailure: func(ctx context.Context, id uuid.UUID, name string, err error) {
			callbackCalled = true
			callbackJobID = id
			callbackJobName = name
			callbackErr = err
		},
	}

	executor := NewRetryExecutor(config, jobID, jobName)

	task := func() error {
		return errors.New("final error")
	}

	err := executor.ExecuteWithRetry(ctx, task)

	if err == nil {
		t.Error("ExecuteWithRetry() error = nil, want error")
	}
	if !callbackCalled {
		t.Error("OnFinalFailure callback was not called")
	}
	if callbackJobID != jobID {
		t.Errorf("callback jobID = %v, want %v", callbackJobID, jobID)
	}
	if callbackJobName != jobName {
		t.Errorf("callback jobName = %s, want %s", callbackJobName, jobName)
	}
	if callbackErr == nil {
		t.Error("callback error = nil, want error")
	}
}
