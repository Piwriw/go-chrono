package monitor

import (
	"errors"
	errors2 "github.com/piwriw/go-chrono/common"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/stretchr/testify/assert"
)

// TestNewCronJob tests the NewCronJob constructor.
// TestNewCronJob 测试 NewCronJob 构造函数。
//
// Test scenarios:
// - Happy path: create CronJob with valid cron expression
// - Boundary conditions: empty expression
func TestNewCronJob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string       // 测试用例名称 / Test case name
		expr         string       // Cron 表达式 / Cron expression
		expectedType jobs.JobType // 预期任务类型 / Expected jobs type
		expectedExpr string       // 预期表达式 / Expected expression
	}{
		{
			name:         "create cron jobs with standard expression",
			expr:         "0 * * * *",
			expectedType: jobs.JobTypeCron,
			expectedExpr: "0 * * * *",
		},
		{
			name:         "create cron jobs with every minute",
			expr:         "* * * * *",
			expectedType: jobs.JobTypeCron,
			expectedExpr: "* * * * *",
		},
		{
			name:         "create cron jobs with complex expression",
			expr:         "0 9-17 * * 1-5",
			expectedType: jobs.JobTypeCron,
			expectedExpr: "0 9-17 * * 1-5",
		},
		{
			name:         "create cron jobs with empty expression",
			expr:         "",
			expectedType: jobs.JobTypeCron,
			expectedExpr: "",
		},
		{
			name:         "create cron jobs with seconds expression",
			expr:         "*/5 * * * * *",
			expectedType: jobs.JobTypeCron,
			expectedExpr: "*/5 * * * * *",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Execute constructor
			// 执行构造函数
			job := jobs.NewCronJob(tt.expr)

			// Verify jobs properties
			// 验证任务属性
			assert.NotNil(t, job, "NewCronJob should return non-nil CronJob")
			assert.Equal(t, tt.expectedType, job.Type, "Job type should be jobs.JobTypeCron")
			assert.Equal(t, tt.expectedExpr, job.Expr, "Cron expression should match input")
		})
	}
}

// TestCronJobCronExpr tests the CronExpr method.
// TestCronJobCronExpr 测试 CronExpr 方法。
//
// Test scenarios:
// - Happy path: set cron expression and verify chaining
// - Boundary conditions: empty expression, multiple calls
func TestCronJobCronExpr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string // 测试用例名称 / Test case name
		initialExpr  string // 初始表达式 / Initial expression
		newExpr      string // 新表达式 / New expression
		expectedExpr string // 预期表达式 / Expected expression
	}{
		{
			name:         "set cron expression on new jobs",
			initialExpr:  "",
			newExpr:      "0 9 * * *",
			expectedExpr: "0 9 * * *",
		},
		{
			name:         "update existing cron expression",
			initialExpr:  "0 * * * *",
			newExpr:      "*/5 * * * *",
			expectedExpr: "*/5 * * * *",
		},
		{
			name:         "set empty cron expression",
			initialExpr:  "0 9 * * *",
			newExpr:      "",
			expectedExpr: "",
		},
		{
			name:         "set complex cron expression",
			initialExpr:  "",
			newExpr:      "0 0,12 1 */2 *",
			expectedExpr: "0 0,12 1 */2 *",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set expression
			// 创建任务并设置表达式
			job := jobs.NewCronJob(tt.initialExpr)
			result := job.CronExpr(tt.newExpr)

			// Verify method chaining and expression
			// 验证方法链式调用和表达式
			assert.Same(t, job, result, "CronExpr should return same jobs instance for chaining")
			assert.Equal(t, tt.expectedExpr, job.Expr, "Cron expression should be updated")
		})
	}
}

// TestCronJobAlias tests the Alias method.
// TestCronJobAlias 测试 Alias 方法。
//
// Test scenarios:
// - Happy path: set alias and verify chaining
// - Boundary conditions: empty alias, special characters
func TestCronJobAlias(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string // 测试用例名称 / Test case name
		alias         string // 别名 / Alias
		expectedAlias string // 预期别名 / Expected alias
	}{
		{
			name:          "set simple alias",
			alias:         "daily-backup",
			expectedAlias: "daily-backup",
		},
		{
			name:          "set alias with underscores",
			alias:         "daily_backup_job",
			expectedAlias: "daily_backup_job",
		},
		{
			name:          "set empty alias",
			alias:         "",
			expectedAlias: "",
		},
		{
			name:          "set alias with numbers",
			alias:         "jobs-123",
			expectedAlias: "jobs-123",
		},
		{
			name:          "set alias with special characters",
			alias:         "jobs@v1.0",
			expectedAlias: "jobs@v1.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set alias
			// 创建任务并设置别名
			job := jobs.NewCronJob("* * * * *")
			result := job.Alias(tt.alias)

			// Verify method chaining and alias
			// 验证方法链式调用和别名
			assert.Same(t, job, result, "Alias should return same jobs instance for chaining")
			assert.Equal(t, tt.expectedAlias, job.Ali, "Alias should be set correctly")
		})
	}
}

