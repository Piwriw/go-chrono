package scheduler

import (
	"errors"
	"testing"

	errors2 "github.com/piwriw/go-chrono/common"
)

// TestErrorDefinitions tests that all error variables are properly defined.
// TestErrorDefinitions 测试所有错误变量是否正确定义
func TestErrorDefinitions(t *testing.T) {
	tests := []struct {
		name string // 测试名称 / Test name
		err  error  // 错误变量 / Error variable
		want string // 期望的错误消息 / Expected error message
	}{
		{
			name: "pkg.ErrInvalidJob message",
			err:  errors2.ErrInvalidJob,
			want: "chrono:invalid jobs",
		},
		{
			name: "pkg.ErrAtTimeDaysNil message",
			err:  errors2.ErrAtTimeDaysNil,
			want: "chrono:at time must have at least one day",
		},
		{
			name: "pkg.ErrTaskFuncNil message",
			err:  errors2.ErrTaskFuncNil,
			want: "chrono:task function cannot be nil",
		},
		{
			name: "pkg.ErrTaskFailed message",
			err:  errors2.ErrTaskFailed,
			want: "chrono:task is failed",
		},
		{
			name: "pkg.ErrFoundAlias message",
			err:  errors2.ErrFoundAlias,
			want: "chrono:can not found  by alias",
		},
		{
			name: "pkg.ErrDisEnableAlias message",
			err:  errors2.ErrDisEnableAlias,
			want: "chrono:alias is disable",
		},
		{
			name: "pkg.ErrScheduleNil message",
			err:  errors2.ErrScheduleNil,
			want: "chrono:schedule cannot be nil",
		},
		{
			name: "pkg.ErrOnceJobNil message",
			err:  errors2.ErrOnceJobNil,
			want: "chrono:once jobs cannot be nil",
		},
		{
			name: "pkg.ErrCronJobNil message",
			err:  errors2.ErrCronJobNil,
			want: "chrono:cron jobs cannot be nil",
		},
		{
			name: "pkg.ErrDailyJobNil message",
			err:  errors2.ErrDailyJobNil,
			want: "chrono:daily jobs cannot be nil",
		},
		{
			name: "pkg.ErrIntervalJobNil message",
			err:  errors2.ErrIntervalJobNil,
			want: "chrono:interval jobs cannot be nil",
		},
		{
			name: "pkg.ErrMonthJobNil message",
			err:  errors2.ErrMonthJobNil,
			want: "chrono:monthly jobs cannot be nil",
		},
		{
			name: "pkg.ErrWeeklyJobNil message",
			err:  errors2.ErrWeeklyJobNil,
			want: "chrono:weekly jobs cannot be nil",
		},
		{
			name: "pkg.ErrJobNotFound message",
			err:  errors2.ErrJobNotFound,
			want: "chrono:jobs not found (missing JobID/Alias/Name)",
		},
		{
			name: "pkg.ErrDisEnableLimit message",
			err:  errors2.ErrDisEnableLimit,
			want: "chrono:limit is disable",
		},
		{
			name: "pkg.ErrMoreLimit message",
			err:  errors2.ErrMoreLimit,
			want: "chrono:limit is reached",
		},
		{
			name: "pkg.ErrJobIDNil message",
			err:  errors2.ErrJobIDNil,
			want: "chrono:jobs id cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("Error variable should not be nil")
			}
			if tt.err.Error() != tt.want {
				t.Errorf("Error message = %v, want %v", tt.err.Error(), tt.want)
			}
		})
	}
}

// TestErrorIsNotNil tests that all error variables are not nil.
// TestErrorIsNotNil 测试所有错误变量不为 nil
func TestErrorIsNotNil(t *testing.T) {
	tests := []struct {
		name string // 测试名称 / Test name
		err  error  // 错误变量 / Error variable
	}{
		{name: "pkg.ErrInvalidJob is not nil", err: errors2.ErrInvalidJob},
		{name: "pkg.ErrAtTimeDaysNil is not nil", err: errors2.ErrAtTimeDaysNil},
		{name: "pkg.ErrTaskFuncNil is not nil", err: errors2.ErrTaskFuncNil},
		{name: "pkg.ErrTaskFailed is not nil", err: errors2.ErrTaskFailed},
		{name: "pkg.ErrFoundAlias is not nil", err: errors2.ErrFoundAlias},
		{name: "pkg.ErrDisEnableAlias is not nil", err: errors2.ErrDisEnableAlias},
		{name: "pkg.ErrScheduleNil is not nil", err: errors2.ErrScheduleNil},
		{name: "pkg.ErrOnceJobNil is not nil", err: errors2.ErrOnceJobNil},
		{name: "pkg.ErrCronJobNil is not nil", err: errors2.ErrCronJobNil},
		{name: "pkg.ErrDailyJobNil is not nil", err: errors2.ErrDailyJobNil},
		{name: "pkg.ErrIntervalJobNil is not nil", err: errors2.ErrIntervalJobNil},
		{name: "pkg.ErrMonthJobNil is not nil", err: errors2.ErrMonthJobNil},
		{name: "pkg.ErrWeeklyJobNil is not nil", err: errors2.ErrWeeklyJobNil},
		{name: "pkg.ErrJobNotFound is not nil", err: errors2.ErrJobNotFound},
		{name: "pkg.ErrDisEnableLimit is not nil", err: errors2.ErrDisEnableLimit},
		{name: "pkg.ErrMoreLimit is not nil", err: errors2.ErrMoreLimit},
		{name: "pkg.ErrJobIDNil is not nil", err: errors2.ErrJobIDNil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("Error variable should not be nil")
			}
		})
	}
}

