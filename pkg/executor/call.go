package executor

import (
	"fmt"
	"reflect"
)

// CallJobFunc calls a jobs function with parameters.
// CallJobFunc 调用带有参数的任务函数。
//
// Parameters:
//
//	jobFunc - The jobs function to call / 要调用的任务函数
//	params  - Variable parameters to pass to the function / 传递给函数的可变参数
//
// Returns:
//
//	error - Error if the call fails / 如果调用失败则返回错误
func CallJobFunc(jobFunc any, params ...any) error {
	if jobFunc == nil {
		return nil
	}
	f := reflect.ValueOf(jobFunc)
	if f.IsZero() {
		return nil
	}
	if len(params) != f.Type().NumIn() {
		return fmt.Errorf("chrono:expected function with %d parameters, got one with %d", f.Type().NumIn(), len(params))
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	returnValues := f.Call(in)
	for _, val := range returnValues {
		i := val.Interface()
		if err, ok := i.(error); ok {
			return err
		}
	}
	return nil
}