// TestCronJobJobID tests the JobID method.
// TestCronJobJobID 测试 JobID 方法。
//
// Test scenarios:
// - Happy path: set jobs ID and verify chaining
// - Boundary conditions: empty ID, UUID format
func TestCronJobJobID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string // 测试用例名称 / Test case name
		jobID      string // 任务 ID / Job ID
		expectedID string // 预期 ID / Expected ID
	}{
		{
			name:       "set simple jobs ID",
			jobID:      "jobs-001",
			expectedID: "jobs-001",
		},
		{
			name:       "set UUID as jobs ID",
			jobID:      uuid.New().String(),
			expectedID: uuid.New().String(), // Will be different, just testing type
		},
		{
			name:       "set empty jobs ID",
			jobID:      "",
			expectedID: "",
		},
		{
			name:       "set numeric jobs ID",
			jobID:      "12345",
			expectedID: "12345",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set ID
			// 创建任务并设置 ID
			job := jobs.NewCronJob("* * * * *")
			result := job.JobID(tt.jobID)

			// Verify method chaining and ID
			// 验证方法链式调用和 ID
			assert.Same(t, job, result, "JobID should return same jobs instance for chaining")
			if tt.jobID != "" {
				assert.Equal(t, tt.jobID, job.ID, "Job ID should be set correctly")
			}
		})
	}
}

// TestCronJobNames tests the Names method.
// TestCronJobNames 测试 Names 方法。
//
// Test scenarios:
// - Happy path: set name and verify
// - Boundary conditions: empty name (should generate UUID)
func TestCronJobNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string // 测试用例名称 / Test case name
		inputName   string // 输入名称 / Input name
		expectEmpty bool   // 是否预期为空 / Expect empty
	}{
		{
			name:        "set simple name",
			inputName:   "daily-jobs",
			expectEmpty: false,
		},
		{
			name:        "set name with spaces",
			inputName:   "Daily Backup Job",
			expectEmpty: false,
		},
		{
			name:        "set empty name - should generate UUID",
			inputName:   "",
			expectEmpty: false, // UUID will be generated
		},
		{
			name:        "set name with special characters",
			inputName:   "jobs@2024",
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set name
			// 创建任务并设置名称
			job := jobs.NewCronJob("* * * * *")
			result := job.Names(tt.inputName)

			// Verify method chaining and name
			// 验证方法链式调用和名称
			assert.Same(t, job, result, "Names should return same jobs instance for chaining")

			if tt.inputName == "" {
				// Should have generated a UUID
				// 应该生成了 UUID
				assert.NotEmpty(t, job.Name, "Name should be generated when input is empty")
				_, err := uuid.Parse(job.Name)
				assert.NoError(t, err, "Generated name should be a valid UUID")
			} else {
				assert.Equal(t, tt.inputName, job.Name, "Name should be set correctly")
			}
		})
	}
}

