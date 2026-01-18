package scheduler

import (
	"testing"

	"github.com/piwriw/go-chrono/jobs"
)

// TestJobTypeString tests JobType string representation.
// TestJobTypeString 测试 JobType 的字符串表示
func TestJobTypeString(t *testing.T) {
	tests := []struct {
		name     string       // 测试名称 / Test name
		jobType  jobs.JobType // 作业类型 / Job type
		expected string       // 期望的字符串值 / Expected string value
	}{
		{
			name:     "jobs.JobTypeOnce returns once",
			jobType:  jobs.JobTypeOnce,
			expected: "once",
		},
		{
			name:     "jobs.JobTypeCron returns cron",
			jobType:  jobs.JobTypeCron,
			expected: "cron",
		},
		{
			name:     "jobs.JobTypeDaily returns daily",
			jobType:  jobs.JobTypeDaily,
			expected: "daily",
		},
		{
			name:     "jobs.JobTypeWeekly returns weekly",
			jobType:  jobs.JobTypeWeekly,
			expected: "weekly",
		},
		{
			name:     "jobs.JobTypeMonthly returns monthly",
			jobType:  jobs.JobTypeMonthly,
			expected: "monthly",
		},
		{
			name:     "JobInterval returns interval",
			jobType:  jobs.JobInterval,
			expected: "interval",
		},
		{
			name:     "jobs.JobTypeUnknown returns unknown",
			jobType:  jobs.JobTypeUnknown,
			expected: "unknown",
		},
		{
			name:     "Custom JobType returns its value",
			jobType:  jobs.JobType("custom"),
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
		name     string       // 测试名称 / Test name
		jobType1 jobs.JobType // 第一个作业类型 / First jobs type
		jobType2 jobs.JobType // 第二个作业类型 / Second jobs type
		wantEq   bool         // 是否期望相等 / Whether equality is expected
	}{
		{
			name:     "Same jobs types are equal",
			jobType1: jobs.JobTypeOnce,
			jobType2: jobs.JobTypeOnce,
			wantEq:   true,
		},
		{
			name:     "Different jobs types are not equal",
			jobType1: jobs.JobTypeCron,
			jobType2: jobs.JobTypeDaily,
			wantEq:   false,
		},
		{
			name:     "Job type and string constant are equal",
			jobType1: jobs.JobTypeWeekly,
			jobType2: jobs.JobType("weekly"),
			wantEq:   true,
		},
		{
			name:     "JobInterval and jobs.JobTypeOnce are not equal",
			jobType1: jobs.JobInterval,
			jobType2: jobs.JobTypeOnce,
			wantEq:   false,
		},
		{
			name:     "jobs.JobTypeUnknown equals itself",
			jobType1: jobs.JobTypeUnknown,
			jobType2: jobs.JobTypeUnknown,
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

// TestJobTypeAllDefined tests that all jobs type constants are defined and not empty.
// TestJobTypeAllDefined 测试所有作业类型常量是否已定义且不为空
func TestJobTypeAllDefined(t *testing.T) {
	tests := []struct {
		name    string       // 测试名称 / Test name
		jobType jobs.JobType // 要检查的作业类型 / Job type to check
		wantVal string       // 期望的字符串值 / Expected string value
	}{
		{
			name:    "jobs.JobTypeOnce is defined",
			jobType: jobs.JobTypeOnce,
			wantVal: "once",
		},
		{
			name:    "jobs.JobTypeCron is defined",
			jobType: jobs.JobTypeCron,
			wantVal: "cron",
		},
		{
			name:    "jobs.JobTypeDaily is defined",
			jobType: jobs.JobTypeDaily,
			wantVal: "daily",
		},
		{
			name:    "jobs.JobTypeWeekly is defined",
			jobType: jobs.JobTypeWeekly,
			wantVal: "weekly",
		},
		{
			name:    "jobs.JobTypeMonthly is defined",
			jobType: jobs.JobTypeMonthly,
			wantVal: "monthly",
		},
		{
			name:    "JobInterval is defined",
			jobType: jobs.JobInterval,
			wantVal: "interval",
		},
		{
			name:    "jobs.JobTypeUnknown is defined",
			jobType: jobs.JobTypeUnknown,
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
