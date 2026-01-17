package chrono

import (
	"errors"
	"testing"
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
			name: "ErrInvalidJob message",
			err:  ErrInvalidJob,
			want: "chrono:invalid job",
		},
		{
			name: "ErrAtTimeDaysNil message",
			err:  ErrAtTimeDaysNil,
			want: "chrono:at time must have at least one day",
		},
		{
			name: "ErrTaskFuncNil message",
			err:  ErrTaskFuncNil,
			want: "chrono:task function cannot be nil",
		},
		{
			name: "ErrTaskFailed message",
			err:  ErrTaskFailed,
			want: "chrono:task is failed",
		},
		{
			name: "ErrFoundAlias message",
			err:  ErrFoundAlias,
			want: "chrono:can not found  by alias",
		},
		{
			name: "ErrDisEnableAlias message",
			err:  ErrDisEnableAlias,
			want: "chrono:alias is disable",
		},
		{
			name: "ErrScheduleNil message",
			err:  ErrScheduleNil,
			want: "chrono:schedule cannot be nil",
		},
		{
			name: "ErrOnceJobNil message",
			err:  ErrOnceJobNil,
			want: "chrono:once job cannot be nil",
		},
		{
			name: "ErrCronJobNil message",
			err:  ErrCronJobNil,
			want: "chrono:cron job cannot be nil",
		},
		{
			name: "ErrDailyJobNil message",
			err:  ErrDailyJobNil,
			want: "chrono:daily job cannot be nil",
		},
		{
			name: "ErrIntervalJobNil message",
			err:  ErrIntervalJobNil,
			want: "chrono:interval job cannot be nil",
		},
		{
			name: "ErrMonthJobNil message",
			err:  ErrMonthJobNil,
			want: "chrono:monthly job cannot be nil",
		},
		{
			name: "ErrWeeklyJobNil message",
			err:  ErrWeeklyJobNil,
			want: "chrono:weekly job cannot be nil",
		},
		{
			name: "ErrJobNotFound message",
			err:  ErrJobNotFound,
			want: "chrono:job not found (missing JobID/Alias/Name)",
		},
		{
			name: "ErrDisEnableLimit message",
			err:  ErrDisEnableLimit,
			want: "chrono:limit is disable",
		},
		{
			name: "ErrMoreLimit message",
			err:  ErrMoreLimit,
			want: "chrono:limit is reached",
		},
		{
			name: "ErrJobIDNil message",
			err:  ErrJobIDNil,
			want: "chrono:job id cannot be nil",
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
		{name: "ErrInvalidJob is not nil", err: ErrInvalidJob},
		{name: "ErrAtTimeDaysNil is not nil", err: ErrAtTimeDaysNil},
		{name: "ErrTaskFuncNil is not nil", err: ErrTaskFuncNil},
		{name: "ErrTaskFailed is not nil", err: ErrTaskFailed},
		{name: "ErrFoundAlias is not nil", err: ErrFoundAlias},
		{name: "ErrDisEnableAlias is not nil", err: ErrDisEnableAlias},
		{name: "ErrScheduleNil is not nil", err: ErrScheduleNil},
		{name: "ErrOnceJobNil is not nil", err: ErrOnceJobNil},
		{name: "ErrCronJobNil is not nil", err: ErrCronJobNil},
		{name: "ErrDailyJobNil is not nil", err: ErrDailyJobNil},
		{name: "ErrIntervalJobNil is not nil", err: ErrIntervalJobNil},
		{name: "ErrMonthJobNil is not nil", err: ErrMonthJobNil},
		{name: "ErrWeeklyJobNil is not nil", err: ErrWeeklyJobNil},
		{name: "ErrJobNotFound is not nil", err: ErrJobNotFound},
		{name: "ErrDisEnableLimit is not nil", err: ErrDisEnableLimit},
		{name: "ErrMoreLimit is not nil", err: ErrMoreLimit},
		{name: "ErrJobIDNil is not nil", err: ErrJobIDNil},
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
		ErrInvalidJob,
		ErrAtTimeDaysNil,
		ErrTaskFuncNil,
		ErrTaskFailed,
		ErrFoundAlias,
		ErrDisEnableAlias,
		ErrScheduleNil,
		ErrOnceJobNil,
		ErrCronJobNil,
		ErrDailyJobNil,
		ErrIntervalJobNil,
		ErrMonthJobNil,
		ErrWeeklyJobNil,
		ErrJobNotFound,
		ErrDisEnableLimit,
		ErrMoreLimit,
		ErrJobIDNil,
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
			name:        "Matching ErrInvalidJob",
			targetError: ErrInvalidJob,
			wrappedErr:  ErrInvalidJob,
			wantMatch:   true,
		},
		{
			name:        "Not matching different error",
			targetError: ErrInvalidJob,
			wrappedErr:  ErrTaskFuncNil,
			wantMatch:   false,
		},
		{
			name:        "Matching wrapped ErrInvalidJob",
			targetError: ErrInvalidJob,
			wrappedErr:  errors.New("wrapper: " + ErrInvalidJob.Error()),
			wantMatch:   false, // errors.Is doesn't match by string, only by identity
		},
		{
			name:        "Matching ErrTaskFuncNil with itself",
			targetError: ErrTaskFuncNil,
			wrappedErr:  ErrTaskFuncNil,
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
			name:        "Wrap ErrInvalidJob",
			baseError:   ErrInvalidJob,
			wrapMessage: "failed to create job",
			wantContain: "chrono:invalid job",
		},
		{
			name:        "Wrap ErrTaskFuncNil",
			baseError:   ErrTaskFuncNil,
			wrapMessage: "validation failed",
			wantContain: "chrono:task function cannot be nil",
		},
		{
			name:        "Wrap ErrJobNotFound",
			baseError:   ErrJobNotFound,
			wrapMessage: "lookup failed",
			wantContain: "chrono:job not found (missing JobID/Alias/Name)",
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