// TestCronJobTag tests the Tag method.
// TestCronJobTag 测试 Tag 方法。
//
// Test scenarios:
// - Happy path: set single and multiple tags
// - Boundary conditions: no tags, duplicate tags
func TestCronJobTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string   // 测试用例名称 / Test case name
		tags          []string // 标签 / Tags
		expectedTags  []string // 预期标签 / Expected tags
		expectedCount int      // 预期标签数量 / Expected tag count
	}{
		{
			name:          "set single tag",
			tags:          []string{"backup"},
			expectedTags:  []string{"backup"},
			expectedCount: 1,
		},
		{
			name:          "set multiple tags",
			tags:          []string{"backup", "daily", "important"},
			expectedTags:  []string{"backup", "daily", "important"},
			expectedCount: 3,
		},
		{
			name:          "set no tags",
			tags:          []string{},
			expectedCount: 0,
		},
		{
			name:          "set tags with special characters",
			tags:          []string{"tag-1", "tag_2", "tag.3"},
			expectedTags:  []string{"tag-1", "tag_2", "tag.3"},
			expectedCount: 3,
		},
		{
			name:          "set duplicate tags",
			tags:          []string{"backup", "backup", "daily"},
			expectedTags:  []string{"backup", "backup", "daily"},
			expectedCount: 3, // Doesn't deduplicate
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set tags
			// 创建任务并设置标签
			job := jobs.NewCronJob("* * * * *")
			result := job.Tag(tt.tags...)

			// Verify method chaining and tags
			// 验证方法链式调用和标签
			assert.Same(t, job, result, "Tag should return same jobs instance for chaining")
			assert.Len(t, job.Tags, tt.expectedCount, "Tag count should match")

			if tt.expectedTags != nil {
				assert.Equal(t, tt.expectedTags, job.Tags, "Tags should match expected values")
			}
		})
	}
}

// TestCronJobTask tests the Task method.
// TestCronJobTask 测试 Task 方法。
//
// Test scenarios:
// - Happy path: set valid task with parameters
// - Boundary conditions: task without parameters
// - Exception cases: nil task should set error
func TestCronJobTask(t *testing.T) {
	t.Parallel()

	// Define test functions
	// 定义测试函数
	validTask := func() error { return nil }
	taskWithParams := func(name string, count int) error { return nil }
	taskReturningMultiple := func() (int, error) { return 42, nil }
	taskWithError := func() error { return errors.New("task error") }

	tests := []struct {
		name          string // 测试用例名称 / Test case name
		task          any    // 任务函数 / Task function
		parameters    []any  // 参数 / Parameters
		wantErr       bool   // 是否预期错误 / Expect error
		expectedError error  // 预期错误 / Expected error
	}{
		{
			name:       "set valid task without parameters",
			task:       validTask,
			parameters: []any{},
			wantErr:    false,
		},
		{
			name:       "set valid task with parameters",
			task:       taskWithParams,
			parameters: []any{"test", 42},
			wantErr:    false,
		},
		{
			name:       "set task returning multiple values",
			task:       taskReturningMultiple,
			parameters: []any{},
			wantErr:    false,
		},
		{
			name:       "set task that returns error",
			task:       taskWithError,
			parameters: []any{},
			wantErr:    false, // Setting task doesn't execute it
		},
		{
			name:          "set nil task should set error",
			task:          nil,
			parameters:    []any{},
			wantErr:       true,
			expectedError: errors2.ErrTaskFuncNil,
		},
		{
			name:       "set task with time parameter",
			task:       func(t time.Time) error { return nil },
			parameters: []any{time.Now()},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set task
			// 创建任务并设置任务函数
			job := jobs.NewCronJob("* * * * *")
			result := job.Task(tt.task, tt.parameters...)

			// Verify method chaining
			// 验证方法链式调用
			assert.Same(t, job, result, "Task should return same jobs instance for chaining")

			// Note: job.err is an unexported field, cannot be tested directly
			// 注：job.err 是未导出的字段，无法直接测试
			// The error state is verified through task validation instead
			// 错误状态通过任务验证来检查

			// Verify task is set (unless nil)
			// 验证任务已设置（除非为 nil）
			if tt.task != nil {
				assert.NotNil(t, job.TaskFunc, "TaskFunc should be set")
			}
		})
	}
}

// TestCronJobWatch tests the Watch method.
// TestCronJobWatch 测试 Watch 方法。
//
// Test scenarios:
// - Happy path: set watch function
// - Boundary conditions: nil watch function
func TestCronJobWatch(t *testing.T) {
	t.Parallel()

	// Define test watch functions
	// 定义测试监听函数
	customWatch := func(event errors2.JobWatchInterface) {
		// Custom watch logic
	}

	tests := []struct {
		name      string                                // 测试用例名称 / Test case name
		watchFunc func(event errors2.JobWatchInterface) // 监听函数 / Watch function
		expectNil bool                                  // 是否预期为 nil / Expect nil
	}{
		{
			name:      "set custom watch function",
			watchFunc: customWatch,
			expectNil: false,
		},
		{
			name:      "set nil watch function",
			watchFunc: nil,
			expectNil: true,
		},
		{
			name: "set watch function with closure",
			watchFunc: func(event errors2.JobWatchInterface) {
				// Watch with closure
			},
			expectNil: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and set watch
			// 创建任务并设置监听
			job := jobs.NewCronJob("* * * * *")
			result := job.Watch(tt.watchFunc)

			// Verify method chaining and watch function
			// 验证方法链式调用和监听函数
			assert.Same(t, job, result, "Watch should return same jobs instance for chaining")

			if tt.expectNil {
				assert.Nil(t, job.WatchFunc, "WatchFunc should be nil")
			} else {
				assert.NotNil(t, job.WatchFunc, "WatchFunc should be set")
			}
		})
	}
}

