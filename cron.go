package chrono

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

const (
	DayTimeType   string = "%d %d * * * "
	WeekTimeType  string = "%d %d * * %d "
	MonthTimeType string = "%d %d * %d * "
)

type CronJob struct {
	ID         string
	Ali        string
	Name       string
	Expr       string
	TaskFunc   any
	Parameters []any
	Hooks      []gocron.EventListener
	WatchFunc  func(event JobWatchInterface)
	timeout    time.Duration
	err        error
}

func NewCronJob(expr string) *CronJob {
	return &CronJob{
		Expr: expr,
	}
}

func (c *CronJob) Error() string {
	return c.err.Error()
}

func (c *CronJob) Alias(alias string) *CronJob {
	c.Ali = alias
	return c
}

func (c *CronJob) JobID(id string) *CronJob {
	c.ID = id
	return c
}

func (c *CronJob) Names(name string) *CronJob {
	if name == "" {
		name = uuid.New().String()
	}
	c.Name = name
	return c
}

func (c *CronJob) Task(task any, parameters ...any) *CronJob {
	if task == nil {
		c.err = errors.Join(c.err, ErrTaskFuncNil)
		return c
	}
	c.Parameters = append(c.Parameters, parameters...)
	return c
}

func (c *CronJob) Timeout(timeout time.Duration) *CronJob {
	if timeout <= 0 {
		c.err = errors.Join(c.err, ErrValidateTimeout)
		return c
	}
	c.timeout = timeout
	return c
}

func (c *CronJob) Watch(watch func(event JobWatchInterface)) *CronJob {
	c.WatchFunc = watch
	return c
}

func (c *CronJob) addHooks(hook ...gocron.EventListener) *CronJob {
	if c.Hooks == nil {
		c.Hooks = make([]gocron.EventListener, 0)
	}
	c.Hooks = append(c.Hooks, hook...)
	return c
}

func (c *CronJob) DefaultHooks() *CronJob {
	return c.addHooks(
		gocron.BeforeJobRuns(defaultBeforeJobRuns),
		gocron.BeforeJobRunsSkipIfBeforeFuncErrors(defaultBeforeJobRunsSkipIfBeforeFuncErrors),
		gocron.AfterJobRuns(defaultAfterJobRuns),
		gocron.AfterJobRunsWithError(defaultAfterJobRunsWithError),
		gocron.AfterJobRunsWithPanic(defaultAfterJobRunsWithPanic),
		gocron.AfterLockError(defaultAfterLockError))
}

// BeforeJobRuns 添加任务运行前的钩子函数
func (c *CronJob) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.addHooks(gocron.BeforeJobRuns(eventListenerFunc))
}

// BeforeJobRunsSkipIfBeforeFuncErrors 添加任务运行前的钩子函数（如果前置函数出错则跳过）
func (c *CronJob) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob {
	return c.addHooks(gocron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc))
}

// AfterJobRuns 添加任务运行后的钩子函数
func (c *CronJob) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.addHooks(gocron.AfterJobRuns(eventListenerFunc))
}

// AfterJobRunsWithError 添加任务运行出错时的钩子函数
func (c *CronJob) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.addHooks(gocron.AfterJobRunsWithError(eventListenerFunc))
}

// AfterJobRunsWithPanic 添加任务运行发生 panic 时的钩子函数
func (c *CronJob) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob {
	return c.addHooks(gocron.AfterJobRunsWithPanic(eventListenerFunc))
}

// AfterLockError 添加任务加锁出错时的钩子函数
func (c *CronJob) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.addHooks(gocron.AfterLockError(eventListenerFunc))
}

type TimeType string

// DayTimeToCron 将 time.Time 转换为 Cron 表达式
func DayTimeToCron(t time.Time) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(DayTimeType, minute, hour)
}

func WeekTimeToCron(t time.Time, week time.Weekday) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(WeekTimeType, minute, hour, week)
}

func MonthTimeToCron(t time.Time, month time.Month) string {
	// 提取时间字段
	minute := t.Minute()
	hour := t.Hour()

	// 返回 Cron 表达式
	return fmt.Sprintf(MonthTimeType, minute, hour, month)
}
