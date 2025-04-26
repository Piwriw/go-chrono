package chrono

import (
	"fmt"
	"reflect"
)

func callJobFunc(jobFunc any, params ...any) error {
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