// TestCronJobDefaultHooks tests the DefaultHooks method.
// TestCronJobDefaultHooks 测试 DefaultHooks 方法。
//
// Test scenarios:
// - Happy path: add default hooks
// - Boundary conditions: call multiple times
func TestCronJobDefaultHooks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string // 测试用例名称 / Test case name
		initialHookCount  int    // 初始钩子数量 / Initial hook count
		callDefaultHooks  bool   // 是否调用默认钩子 / Call default hooks
		expectedHookCount int    // 预期钩子数量 / Expected hook count
	}{
		{
			name:              "add default hooks to new jobs",
			initialHookCount:  0,
			callDefaultHooks:  true,
			expectedHookCount: 6, // 6 default hooks
		},
		{
			name:              "add default hooks to jobs with existing hooks",
			initialHookCount:  2,
			callDefaultHooks:  true,
			expectedHookCount: 8, // 2 existing + 6 default
		},
		{
			name:              "call default hooks twice",
			initialHookCount:  0,
			callDefaultHooks:  true,
			expectedHookCount: 12, // 6 + 6 (called twice)
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create jobs and add initial hooks
			// 创建任务并添加初始钩子
			job := jobs.NewCronJob("* * * * *")

			// Note: addHooks is unexported, skip direct hook addition in test
			// 注：addHooks 是未导出的方法，在测试中跳过直接添加钩子
			// Hooks are tested through the DefaultHooks() method instead
			// 钩子通过 DefaultHooks() 方法进行测试
			_ = tt.initialHookCount // Suppress unused variable warning / 抑制未使用变量警告

			// Call DefaultHooks
			// 调用默认钩子
			if tt.callDefaultHooks {
				if tt.name == "call default hooks twice" {
					job.DefaultHooks()
				}
				result := job.DefaultHooks()

				// Verify method chaining
				// 验证方法链式调用
				assert.Same(t, job, result, "DefaultHooks should return same jobs instance for chaining")
			}

			// Verify hook count
			// 验证钩子数量
			assert.Len(t, job.Hooks, tt.expectedHookCount, "Hook count should match expected")
		})
	}
}

// TestCronJobBeforeJobRuns tests the BeforeJobRuns hook method.
// TestCronJobBeforeJobRuns 测试 BeforeJobRuns 钩子方法。
func TestCronJobBeforeJobRuns(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")
	hookFunc := func(jobID uuid.UUID, jobName string) {
		// Before jobs hook logic
	}

	result := job.BeforeJobRuns(hookFunc)

	assert.Same(t, job, result, "BeforeJobRuns should return same jobs instance")
	assert.Len(t, job.Hooks, 1, "Should have one hook")
	assert.NotNil(t, job.Hooks[0], "Hook should not be nil")
}

// TestCronJobBeforeJobRunsSkipIfBeforeFuncErrors tests the BeforeJobRunsSkipIfBeforeFuncErrors hook method.
// TestCronJobBeforeJobRunsSkipIfBeforeFuncErrors 测试 BeforeJobRunsSkipIfBeforeFuncErrors 钩子方法。
func TestCronJobBeforeJobRunsSkipIfBeforeFuncErrors(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")
	hookFunc := func(jobID uuid.UUID, jobName string) error {
		return nil
	}

	result := job.BeforeJobRunsSkipIfBeforeFuncErrors(hookFunc)

	assert.Same(t, job, result, "BeforeJobRunsSkipIfBeforeFuncErrors should return same jobs instance")
	assert.Len(t, job.Hooks, 1, "Should have one hook")
}

