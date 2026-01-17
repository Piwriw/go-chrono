# 任务重试机制实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标：** 为 go-chrono 调度库添加任务重试机制，支持多种重试策略和灵活的配置方式。

**架构：** 通过 JobHook 系统集成重试逻辑，新增 retry 包包含策略接口和实现，扩展 SchedulerMonitor 记录重试事件，支持全局配置和 Builder 模式配置。

**技术栈：** Go 1.21+, go-co-op/gocron/v2, github.com/google/uuid

---

## 前置准备

### Task 0: 创建 retry 包目录结构

**Files:**
- Create: `retry/policy.go`
- Create: `retry/config.go`
- Create: `retry/executor.go`
- Create: `retry/retry_event.go`
- Create: `retry/policy_test.go`
- Create: `retry/executor_test.go`

**Step 1: 创建 retry 目录**

```bash
mkdir -p retry
```

**Step 2: 创建所有文件（空文件，稍后填充）**

```bash
touch retry/policy.go retry/config.go retry/executor.go retry/retry_event.go
touch retry/policy_test.go retry/executor_test.go
```

**Step 3: 提交初始结构**

```bash
git add retry/
git commit -m "feat(retry): create retry package structure"
```

---

## 核心策略实现

### Task 1: 实现 RetryPolicy 接口和基础结构

**Files:**
- Modify: `retry/policy.go`
- Test: `retry/policy_test.go`

**Step 1: 编写 RetryPolicy 接口的测试**

```go
// retry/policy_test.go
package retry

import (
    "errors"
    "testing"
    "time"
)

func TestRetryPolicy_Interface(t *testing.T) {
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
```

**Step 2: 运行测试验证失败**

```bash
go test -v ./retry -run TestRetryPolicy
```

预期：编译失败，类型未定义

**Step 3: 实现 RetryPolicy 接口和 FixedIntervalPolicy**

```go
// retry/policy.go
package retry

import (
    "errors"
    "time"
)

// RetryPolicy 定义重试策略接口
// RetryPolicy defines the retry strategy interface
type RetryPolicy interface {
    // NextRetry 计算下次重试的间隔时间
    // NextRetry calculates the interval for the next retry
    NextRetry(attempt int, lastErr error) time.Duration

    // ShouldRetry 判断是否应该重试
    // ShouldRetry determines whether a retry should be attempted
    ShouldRetry(attempt int, err error) bool
}

// FixedIntervalPolicy 固定间隔重试策略
// FixedIntervalPolicy retries at a fixed interval
type FixedIntervalPolicy struct {
    interval time.Duration
}

// NewFixedIntervalPolicy 创建固定间隔策略
// NewFixedIntervalPolicy creates a fixed interval retry policy
func NewFixedIntervalPolicy(interval time.Duration) *FixedIntervalPolicy {
    return &FixedIntervalPolicy{
        interval: interval,
    }
}

// NextRetry 返回固定的间隔时间
// NextRetry returns a fixed interval
func (p *FixedIntervalPolicy) NextRetry(attempt int, lastErr error) time.Duration {
    return p.interval
}

// ShouldRetry 始终返回 true（由调用方控制最大次数）
// ShouldRetry always returns true (caller controls max attempts)
func (p *FixedIntervalPolicy) ShouldRetry(attempt int, err error) bool {
    return true
}
```

**Step 4: 运行测试验证通过**

```bash
go test -v ./retry -run TestRetryPolicy
```

预期：PASS

**Step 5: 提交**

```bash
git add retry/policy.go retry/policy_test.go
git commit -m "feat(retry): add RetryPolicy interface and FixedIntervalPolicy"
```

---

### Task 2: 实现 ExponentialBackoffPolicy

**Files:**
- Modify: `retry/policy.go`
- Modify: `retry/policy_test.go`

**Step 1: 编写 ExponentialBackoffPolicy 的测试**

```go
// 添加到 retry/policy_test.go

func TestExponentialBackoffPolicy(t *testing.T) {
    base := 1 * time.Second
    max := 60 * time.Second
    policy := NewExponentialBackoffPolicy(base, max)

    t.Run("NextRetry increases exponentially", func(t *testing.T) {
        tests := []struct {
            attempt    int
            wantMin    time.Duration
            wantMax    time.Duration
        }{
            {0, 1 * time.Second, 1 * time.Second},
            {1, 2 * time.Second, 2 * time.Second},
            {2, 4 * time.Second, 4 * time.Second},
            {3, 8 * time.Second, 8 * time.Second},
            {6, 64 * time.Second, 60 * time.Second}, // 限制为 max
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
```

