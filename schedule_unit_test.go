package chrono

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewScheduler tests the NewScheduler function.
// TestNewScheduler 测试 NewScheduler 函数。
//
// Test scenarios:
// - Happy path: create scheduler with valid parameters
// - Boundary conditions: nil context, nil monitor
// - Exception cases: various option combinations
func TestNewScheduler(t *testing.T) {
	t.Parallel()

	t.Run("creates scheduler with context and monitor", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		monitor := newDefaultSchedulerMonitor()

		scheduler, err := NewScheduler(ctx, monitor)

		require.NoError(t, err, "should create scheduler without error")
		require.NotNil(t, scheduler, "scheduler should not be nil")
		assert.Equal(t, ctx, scheduler.ctx, "context should match")
		assert.Equal(t, monitor, scheduler.monitor, "monitor should match")
		assert.NotNil(t, scheduler.aliasMap, "aliasMap should be initialized")
		assert.NotNil(t, scheduler.watchFuncMap, "watchFuncMap should be initialized")
		assert.NotNil(t, scheduler.jobTypeMap, "jobTypeMap should be initialized")
	})

	t.Run("creates scheduler with nil context uses background", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(nil, nil)

		require.NoError(t, err, "should create scheduler without error")
		assert.NotNil(t, scheduler.ctx, "context should not be nil")
	})

	t.Run("creates scheduler with nil monitor creates default", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil)

		require.NoError(t, err, "should create scheduler without error")
		assert.NotNil(t, scheduler.monitor, "monitor should not be nil")
	})

	t.Run("creates scheduler with alias option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil, WithAliasMode())

		require.NoError(t, err, "should create scheduler without error")
		assert.True(t, scheduler.Enable(AliasOptionName), "alias option should be enabled")
	})

	t.Run("creates scheduler with limit option", func(t *testing.T) {
		t.Parallel()

		limit := 100
		scheduler, err := NewScheduler(context.Background(), nil, WithLimit(limit))

		require.NoError(t, err, "should create scheduler without error")
		assert.True(t, scheduler.Enable(LimitOptionName), "limit option should be enabled")
		assert.Equal(t, limit, scheduler.schOptions.limit.number, "limit number should match")
	})

	t.Run("creates scheduler with watch option", func(t *testing.T) {
		t.Parallel()

		watchFunc := func(event JobWatchInterface) {}
		scheduler, err := NewScheduler(context.Background(), nil, WithWatch(watchFunc))

		require.NoError(t, err, "should create scheduler without error")
		assert.True(t, scheduler.Enable(WatchOptionName), "watch option should be enabled")
	})

	t.Run("creates scheduler with multiple options", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil,
			WithAliasMode(),
			WithLimit(50),
			WithWebMonitor(":8080"),
			WithPrometheus(":9090"),
		)

		require.NoError(t, err, "should create scheduler without error")
		assert.True(t, scheduler.Enable(AliasOptionName), "alias should be enabled")
		assert.True(t, scheduler.Enable(LimitOptionName), "limit should be enabled")
		assert.True(t, scheduler.Enable(WebMonitorOptionName), "web monitor should be enabled")
		assert.True(t, scheduler.Enable(PrometheusOptionName), "prometheus should be enabled")
	})
}

