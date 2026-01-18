package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestDefaultBeforeJobRuns tests the default before jobs runs hook.
// TestDefaultBeforeJobRuns 测试默认的作业启动前钩子
func TestDefaultBeforeJobRuns(t *testing.T) {
	tests := []struct {
		name      string    // 测试名称 / Test name
		jobID     uuid.UUID // 作业 ID / Job ID
		jobName   string    // 作业名称 / Job name
		wantPanic bool      // 是否期望 panic / Whether panic is expected
	}{
		{
			name:      "Valid jobs ID and name",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			jobName:   "test-jobs",
			wantPanic: false,
		},
		{
			name:      "Empty jobs name",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			jobName:   "",
			wantPanic: false,
		},
		{
			name:      "Nil UUID",
			jobID:     uuid.Nil,
			jobName:   "test-jobs",
			wantPanic: false,
		},
		{
			name:      "Special characters in jobs name",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			jobName:   "test-jobs-with-special-chars-中文",
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				defaultBeforeJobRuns(tt.jobID, tt.jobName)
			}()

			if didPanic != tt.wantPanic {
				t.Errorf("defaultBeforeJobRuns() panic = %v, want %v", didPanic, tt.wantPanic)
			}
		})
	}
}

// TestDefaultBeforeJobRunsSkipIfBeforeFuncErrors tests the skippable before jobs runs hook.
// TestDefaultBeforeJobRunsSkipIfBeforeFuncErrors 测试可跳过的作业启动前钩子
func TestDefaultBeforeJobRunsSkipIfBeforeFuncErrors(t *testing.T) {
	tests := []struct {
		name        string    // 测试名称 / Test name
		jobID       uuid.UUID // 作业 ID / Job ID
		jobName     string    // 作业名称 / Job name
		wantErr     bool      // 是否期望错误 / Whether error is expected
		expectedErr error     // 期望的错误类型 / Expected error type
	}{
		{
			name:        "Valid jobs ID and name returns nil",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			jobName:     "test-jobs",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Empty jobs name returns nil",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			jobName:     "",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Nil UUID returns nil",
			jobID:       uuid.Nil,
			jobName:     "test-jobs",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "Long jobs name returns nil",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			jobName:     "this-is-a-very-long-jobs-name-with-lots-of-characters",
			wantErr:     false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := defaultBeforeJobRunsSkipIfBeforeFuncErrors(tt.jobID, tt.jobName)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultBeforeJobRunsSkipIfBeforeFuncErrors() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("defaultBeforeJobRunsSkipIfBeforeFuncErrors() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

// TestDefaultAfterJobRuns tests the default after jobs runs hook.
// TestDefaultAfterJobRuns 测试默认的作业完成后钩子
func TestDefaultAfterJobRuns(t *testing.T) {
	tests := []struct {
		name      string    // 测试名称 / Test name
		jobID     uuid.UUID // 作业 ID / Job ID
		jobName   string    // 作业名称 / Job name
		wantPanic bool      // 是否期望 panic / Whether panic is expected
	}{
		{
			name:      "Successful jobs completion",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			jobName:   "completed-jobs",
			wantPanic: false,
		},
		{
			name:      "Job completion with empty name",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			jobName:   "",
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				defaultAfterJobRuns(tt.jobID, tt.jobName)
			}()

			if didPanic != tt.wantPanic {
				t.Errorf("defaultAfterJobRuns() panic = %v, want %v", didPanic, tt.wantPanic)
			}
		})
	}
}

// TestDefaultAfterJobRunsWithError tests the default after jobs runs with error hook.
// TestDefaultAfterJobRunsWithError 测试默认的作业完成带错误钩子
func TestDefaultAfterJobRunsWithError(t *testing.T) {
	tests := []struct {
		name      string    // 测试名称 / Test name
		jobID     uuid.UUID // 作业 ID / Job ID
		jobName   string    // 作业名称 / Job name
		err       error     // 错误 / Error
		wantPanic bool      // 是否期望 panic / Whether panic is expected
	}{
		{
			name:      "Job with error",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			jobName:   "failed-jobs",
			err:       errors.New("task failed"),
			wantPanic: false,
		},
		{
			name:      "Job with nil error",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			jobName:   "no-error-jobs",
			err:       nil,
			wantPanic: false,
		},
		{
			name:      "Job with wrapped error",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			jobName:   "wrapped-error-jobs",
			err:       errors.New("wrapped: " + "original error"),
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				defaultAfterJobRunsWithError(tt.jobID, tt.jobName, tt.err)
			}()

			if didPanic != tt.wantPanic {
				t.Errorf("defaultAfterJobRunsWithError() panic = %v, want %v", didPanic, tt.wantPanic)
			}
		})
	}
}

