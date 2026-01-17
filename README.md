#  go-chrono
## 支持功能
### 支持的任务类型
- Cron表达式（Cron Job）：任务可以使用 crontab 格式的时间表达式运行
- IntervalJob：任务可以每隔 X 秒、分钟、小时、天、周、月运行
- 每日（DailyJob）：任务可以每隔 X 天在特定的时间运行
- 每周（WeeklyJob）：任务可以每隔 X 周在特定的星期几和时间运行
- 每月（MonthlyJob）：任务可以每隔 X 月在特定的日期和时间运行
- 一次性（OnceJob）：任务可以在特定的时间运行（可以是一次或多次）

### 支持功能
- [x] 允许设置调度器全局配置，优先使用Job自身配置
  - [x] 超时时间
  - [x] Watch监听模式
  - [x] 钩子函数
  - [x] 重试策略
- [x] 允许自定义实现JobClient
  - [ ] Example
- [ ] 统一Time.Format格式
- [x] 设置Schedule最大管控调用任务数量
- [ ] 自定义Logger
- [x] Job标签级别管控
  - 缺少移除func
- [x] 优雅终止关闭
- [x] 监控模式
  - [x]  Prometheus端点监控
  - [x]  Listen Web监控
  - [x]  可以查询最近X次的执行结果
     - [ ]  支持查询指定时间范围的执行结果
     - [ ]  支持查询指定标签的执行结果
     - [ ]  支持查询指定任务的执行结果
     - [ ]  支持查询指定任务的执行结果
     - [x]  新增EventID支持
        - 允许自定义实现Event ID生成器
        - 默认使用JobID+JobName+时间戳
        - 内置支持生成器：
          1. UUID
          2. JobID+JobName+时间戳

## 重试机制 (Retry Mechanism)

go-chrono 提供了灵活的任务重试机制，支持多种重试策略。

### 内置重试策略

1. **固定间隔策略 (FixedIntervalPolicy)**: 每次重试使用固定的时间间隔
2. **指数退避策略 (ExponentialBackoffPolicy)**: 每次重试间隔指数增长（2^n * base）
3. **抖动策略 (JitterPolicy)**: 在基础策略上添加随机抖动，避免雷击效应

### 使用示例

```go
package main

import (
    "github.com/piwriw/go-chrono"
    "github.com/piwriw/go-chrono/retry"
    "time"
)

func main() {
    scheduler := chrono.NewScheduler()

    // 使用固定间隔重试策略
    chrono.NewCronJob(scheduler).
        CronExpr("* * * * *").
        Name("fixed-retry-job").
        Task(myTask).
        WithRetry(3, retry.NewFixedIntervalPolicy(5*time.Second)).
        Add()

    // 使用指数退避重试策略
    chrono.NewCronJob(scheduler).
        CronExpr("* * * * *").
        Name("exponential-retry-job").
        Task(myTask).
        WithRetry(5, retry.NewExponentialBackoffPolicy(1*time.Second, 60*time.Second)).
        Add()

    // 使用抖动策略
    basePolicy := retry.NewFixedIntervalPolicy(10 * time.Second)
    jitterPolicy := retry.NewJitterPolicy(basePolicy, 0.2) // 20% 抖动
    chrono.NewCronJob(scheduler).
        CronExpr("* * * * *").
        Name("jitter-retry-job").
        Task(myTask).
        WithRetryConfig(&retry.RetryConfig{
            MaxRetries: 3,
            Policy:     jitterPolicy,
        }).
        Add()

    scheduler.Start()
}

func myTask() error {
    // 你的任务逻辑
    return nil
}
```

### 查询重试历史

通过 Web 监控端点查询任务的重试历史：

```bash
# 获取指定任务的重试历史
curl http://localhost:8080/jobs/{job_id}/retries
```

### 重试配置选项

- **MaxRetries**: 最大重试次数
- **Policy**: 重试策略（必须实现 RetryPolicy 接口）
- **IsRetryable** (可选): 自定义错误重试判断函数
- **OnRetry** (可选): 每次重试前的回调函数
- **OnFinalFailure** (可选): 最终失败后的回调函数

## Example
