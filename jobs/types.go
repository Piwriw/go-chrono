package jobs

// JobType represents the type of a scheduled jobs.
// JobType 表示计划任务的类型。
type JobType string

const (
	JobTypeOnce    JobType = "once"     // 一次性任务 / One-time task
	JobTypeCron    JobType = "cron"     // 定时任务 / Scheduled task
	JobTypeDaily   JobType = "daily"    // 每日任务 / Daily task
	JobTypeWeekly  JobType = "weekly"   // 每周任务 / Weekly task
	JobTypeMonthly JobType = "monthly"  // 每月任务 / Monthly task
	JobInterval    JobType = "interval" // 间隔任务 / Interval task
	JobTypeUnknown JobType = "unknown"  // 未知任务类型 / Unknown task type
)