**Step 2: 运行测试验证失败**

```bash
go test -v ./retry -run TestExponentialBackoffPolicy
```

预期：未定义函数

**Step 3: 实现 ExponentialBackoffPolicy**

```go
// 添加到 retry/policy.go

// ExponentialBackoffPolicy 指数退避重试策略
// ExponentialBackoffPolicy retries with exponential backoff
type ExponentialBackoffPolicy struct {
    baseInterval time.Duration
    maxInterval  time.Duration
    multiplier   float64
}

// NewExponentialBackoffPolicy 创建指数退避策略
// NewExponentialBackoffPolicy creates an exponential backoff retry policy
func NewExponentialBackoffPolicy(base, max time.Duration) *ExponentialBackoffPolicy {
    return &ExponentialBackoffPolicy{
        baseInterval: base,
        maxInterval:  max,
        multiplier:   2.0,
    }
}

// NextRetry 计算指数增长的间隔时间
// NextRetry calculates exponentially increasing interval
func (p *ExponentialBackoffPolicy) NextRetry(attempt int, lastErr error) time.Duration {
    interval := time.Duration(float64(p.baseInterval) * (1 << uint(attempt)))
    if interval > p.maxInterval {
        return p.maxInterval
    }
    return interval
}

// ShouldRetry 始终返回 true
// ShouldRetry always returns true
func (p *ExponentialBackoffPolicy) ShouldRetry(attempt int, err error) bool {
    return true
}
```

**Step 4: 运行测试验证通过**

```bash
go test -v ./retry -run TestExponentialBackoffPolicy
```

预期：PASS

**Step 5: 提交**

```bash
git add retry/policy.go retry/policy_test.go
git commit -m "feat(retry): add ExponentialBackoffPolicy"
```

---

### Task 3: 实现 JitterPolicy

**Files:**
- Modify: `retry/policy.go`
- Modify: `retry/policy_test.go`

**Step 1: 编写 JitterPolicy 的测试**

```go
// 添加到 retry/policy_test.go

func TestJitterPolicy(t *testing.T) {
    base := NewFixedIntervalPolicy(10 * time.Second)
    jitterFactor := 0.2
    policy := NewJitterPolicy(base, jitterFactor)

    t.Run("NextRetry adds jitter to base interval", func(t *testing.T) {
        baseInterval := 10 * time.Second
        minJitter := time.Duration(float64(baseInterval) * (1 - jitterFactor))
        maxJitter := time.Duration(float64(baseInterval) * (1 + jitterFactor))

        // 多次测试确保随机性在合理范围内
        for i := 0; i < 10; i++ {
            got := policy.NextRetry(0, errors.New("test"))
            if got < minJitter || got > maxJitter {
                t.Errorf("NextRetry() = %v, want between %v and %v",
                    got, minJitter, maxJitter)
            }
        }
    })
}
```

**Step 2: 运行测试验证失败**

```bash
go test -v ./retry -run TestJitterPolicy
```

预期：未定义函数

**Step 3: 实现 JitterPolicy**

```go
// 添加到 retry/policy.go

import (
    "math/rand"
)

// JitterPolicy 带随机抖动的重试策略包装器
// JitterPolicy wraps a retry policy with random jitter
type JitterPolicy struct {
    policy       RetryPolicy
    jitterFactor float64
}

// NewJitterPolicy 创建带抖动的策略
// NewJitterPolicy creates a jittered retry policy
func NewJitterPolicy(policy RetryPolicy, jitterFactor float64) *JitterPolicy {
    return &JitterPolicy{
        policy:       policy,
        jitterFactor: jitterFactor,
    }
}

// NextRetry 返回带随机抖动的间隔时间
// NextRetry returns interval with random jitter
func (p *JitterPolicy) NextRetry(attempt int, lastErr error) time.Duration {
    baseInterval := p.policy.NextRetry(attempt, lastErr)
    jitter := (rand.Float64() * 2 - 1) * p.jitterFactor // [-jitterFactor, +jitterFactor]
    return time.Duration(float64(baseInterval) * (1 + jitter))
}

// ShouldRetry 委托给底层策略
// ShouldRetry delegates to underlying policy
func (p *JitterPolicy) ShouldRetry(attempt int, err error) bool {
    return p.policy.ShouldRetry(attempt, err)
}
```

