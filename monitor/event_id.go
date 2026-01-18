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
	b.WriteString(jobSpec.JobID)
	b.WriteByte('_')
	b.WriteString(jobSpec.JobName)
	b.WriteByte('_')
	if g.timeFormat != "" {
		b.WriteString(time.Now().Format(g.timeFormat))
	}
	b.WriteString(time.Now().Format("20060102150405"))
	return b.String()
}
