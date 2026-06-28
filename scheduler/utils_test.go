package scheduler

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/piwriw/go-chrono/pkg/executor"
	"github.com/piwriw/go-chrono/pkg/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCallJobFunc tests the pkg.CallJobFunc utility function.
// TestCallJobFunc 测试 pkg.CallJobFunc 工具函数。
//
// Test scenarios:
// - Happy path: successful function call with various parameter types
// - Boundary conditions: nil function, zero value function, no parameters
// - Exception cases: parameter count mismatch, function returns error
func TestCallJobFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string // 测试用例名称 / Test case name
		jobFunc  any    // 要调用的函数 / Function to call
		params   []any  // 函数参数 / Function parameters
		wantErr  bool   // 是否预期错误 / Expect error
		errMsg   string // 预期错误信息子字符串 / Expected error message substring
		setup    func() // 测试前设置 / Setup before test
		teardown func() // 测试后清理 / Cleanup after test
	}{
		{
			name:    "successful call with no parameters",
			jobFunc: func() error { return nil },
			params:  []any{},
			wantErr: false,
		},
		{
			name:    "successful call with one parameter",
			jobFunc: func(s string) error { return nil },
			params:  []any{"test"},
			wantErr: false,
		},
		{
			name:    "successful call with multiple parameters",
			jobFunc: func(s string, i int, b bool) error { return nil },
			params:  []any{"test", 42, true},
			wantErr: false,
		},
		{
			name:    "function returns error",
			jobFunc: func() error { return errors.New("task failed") },
			params:  []any{},
			wantErr: true,
			errMsg:  "task failed",
		},
		{
			name:    "nil function",
			jobFunc: nil,
			params:  []any{},
			wantErr: false, // nil function returns nil, not error
		},
		{
			name:    "zero value function",
			jobFunc: (func(string) error)(nil),
			params:  []any{"test"},
			wantErr: false, // zero value function returns nil, not error
		},
		{
			name:    "parameter count mismatch - fewer params",
			jobFunc: func(s string, i int) error { return nil },
			params:  []any{"test"},
			wantErr: true,
			errMsg:  "expected function with 2 parameters, got one with 1",
		},
		{
			name:    "parameter count mismatch - more params",
			jobFunc: func(s string) error { return nil },
			params:  []any{"test", 42},
			wantErr: true,
			errMsg:  "expected function with 1 parameters, got one with 2",
		},
		{
			name:    "function with variadic parameters",
			jobFunc: func(s string, args ...int) error { return nil },
			params:  []any{"test", 1, 2, 3},
			wantErr: true, // reflect doesn't handle variadic directly
		},
		{
			name:    "function with no return value",
			jobFunc: func(s string) {},
			params:  []any{"test"},
			wantErr: false,
		},
		{
			name:    "function with multiple return values, last is error",
			jobFunc: func(s string) (int, error) { return 42, nil },
			params:  []any{"test"},
			wantErr: false,
		},
		{
			name:    "function with multiple return values, error is nil",
			jobFunc: func() (int, string, error) { return 42, "result", nil },
			params:  []any{},
			wantErr: false,
		},
		{
			name:    "function with multiple return values, error is not nil",
			jobFunc: func() (int, string, error) { return 0, "", errors.New("multi error") },
			params:  []any{},
			wantErr: true,
			errMsg:  "multi error",
		},
		{
			name: "successful call with time parameter",
			jobFunc: func(tm time.Time) error {
				assert.False(t, tm.IsZero())
				return nil
			},
			params:  []any{time.Now()},
			wantErr: false,
		},
		{
			name:    "successful call with interface parameter",
			jobFunc: func(a any) error { return nil },
			params:  []any{"test"},
			wantErr: false,
		},
		{
			name:    "successful call with struct parameter",
			jobFunc: func(s struct{ Name string }) error { return nil },
			params:  []any{struct{ Name string }{Name: "test"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup if provided
			// 如果提供了设置函数，则执行
			if tt.setup != nil {
				tt.setup()
			}
			// Ensure teardown is called
			// 确保清理函数被调用
			if tt.teardown != nil {
				defer tt.teardown()
			}

			// Execute the function
			// 执行函数
			err := executor.CallJobFunc(tt.jobFunc, tt.params...)

			// Verify results
			// 验证结果
			if tt.wantErr {
				require.Error(t, err, "pkg.CallJobFunc should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg, "error message should contain expected substring")
				}
			} else {
				assert.NoError(t, err, "pkg.CallJobFunc should not return an error")
			}
		})
	}
}

// BenchmarkCallJobFunc benchmarks the pkg.CallJobFunc function.
// BenchmarkCallJobFunc 对 pkg.CallJobFunc 函数进行基准测试。
func BenchmarkCallJobFunc(b *testing.B) {
	fn := func(s string, i int) error { return nil }
	params := []any{"test", 42}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = executor.CallJobFunc(fn, params...)
	}
}