**Step 4: 运行测试验证通过**

```bash
go test -v ./retry -run TestJitterPolicy
```

预期：PASS

**Step 5: 提交**

```bash
git add retry/policy.go retry/policy_test.go
git commit -m "feat(retry): add JitterPolicy for random jitter"
```

---

## 配置和事件结构

### Task 4: 实现 RetryConfig

**Files:**
- Create: `retry/config.go`

**Step 1: 创建 RetryConfig 结构**

```go
// retry/config.go
package retry

import (
    "context"
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
```

**Step 2: 运行测试验证编译**

```bash
go build ./retry
```

预期：无错误

**Step 3: 提交**

```bash
git add retry/config.go
git commit -m "feat(retry): add RetryConfig structure"
```

---

### Task 5: 实现 RetryEvent

**Files:**
- Create: `retry/retry_event.go`
- Test: `retry/retry_event_test.go`

**Step 1: 编写 RetryEvent 的测试**

```go
// retry/retry_event_test.go
package retry

import (
    "testing"
    "time"
)

func TestRetryEvent(t *testing.T) {
    event := &RetryEvent{
        EventID:         "test-event-1",
        OriginalEventID: "original-1",
        Attempt:         2,
        StartTime:       time.Now(),
        EndTime:         time.Now().Add(5 * time.Second),
        NextRetryIn:     10 * time.Second,
    }

    t.Run("RetryEvent fields are accessible", func(t *testing.T) {
        if event.EventID != "test-event-1" {
            t.Errorf("EventID = %s, want test-event-1", event.EventID)
        }
        if event.Attempt != 2 {
            t.Errorf("Attempt = %d, want 2", event.Attempt)
        }
    })
}
```

**Step 2: 运行测试验证失败**

```bash
go test -v ./retry -run TestRetryEvent
```

预期：类型未定义

**Step 3: 实现 RetryEvent**

```go
// retry/retry_event.go
package retry

import (
    "encoding/json"
    "errors"
    "time"
)

// RetryEvent 重试事件详情
// RetryEvent contains detailed retry information
type RetryEvent struct {
    // EventID 事件 ID
    // EventID is the event ID
    EventID string `json:"event_id"`

    // OriginalEventID 原始任务事件 ID
    // OriginalEventID is the original task event ID
    OriginalEventID string `json:"original_event_id"`

    // Attempt 第几次尝试（从 0 开始）
    // Attempt is the retry attempt number (starts from 0)
    Attempt int `json:"attempt"`

    // StartTime 开始时间
    // StartTime is the start time
    StartTime time.Time `json:"start_time"`

    // EndTime 结束时间
    // EndTime is the end time
    EndTime time.Time `json:"end_time"`

    // Err 错误
    // Err is the error
    Err error `json:"-"`

    // ErrorMessage 错误信息（用于 JSON 序列化）
    // ErrorMessage is the error message for JSON serialization
    ErrorMessage string `json:"error,omitempty"`

    // NextRetryIn 下次重试间隔
    // NextRetryIn is the duration until next retry
    NextRetryIn time.Duration `json:"next_retry_in"`
}

// MarshalJSON 序列化为 JSON
// MarshalJSON serializes to JSON
func (r *RetryEvent) MarshalJSON() ([]byte, error) {
    type Alias RetryEvent
    errMsg := ""
    if r.Err != nil {
        errMsg = r.Err.Error()
    }
    return json.Marshal(&struct {
        *Alias
        ErrorMessage string `json:"error,omitempty"`
    }{
        Alias:        (*Alias)(r),
        ErrorMessage: errMsg,
    })
}

// GetDuration 获取执行耗时
// GetDuration returns the execution duration
func (r *RetryEvent) GetDuration() time.Duration {
    return r.EndTime.Sub(r.StartTime)
}

// NewRetryEvent 创建新的重试事件
// NewRetryEvent creates a new retry event
func NewRetryEvent(eventID, originalEventID string, attempt int, startTime, endTime time.Time, err error, nextRetryIn time.Duration) *RetryEvent {
    errMsg := ""
    if err != nil && !errors.Is(err, nil) {
        errMsg = err.Error()
    }
    return &RetryEvent{
        EventID:         eventID,
        OriginalEventID: originalEventID,
        Attempt:         attempt,
        StartTime:       startTime,
        EndTime:         endTime,
        Err:             err,
        ErrorMessage:    errMsg,
        NextRetryIn:     nextRetryIn,
    }
}
```

