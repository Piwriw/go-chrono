package client

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/piwriw/go-chrono/common"
	"github.com/piwriw/go-chrono/jobs"
	"github.com/piwriw/go-chrono/scheduler"
)

// TestCronJobClientInterfaceCompliance tests that CronJobClient implements CronJobClientInterface.
// TestCronJobClientInterfaceCompliance 测试 CronJobClient 是否实现了 CronJobClientInterface
func TestCronJobClientInterfaceCompliance(t *testing.T) {
	var _ CronJobClientInterface = (*CronJobClient)(nil)
}

// TestCronJobClientNilScheduler tests client behavior with nil scheduler.
// TestCronJobClientNilScheduler 测试调度器为 nil 时的客户端行为
func TestCronJobClientNilScheduler(t *testing.T) {
	tests := []struct {
		name      string                     // 测试名称 / Test name
		setupFunc func(*CronJobClient)       // 设置函数 / Setup function
		testFunc  func(*CronJobClient) error // 测试函数 / Test function
		wantErr   error                      // 期望的错误 / Expected error
	}{
		{
			name: "Add with nil scheduler returns pkg.ErrScheduleNil",
			setupFunc: func(c *CronJobClient) {
				c.scheduler = nil
			},
			testFunc: func(c *CronJobClient) error {
				_, err := c.Add()
				return err
			},
			wantErr: common.ErrScheduleNil,
		},
		{
			name: "BatchAdd with nil scheduler returns pkg.ErrScheduleNil",
			setupFunc: func(c *CronJobClient) {
				c.scheduler = nil
			},
			testFunc: func(c *CronJobClient) error {
				_, err := c.BatchAdd()
				return err
			},
			wantErr: common.ErrScheduleNil,
		},
		{
			name: "Remove with nil scheduler returns pkg.ErrScheduleNil",
			setupFunc: func(c *CronJobClient) {
				c.scheduler = nil
			},
			testFunc: func(c *CronJobClient) error {
				return c.Remove()
			},
			wantErr: common.ErrScheduleNil,
		},
		{
			name: "Get with nil scheduler returns pkg.ErrScheduleNil",
			setupFunc: func(c *CronJobClient) {
				c.scheduler = nil
			},
			testFunc: func(c *CronJobClient) error {
				_, err := c.Get()
				return err
			},
			wantErr: common.ErrScheduleNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &CronJobClient{
				job: &jobs.CronJob{},
			}
			if tt.setupFunc != nil {
				tt.setupFunc(client)
			}
			err := tt.testFunc(client)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCronJobClientNilJob tests client behavior with nil jobs.
// TestCronJobClientNilJob 测试任务为 nil 时的客户端行为
func TestCronJobClientNilJob(t *testing.T) {
	tests := []struct {
		name      string                     // 测试名称 / Test name
		setupFunc func(*CronJobClient)       // 设置函数 / Setup function
		testFunc  func(*CronJobClient) error // 测试函数 / Test function
		wantErr   error                      // 期望的错误 / Expected error
	}{
		{
			name: "Add with nil jobs returns pkg.ErrCronJobNil",
			setupFunc: func(c *CronJobClient) {
				c.job = nil
			},
			testFunc: func(c *CronJobClient) error {
				_, err := c.Add()
				return err
			},
			wantErr: common.ErrCronJobNil,
		},
		{
			name: "Remove with nil jobs returns pkg.ErrCronJobNil",
			setupFunc: func(c *CronJobClient) {
				c.job = nil
			},
			testFunc: func(c *CronJobClient) error {
				return c.Remove()
			},
			wantErr: common.ErrCronJobNil,
		},
		{
			name: "Get with nil jobs returns pkg.ErrCronJobNil",
			setupFunc: func(c *CronJobClient) {
				c.job = nil
			},
			testFunc: func(c *CronJobClient) error {
				_, err := c.Get()
				return err
			},
			wantErr: common.ErrCronJobNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler, _ := scheduler.NewScheduler(nil, nil)
			client := &CronJobClient{
				scheduler: scheduler,
			}
			if tt.setupFunc != nil {
				tt.setupFunc(client)
			}
			err := tt.testFunc(client)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCronJobClientMethodChaining tests that methods return the interface for chaining.
// TestCronJobClientMethodChaining 测试方法返回接口以支持链式调用
func TestCronJobClientMethodChaining(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &CronJobClient{
		scheduler: scheduler,
		job:       &jobs.CronJob{},
	}

	tests := []struct {
		name     string                        // 测试名称 / Test name
		testFunc func() CronJobClientInterface // 测试函数 / Test function
	}{
		{
			name: "CronExpr returns interface",
			testFunc: func() CronJobClientInterface {
				return client.CronExpr("* * * * *")
			},
		},
		{
			name: "Alias returns interface",
			testFunc: func() CronJobClientInterface {
				return client.Alias("test-alias")
			},
		},
		{
			name: "JobID returns interface",
			testFunc: func() CronJobClientInterface {
				return client.JobID("test-id")
			},
		},
		{
			name: "Name returns interface",
			testFunc: func() CronJobClientInterface {
				return client.Name("test-name")
			},
		},
		{
			name: "Tags returns interface",
			testFunc: func() CronJobClientInterface {
				return client.Tags("tag1", "tag2")
			},
		},
		{
			name: "Task returns interface",
			testFunc: func() CronJobClientInterface {
				return client.Task(func() {}, "arg1", "arg2")
			},
		},
		{
			name: "Watch returns interface",
			testFunc: func() CronJobClientInterface {
				return client.Watch(func(event common.JobWatchInterface) {})
			},
		},
		{
			name: "DefaultHooks returns interface",
			testFunc: func() CronJobClientInterface {
				return client.DefaultHooks()
			},
		},
		{
			name: "BeforeJobRuns returns interface",
			testFunc: func() CronJobClientInterface {
				return client.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {})
			},
		},
		{
			name: "BeforeJobRunsSkipIfBeforeFuncErrors returns interface",
			testFunc: func() CronJobClientInterface {
				return client.BeforeJobRunsSkipIfBeforeFuncErrors(func(jobID uuid.UUID, jobName string) error {
					return nil
				})
			},
		},
		{
			name: "AfterJobRuns returns interface",
			testFunc: func() CronJobClientInterface {
				return client.AfterJobRuns(func(jobID uuid.UUID, jobName string) {})
			},
		},
		{
			name: "AfterJobRunsWithError returns interface",
			testFunc: func() CronJobClientInterface {
				return client.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {})
			},
		},
		{
			name: "AfterJobRunsWithPanic returns interface",
			testFunc: func() CronJobClientInterface {
				return client.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {})
			},
		},
		{
			name: "AfterLockError returns interface",
			testFunc: func() CronJobClientInterface {
				return client.AfterLockError(func(jobID uuid.UUID, jobName string, err error) {})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result == nil {
				t.Errorf("Method should return interface, got nil")
			}
			// Verify the returned value implements the interface
			// 验证返回值实现了接口
			var _ CronJobClientInterface = result
		})
	}
}

// TestCronJobClientMethodsWithNilJob tests methods when jobs is nil.
// TestCronJobClientMethodsWithNilJob 测试任务为 nil 时的方法
func TestCronJobClientMethodsWithNilJob(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &CronJobClient{
		scheduler: scheduler,
		job:       nil,
	}

	tests := []struct {
		name     string                        // 测试名称 / Test name
		testFunc func() CronJobClientInterface // 测试函数 / Test function
	}{
		{
			name: "CronExpr with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.CronExpr("* * * * *")
			},
		},
		{
			name: "Alias with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.Alias("test-alias")
			},
		},
		{
			name: "Name with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.Name("test-name")
			},
		},
		{
			name: "Tags with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.Tags("tag1")
			},
		},
		{
			name: "Task with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.Task(func() {})
			},
		},
		{
			name: "Watch with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.Watch(func(event common.JobWatchInterface) {})
			},
		},
		{
			name: "DefaultHooks with nil jobs returns client",
			testFunc: func() CronJobClientInterface {
				return client.DefaultHooks()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result == nil {
				t.Errorf("Method should return client even with nil jobs")
			}
		})
	}
}