// TestDefaultAfterJobRunsWithPanic tests the default after jobs runs with panic hook.
// TestDefaultAfterJobRunsWithPanic 测试默认的作业 panic 钩子
func TestDefaultAfterJobRunsWithPanic(t *testing.T) {
	tests := []struct {
		name        string    // 测试名称 / Test name
		jobID       uuid.UUID // 作业 ID / Job ID
		jobName     string    // 作业名称 / Job name
		recoverData any       // panic 恢复数据 / Panic recover data
		wantPanic   bool      // 是否期望 panic / Whether panic is expected
	}{
		{
			name:        "Job with string panic",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			jobName:     "panic-jobs",
			recoverData: "panic message",
			wantPanic:   false,
		},
		{
			name:        "Job with nil panic data",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			jobName:     "nil-panic-jobs",
			recoverData: nil,
			wantPanic:   false,
		},
		{
			name:        "Job with error panic",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			jobName:     "error-panic-jobs",
			recoverData: errors.New("panic error"),
			wantPanic:   false,
		},
		{
			name:        "Job with struct panic",
			jobID:       uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			jobName:     "struct-panic-jobs",
			recoverData: struct{ Field string }{Field: "value"},
			wantPanic:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				defaultAfterJobRunsWithPanic(tt.jobID, tt.jobName, tt.recoverData)
			}()

			if didPanic != tt.wantPanic {
				t.Errorf("defaultAfterJobRunsWithPanic() panic = %v, want %v", didPanic, tt.wantPanic)
			}
		})
	}
}

// TestDefaultAfterLockError tests the default after lock error hook.
// TestDefaultAfterLockError 测试默认的锁失败钩子
func TestDefaultAfterLockError(t *testing.T) {
	tests := []struct {
		name      string    // 测试名称 / Test name
		jobID     uuid.UUID // 作业 ID / Job ID
		jobName   string    // 作业名称 / Job name
		err       error     // 错误 / Error
		wantPanic bool      // 是否期望 panic / Whether panic is expected
	}{
		{
			name:      "Lock error",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			jobName:   "locked-jobs",
			err:       errors.New("lock acquisition failed"),
			wantPanic: false,
		},
		{
			name:      "Nil lock error",
			jobID:     uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			jobName:   "no-lock-error-jobs",
			err:       nil,
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				defaultAfterLockError(tt.jobID, tt.jobName, tt.err)
			}()

			if didPanic != tt.wantPanic {
				t.Errorf("defaultAfterLockError() panic = %v, want %v", didPanic, tt.wantPanic)
			}
		})
	}
}

// TestEmptyWatchFunc tests the empty watch function.
// TestEmptyWatchFunc 测试空的监控函数
func TestEmptyWatchFunc(t *testing.T) {
	tests := []struct {
		name  string            // 测试名称 / Test name
		event JobWatchInterface // 作业事件 / Job event
	}{
		{
			name:  "Nil event",
			event: nil,
		},
		{
			name:  "Mock event",
			event: &mockJobWatchInterface{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				EmptyWatchFunc(tt.event)
			}()

			if didPanic {
				t.Errorf("EmptyWatchFunc() should not panic")
			}
		})
	}
}

// mockJobWatchInterface is a mock implementation of JobWatchInterface for testing.
// mockJobWatchInterface 是 JobWatchInterface 的模拟实现，用于测试
type mockJobWatchInterface struct{}

func (m *mockJobWatchInterface) GetJobID() string {
	return "00000000-0000-0000-0000-000000000001"
}

func (m *mockJobWatchInterface) GetJobName() string {
	return "mock-jobs"
}

func (m *mockJobWatchInterface) GetStartTime() time.Time {
	return time.Time{}
}

func (m *mockJobWatchInterface) GetEndTime() time.Time {
	return time.Time{}
}

func (m *mockJobWatchInterface) GetStatus() int {
	return 0
}

func (m *mockJobWatchInterface) GetTags() []string {
	return []string{"test-tag"}
}

func (m *mockJobWatchInterface) Error() error {
	return nil
}

func (m *mockJobWatchInterface) GetCurrentEvent() *JobEvent {
	return &JobEvent{}
}