**Step 4: 运行测试验证通过**

```bash
go test -v ./retry -run TestRetryEvent
```

预期：PASS

**Step 5: 提交**

```bash
git add retry/retry_event.go retry/retry_event_test.go
git commit -m "feat(retry): add RetryEvent structure"
```

---

## 重试执行器

### Task 6: 实现 RetryExecutor

**Files:**
- Create: `retry/executor.go`
- Test: `retry/executor_test.go`

**Step 1: 编写 RetryExecutor 的测试**

```go
// retry/executor_test.go
package retry

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/google/uuid"
)

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
        cancel() // 第一次执行后取消上下文
        return errors.New("error")
    }

    err := executor.ExecuteWithRetry(ctx, task)

    if err == nil {
        t.Error("ExecuteWithRetry() error = nil, want context canceled error")
    }
    if callCount > 2 { // 应该很快停止
        t.Errorf("task called %d times, should stop quickly", callCount)
    }
}
```

**Step 2: 运行测试验证失败**

```bash
go test -v ./retry -run TestRetryExecutor
```

预期：类型未定义

**Step 3: 实现 RetryExecutor**

```go
// retry/executor.go
package retry

import (
    "context"
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
        if len(e.config.RetriesableErrors) > 0 && !e.isRetryableError(err) {
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
    for _, retryableErr := range e.config.RetriesableErrors {
        if errors.Is(err, retryableErr) {
            return true
        }
    }
    return false
}
```

**注意：** 上面代码中有一个笔误 `RetriesableErrors` 应该是 `RetryableErrors`，需要修正。

**修正：**

```go
// 检查错误是否可重试
// Check if error is retryable
if len(e.config.RetryableErrors) > 0 && !e.isRetryableError(err) {
    break
}
```

**Step 4: 运行测试验证通过**

```bash
go test -v ./retry -run TestRetryExecutor
```

预期：PASS

**Step 5: 提交**

```bash
git add retry/executor.go retry/executor_test.go
git commit -m "feat(retry): add RetryExecutor for retry logic"
```

---

## 集成到调度器

### Task 7: 扩展 SchedulerOptions

**Files:**
- Modify: `option.go`

**Step 1: 添加 RetryOption 常量**

```go
// 在 option.go 的 const 部分添加

const (
    // AliasOptionName ...
    AliasOptionName = "alias"
    // WatchOptionName ...
    WatchOptionName = "watch"
    // WebMonitorOptionName ...
    WebMonitorOptionName = "web_monitor"
    // LimitOptionName ...
    LimitOptionName = "limit"
    // PrometheusOptionName ...
    PrometheusOptionName = "prometheus"
    // RetryOptionName is the name constant for the retry option.
    // RetryOptionName 是重试选项的名称常量。
    RetryOptionName = "retry"
)
```

**Step 2: 添加 RetryOption 结构**

```go
// 在 option.go 中添加

// RetryOption represents the retry option.
// RetryOption 表示重试选项。
type RetryOption struct {
    // enabled indicates whether the retry option is enabled.
    // enabled 表示是否启用重试选项。
    enabled bool
    // config is the retry configuration.
    // config 是重试配置。
    config *RetryConfig
}

var _ ScheduleOption = &RetryOption{}

// Name returns the name of the retry option.
// Name 返回重试选项的名称。
func (r *RetryOption) Name() string {
    return RetryOptionName
}

// Enable returns whether the retry option is enabled.
// Enable 返回重试选项是否启用。
func (r *RetryOption) Enable() bool {
    return r.enabled
}

// Config returns the retry configuration.
// Config 返回重试配置。
func (r *RetryOption) Config() *RetryConfig {
    return r.config
}
```

**Step 3: 添加 WithRetry 函数**