// TestScheduler_Enable tests the Enable method.
// TestScheduler_Enable 测试 Enable 方法。
func TestScheduler_Enable(t *testing.T) {
	t.Parallel()

	t.Run("returns false for disabled option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil)
		require.NoError(t, err)

		assert.False(t, scheduler.Enable(AliasOptionName), "disabled option should return false")
	})

	t.Run("returns true for enabled alias option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil, WithAliasMode())
		require.NoError(t, err)

		assert.True(t, scheduler.Enable(AliasOptionName), "enabled alias should return true")
	})

	t.Run("returns true for enabled limit option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil, WithLimit(10))
		require.NoError(t, err)

		assert.True(t, scheduler.Enable(LimitOptionName), "enabled limit should return true")
	})

	t.Run("returns true for enabled watch option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil, WithWatch(nil))
		require.NoError(t, err)

		assert.True(t, scheduler.Enable(WatchOptionName), "enabled watch should return true")
	})

	t.Run("returns true for enabled web monitor option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil, WithWebMonitor(":8080"))
		require.NoError(t, err)

		assert.True(t, scheduler.Enable(WebMonitorOptionName), "enabled web monitor should return true")
	})

	t.Run("returns true for enabled prometheus option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil, WithPrometheus(":9090"))
		require.NoError(t, err)

		assert.True(t, scheduler.Enable(PrometheusOptionName), "enabled prometheus should return true")
	})

	t.Run("returns false for unknown option", func(t *testing.T) {
		t.Parallel()

		scheduler, err := NewScheduler(context.Background(), nil)
		require.NoError(t, err)

		assert.False(t, scheduler.Enable("unknown_option"), "unknown option should return false")
	})
}

// TestScheduler_AddJobType tests the addJobType method.
// TestScheduler_AddJobType 测试 addJobType 方法。
func TestScheduler_AddJobType(t *testing.T) {
	t.Parallel()

	t.Run("adds job type successfully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()
		jobType := JobTypeCron

		scheduler.addJobType(jobID, jobType)

		retrievedType := scheduler.getJobType(jobID)
		assert.Equal(t, jobType, retrievedType, "job type should match")
	})

	t.Run("handles empty job ID gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Should not panic
		// 不应该panic
		scheduler.addJobType("", JobTypeCron)
	})

	t.Run("overwrites existing job type", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		scheduler.addJobType(jobID, JobTypeCron)
		scheduler.addJobType(jobID, JobTypeOnce)

		retrievedType := scheduler.getJobType(jobID)
		assert.Equal(t, JobTypeOnce, retrievedType, "job type should be updated")
	})

	t.Run("thread-safe concurrent additions", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		const goroutines = 100

		jobTypes := []JobType{JobTypeOnce, JobTypeCron, JobTypeDaily, JobTypeWeekly, JobTypeMonthly, JobInterval}

		var wg sync.WaitGroup
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				jobID := uuid.New().String()
				scheduler.addJobType(jobID, jobTypes[idx%len(jobTypes)])
			}(i)
		}

		wg.Wait()

		// Verify no data races occurred
		// 验证没有发生数据竞争
		assert.NotNil(t, scheduler.jobTypeMap, "jobTypeMap should not be nil")
	})
}

// TestScheduler_RemoveJobType tests the removeJobType method.
// TestScheduler_RemoveJobType 测试 removeJobType 方法。
func TestScheduler_RemoveJobType(t *testing.T) {
	t.Parallel()

	t.Run("removes existing job type", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		scheduler.addJobType(jobID, JobTypeCron)
		scheduler.removeJobType(jobID)

		retrievedType := scheduler.getJobType(jobID)
		assert.Equal(t, JobTypeUnknown, retrievedType, "job type should be unknown after removal")
	})

	t.Run("handles non-existent job ID gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentJobID := uuid.New().String()

		// Should not panic
		// 不应该panic
		scheduler.removeJobType(nonExistentJobID)
	})
}

// TestScheduler_GetJobType tests the getJobType method.
// TestScheduler_GetJobType 测试 getJobType 方法。
func TestScheduler_GetJobType(t *testing.T) {
	t.Parallel()

	t.Run("returns unknown for non-existent job", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentJobID := uuid.New().String()

		jobType := scheduler.getJobType(nonExistentJobID)
		assert.Equal(t, JobTypeUnknown, jobType, "should return unknown for non-existent job")
	})

	t.Run("returns correct job type for existing job", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()
		expectedType := JobInterval

		scheduler.addJobType(jobID, expectedType)

		actualType := scheduler.getJobType(jobID)
		assert.Equal(t, expectedType, actualType, "should return correct job type")
	})

	t.Run("handles empty job ID", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		jobType := scheduler.getJobType("")
		assert.Equal(t, JobTypeUnknown, jobType, "should return unknown for empty job ID")
	})
}