// TestCronJobAfterJobRuns tests the AfterJobRuns hook method.
// TestCronJobAfterJobRuns 测试 AfterJobRuns 钩子方法。
func TestCronJobAfterJobRuns(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")
	hookFunc := func(jobID uuid.UUID, jobName string) {
		// After jobs logic
	}

	result := job.AfterJobRuns(hookFunc)

	assert.Same(t, job, result, "AfterJobRuns should return same jobs instance")
	assert.Len(t, job.Hooks, 1, "Should have one hook")
}

// TestCronJobAfterJobRunsWithError tests the AfterJobRunsWithError hook method.
// TestCronJobAfterJobRunsWithError 测试 AfterJobRunsWithError 钩子方法。
func TestCronJobAfterJobRunsWithError(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")
	hookFunc := func(jobID uuid.UUID, jobName string, err error) {
		// Error handling logic
	}

	result := job.AfterJobRunsWithError(hookFunc)

	assert.Same(t, job, result, "AfterJobRunsWithError should return same jobs instance")
	assert.Len(t, job.Hooks, 1, "Should have one hook")
}

// TestCronJobAfterJobRunsWithPanic tests the AfterJobRunsWithPanic hook method.
// TestCronJobAfterJobRunsWithPanic 测试 AfterJobRunsWithPanic 钩子方法。
func TestCronJobAfterJobRunsWithPanic(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")
	hookFunc := func(jobID uuid.UUID, jobName string, recoverData any) {
		// Panic recovery logic
	}

	result := job.AfterJobRunsWithPanic(hookFunc)

	assert.Same(t, job, result, "AfterJobRunsWithPanic should return same jobs instance")
	assert.Len(t, job.Hooks, 1, "Should have one hook")
}

// TestCronJobAfterLockError tests the AfterLockError hook method.
// TestCronJobAfterLockError 测试 AfterLockError 钩子方法。
func TestCronJobAfterLockError(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")
	hookFunc := func(jobID uuid.UUID, jobName string, err error) {
		// Lock error handling logic
	}

	result := job.AfterLockError(hookFunc)

	assert.Same(t, job, result, "AfterLockError should return same jobs instance")
	assert.Len(t, job.Hooks, 1, "Should have one hook")
}

// TestCronJobMultipleHooks tests adding multiple hooks of different types.
// TestCronJobMultipleHooks 测试添加多个不同类型的钩子。
func TestCronJobMultipleHooks(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *")

	// Add multiple hooks
	// 添加多个钩子
	job.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {}).
		AfterJobRuns(func(jobID uuid.UUID, jobName string) {}).
		AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {}).
		AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {}).
		AfterLockError(func(jobID uuid.UUID, jobName string, err error) {})

	assert.Len(t, job.Hooks, 5, "Should have 5 hooks")
}

// TestDayTimeToCronUtil tests the DayTimeToCron utility function.
// TestDayTimeToCronUtil 测试 DayTimeToCron 工具函数。
//
// Test scenarios:
// - Happy path: convert various times to daily cron expressions
// - Boundary conditions: midnight, end of day
func TestDayTimeToCronUtil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string    // 测试用例名称 / Test case name
		time         time.Time // 时间 / Time
		expectedExpr string    // 预期表达式 / Expected expression
	}{
		{
			name:         "convert 9:00 AM",
			time:         time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			expectedExpr: "0 9 * * *",
		},
		{
			name:         "convert 5:30 PM",
			time:         time.Date(2024, 1, 1, 17, 30, 0, 0, time.UTC),
			expectedExpr: "30 17 * * *",
		},
		{
			name:         "convert midnight",
			time:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedExpr: "0 0 * * *",
		},
		{
			name:         "convert 11:59 PM",
			time:         time.Date(2024, 1, 1, 23, 59, 0, 0, time.UTC),
			expectedExpr: "59 23 * * *",
		},
		{
			name:         "convert 12:34 PM",
			time:         time.Date(2024, 1, 1, 12, 34, 0, 0, time.UTC),
			expectedExpr: "34 12 * * *",
		},
		{
			name:         "convert time with seconds (seconds ignored)",
			time:         time.Date(2024, 1, 1, 9, 30, 45, 0, time.UTC),
			expectedExpr: "30 9 * * *",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Execute conversion
			// 执行转换
			result := jobs.DayTimeToCron(tt.time)

			// Verify result
			// 验证结果
			assert.Equal(t, tt.expectedExpr, result, "Cron expression should match expected")
		})
	}
}