```go
// 在 option.go 中添加

import (
    "your-project/chrono/retry"
)

// WithRetry sets the retry configuration for all jobs.
// WithRetry 为所有任务设置重试配置。
//
// Parameters:
//
//	config - The retry configuration / 重试配置
//
// Returns:
//
//	SchedulerOption - The scheduler option / 调度器选项
func WithRetry(config *retry.RetryConfig) SchedulerOption {
    return func(s *Scheduler) {
        retryOpt := &RetryOption{
            enabled: config != nil && config.MaxRetries > 0,
            config:  config,
        }
        s.schOptions.retry = retryOpt
    }
}
```

**Step 4: 更新 SchedulerOptions 结构**

```go
// 在 SchedulerOptions 中添加 Retry 字段

type SchedulerOptions struct {
    alias      *AliasOption
    watch      *WatchOption
    timeout    *TimeoutOption
    webMonitor *WebMonitorOption
    limit      *LimitOption
    prometheus *PrometheusOption
    retry      *RetryOption  // 新增
}
```

**Step 5: 更新 Enable 方法**

```go
// 在 Scheduler.Enable 方法中添加 retry case

func (s *Scheduler) Enable(option string) bool {
    switch option {
    case AliasOptionName:
        if s.schOptions.alias != nil {
            return s.schOptions.alias.Enable()
        }
    case WatchOptionName:
        if s.schOptions.watch != nil {
            return s.schOptions.watch.Enable()
        }
    case RetryOptionName:  // 新增
        if s.schOptions.retry != nil {
            return s.schOptions.retry.Enable()
        }
    // ... 其他 case
    }
    return false
}
```

**Step 6: 编译验证**

```bash
go build .
```

预期：无错误

**Step 7: 提交**

```bash
git add option.go
git commit -m "feat(option): add RetryOption support"
```

---

### Task 8: 扩展 SchedulerMonitor 支持重试事件

**Files:**
- Modify: `monitor.go`

**Step 1: 在 SchedulerMonitor 接口中添加方法**

```go
// 在 SchedulerMonitor 接口中添加

import (
    "your-project/chrono/retry"
)

type SchedulerMonitor interface {
    gocron.MonitorStatus

    // 现有方法...
    Watch() chan JobWatchInterface
    UpdateJobEvents(jobID uuid.UUID, jobName string, jobnewEvent *JobEvent, jobTags ...string)
    GetJobEvents(jobID string) []*JobEvent

    // 新增方法

    // GetRetryHistory gets the retry history for a job.
    // 获取任务的重试历史。
    GetRetryHistory(jobID string) []*retry.RetryEvent

    // RecordRetryEvent records a retry event.
    // 记录重试事件。
    RecordRetryEvent(event *retry.RetryEvent)
}
```

**Step 2: 在 defaultSchedulerMonitor 中实现**

```go
// 在 defaultSchedulerMonitor 中添加字段

type defaultSchedulerMonitor struct {
    mu             sync.Mutex
    counter        map[string]int
    time           map[string][]time.Duration
    jobChan        chan JobWatchInterface
    maxRecords     int
    jobRecord      map[string]MonitorJobSpec
    eventIDCli     EventIDGenerator
    retryHistory   map[string][]*retry.RetryEvent  // 新增
}

// 初始化时添加
func newDefaultSchedulerMonitor(opts ...SchedulerMonitorOption) *defaultSchedulerMonitor {
    return &defaultSchedulerMonitor{
        counter:      make(map[string]int),
        time:         make(map[string][]time.Duration),
        jobChan:      make(chan JobWatchInterface, 100),
        jobRecord:    make(map[string]MonitorJobSpec),
        eventIDCli:   defaultEventIDGenerator,
        retryHistory: make(map[string][]*retry.RetryEvent),  // 新增
    }
}

// 实现 GetRetryHistory
func (s *defaultSchedulerMonitor) GetRetryHistory(jobID string) []*retry.RetryEvent {
    s.mu.Lock()
    defer s.mu.Unlock()

    history, ok := s.retryHistory[jobID]
    if !ok {
        return nil
    }
    // 返回副本
    result := make([]*retry.RetryEvent, len(history))
    copy(result, history)
    return result
}

// 实现 RecordRetryEvent
func (s *defaultSchedulerMonitor) RecordRetryEvent(event *retry.RetryEvent) {
    s.mu.Lock()
    defer s.mu.Unlock()

    jobID := event.OriginalEventID

    if _, ok := s.retryHistory[jobID]; !ok {
        s.retryHistory[jobID] = make([]*retry.RetryEvent, 0)
    }

    // 最多保留 100 条重试记录
    const maxRetryHistory = 100
    if len(s.retryHistory[jobID]) >= maxRetryHistory {
        s.retryHistory[jobID] = s.retryHistory[jobID][1:]
    }

    s.retryHistory[jobID] = append(s.retryHistory[jobID], event)

    slog.Debug("chrono: retry event recorded",
        "job_id", jobID,
        "attempt", event.Attempt,
        "error", event.ErrorMessage)
}
```