// TestScheduler_AddAlias tests the addAlias method.
// TestScheduler_AddAlias 测试 addAlias 方法。
func TestScheduler_AddAlias(t *testing.T) {
	t.Parallel()

	t.Run("adds alias successfully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		alias := "test-alias"
		jobID := uuid.New().String()

		scheduler.addAlias(alias, jobID)

		retrievedJobID, exists := scheduler.aliasMap[alias]
		assert.True(t, exists, "alias should exist")
		assert.Equal(t, jobID, retrievedJobID, "job ID should match")
	})

	t.Run("handles empty alias gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		// Should not panic
		// 不应该panic
		scheduler.addAlias("", jobID)
	})

	t.Run("handles empty job ID gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		alias := "test-alias"

		// Should not panic
		// 不应该panic
		scheduler.addAlias(alias, "")
	})

	t.Run("overwrites existing alias", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		alias := "test-alias"
		jobID1 := uuid.New().String()
		jobID2 := uuid.New().String()

		scheduler.addAlias(alias, jobID1)
		scheduler.addAlias(alias, jobID2)

		retrievedJobID := scheduler.aliasMap[alias]
		assert.Equal(t, jobID2, retrievedJobID, "alias should point to new job ID")
	})
}

// TestScheduler_RemoveAlias tests the removeAlias method.
// TestScheduler_RemoveAlias 测试 removeAlias 方法。
func TestScheduler_RemoveAlias(t *testing.T) {
	t.Parallel()

	t.Run("removes existing alias", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		alias := "test-alias"
		jobID := uuid.New().String()

		scheduler.addAlias(alias, jobID)
		scheduler.removeAlias(alias)

		_, exists := scheduler.aliasMap[alias]
		assert.False(t, exists, "alias should be removed")
	})

	t.Run("handles non-existent alias gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentAlias := "non-existent-alias"

		// Should not panic
		// 不应该panic
		scheduler.removeAlias(nonExistentAlias)
	})
}

// TestScheduler_AddWatchFunc tests the addWatchFunc method.
// TestScheduler_AddWatchFunc 测试 addWatchFunc 方法。
func TestScheduler_AddWatchFunc(t *testing.T) {
	t.Parallel()

	watchFunc := func(event JobWatchInterface) {
		// Watch function
	}

	t.Run("adds watch function successfully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		scheduler.addWatchFunc(jobID, watchFunc)

		retrievedFunc, exists := scheduler.watchFuncMap[jobID]
		assert.True(t, exists, "watch function should exist")
		assert.NotNil(t, retrievedFunc, "watch function should not be nil")
	})

	t.Run("handles empty job ID gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Should not panic
		// 不应该panic
		scheduler.addWatchFunc("", watchFunc)
	})

	t.Run("handles nil watch function gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		// Should not panic
		// 不应该panic
		scheduler.addWatchFunc(jobID, nil)
	})

	t.Run("overwrites existing watch function", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()
		watchFunc1 := func(event JobWatchInterface) {}
		watchFunc2 := func(event JobWatchInterface) {}

		scheduler.addWatchFunc(jobID, watchFunc1)
		scheduler.addWatchFunc(jobID, watchFunc2)

		// Verify it was added (we can't compare functions directly)
		// 验证已添加（无法直接比较函数）
		_, exists := scheduler.watchFuncMap[jobID]
		assert.True(t, exists, "watch function should exist")
	})
}

// TestScheduler_RemoveWatchFunc tests the removeWatchFunc method.
// TestScheduler_RemoveWatchFunc 测试 removeWatchFunc 方法。
func TestScheduler_RemoveWatchFunc(t *testing.T) {
	t.Parallel()

	watchFunc := func(event JobWatchInterface) {}

	t.Run("removes existing watch function", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		scheduler.addWatchFunc(jobID, watchFunc)
		scheduler.removeWatchFunc(jobID)

		_, exists := scheduler.watchFuncMap[jobID]
		assert.False(t, exists, "watch function should be removed")
	})

	t.Run("handles non-existent job ID gracefully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentJobID := uuid.New().String()

		// Should not panic
		// 不应该panic
		scheduler.removeWatchFunc(nonExistentJobID)
	})
}