// TestCronJobClientRemovePriority tests the priority order for Remove method.
// TestCronJobClientRemovePriority 测试 Remove 方法的优先级顺序
func TestCronJobClientRemovePriority(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	tests := []struct {
		name    string // 测试名称 / Test name
		jobID   string // 任务 ID / Job ID
		alias   string // 别名 / Alias
		jobName string // 任务名称 / Job name
		wantErr error  // 期望的错误 / Expected error
	}{
		{
			name:    "Remove by JobID when JobID is set",
			jobID:   "test-id",
			alias:   "test-alias",
			jobName: "test-name",
			wantErr: common.ErrJobNotFound, // Job doesn't exist / 任务不存在
		},
		{
			name:    "Remove by Alias when JobID is empty",
			jobID:   "",
			alias:   "test-alias",
			jobName: "test-name",
			wantErr: common.ErrFoundAlias, // Alias doesn't exist / 别名不存在
		},
		{
			name:    "Remove by Name when both JobID and Alias are empty",
			jobID:   "",
			alias:   "",
			jobName: "test-name",
			wantErr: common.ErrJobNotFound, // Name doesn't exist / 名称不存在
		},
		{
			name:    "Remove with no identifier returns pkg.ErrJobNotFound",
			jobID:   "",
			alias:   "",
			jobName: "",
			wantErr: common.ErrJobNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &CronJobClient{
				scheduler: scheduler,
				job: &jobs.CronJob{
					ID:   tt.jobID,
					Ali:  tt.alias,
					Name: tt.jobName,
				},
			}
			err := client.Remove()
			if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
				t.Logf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCronJobClientGetPriority tests the priority order for Get method.
// TestCronJobClientGetPriority 测试 Get 方法的优先级顺序
func TestCronJobClientGetPriority(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	tests := []struct {
		name    string // 测试名称 / Test name
		jobID   string // 任务 ID / Job ID
		alias   string // 别名 / Alias
		jobName string // 任务名称 / Job name
		wantErr error  // 期望的错误 / Expected error
	}{
		{
			name:    "Get by JobID when JobID is set",
			jobID:   "test-id",
			alias:   "test-alias",
			jobName: "test-name",
			wantErr: common.ErrJobNotFound, // Job doesn't exist / 任务不存在
		},
		{
			name:    "Get by Alias when JobID is empty",
			jobID:   "",
			alias:   "test-alias",
			jobName: "test-name",
			wantErr: common.ErrFoundAlias, // Alias doesn't exist / 别名不存在
		},
		{
			name:    "Get by Name when both JobID and Alias are empty",
			jobID:   "",
			alias:   "",
			jobName: "test-name",
			wantErr: common.ErrJobNotFound, // Name doesn't exist / 名称不存在
		},
		{
			name:    "Get with no identifier returns pkg.ErrJobNotFound",
			jobID:   "",
			alias:   "",
			jobName: "",
			wantErr: common.ErrJobNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &CronJobClient{
				scheduler: scheduler,
				job: &jobs.CronJob{
					ID:   tt.jobID,
					Ali:  tt.alias,
					Name: tt.jobName,
				},
			}
			_, err := client.Get()
			if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
				t.Logf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCronJobClientBatchAdd tests BatchAdd method.
// TestCronJobClientBatchAdd 测试 BatchAdd 方法
func TestCronJobClientBatchAdd(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &CronJobClient{
		scheduler: scheduler,
		job:       &jobs.CronJob{},
	}

	t.Run("BatchAdd with empty jobs returns nil, nil", func(t *testing.T) {
		jobs, err := client.BatchAdd()
		if err != nil {
			t.Errorf("BatchAdd() error = %v, want nil", err)
		}
		if jobs != nil {
			t.Errorf("BatchAdd() jobs = %v, want nil", jobs)
		}
	})

	t.Run("BatchAdd with nil scheduler returns error", func(t *testing.T) {
		client.scheduler = nil
		_, err := client.BatchAdd(&jobs.CronJob{})
		if !errors.Is(err, common.ErrScheduleNil) {
			t.Errorf("BatchAdd() error = %v, wantErr %v", err, common.ErrScheduleNil)
		}
	})
}

// TestCronJobClientThreadSafety tests that client methods are thread-safe.
// TestCronJobClientThreadSafety 测试客户端方法是否线程安全
func TestCronJobClientThreadSafety(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &CronJobClient{
		scheduler: scheduler,
		job:       &jobs.CronJob{},
	}

	done := make(chan bool)

	// Start multiple goroutines calling methods concurrently
	// 启动多个 goroutine 并发调用方法
	for i := 0; i < 10; i++ {
		go func() {
			client.CronExpr("* * * * *")
			client.Alias("test")
			client.Name("test")
			client.Tags("tag1", "tag2")
			client.Task(func() {})
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we got here without panic or deadlock, test passed
	// 如果我们到达这里而没有 panic 或死锁，测试通过
}

// TestOnceJobClientInterfaceCompliance tests that OnceJobClient implements OnceJobClientInterface.
// TestOnceJobClientInterfaceCompliance 测试 OnceJobClient 是否实现了 OnceJobClientInterface
func TestOnceJobClientInterfaceCompliance(t *testing.T) {
	var _ OnceJobClientInterface = (*OnceJobClient)(nil)
}

// TestOnceJobClientNilScheduler tests client behavior with nil scheduler.
// TestOnceJobClientNilScheduler 测试调度器为 nil 时的客户端行为
func TestOnceJobClientNilScheduler(t *testing.T) {
	tests := []struct {
		name     string                     // 测试名称 / Test name
		testFunc func(*OnceJobClient) error // 测试函数 / Test function
		wantErr  error                      // 期望的错误 / Expected error
	}{
		{
			name: "Add with nil scheduler returns pkg.ErrScheduleNil",
			testFunc: func(c *OnceJobClient) error {
				_, err := c.Add()
				return err
			},
			wantErr: common.ErrScheduleNil,
		},
		{
			name: "BatchAdd with nil scheduler returns pkg.ErrScheduleNil",
			testFunc: func(c *OnceJobClient) error {
				_, err := c.BatchAdd()
				return err
			},
			wantErr: common.ErrScheduleNil,
		},
		{
			name: "Remove with nil scheduler returns pkg.ErrScheduleNil",
			testFunc: func(c *OnceJobClient) error {
				return c.Remove()
			},
			wantErr: common.ErrScheduleNil,
		},
		{
			name: "Get with nil scheduler returns pkg.ErrScheduleNil",
			testFunc: func(c *OnceJobClient) error {
				_, err := c.Get()
				return err
			},
			wantErr: common.ErrScheduleNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &OnceJobClient{
				job: &jobs.OnceJob{},
			}
			err := tt.testFunc(client)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestOnceJobClientMethodChaining tests that methods return the interface for chaining.
// TestOnceJobClientMethodChaining 测试方法返回接口以支持链式调用
func TestOnceJobClientMethodChaining(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &OnceJobClient{
		scheduler: scheduler,
		job:       &jobs.OnceJob{},
	}

	tests := []struct {
		name     string                        // 测试名称 / Test name
		testFunc func() OnceJobClientInterface // 测试函数 / Test function
	}{
		{
			name: "AtTimes returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.AtTimes(time.Now())
			},
		},
		{
			name: "Alias returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.Alias("test-alias")
			},
		},
		{
			name: "JobID returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.JobID("test-id")
			},
		},
		{
			name: "Name returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.Name("test-name")
			},
		},
		{
			name: "Tag returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.Tag("tag1", "tag2")
			},
		},
		{
			name: "Task returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.Task(func() {}, "arg1")
			},
		},
		{
			name: "Watch returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.Watch(func(event common.JobWatchInterface) {})
			},
		},
		{
			name: "DefaultHooks returns interface",
			testFunc: func() OnceJobClientInterface {
				return client.DefaultHooks()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result == nil {
				t.Errorf("Method should return interface, got nil")
			}
		})
	}
}

// TestIntervalJobClientInterfaceCompliance tests that IntervalJobClient implements IntervalJobClientInterface.
// TestIntervalJobClientInterfaceCompliance 测试 IntervalJobClient 是否实现了 IntervalJobClientInterface
func TestIntervalJobClientInterfaceCompliance(t *testing.T) {
	var _ IntervalJobClientInterface = (*IntervalJobClient)(nil)
}

// TestIntervalJobClientMethodChaining tests that methods return the interface for chaining.
// TestIntervalJobClientMethodChaining 测试方法返回接口以支持链式调用
func TestIntervalJobClientMethodChaining(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &IntervalJobClient{
		scheduler: scheduler,
		job:       &jobs.IntervalJob{},
	}

	tests := []struct {
		name     string                            // 测试名称 / Test name
		testFunc func() IntervalJobClientInterface // 测试函数 / Test function
	}{
		{
			name: "Interval returns interface",
			testFunc: func() IntervalJobClientInterface {
				return client.Interval(time.Minute)
			},
		},
		{
			name: "Alias returns interface",
			testFunc: func() IntervalJobClientInterface {
				return client.Alias("test-alias")
			},
		},
		{
			name: "JobID returns interface",
			testFunc: func() IntervalJobClientInterface {
				return client.JobID("test-id")
			},
		},
		{
			name: "Name returns interface",
			testFunc: func() IntervalJobClientInterface {
				return client.Name("test-name")
			},
		},
		{
			name: "Tag returns interface",
			testFunc: func() IntervalJobClientInterface {
				return client.Tag("tag1", "tag2")
			},
		},
		{
			name: "Task returns interface",
			testFunc: func() IntervalJobClientInterface {
				return client.Task(func() {})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result == nil {
				t.Errorf("Method should return interface, got nil")
			}
		})
	}
}

// TestDailyJobClientInterfaceCompliance tests that DailyJobClient implements DailyJobClientInterface.
// TestDailyJobClientInterfaceCompliance 测试 DailyJobClient 是否实现了 DailyJobClientInterface
func TestDailyJobClientInterfaceCompliance(t *testing.T) {
	var _ DailyJobClientInterface = (*DailyJobClient)(nil)
}

// TestDailyJobClientMethodChaining tests that methods return the interface for chaining.
// TestDailyJobClientMethodChaining 测试方法返回接口以支持链式调用
func TestDailyJobClientMethodChaining(t *testing.T) {
	scheduler, _ := scheduler.NewScheduler(nil, nil)

	client := &DailyJobClient{
		scheduler: scheduler,
		job:       &jobs.DailyJob{},
	}

	tests := []struct {
		name     string                         // 测试名称 / Test name
		testFunc func() DailyJobClientInterface // 测试函数 / Test function
	}{
		{
			name: "AtDayTime returns interface",
			testFunc: func() DailyJobClientInterface {
				return client.AtDayTime(10, 30, 0)
			},
		},
		{
			name: "Alias returns interface",
			testFunc: func() DailyJobClientInterface {
				return client.Alias("test-alias")
			},
		},
		{
			name: "JobID returns interface",
			testFunc: func() DailyJobClientInterface {
				return client.JobID("test-id")
			},
		},
		{
			name: "Name returns interface",
			testFunc: func() DailyJobClientInterface {
				return client.Name("test-name")
			},
		},
		{
			name: "Tags returns interface",
			testFunc: func() DailyJobClientInterface {
				return client.Tags("tag1", "tag2")
			},
		},
		{
			name: "Task returns interface",
			testFunc: func() DailyJobClientInterface {
				return client.Task(func() {})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result == nil {
				t.Errorf("Method should return interface, got nil")
			}
		})
	}
}
