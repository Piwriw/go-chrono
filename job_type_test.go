package chrono

import (
	"testing"
)

// TestJobTypeString tests JobType string representation.
// TestJobTypeString 测试 JobType 的字符串表示
func TestJobTypeString(t *testing.T) {
	tests := []struct {
		name     string  // 测试名称 / Test name
		jobType  JobType // 作业类型 / Job type
		expected string  // 期望的字符串值 / Expected string value
	}{
		{
			name:     "JobTypeOnce returns once",
			jobType:  JobTypeOnce,
			expected: "once",
		},
		{
			name:     "JobTypeCron returns cron",
			jobType:  JobTypeCron,
			expected: "cron",
		},
		{
			name:     "JobTypeDaily returns daily",
			jobType:  JobTypeDaily,
			expected: "daily",
		},
		{
			name:     "JobTypeWeekly returns weekly",
			jobType:  JobTypeWeekly,
			expected: "weekly",
		},
		{
			name:     "JobTypeMonthly returns monthly",
			jobType:  JobTypeMonthly,
			expected: "monthly",
		},
		{
			name:     "JobInterval returns interval",
			jobType:  JobInterval,
			expected: "interval",
		},
		{
			name:     "JobTypeUnknown returns unknown",
			jobType:  JobTypeUnknown,
			expected: "unknown",
		},
		{
			name:     "Custom JobType returns its value",
			jobType:  JobType("custom"),
			expected: "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.jobType)
			if got != tt.expected {
				t.Errorf("JobType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestJobTypeEquality tests JobType equality comparison.
// TestJobTypeEquality 测试 JobType 的相等性比较
func TestJobTypeEquality(t *testing.T) {
	tests := []struct {
		name     string  // 测试名称 / Test name
		jobType1 JobType // 第一个作业类型 / First job type
		jobType2 JobType // 第二个作业类型 / Second job type
		wantEq   bool    // 是否期望相等 / Whether equality is expected
	}{
		{
			name:     "Same job types are equal",
			jobType1: JobTypeOnce,
			jobType2: JobTypeOnce,
			wantEq:   true,
		},
		{
			name:     "Different job types are not equal",
			jobType1: JobTypeCron,
			jobType2: JobTypeDaily,
			wantEq:   false,
		},
		{
			name:     "Job type and string constant are equal",
			jobType1: JobTypeWeekly,
			jobType2: JobType("weekly"),
			wantEq:   true,
		},
		{
			name:     "JobInterval and JobTypeOnce are not equal",
			jobType1: JobInterval,
			jobType2: JobTypeOnce,
			wantEq:   false,
		},
		{
			name:     "JobTypeUnknown equals itself",
			jobType1: JobTypeUnknown,
			jobType2: JobTypeUnknown,
			wantEq:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEq := tt.jobType1 == tt.jobType2
			if gotEq != tt.wantEq {
				t.Errorf("JobType equality = %v, want %v (jobType1=%v, jobType2=%v)",
					gotEq, tt.wantEq, tt.jobType1, tt.jobType2)
			}
		})
	}
}

// TestJobTypeAllDefined tests that all job type constants are defined and not empty.
// TestJobTypeAllDefined 测试所有作业类型常量是否已定义且不为空
func TestJobTypeAllDefined(t *testing.T) {
	tests := []struct {
		name    string  // 测试名称 / Test name
		jobType JobType // 要检查的作业类型 / Job type to check
		wantVal string  // 期望的字符串值 / Expected string value
	}{
		{
			name:    "JobTypeOnce is defined",
			jobType: JobTypeOnce,
			wantVal: "once",
		},
		{
			name:    "JobTypeCron is defined",
			jobType: JobTypeCron,
			wantVal: "cron",
		},
		{
			name:    "JobTypeDaily is defined",
			jobType: JobTypeDaily,
			wantVal: "daily",
		},
		{
			name:    "JobTypeWeekly is defined",
			jobType: JobTypeWeekly,
			wantVal: "weekly",
		},
		{
			name:    "JobTypeMonthly is defined",
			jobType: JobTypeMonthly,
			wantVal: "monthly",
		},
		{
			name:    "JobInterval is defined",
			jobType: JobInterval,
			wantVal: "interval",
		},
		{
			name:    "JobTypeUnknown is defined",
			jobType: JobTypeUnknown,
			wantVal: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.jobType) != tt.wantVal {
				t.Errorf("JobType %v = %v, want %v", tt.name, tt.jobType, tt.wantVal)
			}
		})
	}
}