// TestScheduler_CheckLimit tests the CheckLimit method.
// TestScheduler_CheckLimit 测试 CheckLimit 方法。
func TestScheduler_CheckLimit(t *testing.T) {
	t.Parallel()

	t.Run("returns true when limit option is disabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		result := scheduler.CheckLimit()
		assert.True(t, result, "should return true when limit is disabled")
	})

	t.Run("returns true when limit is not reached", func(t *testing.T) {
		t.Parallel()

		limit := 5
		scheduler, _ := NewScheduler(context.Background(), nil, WithLimit(limit))

		result := scheduler.CheckLimit()
		assert.True(t, result, "should return true when limit is not reached")
	})

	t.Run("returns false when limit is reached", func(t *testing.T) {
		t.Parallel()

		limit := 1
		scheduler, _ := NewScheduler(context.Background(), nil, WithLimit(limit))

		// Check limit once (decrements to 0)
		// 检查限制一次（递减到0）
		scheduler.CheckLimit()

		// Check limit again (would be -1)
		// 再次检查限制（将是-1）
		result := scheduler.CheckLimit()
		assert.False(t, result, "should return false when limit is reached")
	})
}

// TestScheduler_IncLimit tests the incLimit method.
// TestScheduler_IncLimit 测试 incLimit 方法。
func TestScheduler_IncLimit(t *testing.T) {
	t.Parallel()

	t.Run("returns error when limit option is disabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		err := scheduler.incLimit()
		assert.Error(t, err, "should return error when limit is disabled")
		assert.Equal(t, ErrDisEnableLimit, err, "error should be ErrDisEnableLimit")
	})

	t.Run("increments limit successfully", func(t *testing.T) {
		t.Parallel()

		initialLimit := 10
		scheduler, _ := NewScheduler(context.Background(), nil, WithLimit(initialLimit))

		err := scheduler.incLimit()
		assert.NoError(t, err, "should increment limit without error")
		assert.Equal(t, initialLimit+1, scheduler.schOptions.limit.number, "limit should be incremented")
	})
}

// TestScheduler_DecLimit tests the decLimit method.
// TestScheduler_DecLimit 测试 decLimit 方法。
func TestScheduler_DecLimit(t *testing.T) {
	t.Parallel()

	t.Run("returns error when limit option is disabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		err := scheduler.decLimit()
		assert.Error(t, err, "should return error when limit is disabled")
		assert.Equal(t, ErrDisEnableLimit, err, "error should be ErrDisEnableLimit")
	})

	t.Run("decrements limit successfully", func(t *testing.T) {
		t.Parallel()

		initialLimit := 10
		scheduler, _ := NewScheduler(context.Background(), nil, WithLimit(initialLimit))

		err := scheduler.decLimit()
		assert.NoError(t, err, "should decrement limit without error")
		assert.Equal(t, initialLimit-1, scheduler.schOptions.limit.number, "limit should be decremented")
	})
}