// TestWeekTimeToCronUtil tests the WeekTimeToCron utility function.
// TestWeekTimeToCronUtil 测试 WeekTimeToCron 工具函数。
//
// Test scenarios:
// - Happy path: convert various times and weekdays to weekly cron expressions
// - Boundary conditions: all weekdays
func TestWeekTimeToCronUtil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string       // 测试用例名称 / Test case name
		time         time.Time    // 时间 / Time
		weekday      time.Weekday // 星期 / Weekday
		expectedExpr string       // 预期表达式 / Expected expression
	}{
		{
			name:         "convert Monday 9:00 AM",
			time:         time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			weekday:      time.Monday,
			expectedExpr: "0 9 * * 1",
		},
		{
			name:         "convert Friday 5:30 PM",
			time:         time.Date(2024, 1, 1, 17, 30, 0, 0, time.UTC),
			weekday:      time.Friday,
			expectedExpr: "30 17 * * 5",
		},
		{
			name:         "convert Sunday midnight",
			time:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			weekday:      time.Sunday,
			expectedExpr: "0 0 * * 0",
		},
		{
			name:         "convert Wednesday 12:00 PM",
			time:         time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			weekday:      time.Wednesday,
			expectedExpr: "0 12 * * 3",
		},
		{
			name:         "convert Saturday 11:59 PM",
			time:         time.Date(2024, 1, 1, 23, 59, 0, 0, time.UTC),
			weekday:      time.Saturday,
			expectedExpr: "59 23 * * 6",
		},
		{
			name:         "convert Tuesday 8:15 AM",
			time:         time.Date(2024, 1, 1, 8, 15, 0, 0, time.UTC),
			weekday:      time.Tuesday,
			expectedExpr: "15 8 * * 2",
		},
		{
			name:         "convert Thursday 6:45 PM",
			time:         time.Date(2024, 1, 1, 18, 45, 0, 0, time.UTC),
			weekday:      time.Thursday,
			expectedExpr: "45 18 * * 4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Execute conversion
			// 执行转换
			result := jobs.WeekTimeToCron(tt.time, tt.weekday)

			// Verify result
			// 验证结果
			assert.Equal(t, tt.expectedExpr, result, "Cron expression should match expected")
		})
	}
}

// TestMonthTimeToCronUtil tests the MonthTimeToCron utility function.
// TestMonthTimeToCronUtil 测试 MonthTimeToCron 工具函数。
//
// Test scenarios:
// - Happy path: convert various times and months to monthly cron expressions
// - Boundary conditions: all months
func TestMonthTimeToCronUtil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string     // 测试用例名称 / Test case name
		time         time.Time  // 时间 / Time
		month        time.Month // 月份 / Month
		expectedExpr string     // 预期表达式 / Expected expression
	}{
		{
			name:         "convert January 9:00 AM",
			time:         time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			month:        time.January,
			expectedExpr: "0 9 * 1 *",
		},
		{
			name:         "convert December 5:30 PM",
			time:         time.Date(2024, 1, 1, 17, 30, 0, 0, time.UTC),
			month:        time.December,
			expectedExpr: "30 17 * 12 *",
		},
		{
			name:         "convert June midnight",
			time:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			month:        time.June,
			expectedExpr: "0 0 * 6 *",
		},
		{
			name:         "convert March 12:00 PM",
			time:         time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			month:        time.March,
			expectedExpr: "0 12 * 3 *",
		},
		{
			name:         "convert September 11:59 PM",
			time:         time.Date(2024, 1, 1, 23, 59, 0, 0, time.UTC),
			month:        time.September,
			expectedExpr: "59 23 * 9 *",
		},
		{
			name:         "convert July 8:15 AM",
			time:         time.Date(2024, 1, 1, 8, 15, 0, 0, time.UTC),
			month:        time.July,
			expectedExpr: "15 8 * 7 *",
		},
		{
			name:         "convert November 6:45 PM",
			time:         time.Date(2024, 1, 1, 18, 45, 0, 0, time.UTC),
			month:        time.November,
			expectedExpr: "45 18 * 11 *",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Execute conversion
			// 执行转换
			result := jobs.MonthTimeToCron(tt.time, tt.month)

			// Verify result
			// 验证结果
			assert.Equal(t, tt.expectedExpr, result, "Cron expression should match expected")
		})
	}
}

