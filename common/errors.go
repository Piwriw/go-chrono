package common

import "errors"

var (
	ErrInvalidJob     = errors.New("chrono:invalid jobs")
	ErrAtTimeDaysNil  = errors.New("chrono:at time must have at least one day")
	ErrTaskFuncNil    = errors.New("chrono:task function cannot be nil")
	ErrTaskFailed     = errors.New("chrono:task is failed")
	ErrFoundAlias     = errors.New("chrono:can not found  by alias")
	ErrDisEnableAlias = errors.New("chrono:alias is disable")
	ErrScheduleNil    = errors.New("chrono:schedule cannot be nil")
	ErrOnceJobNil     = errors.New("chrono:once jobs cannot be nil")
	ErrCronJobNil     = errors.New("chrono:cron jobs cannot be nil")
	ErrDailyJobNil    = errors.New("chrono:daily jobs cannot be nil")
	ErrIntervalJobNil = errors.New("chrono:interval jobs cannot be nil")
	ErrMonthJobNil    = errors.New("chrono:monthly jobs cannot be nil")
	ErrWeeklyJobNil   = errors.New("chrono:weekly jobs cannot be nil")
	ErrJobNotFound    = errors.New("chrono:jobs not found (missing JobID/Alias/Name)")
	ErrDisEnableLimit = errors.New("chrono:limit is disable")
	ErrMoreLimit      = errors.New("chrono:limit is reached")
	ErrJobIDNil       = errors.New("chrono:jobs id cannot be nil")
)
