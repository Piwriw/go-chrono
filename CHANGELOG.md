# Changelog

All notable changes to go-chrono will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - 2025-01-18

### Added

- **Retry History Configuration** - Make retry history storage size configurable
- **Integration Tests** - Comprehensive test coverage for retry mechanism scenarios
- **Web API Endpoint** - `GET /jobs/{job_id}/retries` for querying retry history

### Changed

- **Code Formatting** - Fixed formatting issues across the codebase
- **Documentation** - Improved README structure and clarity
- **Comments** - Enhanced bilingual documentation for retry-related features

### Fixed

- **Thread Safety** - Added mutex protection to monitor methods
- **Comment Format** - Corrected bilingual comment formatting in retry_event.go

## [2.0.0] - 2025-01-18

> **⚠️ BREAKING CHANGES** - Major version upgrade with significant package restructuring.
> **⚠️ 重大变更** - 包含大量包结构重构的重大版本升级。

**Current Branch:** `v2.0.0` | **Main Branch:** `v1.0.0`

### Migration Guide / 迁移指南

**Import Path Changes / 导入路径变更:**

| Before (v1.0.0) | After (v2.0.0) |
|------------------|-----------------|
| `github.com/piwriw/go-chrono/pkg`.JobType | `github.com/piwriw/go-chrono/jobs`.JobType |
| `github.com/piwriw/go-chrono/pkg`.CronJob | `github.com/piwriw/go-chrono/jobs`.CronJob |
| `github.com/piwriw/go-chrono/pkg`.IntervalJob | `github.com/piwriw/go-chrono/jobs`.IntervalJob |
| `github.com/piwriw/go-chrono/pkg`.DailyJob | `github.com/piwriw/go-chrono/jobs`.DailyJob |
| `github.com/piwriw/go-chrono/pkg`.WeeklyJob | `github.com/piwriw/go-chrono/jobs`.WeeklyJob |
| `github.com/piwriw/go-chrono/pkg`.MonthlyJob | `github.com/piwriw/go-chrono/jobs`.MonthlyJob |
| `github.com/piwriw/go-chrono/pkg`.OnceJob | `github.com/piwriw/go-chrono/jobs`.OnceJob |

**Example / 示例:**

```go
// v1.0.0 (OLD)
import "github.com/piwriw/go-chrono/pkg"

jobType := pkg.JobTypeCron

// v2.0.0 (NEW)
import "github.com/piwriw/go-chrono/jobs"

jobType := jobs.JobTypeCron
```

### Added

- **Retry Mechanism** - 完整的重试机制支持
  - 添加 `RetryExecutor` 用于执行带重试逻辑的任务
  - 支持多种重试策略：`FixedIntervalPolicy`、`ExponentialBackoffPolicy`、`JitterPolicy`
  - 所有 JobClient 支持 `WithRetry()` 和 `WithRetryConfig()` 方法
  - Web 监控器新增 `/jobs/{job_id}/retries` API 端点查询重试历史

### Changed

- **包结构重构** - 优化项目代码组织
  - 将 `JobType` 移至 `jobs/types.go`
  - 将错误定义移至 `common/errors.go`
  - 将钩子函数移至 `common/hooks.go`
  - 将选项定义移至 `common/options.go`
  - 创建独立的 `client/`、`jobs/`、`monitor/`、`pkg/` 子目录

### Fixed

- **循环导入** - 修复包之间的循环依赖问题
  - 移除 `common/options.go` 对 `monitor` 包的导入
  - 使用 `interface{}` 类型避免类型依赖
  - 修复 `monitor/web.go` 中的无效包导入
  - 更新所有测试文件的包引用

- **测试修复** - 修复因包重组导致的测试失败
  - 更新 `scheduler/` 包中的所有测试文件
  - 修复 `monitor/` 包中的类型引用
  - 修正字段名大小写问题（`enabled` → `Enabled`）

### Technical Details

#### 修复的循环依赖路径

修复前：
```
client → common → monitor → jobs → common (循环!)
```

修复后：
```
client → common (基础层)
client → jobs (业务层)
client → monitor (业务层)
jobs → common
monitor → common
scheduler → common, jobs, monitor
```

#### 修改的文件列表

**核心修复：**
- `common/options.go` - 移除 monitor 导入，WatchFunc 改用 interface{}
- `monitor/web.go` - 移除无效的 pkg 导入
- `scheduler/scheduler.go` - 添加类型断言处理

**测试文件：**
- `scheduler/job_type_test.go`
- `scheduler/schedule_unit_test.go`
- `scheduler/option_test.go`
- `monitor/cron_test.go`
- `monitor/retry_integration_test.go`
- `monitor/monitor_test.go`

### Documentation

- 添加代码审查报告 `docs/review-reports/code-review-2025-01-18.md`
- 更新项目文档以反映新的包结构

---

## [0.2.0] - 2025-01-XX

### Added

- **重试机制** - 完整的作业重试功能
  - 新增 `retry` 包，包含重试执行器和策略
  - 支持固定间隔、指数退避和抖动重试策略
  - 可配置最大重试次数和重试回调
  - 重试历史记录和查询功能

- **Web 监控增强**
  - 新增重试历史查询 API 端点
  - 支持按 job ID 查询重试记录

- **配置选项**
  - 新增 `WithRetry()` 和 `WithRetryConfig()` 选项
  - 所有 JobClient 支持重试配置

### Changed

- 优化重试历史的内存管理
- 改进监控器的线程安全性

---

## [1.0.0] - 2025-01-XX

### Added

- **基础调度功能**
  - 支持 Cron 表达式调度
  - 支持间隔调度
  - 支持一次性任务
  - 支持每日、每周、每月任务

- **监控功能**
  - 作业执行事件记录
  - Web 监控仪表板
  - Prometheus 指标支持

- **任务管理**
  - 任务别名
  - 任务标签
  - 任务限制（并发限制）

---

## 版本说明

- **[2.0.0]** - v2.0.0 分支（当前分支）- 重大版本升级，包含不兼容的包结构变更
- **[1.0.0]** - main 分支 - 稳定的 v1.0.0 版本

**Branch Information / 分支信息:**
- `v2.0.0` - Current development branch with refactored package structure
- `main` - Stable v1.0.0 release branch

**Versioning / 版本规则:**
- Major version (X.0.0) - Breaking changes / 不兼容的重大变更
- Minor version (0.Y.0) - New features / 新功能
- Patch version (0.0.Z) - Bug fixes / 错误修复

### 变更类型说明

- **Added** - 新增功能
- **Changed** - 对现有功能的变更
- **Deprecated** - 即将移除的功能
- **Removed** - 已移除的功能
- **Fixed** - 错误修复
- **Security** - 安全相关的修复或改进
