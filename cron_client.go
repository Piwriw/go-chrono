package chrono

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type CronJobClientInterface interface {
	// CronExpr Linux Cron 表达式必须填写
	CronExpr(expr string) CronJobClientInterface
	Alias(alias string) CronJobClientInterface
	JobID(id string) CronJobClientInterface
	Name(name string) CronJobClientInterface
	Tags(tags ...string) CronJobClientInterface
	Task(task any, parameters ...any) CronJobClientInterface
	Timeout(timeout time.Duration) CronJobClientInterface
	Watch(watch func(event JobWatchInterface)) CronJobClientInterface
	DefaultHooks() CronJobClientInterface
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) CronJobClientInterface
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) CronJobClientInterface
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface
	Add() (gocron.Job, error)
	BatchAdd(cronJobs ...*CronJob) ([]gocron.Job, error)
	Remove() error
	Get() (gocron.Job, error)
}

type CronJobClient struct {
	scheduler *Scheduler
	job       *CronJob
}

var _ CronJobClientInterface = (*CronJobClient)(nil)

func (c *CronJobClient) CronExpr(expr string) CronJobClientInterface {
	c.job.CronExpr(expr)
	return c
}

func (c *CronJobClient) Alias(alias string) CronJobClientInterface {
	c.job.Alias(alias)
	return c
}

func (c *CronJobClient) JobID(id string) CronJobClientInterface {
	c.job.JobID(id)
	return c
}

func (c *CronJobClient) Name(name string) CronJobClientInterface {
	c.job.Names(name)
	return c
}

func (c *CronJobClient) Tags(tags ...string) CronJobClientInterface {
	c.job.Tag(tags...)
	return c
}

func (c *CronJobClient) Task(task any, parameters ...any) CronJobClientInterface {
	c.job.Task(task, parameters...)
	return c
}

func (c *CronJobClient) Timeout(timeout time.Duration) CronJobClientInterface {
	c.job.Timeout(timeout)
	return c
}

func (c *CronJobClient) Watch(watch func(event JobWatchInterface)) CronJobClientInterface {
	c.job.Watch(watch)
	return c
}

func (c *CronJobClient) DefaultHooks() CronJobClientInterface {
	c.job.DefaultHooks()
	return c
}

func (c *CronJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface {
	c.job.BeforeJobRuns(eventListenerFunc)
	return c
}

func (c *CronJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) CronJobClientInterface {
	c.job.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
	return c
}

func (c *CronJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) CronJobClientInterface {
	c.job.AfterJobRuns(eventListenerFunc)
	return c
}

func (c *CronJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface {
	c.job.AfterJobRunsWithError(eventListenerFunc)
	return c
}

func (c *CronJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) CronJobClientInterface {
	c.job.AfterJobRunsWithPanic(eventListenerFunc)
	return c
}

func (c *CronJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) CronJobClientInterface {
	c.job.AfterLockError(eventListenerFunc)
	return c
}

// Add 添加任务
func (c *CronJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrCronJobNil
	}
	return c.scheduler.AddCronJob(c.job)
}

func (c *CronJobClient) BatchAdd(cronJobs ...*CronJob) ([]gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	return c.scheduler.AddCronJobs(cronJobs...)
}

// Remove 删除任务
func (c *CronJobClient) Remove() error {
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.job == nil {
		return ErrCronJobNil
	}

	// 1. 优先 JobID
	if c.job.ID != "" {
		return c.scheduler.RemoveJob(c.job.ID)
	}
	// 2. 其次 Alias
	if c.job.Ali != "" {
		return c.scheduler.RemoveJobByAlias(c.job.Ali)
	}
	// 3. Name 或多条件才遍历
	if c.job.Name != "" {
		return c.scheduler.RemoveJobByName(c.job.Name)
	}
	return ErrJobNotFound
}

// Get 获取任务
func (c *CronJobClient) Get() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.job == nil {
		return nil, ErrCronJobNil
	}
	if c.job.ID != "" {
		return c.scheduler.GetJobByID(c.job.ID)
	}
	if c.job.Ali != "" {
		return c.scheduler.GetJobByAlias(c.job.Ali)
	}
	if c.job.Name != "" {
		return c.scheduler.GetJobByName(c.job.Name)
	}
	return nil, ErrJobNotFound
}
