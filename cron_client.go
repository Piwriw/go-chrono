package chrono

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"time"
)

type CronJobClient struct {
	scheduler *Scheduler
	cron      *CronJob
}

type CronJobInterface interface {
	Alias(alias string) *CronJob
	JobID(id string) *CronJob
	Name(name string) *CronJob
	Task(task any, parameters ...any) *CronJob
	Timeout(timeout time.Duration) *CronJob
	Watch(watch func(event JobWatchInterface)) *CronJob
	DefaultHooks() *CronJob
	BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob
	BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob
	AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob
	AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob
	AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob
	AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob
}

func (c *CronJobClient) Alias(alias string) *CronJob {
	return c.cron.Alias(alias)
}

func (c *CronJobClient) JobID(id string) *CronJob {
	return c.cron.JobID(id)
}

func (c *CronJobClient) Name(name string) *CronJob {
	return c.cron.Names(name)
}

func (c *CronJobClient) Task(task any, parameters ...any) *CronJob {
	return c.cron.Task(task, parameters...)
}

func (c *CronJobClient) Timeout(timeout time.Duration) *CronJob {
	return c.cron.Timeout(timeout)
}

func (c *CronJobClient) Watch(watch func(event JobWatchInterface)) *CronJob {
	return c.cron.Watch(watch)
}

func (c *CronJobClient) DefaultHooks() *CronJob {
	return c.cron.DefaultHooks()
}

func (c *CronJobClient) BeforeJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.cron.BeforeJobRuns(eventListenerFunc)
}

func (c *CronJobClient) BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc func(jobID uuid.UUID, jobName string) error) *CronJob {
	return c.cron.BeforeJobRunsSkipIfBeforeFuncErrors(eventListenerFunc)
}

func (c *CronJobClient) AfterJobRuns(eventListenerFunc func(jobID uuid.UUID, jobName string)) *CronJob {
	return c.cron.AfterJobRuns(eventListenerFunc)
}

func (c *CronJobClient) AfterJobRunsWithError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.cron.AfterJobRunsWithError(eventListenerFunc)
}

func (c *CronJobClient) AfterJobRunsWithPanic(eventListenerFunc func(jobID uuid.UUID, jobName string, recoverData any)) *CronJob {
	return c.cron.AfterJobRunsWithPanic(eventListenerFunc)
}

func (c *CronJobClient) AfterLockError(eventListenerFunc func(jobID uuid.UUID, jobName string, err error)) *CronJob {
	return c.cron.AfterLockError(eventListenerFunc)
}

// Add 添加任务
func (c *CronJobClient) Add() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.cron == nil {
		return nil, ErrCronJobNil
	}
	return c.scheduler.AddCronJob(c.cron)
}

// Remove 删除任务
func (c *CronJobClient) Remove() error {
	if c.scheduler == nil {
		return ErrScheduleNil
	}
	if c.cron == nil {
		return ErrCronJobNil
	}

	// 1. 优先 JobID
	if c.cron.ID != "" {
		return c.scheduler.RemoveJob(c.cron.ID)
	}
	// 2. 其次 Alias
	if c.cron.Ali != "" {
		return c.scheduler.RemoveJobByAlias(c.cron.Ali)
	}
	// 3. Name 或多条件才遍历
	if c.cron.Name != "" {
		jobs, err := c.scheduler.GetJobs()
		if err != nil {
			return err
		}
		for _, job := range jobs {
			if job.Name() == c.cron.Name {
				return c.scheduler.RemoveJob(job.ID().String())
			}
		}
		return fmt.Errorf("job with name %s not found", c.cron.Name)
	}
	return ErrJobNotFound
}

// Get 获取任务
func (c *CronJobClient) Get() (gocron.Job, error) {
	if c.scheduler == nil {
		return nil, ErrScheduleNil
	}
	if c.cron == nil {
		return nil, ErrCronJobNil
	}
	if c.cron.ID != "" {
		return c.scheduler.GetJobByID(c.cron.ID)
	}
	if c.cron.Ali != "" {
		return c.scheduler.GetJobByAlias(c.cron.Ali)
	}
	if c.cron.Name != "" {
		return c.scheduler.GetJobByName(c.cron.Name)
	}
	return nil, ErrJobNotFound
}