// ExampleCallJobFunc demonstrates the usage of pkg.CallJobFunc.
// ExampleCallJobFunc 展示 pkg.CallJobFunc 的使用方法。
func ExampleCallJobFunc() {
	// Define a function to call
	// 定义要调用的函数
	task := func(name string, count int) error {
		fmt.Printf("Task: %s, Count: %d\n", name, count)
		return nil
	}

	// Call the function with parameters
	// 使用参数调用函数
	err := executor.CallJobFunc(task, "example", 10)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Output: Task: example, Count: 10
}

// TestValidateURLAddr tests the pkg.ValidateURLAddr utility function.
// TestValidateURLAddr 测试 pkg.ValidateURLAddr 工具函数。
//
// Test scenarios:
// - Happy path: valid localhost, IP addresses, host:port combinations
// - Boundary conditions: empty string, port only
// - Exception cases: invalid address formats
func TestValidateURLAddr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string // 测试用例名称 / Test case name
		addr    string // 要验证的地址 / Address to validate
		wantErr bool   // 是否预期错误 / Expect error
		errMsg  string // 预期错误信息子字符串 / Expected error message substring
	}{
		// Happy path - valid addresses
		// 正常场景 - 有效地址
		{
			name:    "localhost address",
			addr:    "localhost",
			wantErr: false,
		},
		{
			name:    "localhost with port",
			addr:    "localhost:8080",
			wantErr: false,
		},
		{
			name:    "127.0.0.1 address",
			addr:    "127.0.0.1",
			wantErr: false,
		},
		{
			name:    "127.0.0.1 with port",
			addr:    "127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "valid IPv4 address",
			addr:    "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "valid IPv4 with port",
			addr:    "192.168.1.1:9090",
			wantErr: false,
		},
		{
			name:    "valid IPv6 address",
			addr:    "::1",
			wantErr: false,
		},
		{
			name:    "valid IPv6 with port",
			addr:    "[::1]:8080",
			wantErr: false,
		},
		{
			name:    "port only - colon prefix",
			addr:    ":8080",
			wantErr: false,
		},
		{
			name:    "valid hostname with port",
			addr:    "example.com:8080",
			wantErr: false,
		},
		{
			name:    "valid hostname with port - underscore",
			addr:    "example_host.com:8080",
			wantErr: false,
		},

		// Boundary conditions
		// 边界条件
		{
			name:    "empty string",
			addr:    "",
			wantErr: true,
			errMsg:  "missing port in address",
		},
		{
			name:    "port only - numeric",
			addr:    ":0",
			wantErr: false,
		},
		{
			name:    "port only - max port",
			addr:    ":65535",
			wantErr: false,
		},
		{
			name:    "port only - invalid high port",
			addr:    ":99999",
			wantErr: false, // validation doesn't check port range
		},

		// Exception cases - invalid addresses
		// 异常场景 - 无效地址
		{
			name:    "invalid hostname without port",
			addr:    "example.com",
			wantErr: true,
			errMsg:  "missing port",
		},
		{
			name:    "invalid characters",
			addr:    "invalid@address:8080",
			wantErr: true,
			errMsg:  "invalid address",
		},
		{
			name:    "invalid port format",
			addr:    "localhost:abc",
			wantErr: true,
			errMsg:  "invalid address",
		},
		{
			name:    "space in address",
			addr:    "localhost :8080",
			wantErr: true,
			errMsg:  "invalid address",
		},
		{
			name:    "multiple colons without port",
			addr:    "localhost::8080",
			wantErr: true,
			errMsg:  "invalid address",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Execute validation
			// 执行验证
			err := url.ValidateURLAddr(tt.addr)

			// Verify results
			// 验证结果
			if tt.wantErr {
				require.Error(t, err, "pkg.ValidateURLAddr should return an error")
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg, "error message should contain expected substring")
				}
			} else {
				assert.NoError(t, err, "pkg.ValidateURLAddr should not return an error")
			}
		})
	}
}

// BenchmarkValidateURLAddr benchmarks the pkg.ValidateURLAddr function.
// BenchmarkValidateURLAddr 对 pkg.ValidateURLAddr 函数进行基准测试。
func BenchmarkValidateURLAddr(b *testing.B) {
	validAddr := "localhost:8080"

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = url.ValidateURLAddr(validAddr)
	}
}

// ExampleValidateURLAddr demonstrates the usage of pkg.ValidateURLAddr.
// ExampleValidateURLAddr 展示 pkg.ValidateURLAddr 的使用方法。
func ExampleValidateURLAddr() {
	// Valid addresses
	// 有效地址
	err := url.ValidateURLAddr("localhost:8080")
	fmt.Println("localhost:8080:", err == nil)

	err = url.ValidateURLAddr("192.168.1.1:9090")
	fmt.Println("192.168.1.1:9090:", err == nil)

	// Invalid address
	// 无效地址
	err = url.ValidateURLAddr("example.com")
	fmt.Println("example.com:", err)

	// Output:
	// localhost:8080: true
	// 192.168.1.1:9090: true
	// example.com: missing port in address
}
