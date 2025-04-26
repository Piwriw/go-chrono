package chrono

import "errors"

var (
	ErrInvalidJob      = errors.New("chrono:invalid job")
	ErrTaskFuncNil     = errors.New("chrono:task function cannot be nil")
	ErrTaskTimeout     = errors.New("chrono:task is timeout")
	ErrValidateTimeout = errors.New("chrono:task timeout must be greater than 0")
	ErrTaskFailed      = errors.New("chrono:task is failed")
	ErrFoundAlias      = errors.New("chrono:can not found  by alias")
)
