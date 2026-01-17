# 任务重试机制设计文档

## 概述

为 go-chrono 调度库添加任务重试机制，当任务执行失败时自动重试，支持多种重试策略和灵活的配置方式。

**设计日期：** 2025-01-17

---

## 核心架构

### 1. RetryPolicy（重试策略接口）

```go
// RetryPolicy 定义重试策略接口
// RetryPolicy defines the retry strategy interface
type RetryPolicy interface {
    // NextRetry 计算下次重试的间隔时间
    // NextRetry calculates the interval for the next retry
    //
    // Parameters:
    //   attempt - 当前重试次数（从 0 开始） / Current retry attempt (starts from 0)
    //   lastErr - 上次执行的错误 / Last execution error
    //
    // Returns:
    //   time.Duration - 下次重试的间隔时间 / Interval until next retry
    NextRetry(attempt int, lastErr error) time.Duration

    // ShouldRetry 判断是否应该重试
    // ShouldRetry determines whether a retry should be attempted
    //
    // Parameters:
    //   attempt - 当前重试次数 / Current retry attempt
    //   err     - 执行错误 / Execution error
    //
    // Returns:
    //   bool - true 表示应该重试 / true means should retry
    ShouldRetry(attempt int, err error) bool
}
```

### 2. 内置策略实现

#### FixedIntervalPolicy（固定间隔）

```go
type FixedIntervalPolicy struct {
    interval time.Duration
}

func NewFixedIntervalPolicy(interval time.Duration) *FixedIntervalPolicy
```

#### ExponentialBackoffPolicy（指数退避）

```go
type ExponentialBackoffPolicy struct {
    baseInterval   time.Duration
    maxInterval    time.Duration
    multiplier     float64
}

func NewExponentialBackoffPolicy(base, max time.Duration) *ExponentialBackoffPolicy
```

#### JitterPolicy（随机抖动）

```go
type JitterPolicy struct {
    policy    RetryPolicy
    jitterFactor float64
}

func NewJitterPolicy(policy RetryPolicy, jitterFactor float64) *JitterPolicy
```

### 3. RetryConfig（重试配置）

```go
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
    // RetryableErrors are the error types that can be retried (empty means all errors retry)
    RetryableErrors []error

    // OnFinalFailure 最终失败回调
    // OnFinalFailure is the callback when all retries fail
    OnFinalFailure func(ctx context.Context, jobID uuid.UUID, jobName string, err error)
}
```

### 4. 集成点

重试逻辑在 `job_hook.go` 的 `wrapJobWithHooks` 中实现，当任务返回 error 时，根据配置决定是否重试。

---

## 配置方式

### 1. 全局配置（Option 模式）

```go
// WithRetry 设置全局重试配置
// WithRetry sets the global retry configuration
func WithRetry(config *RetryConfig) SchedulerOption

// 使用示例
scheduler, err := chrono.NewScheduler(
    ctx,
    monitor,
    chrono.WithRetry(&chrono.RetryConfig{
        MaxRetries: 3,
        Policy:     chrono.NewExponentialBackoffPolicy(time.Second, 2*time.Minute),
    }),
)
```

### 2. Builder 模式（JobClient 级别）

所有 JobClient 添加 `WithRetry()` 方法：

```go
// CronJobClient
func (c *CronJobClient) WithRetry(maxRetries int, policy RetryPolicy) *CronJobClient

// OnceJobClient
func (c *OnceJobClient) WithRetry(maxRetries int, policy RetryPolicy) *OnceJobClient

// IntervalJobClient
func (c *IntervalJobClient) WithRetry(maxRetries int, policy RetryPolicy) *IntervalJobClient

// DailyJobClient
func (c *DailyJobClient) WithRetry(maxRetries int, policy RetryPolicy) *DailyJobClient

// WeeklyJobClient
func (c *WeeklyJobClient) WithRetry(maxRetries int, policy RetryPolicy) *WeeklyJobClient

// MonthlyJobClient
func (c *MonthlyJobClient) WithRetry(maxRetries int, policy RetryPolicy) *MonthlyJobClient
```

