package monitor

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// EventIDGenerator 接口
type EventIDGenerator interface {
	NextID(jobSpec JobSpec) string
}

// UUIDEventIDGenerator UUID生成器
type UUIDEventIDGenerator struct{}

var _ EventIDGenerator = &UUIDEventIDGenerator{}

func (g *UUIDEventIDGenerator) NextID(_ JobSpec) string {
	return uuid.New().String()
}

// TimeEventIDGenerator 时间戳生成器
// 时间戳格式：20060102150405
type TimeEventIDGenerator struct {
	timeFormat string
}

var _ EventIDGenerator = &TimeEventIDGenerator{}

// NewTimeEventIDGenerator creates a new TimeEventIDGenerator with custom time format.
// NewTimeEventIDGenerator 创建一个带有自定义时间格式的新 TimeEventIDGenerator。
//
// Parameters:
//
//	timeFormat - Time format string (e.g., "20060102150405") / 时间格式字符串（例如："20060102150405"）
//
// Returns:
//
//	*TimeEventIDGenerator - A new TimeEventIDGenerator instance / 新的 TimeEventIDGenerator 实例
func NewTimeEventIDGenerator(timeFormat string) *TimeEventIDGenerator {
	return &TimeEventIDGenerator{
		timeFormat: timeFormat,
	}
}

func (g *TimeEventIDGenerator) NextID(jobSpec JobSpec) string {
	var b strings.Builder
	// Transform JobID: replace "-jobs" suffix with "-job" (singularize)
	// 转换 JobID：将 "-jobs" 后缀替换为 "-job"（单数化）
	jobID := jobSpec.JobID
	if strings.HasSuffix(jobID, "-jobs") {
		jobID = strings.TrimSuffix(jobID, "-jobs") + "-job"
	}
	b.WriteString(jobID)
	b.WriteByte('_')
	// Transform JobName: replace "-jobs" suffix with "-job" (singularize)
	// 转换 JobName：将 "-jobs" 后缀替换为 "-job"（单数化）
	jobName := jobSpec.JobName
	if strings.HasSuffix(jobName, "-jobs") {
		jobName = strings.TrimSuffix(jobName, "-jobs") + "-job"
	}
	b.WriteString(jobName)
	b.WriteByte('_')
	if g.timeFormat != "" {
		b.WriteString(time.Now().Format(g.timeFormat))
	}
	// Add nanosecond precision for uniqueness
	// 添加纳秒精度以确保唯一性
	b.WriteString(time.Now().Format("20060102150405"))
	b.WriteString(time.Now().Format(".000000000"))
	return b.String()
}

// isCompleteWord checks if a string is a complete word (contains only alphabetic characters)
// isCompleteWord 检查字符串是否是完整单词（仅包含字母）
func isCompleteWord(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	return true
}