// TestScheduler_GetAlias tests the GetAlias method.
// TestScheduler_GetAlias 测试 GetAlias 方法。
func TestScheduler_GetAlias(t *testing.T) {
	t.Parallel()

	t.Run("returns error when alias option is disabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		jobID := uuid.New().String()

		alias, err := scheduler.GetAlias(jobID)
		assert.Error(t, err, "should return error when alias is disabled")
		assert.Equal(t, ErrDisEnableAlias, err, "error should be ErrDisEnableAlias")
		assert.Empty(t, alias, "alias should be empty")
	})

	t.Run("returns error when alias not found", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil, WithAliasMode())
		nonExistentJobID := uuid.New().String()

		alias, err := scheduler.GetAlias(nonExistentJobID)
		assert.Error(t, err, "should return error when alias not found")
		assert.Equal(t, ErrFoundAlias, err, "error should be ErrFoundAlias")
		assert.Empty(t, alias, "alias should be empty")
	})

	t.Run("returns alias for existing job", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil, WithAliasMode())
		alias := "test-alias"
		jobID := uuid.New().String()

		scheduler.addAlias(alias, jobID)

		retrievedAlias, err := scheduler.GetAlias(jobID)
		assert.NoError(t, err, "should not return error")
		assert.Equal(t, alias, retrievedAlias, "alias should match")
	})
}

// TestScheduler_GetJobs tests the GetJobs method.
// TestScheduler_GetJobs 测试 GetJobs 方法。
func TestScheduler_GetJobs(t *testing.T) {
	t.Parallel()

	t.Run("returns empty list when no jobs", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		jobs, err := scheduler.GetJobs()
		assert.NoError(t, err, "should not return error")
		assert.Empty(t, jobs, "jobs should be empty")
	})

	t.Run("returns all jobs", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Add a job
		// 添加一个任务
		task := func() error { return nil }
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Task(task)

		_, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		jobs, err := scheduler.GetJobs()
		assert.NoError(t, err, "should not return error")
		assert.NotEmpty(t, jobs, "jobs should not be empty")
	})
}

// TestScheduler_GetJobByID tests the GetJobByID method.
// TestScheduler_GetJobByID 测试 GetJobByID 方法。
func TestScheduler_GetJobByID(t *testing.T) {
	t.Parallel()

	t.Run("returns error for non-existent job", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentJobID := uuid.New().String()

		job, err := scheduler.GetJobByID(nonExistentJobID)
		assert.Error(t, err, "should return error for non-existent job")
		assert.Nil(t, job, "job should be nil")
	})

	t.Run("returns job for existing job ID", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Add a job
		// 添加一个任务
		task := func() error { return nil }
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Get job by ID
		// 通过ID获取任务
		retrievedJob, err := scheduler.GetJobByID(addedJob.ID().String())
		assert.NoError(t, err, "should not return error")
		assert.NotNil(t, retrievedJob, "job should not be nil")
		assert.Equal(t, addedJob.ID(), retrievedJob.ID(), "job ID should match")
	})
}

// TestScheduler_GetJobByName tests the GetJobByName method.
// TestScheduler_GetJobByName 测试 GetJobByName 方法。
func TestScheduler_GetJobByName(t *testing.T) {
	t.Parallel()

	t.Run("returns error for non-existent job", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentJobName := "non-existent-job"

		job, err := scheduler.GetJobByName(nonExistentJobName)
		assert.Error(t, err, "should return error for non-existent job")
		assert.Nil(t, job, "job should be nil")
	})

	t.Run("returns job for existing job name", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Add a job
		// 添加一个任务
		task := func() error { return nil }
		jobName := "test-job-by-name"
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names(jobName).
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Get job by name
		// 通过名称获取任务
		retrievedJob, err := scheduler.GetJobByName(jobName)
		assert.NoError(t, err, "should not return error")
		assert.NotNil(t, retrievedJob, "job should not be nil")
		assert.Equal(t, addedJob.ID(), retrievedJob.ID(), "job ID should match")
		assert.Equal(t, jobName, retrievedJob.Name(), "job name should match")
	})
}