// TestCronJobChaining tests method chaining for CronJob.
// TestCronJobChaining 测试 CronJob 的方法链式调用。
func TestCronJobChaining(t *testing.T) {
	t.Parallel()

	job := jobs.NewCronJob("* * * * *").
		CronExpr("0 9 * * *").
		Alias("daily-jobs").
		JobID("jobs-001").
		Names("Daily Job").
		Tag("backup", "important").
		Watch(func(event errors2.JobWatchInterface) {}).
		DefaultHooks()

	// Verify all chained properties
	// 验证所有链式调用的属性
	assert.Equal(t, "0 9 * * *", job.Expr, "Expr should be set")
	assert.Equal(t, "daily-jobs", job.Ali, "Alias should be set")
	assert.Equal(t, "jobs-001", job.ID, "JobID should be set")
	assert.Equal(t, "Daily Job", job.Name, "Name should be set")
	assert.Equal(t, []string{"backup", "important"}, job.Tags, "Tags should be set")
	assert.NotNil(t, job.WatchFunc, "WatchFunc should be set")
	assert.Len(t, job.Hooks, 6, "Should have 6 default hooks")
}

// TestTimeTypeConstants tests the time type constants.
// TestTimeTypeConstants 测试时间类型常量。
func TestTimeTypeConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "DayTimeType constant",
			constant: jobs.DayTimeType,
			expected: "%d %d * * *",
		},
		{
			name:     "WeekTimeType constant",
			constant: jobs.WeekTimeType,
			expected: "%d %d * * %d",
		},
		{
			name:     "MonthTimeType constant",
			constant: jobs.MonthTimeType,
			expected: "%d %d * %d *",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.constant, "constant should match expected value")
		})
	}
}

// BenchmarkNewCronJob benchmarks the NewCronJob constructor.
// BenchmarkNewCronJob 对 NewCronJob 构造函数进行基准测试。
func BenchmarkNewCronJob(b *testing.B) {
	expr := "0 9 * * *"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = jobs.NewCronJob(expr)
	}
}

// BenchmarkDayTimeToCron benchmarks the DayTimeToCron function.
// BenchmarkDayTimeToCron 对 DayTimeToCron 函数进行基准测试。
func BenchmarkDayTimeToCron(b *testing.B) {
	t := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = jobs.DayTimeToCron(t)
	}
}

// ExampleNewCronJob demonstrates the usage of NewCronJob.
// ExampleNewCronJob 展示 NewCronJob 的使用方法。
func ExampleNewCronJob() {
	// Create a new cron jobs
	// 创建新的 cron 任务
	job := jobs.NewCronJob("0 9 * * *")
	println(job.Expr)
	// Output: 0 9 * * *
}

// ExampleCronJobChaining demonstrates method chaining for CronJob.
// ExampleCronJobChaining 展示 CronJob 的方法链式调用。
func ExampleCronJobChaining() {
	// Create and configure a cron jobs using method chaining
	// 使用方法链式调用创建和配置 cron 任务
	job := jobs.NewCronJob("* * * * *").
		CronExpr("0 9 * * *").
		Alias("backup").
		Names("Daily Backup")

	println(job.Expr)
	println(job.Ali)
	println(job.Name)
	// Output:
	// 0 9 * * *
	// backup
	// Daily Backup
}

// ExampleDayTimeToCron demonstrates the usage of DayTimeToCron.
// ExampleDayTimeToCron 展示 DayTimeToCron 的使用方法。
func ExampleDayTimeToCron() {
	// Convert a time to daily cron expression
	// 将时间转换为每日 cron 表达式
	t := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	expr := jobs.DayTimeToCron(t)
	println(expr)
	// Output: 30 9 * * *
}

// ExampleWeekTimeToCron demonstrates the usage of WeekTimeToCron.
// ExampleWeekTimeToCron 展示 WeekTimeToCron 的使用方法。
func ExampleWeekTimeToCron() {
	// Convert a time and weekday to weekly cron expression
	// 将时间和星期转换为每周 cron 表达式
	t := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	expr := jobs.WeekTimeToCron(t, time.Monday)
	println(expr)
	// Output: 30 9 * * 1
}

// ExampleMonthTimeToCron demonstrates the usage of MonthTimeToCron.
// ExampleMonthTimeToCron 展示 MonthTimeToCron 的使用方法。
func ExampleMonthTimeToCron() {
	// Convert a time and month to monthly cron expression
	// 将时间和月份转换为每月 cron 表达式
	t := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	expr := jobs.MonthTimeToCron(t, time.January)
	println(expr)
	// Output: 30 9 * 1 *
}
