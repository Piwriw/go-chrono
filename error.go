package chrono

import "errors"

var (
	ErrInvalidJob      = errors.New("chrono:invalid job")
	ErrAtTimeDaysNil   = errors.New("chrono:at time must have at least one day")
	ErrTaskFuncNil     = errors.New("chrono:task function cannot be nil")
	ErrTaskTimeout     = errors.New("chrono:task is timeout")
	ErrValidateTimeout = errors.New("chrono:task timeout must be greater than 0")
	ErrTaskFailed      = errors.New("chrono:task is failed")
	ErrFoundAlias      = errors.New("chrono:can not found  by alias")
	ErrDisEnableAlias  = errors.New("chrono:alias is disable")
	ErrScheduleNil     = errors.New("chrono:schedule cannot be nil")
	ErrOnceJobNil      = errors.New("chrono:once job cannot be nil")
	ErrCronJobNil      = errors.New("chrono:cron job cannot be nil")
	ErrDailyJobNil     = errors.New("chrono:daily job cannot be nil")
	ErrIntervalJobNil  = errors.New("chrono:interval job cannot be nil")
	ErrMonthJobNil     = errors.New("chrono:monthly job cannot be nil")
	ErrWeeklyJobNil    = errors.New("chrono:weekly job cannot be nil")
	ErrJobNotFound     = errors.New("chrono:job not found (missing JobID/Alias/Name)")
	ErrDisEnableLimit  = errors.New("chrono:limit is disable")
	ErrMoreLimit       = errors.New("chrono:limit is reached")
)
