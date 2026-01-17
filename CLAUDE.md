# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**go-chrono** is an enhanced job scheduling library for Go built on top of `go-co-op/gocron/v2`. It provides a fluent API for scheduling various job types with monitoring, web UI, Prometheus metrics, and flexible job lifecycle management.

## Common Commands

### Build and Test
```bash
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v -run TestFunctionName ./...
```

### Linting and Formatting
```bash
# Run linter (automatically installs golangci-lint v1.61.0 if needed)
make lint

# Format code with goimports and gofmt
make format

# Clean temporary files
make clean
```

## Architecture

### Core Components

**Scheduler** (`schedule.go`)
- Main orchestrator extending `gocron.Scheduler`
- Manages job lifecycle via job lookup by ID, name, or alias
- Coordinates monitoring and web interfaces
- Thread-safe with mutex-protected maps for aliases, watch functions, and job types

**Job Type System** (`job_type.go`)
- Six job types: `JobTypeOnce`, `JobTypeCron`, `JobTypeDaily`, `JobTypeWeekly`, `JobTypeMonthly`, `JobInterval`
- Each job type has a dedicated client for fluent configuration
- Clients follow builder pattern with method chaining

**Job Clients**
- `CronJobClient` (`cron_client.go`)
- `OnceJobClient` (`once_client.go`)
- `IntervalJobClient` (`interval_job_client.go`)
- `DailyJobClient` (`daliy_client.go`)
- `WeeklyJobClient` (`weekly_job_client.go`)
- `MonthlyJobClient` (`monthly_job_client.go`)

**Monitor System** (`monitor.go`)
- `SchedulerMonitor`: Interface for tracking job executions
- `defaultSchedulerMonitor`: Implementation with circular buffer for recent events
- `JobWatchInterface`: Provides job event observation hooks
- Configurable event ID generation (UUID or JobID+Name+Timestamp)

**Options System** (`option.go`)
- Functional options pattern for flexible scheduler configuration
- `WithAliasMode()`: Enable human-readable job aliases
- `WithWatch()`: Enable event watching
- `WithWebMonitor()`: Enable HTTP dashboard
- `WithLimit()`: Set max concurrent jobs
- `WithPrometheus()`: Enable metrics collection

**Retry System** (`retry/` package)
- `RetryPolicy`: Interface defining retry behavior
  - `FixedIntervalPolicy`: Fixed time interval between retries
  - `ExponentialBackoffPolicy`: Exponential backoff with max interval cap
  - `JitterPolicy`: Adds randomness to base policy to prevent thundering herd
- `RetryConfig`: Configuration object for retry behavior (max retries, policy, callbacks)
- `RetryExecutor`: Executes tasks with retry logic, records retry events
- `RetryEvent`: Represents a single retry attempt with timing and error information
- All job clients support `WithRetry()` and `WithRetryConfig()` methods

### Design Patterns

1. **Builder Pattern**: Job clients use method chaining (`Names().CronExpr().Task().Tag().Build()`)
2. **Option Pattern**: Functional options for scheduler configuration
3. **Observer Pattern**: Job event monitoring and hooks
4. **Strategy Pattern**: Configurable event ID generation and retry policies

## Code Conventions

### Bilingual Documentation (Required)
All functions must have both English and Chinese comments following this exact format:

```go
// FunctionName performs an action.
// 函数名称执行一个操作。
//
// Parameters:
//	param - Parameter description / 参数描述
//
// Returns:
//	ReturnType - Return value description / 返回值描述
func FunctionName(param string) (string, error) {
    // implementation
}
```

- Function comments start with function name, not struct name
- Omit Parameters/Returns sections when empty
- Add comments for key business logic and algorithms

### Error Handling
```go
// Always wrap errors with context
return fmt.Errorf("operation failed: %w", err)

// Define custom error types in error.go
type MyError struct {
    Code    int
    Message string
}
```

### Concurrency
- Use `sync.Mutex` to protect shared state (maps, slices)
- Always lock before accessing `aliasMap`, `watchFuncMap`, `jobTypeMap`

### Testing
- Test files: `schedule_test.go`, `retry_integration_test.go`
- Use `github.com/google/uuid` for test IDs
- Mock `SchedulerMonitor` and `JobWatchInterface` as needed
- Retry mechanism tests verify policies, configs, events, and monitor integration

## Dependencies

- `github.com/go-co-op/gocron/v2` (v2.16.2): Core scheduling
- `github.com/google/uuid` (v1.6.0): UUID generation
- `github.com/prometheus/client_golang` (v1.22.0): Metrics
- `gopkg.in/natefinch/lumberjack.v2` (v2.2.1): Log rotation

## Linting Configuration

The project uses `golangci-lint` v1.61.0 with custom rules:
- Max line length: 200 characters
- Enabled linters: errcheck, gofmt, goimports, gosimple, govet, bodyclose, goconst, ineffassign, staticcheck
- Custom formatters: gci, gofmt, gofumpt, goimports, golines
- Auto-rewrite rules: `interface{}` → `any`, `a[b:len(a)]` → `a[b:]`, var declaration simplification

## Web Monitor

When `WithWebMonitor()` is enabled, the scheduler exposes an HTTP dashboard for:
- Viewing job status and recent executions
- Managing jobs (start/stop/remove)
- Real-time monitoring via WebSocket
- Querying retry history: `GET /jobs/{job_id}/retries`

## Graceful Shutdown

Call `scheduler.Shutdown()` to terminate all scheduled jobs cleanly. The scheduler context should be managed by the caller.