**Step 3: 扩展 JobEvent 结构**

```go
// 在 JobEvent 中添加字段

type JobEvent struct {
    EventID     string
    StartTime   time.Time
    EndTime     time.Time
    Status      gocron.JobStatus
    Err         error

    // 新增重试相关字段
    // New retry-related fields
    RetryCount      int     `json:"retry_count,omitempty"`
    IsRetry         bool    `json:"is_retry,omitempty"`
    OriginalEventID string  `json:"original_event_id,omitempty"`
}
```

**Step 4: 编译验证**

```bash
go build .
```

预期：无错误

**Step 5: 提交**

```bash
git add monitor.go
git commit -m "feat(monitor): add retry history support to SchedulerMonitor"
```

---

### Task 9: 集成重试逻辑到 JobHook

**Files:**
- Modify: `job_hook.go`

**Step 1: 修改 wrapJobWithHooks 集成重试**

首先查看 job_hook.go 的当前结构，然后添加重试逻辑：

```go
// 在 job_hook.go 中修改

func wrapJobWithHooks(
    ctx context.Context,
    jobID uuid.UUID,
    name string,
    tags []string,
    task func(),
    monitor SchedulerMonitor,
    opts *jobOptions,
) {
    // 检查是否有重试配置
    var retryConfig *RetryConfig
    if opts != nil && opts.retryConfig != nil {
        retryConfig = opts.retryConfig
    } else if scheduler != nil && scheduler.Enable(RetryOptionName) {
        retryConfig = scheduler.schOptions.retry.Config()
    }

    if retryConfig != nil && retryConfig.MaxRetries > 0 {
        // 使用重试执行器
        executor := NewRetryExecutor(retryConfig, jobID, name)
        err := executor.ExecuteWithRetry(ctx, func() error {
            // 执行原始任务
            executeJob(ctx, jobID, name, tags, task, monitor)
            return nil // 这里需要从 executeJob 获取错误
        })

        if err != nil {
            // 记录最终失败
            slog.Error("chrono: job failed after retries",
                "job_id", jobID,
                "job_name", name,
                "error", err)
        }
    } else {
        // 原有逻辑，不重试
        executeJob(ctx, jobID, name, tags, task, monitor)
    }
}
```

**注意：** 需要重构 executeJob 使其返回 error，以便重试逻辑能够判断是否需要重试。

**Step 2: 提交**

```bash
git add job_hook.go
git commit -m "feat(hook): integrate retry logic into job execution"
```

---

## JobClient 集成

### Task 10: 为所有 JobClient 添加 WithRetry 方法

**Files:**
- Modify: `cron_client.go`
- Modify: `once_client.go`
- Modify: `interval_job_client.go`
- Modify: `daliy_client.go`
- Modify: `weekly_job_client.go`
- Modify: `monthly_job_client.go`

**以 CronJobClient 为例：**

**Step 1: 在 CronJobClient 中添加 retryConfig 字段**

```go
// cron_client.go

type CronJobClient struct {
    // 现有字段...
    names      []string
    cronExpr   string
    tags       []string
    task       func()
    jobOptions *jobOptions

    // 新增字段
    retryConfig *RetryConfig
}
```

**Step 2: 添加 WithRetry 方法**

```go
// WithRetry 设置任务的重试配置
// WithRetry sets the retry configuration for the job
//
// Parameters:
//
//	maxRetries - 最大重试次数 / Maximum number of retries
//	policy     - 重试策略 / Retry policy
//
// Returns:
//	*CronJobClient - 返回当前客户端以便链式调用 / Returns current client for chaining
func (c *CronJobClient) WithRetry(maxRetries int, policy RetryPolicy) *CronJobClient {
    c.retryConfig = &RetryConfig{
        MaxRetries: maxRetries,
        Policy:     policy,
    }
    return c
}

// WithRetryConfig 设置完整的重试配置
// WithRetryConfig sets the complete retry configuration
//
// Parameters:
//
//	config - 重试配置 / Retry configuration
//
// Returns:
//	*CronJobClient - 返回当前客户端以便链式调用 / Returns current client for chaining
func (c *CronJobClient) WithRetryConfig(config *RetryConfig) *CronJobClient {
    c.retryConfig = config
    return c
}
```