// TestErrorUniqueness tests that all error variables are unique.
// TestErrorUniqueness 测试所有错误变量是否唯一
func TestErrorUniqueness(t *testing.T) {
	errors := []error{
		errors2.ErrInvalidJob,
		errors2.ErrAtTimeDaysNil,
		errors2.ErrTaskFuncNil,
		errors2.ErrTaskFailed,
		errors2.ErrFoundAlias,
		errors2.ErrDisEnableAlias,
		errors2.ErrScheduleNil,
		errors2.ErrOnceJobNil,
		errors2.ErrCronJobNil,
		errors2.ErrDailyJobNil,
		errors2.ErrIntervalJobNil,
		errors2.ErrMonthJobNil,
		errors2.ErrWeeklyJobNil,
		errors2.ErrJobNotFound,
		errors2.ErrDisEnableLimit,
		errors2.ErrMoreLimit,
		errors2.ErrJobIDNil,
	}

	// Check for duplicate error messages
	// 检查是否有重复的错误消息
	messages := make(map[string]int)
	for _, err := range errors {
		msg := err.Error()
		messages[msg]++
		if messages[msg] > 1 {
			t.Errorf("Duplicate error message found: %s", msg)
		}
	}
}

// TestErrorIs tests error comparison using errors.Is.
// TestErrorIs 测试使用 errors.Is 进行错误比较
func TestErrorIs(t *testing.T) {
	tests := []struct {
		name        string // 测试名称 / Test name
		targetError error  // 目标错误 / Target error
		wrappedErr  error  // 包装后的错误 / Wrapped error
		wantMatch   bool   // 是否期望匹配 / Whether match is expected
	}{
		{
			name:        "Matching pkg.ErrInvalidJob",
			targetError: errors2.ErrInvalidJob,
			wrappedErr:  errors2.ErrInvalidJob,
			wantMatch:   true,
		},
		{
			name:        "Not matching different error",
			targetError: errors2.ErrInvalidJob,
			wrappedErr:  errors2.ErrTaskFuncNil,
			wantMatch:   false,
		},
		{
			name:        "Matching wrapped pkg.ErrInvalidJob",
			targetError: errors2.ErrInvalidJob,
			wrappedErr:  errors.New("wrapper: " + errors2.ErrInvalidJob.Error()),
			wantMatch:   false, // errors.Is doesn't match by string, only by identity
		},
		{
			name:        "Matching pkg.ErrTaskFuncNil with itself",
			targetError: errors2.ErrTaskFuncNil,
			wrappedErr:  errors2.ErrTaskFuncNil,
			wantMatch:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch := errors.Is(tt.wrappedErr, tt.targetError)
			if gotMatch != tt.wantMatch {
				t.Errorf("errors.Is() = %v, want %v (target=%v, wrapped=%v)",
					gotMatch, tt.wantMatch, tt.targetError, tt.wrappedErr)
			}
		})
	}
}

// TestErrorWrapping tests error wrapping behavior.
// TestErrorWrapping 测试错误包装行为
func TestErrorWrapping(t *testing.T) {
	tests := []struct {
		name        string // 测试名称 / Test name
		baseError   error  // 基础错误 / Base error
		wrapMessage string // 包装消息 / Wrap message
		wantContain string // 期望包含的字符串 / Expected contained string
	}{
		{
			name:        "Wrap pkg.ErrInvalidJob",
			baseError:   errors2.ErrInvalidJob,
			wrapMessage: "failed to create jobs",
			wantContain: "chrono:invalid jobs",
		},
		{
			name:        "Wrap pkg.ErrTaskFuncNil",
			baseError:   errors2.ErrTaskFuncNil,
			wrapMessage: "validation failed",
			wantContain: "chrono:task function cannot be nil",
		},
		{
			name:        "Wrap pkg.ErrJobNotFound",
			baseError:   errors2.ErrJobNotFound,
			wrapMessage: "lookup failed",
			wantContain: "chrono:jobs not found (missing JobID/Alias/Name)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedErr := errors.New(tt.wrapMessage + ": " + tt.baseError.Error())
			if !errors.Is(wrappedErr, tt.baseError) {
				// This is expected since errors.Is doesn't match by string
				// 这是预期的，因为 errors.Is 不按字符串匹配
				t.Logf("Wrapped error doesn't match base error (expected)")
			}
			if !contains(wrappedErr.Error(), tt.wantContain) {
				t.Errorf("Wrapped error should contain %v", tt.wantContain)
			}
		})
	}
}

// contains is a helper function to check if a string contains a substring.
// contains 是一个辅助函数，用于检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr ||
		contains(s[1:], substr) ||
		s[len(s)-len(substr):] == substr ||
		containsInner(s, substr)))
}

// containsInner checks if substr is contained in s starting from position 1.
// containsInner 检查 substr 是否包含在 s 中，从位置 1 开始
func containsInner(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
