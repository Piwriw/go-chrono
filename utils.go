package chrono

import (
	"fmt"
	"net"
	"reflect"
	"strings"
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

func validateURLAddr(addr string) error {
	// 如果是本地地址，例如 "localhost" 或 "127.0.0.1"
	if strings.HasPrefix(addr, "localhost") || strings.HasPrefix(addr, "127.0.0.1") {
		return nil
	}

	// 使用 net.ParseIP 检查是否是有效的 IP 地址
	ip := net.ParseIP(addr)
	if ip != nil {
		return nil
	}

	// 如果是端口检查 (例如 :8080)
	if strings.HasPrefix(addr, ":") {
		return nil
	}

	// 如果是主机名和端口号的组合，例如 localhost:8080
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %v", err)
	}

	return nil
}