### 3. 优先级规则

| 配置级别 | 优先级 |
|---------|--------|
| JobClient Builder 模式 | 高 |
| 全局 Option | 低 |
| 未配置 | 不启用重试 |

---

## 事件和监控

### 1. 扩展 JobEvent

```go
type JobEvent struct {
    EventID     string
    StartTime   time.Time
    EndTime     time.Time
    Status      gocron.JobStatus
    Err         error

    // 新增重试相关字段 / New retry-related fields
    RetryCount      int           // 本次执行重试次数 / Retry count for this execution
    IsRetry         bool          // 是否是重试执行 / Whether this is a retry execution
    OriginalEventID string        // 原始事件 ID / Original event ID
}
```

### 2. 新增 RetryEvent

```go
// RetryEvent 重试事件详情
// RetryEvent contains detailed retry information
type RetryEvent struct {
    // EventID 事件 ID
    // EventID is the event ID
    EventID string

    // OriginalEventID 原始任务事件 ID
    // OriginalEventID is the original task event ID
    OriginalEventID string

    // Attempt 第几次尝试（从 0 开始）
    // Attempt is the retry attempt number (starts from 0)
    Attempt int

    // StartTime 开始时间
    // StartTime is the start time
    StartTime time.Time

    // EndTime 结束时间
    // EndTime is the end time
    EndTime time.Time

    // Err 错误
    // Err is the error
    Err error

    // NextRetryIn 下次重试间隔
    // NextRetryIn is the duration until next retry
    NextRetryIn time.Duration
}
```

### 3. SchedulerMonitor 扩展

```go
type SchedulerMonitor interface {
    // 现有方法...
    // Existing methods...

    // GetRetryHistory 获取任务的重试历史
    // GetRetryHistory gets the retry history for a job
    //
    // Parameters:
    //   jobID - 任务 ID / Job ID
    //
    // Returns:
    //   []*RetryEvent - 重试事件列表 / List of retry events
    GetRetryHistory(jobID string) []*RetryEvent

    // RecordRetryEvent 记录重试事件
    // RecordRetryEvent records a retry event
    //
    // Parameters:
    //   event - 重试事件 / Retry event
    RecordRetryEvent(event *RetryEvent)
}
```

### 4. Web API 扩展

```
GET /jobs/{job_id}/retries - 返回指定任务的重试历史
```

---

## 实现细节

### 重试执行流程

```
┌─────────────────┐
│   任务执行       │
│   Execute Job   │
└────────┬────────┘
         │
         ▼
    ┌────────┐
    │ 返回error? │
    │  Error?  │
    └───┬────┬─┘
        │    │
   No   │    │ Yes
        │    │
        │    ▼
        │ ┌─────────────────┐
        │ │ 启用重试?        │
        │ │ Retry Enabled?  │
        │ └───┬───────┬─────┘
        │     │       │
        │    No       │ Yes
        │     │       │
        │     ▼       ▼
        │  ┌────────────────┐
        │  │ 重试次数已用尽?   │
        │  │ Max Retries?   │
        │  └───┬───────┬────┘
        │      │       │
        │     Yes      │ No
        │      │       │
        │      ▼       ▼
        │  ┌─────────┐ ┌────────────────┐
        │  │ 记录失败 │ │ 错误可重试?     │
        │  │ 回调    │ │ Retryable?     │
        │  └────┬────┘ └───┬────────┬───┘
        │       │          │        │
        │       │        No        │ Yes
        │       │          │        │
        │       │          ▼        ▼
        │       │    ┌─────────┐ ┌────────────┐
        │       │    │ 记录失败 │ │计算下次重试 │
        │       │    │ 回调    │ │ 等待      │
        │       │    └─────────┘ └─────┬──────┘
        │       │                     │
        │       └─────────────────────┘
        │                         │
        └─────────────────────────┼──────────┐
                                  │          │
                                  ▼          │
                            ┌──────────┐    │
                            │ 再次执行  │◄───┘
                            │ Retry    │
                            └──────────┘
```