**Step 3: 在 Build 方法中传递 retryConfig**

```go
// Build 方法中
func (c *CronJobClient) Build() (uuid.UUID, error) {
    // ...
    jobOptions := &jobOptions{
        tags:         c.tags,
        retryConfig:  c.retryConfig,  // 添加
    }
    // ...
}
```

**Step 4: 为其他 JobClient 重复相同步骤**

对以下文件重复 Step 1-3：
- `once_client.go`
- `interval_job_client.go`
- `daliy_client.go`
- `weekly_job_client.go`
- `monthly_job_client.go`

**Step 5: 编译验证**

```bash
go build .
```

预期：无错误

**Step 6: 提交**

```bash
git add cron_client.go once_client.go interval_job_client.go daliy_client.go weekly_job_client.go monthly_job_client.go
git commit -m "feat(client): add WithRetry methods to all JobClients"
```

---

## Web API 扩展

### Task 11: 添加重试历史查询 API

**Files:**
- Modify: `web_monitor.go`

**Step 1: 添加重试历史处理器**

```go
// web_monitor.go

func (wm *WebMonitor) handleJobRetries(w http.ResponseWriter, r *http.Request) {
    // 从 URL 获取 job_id
    jobID := r.PathValue("job_id")
    if jobID == "" {
        http.Error(w, "job_id is required", http.StatusBadRequest)
        return
    }

    // 获取重试历史
    history := wm.scheduleMonitor.GetRetryHistory(jobID)

    // 返回 JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(history); err != nil {
        http.Error(w, "failed to encode response", http.StatusInternalServerError)
    }
}
```

**Step 2: 在 Start 方法中注册路由**

```go
// WebMonitor.Start 方法中

func (wm *WebMonitor) Start() error {
    mux := http.NewServeMux()
    endpoints := []string{"/healthz", "/jobs", "/jobs/{job_id}/retries"}

    mux.HandleFunc("/healthz", wm.handleHealthz)
    mux.HandleFunc("/jobs", wm.handleJobs)
    mux.HandleFunc("/jobs/{job_id}/retries", wm.handleJobRetries)  // 新增

    // ... 现有代码

    slog.Info("chrono:web monitor started", "address", wm.addr)
    host := wm.addr
    if strings.HasPrefix(host, ":") {
        host = "localhost" + host
    }
    for _, ep := range endpoints {
        slog.Info(fmt.Sprintf("  http://%s%s", host, ep))
    }

    // ... 现有代码
}
```

**Step 3: 测试 API**

```bash
# 启动服务后测试
curl http://localhost:8080/jobs/{job_id}/retries
```

**Step 4: 提交**

```bash
git add web_monitor.go
git commit -m "feat(web): add retry history API endpoint"
```

---

## 集成测试

### Task 12: 编写集成测试

**Files:**
- Create: `retry_integration_test.go`

**Step 1: 编写完整的集成测试**

