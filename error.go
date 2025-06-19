package chrono

import "errors"

var (
	ErrInvalidJob      = errors.New("chrono:invalid job")
	ErrTaskFuncNil     = errors.New("chrono:task function cannot be nil")
	ErrTaskTimeout     = errors.New("chrono:task is timeout")
	ErrValidateTimeout = errors.New("chrono:task timeout must be greater than 0")
	ErrTaskFailed      = errors.New("chrono:task is failed")
	ErrFoundAlias      = errors.New("chrono:can not found  by alias")
	ErrDisEnableAlias  = errors.New("chrono:alias is disable")
	ErrScheduleNil     = errors.New("chrono:schedule cannot be nil")
	ErrCronJobNil      = errors.New("chrono:cron job cannot be nil")
	ErrJobNotFound     = errors.New("chrono:job not found (missing JobID/Alias/Name)")
)
