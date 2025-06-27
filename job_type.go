package chrono

type JobType string

const (
	JobTypeOnce    JobType = "once"     // 一次性任务
	JobTypeCron    JobType = "cron"     // 定时任务
	JobTypeDaily   JobType = "daily"    // 每日任务
	JobTypeWeekly  JobType = "weekly"   // 每周任务
	JobTypeMonthly JobType = "monthly"  // 每月任务
	JobInterval    JobType = "interval" // 间隔任务
	JobTypeUnknown JobType = "unknown"  // 未知任务类型
)
