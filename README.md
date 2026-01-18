# go-chrono

> 一个功能强大的 Go 任务调度库，基于 [go-co-op/gocron/v2](https://github.com/go-co-op/gocron) 构建，提供流畅的 API、监控、重试机制和灵活的任务生命周期管理。

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-2.0.0-orange.svg)](https://github.com/piwriw/go-chrono/tree/v2.0.0)

> **⚠️ 注意：** 当前文档为 v2.0.0 版本。如果您使用 v1.0.0，请切换到 [main 分支](https://github.com/piwriw/go-chrono/tree/main)。

## 版本与分支

| 分支 | 版本 | 说明 |
|------|------|------|
| **[v2.0.0](https://github.com/piwriw/go-chrono/tree/v2.0.0)** | `v2.0.0` | 🔄 当前开发分支 - 包含包结构重构和重试机制 |
| **[main](https://github.com/piwriw/go-chrono/tree/main)** | `v1.0.0` | ✅ 稳定版本 - 生产环境推荐 |

**版本升级说明 (v1.0.0 → v2.0.0):**

v2.0.0 是一个重大版本升级，包含不兼容的包结构重构。主要变更：
- 📦 包结构重组：`pkg.JobType` → `jobs.JobType`
- 🔄 新增重试机制和重试历史查询
- 🏗️ 修复循环导入问题，建立清晰的依赖层次

**查看 [CHANGELOG.md](CHANGELOG.md) 了解详细的变更内容和迁移指南。**

---

## 特性

- 🕐 **多种任务类型** - 支持 Cron、定时、一次性、间隔等多种调度方式
- 🔄 **重试机制** - 内置固定间隔、指数退避、抖动等重试策略
- 📊 **监控面板** - 内置 Web 监控和 Prometheus 指标支持
- 🎣 **钩子函数** - 支持任务生命周期的各个阶段回调
- 🏷️ **标签管理** - 任务标签级别的管控和查询
- 🔒 **分布式锁** - 支持分布式环境下的任务去重
- 🎨 **流式 API** - Builder 模式提供的优雅链式调用

## 安装

```bash
go get github.com/piwriw/go-chrono
```

## 快速开始

```go
package main

import (
    "github.com/piwriw/go-chrono"
)

func main() {
    // 创建调度器
    scheduler := chrono.NewScheduler()

    // 创建一个每分钟执行的定时任务
    chrono.NewCronJob(scheduler).
        CronExpr("* * * * *").
        Name("my-jobs").
        Task(func() {
            println("Hello, go-chrono!")
        }).
        Add()

    // 启动调度器
    scheduler.Start()

    // 程序退出时关闭
    defer scheduler.Stop()
}
```

## 任务类型

### Cron 任务

使用 Cron 表达式定义执行时间：

```go
chrono.NewCronJob(scheduler).
    CronExpr("0 */5 * * *").  // 每 5 分钟
    Name("cron-jobs").
    Task(myTask).
    Add()
```

### 定时任务

按固定间隔执行任务：

```go
chrono.NewIntervalJob(scheduler).
    Interval(10 * time.Second).
    Name("interval-jobs").
    Task(myTask).
    Add()
```

### 每日任务

在每天的特定时间执行：

```go
chrono.NewDailyJob(scheduler).
    AtTime(14, 30, 0).  // 每天 14:30:00
    Name("daily-jobs").
    Task(myTask).
    Add()
```

### 每周任务

在每周的特定日期和时间执行：

```go
chrono.NewWeeklyJob(scheduler).
    AtTime([]time.Weekday{time.Monday, time.Friday}, 10, 0, 0).
    Name("weekly-jobs").
    Task(myTask).
    Add()
```

### 每月任务

在每月的特定日期和时间执行：

```go
chrono.NewMonthlyJob(scheduler).
    AtTime([]int{1, 15}, 9, 0, 0).  // 每月 1 号和 15 号 9:00
    Name("monthly-jobs").
    Task(myTask).
    Add()
```

### 一次性任务

在指定时间执行一次：

```go
chrono.NewOnceJob(scheduler).
    At(time.Now().Add(1 * time.Hour)).
    Name("once-jobs").
    Task(myTask).
    Add()
```

## 重试机制

go-chrono 提供灵活的任务重试机制，支持多种重试策略。

### 内置策略

| 策略 | 说明 |
|------|------|
| `FixedIntervalPolicy` | 固定时间间隔重试 |
| `ExponentialBackoffPolicy` | 指数退避（2^n × base，带上限） |
| `JitterPolicy` | 在基础策略上添加随机抖动 |

### 使用示例

```go
import "github.com/piwriw/go-chrono/retry"

// 固定间隔重试
chrono.NewCronJob(scheduler).
    CronExpr("* * * * *").
    Name("fixed-retry-jobs").
    Task(myTask).
    WithRetry(3, retry.NewFixedIntervalPolicy(5*time.Second)).
    Add()

// 指数退避重试
chrono.NewCronJob(scheduler).
    CronExpr("* * * * *").
    Name("exponential-retry-jobs").
    Task(myTask).
    WithRetry(5, retry.NewExponentialBackoffPolicy(1*time.Second, 60*time.Second)).
    Add()

// 带抖动的重试（避免雷击效应）
basePolicy := retry.NewFixedIntervalPolicy(10 * time.Second)
jitterPolicy := retry.NewJitterPolicy(basePolicy, 0.2)  // 20% 抖动
chrono.NewCronJob(scheduler).
    CronExpr("* * * * *").
    Name("jitter-retry-jobs").
    Task(myTask).
    WithRetryConfig(&retry.RetryConfig{
        MaxRetries: 3,
        Policy:     jitterPolicy,
    }).
    Add()
```

### 查询重试历史

```bash
curl http://localhost:8080/jobs/{job_id}/retries
```

## 调度器选项

### 启用别名模式

```go
scheduler := chrono.NewScheduler(
    chrono.WithAliasMode(),
)
```

### 启用 Web 监控

```go
scheduler := chrono.NewScheduler(
    chrono.WithWebMonitor(":8080"),
)
// 访问 http://localhost:8080 查看监控面板
```

### 启用 Prometheus 指标

```go
scheduler := chrono.NewScheduler(
    chrono.WithPrometheus(":9090"),
)
```

### 设置并发限制

```go
scheduler := chrono.NewScheduler(
    chrono.WithLimit(100),  // 最多 100 个并发任务
)
```

## 钩子函数

```go
chrono.NewCronJob(scheduler).
    CronExpr("* * * * *").
    Name("jobs-with-hooks").
    Task(myTask).
    BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
        fmt.Printf("任务 %s 即将执行\n", jobName)
    }).
    AfterJobRuns(func(jobID uuid.UUID, jobName string) {
        fmt.Printf("任务 %s 执行成功\n", jobName)
    }).
    AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
        fmt.Printf("任务 %s 执行失败: %v\n", jobName, err)
    }).
    Add()
```

## 任务管理

### 按别名移除任务

```go
scheduler.RemoveJobByAlias("my-jobs-alias")
```

### 按名称移除任务

```go
scheduler.RemoveJobByName("my-jobs")
```

### 按移除任务

```go
scheduler.RemoveJob(jobID)
```

### 获取所有任务

```go
jobs, err := scheduler.GetJobs()
```

## Web 监控端点

启用 Web 监控后，可使用以下端点：

| 端点 | 说明 |
|------|------|
| `GET /healthz` | 健康检查 |
| `GET /jobs` | 获取所有任务列表 |
| `GET /jobs/{job_id}/retries` | 获取任务的重试历史 |

## License

MIT License
