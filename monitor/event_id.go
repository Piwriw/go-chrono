package monitor

import (
	"strings"
	"sync/atomic"
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
	counter    atomic.Uint64
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
	now := time.Now()
	if g.timeFormat != "" {
		b.WriteString(now.Format(g.timeFormat))
	}
	// Add nanosecond precision for uniqueness
	// 添加纳秒精度以确保唯一性
	b.WriteString(now.Format("20060102150405"))
	b.WriteString(now.Format(".000000000"))
	// Append monotonic counter to guarantee uniqueness under concurrency:
	// multiple goroutines may produce the same timestamp within the same nanosecond.
	// 追加单调递增计数器以确保并发下的唯一性：
	// 多个 goroutine 可能在同一纳秒内产生相同的时间戳。
	b.WriteByte('_')
	b.WriteString(itoaUint64(g.counter.Add(1)))
	return b.String()
}

// itoaUint64 converts a uint64 to its decimal string representation without using strconv,
	// avoiding the import just for this one call site.
// itoaUint64 将 uint64 转换为十进制字符串表示，不使用 strconv，
// 避免仅为这一处调用而引入额外的包。
func itoaUint64(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