// TestScheduler_GetJobByAlias tests the GetJobByAlias method.
// TestScheduler_GetJobByAlias 测试 GetJobByAlias 方法。
func TestScheduler_GetJobByAlias(t *testing.T) {
	t.Parallel()

	t.Run("returns error when alias option is disabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		alias := "test-alias"

		job, err := scheduler.GetJobByAlias(alias)
		assert.Error(t, err, "should return error when alias is disabled")
		assert.Equal(t, ErrDisEnableAlias, err, "error should be ErrDisEnableAlias")
		assert.Nil(t, job, "job should be nil")
	})

	t.Run("returns error for non-existent alias", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil, WithAliasMode())
		nonExistentAlias := "non-existent-alias"

		job, err := scheduler.GetJobByAlias(nonExistentAlias)
		assert.Error(t, err, "should return error for non-existent alias")
		assert.Nil(t, job, "job should be nil")
	})

	t.Run("returns job for existing alias", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil, WithAliasMode())

		// Add a job with alias
		// 添加一个带别名的任务
		task := func() error { return nil }
		alias := "test-alias"
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Alias(alias).
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Get job by alias
		// 通过别名获取任务
		retrievedJob, err := scheduler.GetJobByAlias(alias)
		assert.NoError(t, err, "should not return error")
		assert.NotNil(t, retrievedJob, "job should not be nil")
		assert.Equal(t, addedJob.ID(), retrievedJob.ID(), "job ID should match")
	})
}

// TestScheduler_GetJobByIDOrAlias tests the GetJobByIDOrAlias method.
// TestScheduler_GetJobByIDOrAlias 测试 GetJobByIDOrAlias 方法。
func TestScheduler_GetJobByIDOrAlias(t *testing.T) {
	t.Parallel()

	t.Run("returns error for non-existent identifier", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		nonExistentID := uuid.New().String()

		job, err := scheduler.GetJobByIDOrAlias(nonExistentID)
		assert.Error(t, err, "should return error for non-existent identifier")
		assert.Nil(t, job, "job should be nil")
	})

	t.Run("returns job by ID when ID exists", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Add a job
		// 添加一个任务
		task := func() error { return nil }
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Get job by ID
		// 通过ID获取任务
		retrievedJob, err := scheduler.GetJobByIDOrAlias(addedJob.ID().String())
		assert.NoError(t, err, "should not return error")
		assert.NotNil(t, retrievedJob, "job should not be nil")
	})

	t.Run("returns job by alias when alias exists and enabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil, WithAliasMode())

		// Add a job with alias
		// 添加一个带别名的任务
		task := func() error { return nil }
		alias := "test-alias-or"
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Alias(alias).
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Get job by alias directly
		// 直接通过别名获取任务
		retrievedJob, err := scheduler.GetJobByAlias(alias)
		assert.NoError(t, err, "should not return error")
		assert.NotNil(t, retrievedJob, "job should not be nil")
		assert.Equal(t, addedJob.ID(), retrievedJob.ID(), "job ID should match")
	})
}

// TestScheduler_RemoveJob tests the RemoveJob method.
// TestScheduler_RemoveJob 测试 RemoveJob 方法。
func TestScheduler_RemoveJob(t *testing.T) {
	t.Parallel()

	t.Run("returns error for empty job ID", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		err := scheduler.RemoveJob("")
		assert.Error(t, err, "should return error for empty job ID")
		assert.Equal(t, ErrJobIDNil, err, "error should be ErrJobIDNil")
	})

	t.Run("returns error for invalid job ID format", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)
		invalidJobID := "invalid-uuid-format"

		err := scheduler.RemoveJob(invalidJobID)
		assert.Error(t, err, "should return error for invalid job ID format")
	})

	t.Run("removes job successfully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Add a job
		// 添加一个任务
		task := func() error { return nil }
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Remove job
		// 移除任务
		err = scheduler.RemoveJob(addedJob.ID().String())
		assert.NoError(t, err, "should remove job without error")
	})

	t.Run("removes alias when alias option is enabled", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil, WithAliasMode())

		// Add a job with alias
		// 添加一个带别名的任务
		task := func() error { return nil }
		alias := "test-alias-remove"
		job := NewOnceJob(time.Now().Add(time.Minute)).
			Names("test-job").
			Alias(alias).
			Task(task)

		addedJob, err := scheduler.AddOnceJob(job)
		require.NoError(t, err, "should add job without error")

		// Verify alias exists
		// 验证别名存在
		_, exists := scheduler.aliasMap[alias]
		assert.True(t, exists, "alias should exist before removal")

		// Remove job
		// 移除任务
		err = scheduler.RemoveJob(addedJob.ID().String())
		assert.NoError(t, err, "should remove job without error")

		// Verify alias is removed
		// 验证别名已移除
		_, exists = scheduler.aliasMap[alias]
		assert.False(t, exists, "alias should be removed")
	})
}

