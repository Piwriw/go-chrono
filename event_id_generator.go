package chrono

import (
	"github.com/google/uuid"
	"strings"
	"time"
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