### 并发安全

- 重试状态使用 `sync.Mutex` 保护
- 每个任务独立的重试上下文
- 任务取消时重试也会被取消

### 上下文传递

```go
func (r *retryExecutor) executeWithRetry(ctx context.Context, job func() error) error {
    for attempt := 0; attempt <= r.maxRetries; attempt++ {
        err := job()
        if err == nil {
            return nil
        }
        if !r.shouldRetry(attempt, err) {
            return r.finalFailure(ctx, err)
        }
        select {
        case <-time.After(r.policy.NextRetry(attempt, err)):
            // 继续重试
            // Continue retry
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

---

## 文件结构

```
chrono/
├── retry/
│   ├── policy.go           # 重试策略接口和实现
│   ├── config.go           # 重试配置
│   ├── executor.go         # 重试执行器
│   └── retry_event.go      # 重试事件定义
├── job_hook.go             # 修改：集成重试逻辑
├── monitor.go              # 修改：扩展 SchedulerMonitor
├── option.go               # 修改：添加 WithRetry 选项
├── web_monitor.go          # 修改：添加 /retries API
├── cron_client.go          # 修改：添加 WithRetry 方法
├── once_client.go          # 修改：添加 WithRetry 方法
├── interval_job_client.go  # 修改：添加 WithRetry 方法
├── daliy_client.go         # 修改：添加 WithRetry 方法
├── weekly_job_client.go    # 修改：添加 WithRetry 方法
├── monthly_job_client.go   # 修改：添加 WithRetry 方法
└── retry_test.go           # 测试文件
```

---

## 测试计划

### 单元测试

- 各策略的 `NextRetry()` 计算逻辑
- `ShouldRetry()` 判断逻辑
- Jitter 随机性验证

### 集成测试

- 完整重试流程（成功/失败）
- 全局配置优先级
- Builder 覆盖全局配置

### 边界测试

- 重试 0 次
- 重试 1 次
- 重试最大次数
- 任务取消时的重试行为

### 并发测试

- 多个任务同时重试
- 重试过程中的状态一致性

---

## API 使用示例

### 示例 1：固定间隔重试

```go
scheduler.Cron().
    Names("cleanup").
    CronExpr("0 2 * * *").
    WithRetry(3, chrono.NewFixedIntervalPolicy(5*time.Second)).
    Task(cleanupFunc).
    Build()
```

### 示例 2：指数退避重试

```go
scheduler.Cron().
    Names("api-sync").
    CronExpr("*/10 * * * *").
    WithRetry(5, chrono.NewExponentialBackoffPolicy(time.Second, time.Minute)).
    Task(syncAPI).
    Build()
```

### 示例 3：带随机抖动的重试

```go
basePolicy := chrono.NewExponentialBackoffPolicy(time.Second, 30*time.Second)
jitterPolicy := chrono.NewJitterPolicy(basePolicy, 0.2)

scheduler.Cron().
    Names("db-backup").
    CronExpr("0 3 * * *").
    WithRetry(3, jitterPolicy).
    Task(backupDB).
    Build()
```

### 示例 4：带最终失败回调

```go
retryConfig := &chrono.RetryConfig{
    MaxRetries: 3,
    Policy:     chrono.NewFixedIntervalPolicy(5*time.Second),
    OnFinalFailure: func(ctx context.Context, jobID uuid.UUID, jobName string, err error) {
        log.Error("Job failed after all retries", "job", jobName, "error", err)
        // 发送告警通知
        // Send alert notification
    },
}

scheduler.Cron().
    Names("critical-job").
    CronExpr("*/5 * * * *").
    WithRetryConfig(retryConfig).
    Task(criticalTask).
    Build()
```

---

## 依赖关系

- 无新增外部依赖
- 复用现有的 `uuid`、`time`、`context` 等标准库
- 集成现有的 JobHook 和 Monitor 系统