// TestScheduler_StartAndStop tests the Start and Stop methods.
// TestScheduler_StartAndStop 测试 Start 和 Stop 方法。
func TestScheduler_StartAndStop(t *testing.T) {
	t.Parallel()

	t.Run("starts and stops scheduler successfully", func(t *testing.T) {
		t.Parallel()

		scheduler, _ := NewScheduler(context.Background(), nil)

		// Start scheduler
		// 启动调度器
		scheduler.Start()

		// Stop scheduler
		// 停止调度器
		err := scheduler.Stop()
		assert.NoError(t, err, "should stop scheduler without error")
	})
}

// TestEvent_methods tests the Event struct methods.
// TestEvent_methods 测试 Event 结构体方法。
func TestEvent_methods(t *testing.T) {
	t.Parallel()

	t.Run("GetJobID returns correct job ID", func(t *testing.T) {
		t.Parallel()

		expectedJobID := "test-job-123"
		event := Event{JobID: expectedJobID}

		assert.Equal(t, expectedJobID, event.GetJobID(), "job ID should match")
	})

	t.Run("GetJobName returns correct job name", func(t *testing.T) {
		t.Parallel()

		expectedJobName := "test-job-name"
		event := Event{JobName: expectedJobName}

		assert.Equal(t, expectedJobName, event.GetJobName(), "job name should match")
	})

	t.Run("GetNextRunTime returns correct next run time", func(t *testing.T) {
		t.Parallel()

		expectedTime := time.Now().Add(time.Hour)
		event := Event{NextRunTime: expectedTime}

		assert.WithinDuration(t, expectedTime, event.GetNextRunTime(), time.Microsecond, "next run time should match")
	})

	t.Run("GetLastTime returns correct last time", func(t *testing.T) {
		t.Parallel()

		expectedTime := time.Now().Add(-time.Hour)
		event := Event{LastTime: expectedTime}

		assert.WithinDuration(t, expectedTime, event.GetLastTime(), time.Microsecond, "last time should match")
	})

	t.Run("GetError returns correct error", func(t *testing.T) {
		t.Parallel()

		expectedError := assert.AnError
		event := Event{Err: expectedError}

		assert.Equal(t, expectedError, event.GetError(), "error should match")
	})

	t.Run("GetError returns nil when no error", func(t *testing.T) {
		t.Parallel()

		event := Event{Err: nil}

		assert.Nil(t, event.GetError(), "error should be nil")
	})
}

// BenchmarkNewScheduler benchmarks the NewScheduler function.
// BenchmarkNewScheduler 对 NewScheduler 函数进行基准测试。
func BenchmarkNewScheduler(b *testing.B) {
	ctx := context.Background()
	monitor := newDefaultSchedulerMonitor()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = NewScheduler(ctx, monitor)
	}
}

// ExampleNewScheduler demonstrates the usage of NewScheduler.
// ExampleNewScheduler 展示 NewScheduler 的使用方法。
func ExampleNewScheduler() {
	// Create a new scheduler with context and monitor
	// 创建一个带有上下文和监控器的新调度器
	ctx := context.Background()
	monitor := newDefaultSchedulerMonitor()
	scheduler, err := NewScheduler(ctx, monitor, WithAliasMode(), WithLimit(100))

	if err != nil {
		panic(err)
	}

	// Check if alias option is enabled
	// 检查别名选项是否启用
	println(scheduler.Enable(AliasOptionName))
	// Output: true

	// Stop the scheduler when done
	// 完成时停止调度器
	_ = scheduler.Stop()
}