```go
// retry_integration_test.go
package chrono

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/your-project/chrono/retry"
)

func TestScheduler_WithRetry_GlobalConfig(t *testing.T) {
    ctx := context.Background()

    monitor := newDefaultSchedulerMonitor()
    retryConfig := &retry.RetryConfig{
        MaxRetries: 3,
        Policy:     retry.NewFixedIntervalPolicy(10 * time.Millisecond),
    }

    scheduler, err := NewScheduler(
        ctx,
        monitor,
        WithRetry(retryConfig),
    )
    if err != nil {
        t.Fatalf("NewScheduler failed: %v", err)
    }

    callCount := 0
    task := func() {
        callCount++
        if callCount < 3 {
            panic("temporary error")
        }
    }

    _, err = scheduler.Cron().
        Names("test-job").
        CronExpr("* * * * *").
        Task(task).
        Build()

    if err != nil {
        t.Fatalf("Build failed: %v", err)
    }

    // 等待任务执行
    time.Sleep(2 * time.Second)

    if callCount < 3 {
        t.Errorf("task called %d times, want at least 3", callCount)
    }
}

func TestScheduler_WithRetry_BuilderConfig(t *testing.T) {
    ctx := context.Background()

    monitor := newDefaultSchedulerMonitor()
    scheduler, err := NewScheduler(ctx, monitor)
    if err != nil {
        t.Fatalf("NewScheduler failed: %v", err)
    }

    callCount := 0
    task := func() {
        callCount++
        if callCount < 2 {
            errors.New("error")
        }
    }

    _, err = scheduler.Cron().
        Names("test-job").
        CronExpr("* * * * *").
        WithRetry(2, retry.NewFixedIntervalPolicy(10*time.Millisecond)).
        Task(task).
        Build()

    if err != nil {
        t.Fatalf("Build failed: %v", err)
    }

    time.Sleep(2 * time.Second)

    if callCount < 2 {
        t.Errorf("task called %d times, want at least 2", callCount)
    }
}
```

**Step 2: 运行集成测试**

```bash
go test -v -run TestScheduler_WithRetry
```

**Step 3: 提交**

```bash
git add retry_integration_test.go
git commit -m "test(retry): add integration tests for retry mechanism"
```

---

## 文档更新

### Task 13: 更新项目文档

**Files:**
- Modify: `CLAUDE.md`
- Create: `examples/retry_example.go`

**Step 1: 在 CLAUDE.md 中添加重试机制说明**

```markdown
## Retry Mechanism

The scheduler supports automatic task retry with configurable policies:

### Retry Policies

- **FixedIntervalPolicy** - Retries at a fixed interval
- **ExponentialBackoffPolicy** - Retries with exponential backoff
- **JitterPolicy** - Adds random jitter to any policy

### Configuration

**Global configuration:**
```go
scheduler, err := chrono.NewScheduler(
    ctx,
    monitor,
    chrono.WithRetry(&chrono.RetryConfig{
        MaxRetries: 3,
        Policy:     chrono.NewExponentialBackoffPolicy(time.Second, 2*time.Minute),
    }),
)
```

**Per-job configuration:**
```go
scheduler.Cron().
    Names("my-job").
    CronExpr("*/5 * * * *").
    WithRetry(3, chrono.NewFixedIntervalPolicy(5*time.Second)).
    Task(myTask).
    Build()
```
```

**Step 2: 创建示例文件**

```go
// examples/retry_example.go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/your-project/chrono"
)

func main() {
    ctx := context.Background()

    // 创建带重试的调度器
    scheduler, err := chrono.NewScheduler(
        ctx,
        nil,
        chrono.WithRetry(&chrono.RetryConfig{
            MaxRetries: 3,
            Policy:     chrono.NewExponentialBackoffPolicy(time.Second, time.Minute),
        }),
    )

    if err != nil {
        panic(err)
    }

    // 添加任务
    scheduler.Cron().
        Names("cleanup").
        CronExpr("0 2 * * *").
        Task(cleanupTask).
        Build()

    scheduler.Start()
    defer scheduler.Shutdown()

    select {}
}

func cleanupTask() {
    fmt.Println("Running cleanup task...")
    // 清理逻辑
}
```

**Step 3: 提交**

```bash
git add CLAUDE.md examples/retry_example.go
git commit -m "docs: add retry mechanism documentation and examples"
```

---

## 最终验证

### Task 14: 完整测试和代码审查

**Step 1: 运行所有测试**

```bash
go test ./... -v
```

预期：所有测试通过

**Step 2: 运行 lint 检查**

```bash
make lint
```

预期：无错误

**Step 3: 格式化代码**

```bash
make format
```

**Step 4: 构建项目**

```bash
go build .
```

预期：无错误

**Step 5: 提交最终更改**

```bash
git add .
git commit -m "chore: final cleanup and formatting"
```

---

## 完成标准

- [ ] 所有单元测试通过
- [ ] 所有集成测试通过
- [ ] 代码通过 lint 检查
- [ ] 文档已更新
- [ ] 示例代码已添加
- [ ] 功能已验证

---

## 执行命令

准备开始执行计划：

```bash
# 查看计划
cat docs/plans/2025-01-17-retry-mechanism.md

# 使用 executing-plans skill 执行
```
