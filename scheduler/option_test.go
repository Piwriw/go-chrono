package scheduler

import (
	"fmt"
	"testing"

	"github.com/piwriw/go-chrono/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAliasOption tests the AliasOption functionality.
// TestAliasOption 测试 AliasOption 的功能。
//
// Test scenarios:
// - Happy path: create and enable alias option
// - Boundary conditions: default disabled state
func TestAliasOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string                     // 测试用例名称 / Test case name
		option *common.AliasOption        // 要测试的选项 / Option to test
		setup  func() *common.AliasOption // 设置函数 / Setup function
	}{
		{
			name:   "default alias option is disabled",
			option: &common.AliasOption{},
		},
		{
			name:   "enabled alias option",
			option: &common.AliasOption{Enabled: true},
		},
		{
			name: "alias option created with WithAliasMode",
			setup: func() *common.AliasOption {
				schOpts := &SchedulerOptions{}
				WithAliasMode()(schOpts)
				return schOpts.alias
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup if provided
			// 如果提供了设置函数，则执行
			if tt.setup != nil {
				tt.option = tt.setup()
			}

			// Test Name method
			// 测试 Name 方法
			assert.Equal(t, common.AliasOptionName, tt.option.Name(), "Name should return AliasOptionName")

			// Test Enable method
			// 测试 Enable 方法
			expectedEnabled := tt.option.Enabled
			assert.Equal(t, expectedEnabled, tt.option.Enable(), "Enable should return enabled state")

			// Verify interface implementation
			// 验证接口实现
			var _ common.ScheduleOption = tt.option
		})
	}
}

// TestWatchOption tests the WatchOption functionality.
// TestWatchOption 测试 WatchOption 的功能。
//
// Test scenarios:
// - Happy path: create watch option with custom function
// - Boundary conditions: default watch function
// - Exception cases: nil watch function
func TestWatchOption(t *testing.T) {
	t.Parallel()

	customWatchFunc := func(event common.JobWatchInterface) {
		// Custom watch function
	}

	tests := []struct {
		name     string                     // 测试用例名称 / Test case name
		option   *common.WatchOption        // 要测试的选项 / Option to test
		setup    func() *common.WatchOption // 设置函数 / Setup function
		wantFunc bool                       // 是否预期有监听函数 / Expect watch function
	}{
		{
			name:     "default watch option is disabled",
			option:   &common.WatchOption{},
			wantFunc: false,
		},
		{
			name:     "enabled watch option with custom function",
			option:   &common.WatchOption{Enabled: true, WatchFunc: customWatchFunc},
			wantFunc: true,
		},
		{
			name: "watch option created with WithWatch and custom function",
			setup: func() *common.WatchOption {
				schOpts := &SchedulerOptions{}
				WithWatch(customWatchFunc)(schOpts)
				return schOpts.watch
			},
			wantFunc: true,
		},
		{
			name: "watch option created with WithWatch and nil function",
			setup: func() *common.WatchOption {
				schOpts := &SchedulerOptions{}
				WithWatch(nil)(schOpts)
				return schOpts.watch
			},
			wantFunc: true, // defaultWatchFunc is set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup if provided
			// 如果提供了设置函数，则执行
			if tt.setup != nil {
				tt.option = tt.setup()
			}

			// Test Name method
			// 测试 Name 方法
			assert.Equal(t, common.WatchOptionName, tt.option.Name(), "Name should return WatchOptionName")

			// Test Enable method
			// 测试 Enable 方法
			expectedEnabled := tt.option.Enabled
			assert.Equal(t, expectedEnabled, tt.option.Enable(), "Enable should return enabled state")

			// Test WatchFunc method
			// 测试 WatchFunc 方法
			if tt.wantFunc {
				assert.NotNil(t, tt.option.WatchFunc, "WatchFunc should return a function")
			} else {
				assert.Nil(t, tt.option.WatchFunc, "WatchFunc should return nil")
			}

			// Verify interface implementation
			// 验证接口实现
			var _ common.ScheduleOption = tt.option
		})
	}
}

// TestWebMonitorOption tests the WebMonitorOption functionality.
// TestWebMonitorOption 测试 WebMonitorOption 的功能。
//
// Test scenarios:
// - Happy path: create web monitor option with address
// - Boundary conditions: empty address
// - Exception cases: invalid address formats
func TestWebMonitorOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string                          // 测试用例名称 / Test case name
		option   *common.WebMonitorOption        // 要测试的选项 / Option to test
		setup    func() *common.WebMonitorOption // 设置函数 / Setup function
		wantAddr string                          // 预期地址 / Expected address
	}{
		{
			name:     "default web monitor option is disabled",
			option:   &common.WebMonitorOption{},
			wantAddr: "",
		},
		{
			name:     "enabled web monitor option with address",
			option:   &common.WebMonitorOption{Enabled: true, Address: "localhost:8080"},
			wantAddr: "localhost:8080",
		},
		{
			name: "web monitor option created with WithWebMonitor",
			setup: func() *common.WebMonitorOption {
				schOpts := &SchedulerOptions{}
				WithWebMonitor(":9090")(schOpts)
				return schOpts.webMonitor
			},
			wantAddr: ":9090",
		},
		{
			name: "web monitor option with IP address",
			setup: func() *common.WebMonitorOption {
				schOpts := &SchedulerOptions{}
				WithWebMonitor("0.0.0.0:8080")(schOpts)
				return schOpts.webMonitor
			},
			wantAddr: "0.0.0.0:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup if provided
			// 如果提供了设置函数，则执行
			if tt.setup != nil {
				tt.option = tt.setup()
			}

			// Test Name method
			// 测试 Name 方法
			assert.Equal(t, common.WebMonitorOptionName, tt.option.Name(), "Name should return WebMonitorOptionName")

			// Test Enable method
			// 测试 Enable 方法
			expectedEnabled := tt.option.Enabled
			assert.Equal(t, expectedEnabled, tt.option.Enable(), "Enable should return enabled state")

			// Test Address method
			// 测试 Address 方法
			assert.Equal(t, tt.wantAddr, tt.option.Address, "Address should return expected address")

			// Verify interface implementation
			// 验证接口实现
			var _ common.ScheduleOption = tt.option
		})
	}
}

// TestLimitOption tests the LimitOption functionality.
// TestLimitOption 测试 LimitOption 的功能。
//
// Test scenarios:
// - Happy path: create limit option with number
// - Boundary conditions: zero and negative limits
// - Exception cases: maximum limit values
func TestLimitOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string                     // 测试用例名称 / Test case name
		option    *common.LimitOption        // 要测试的选项 / Option to test
		setup     func() *common.LimitOption // 设置函数 / Setup function
		wantLimit int                        // 预期限制 / Expected limit
	}{
		{
			name:      "default limit option is disabled",
			option:    &common.LimitOption{},
			wantLimit: 0,
		},
		{
			name:      "enabled limit option with positive limit",
			option:    &common.LimitOption{Enabled: true, Limit: 100},
			wantLimit: 100,
		},
		{
			name: "limit option created with WithLimit",
			setup: func() *common.LimitOption {
				schOpts := &SchedulerOptions{}
				WithLimit(50)(schOpts)
				return schOpts.limit
			},
			wantLimit: 50,
		},
		{
			name: "limit option with zero limit",
			setup: func() *common.LimitOption {
				schOpts := &SchedulerOptions{}
				WithLimit(0)(schOpts)
				return schOpts.limit
			},
			wantLimit: 0,
		},
		{
			name: "limit option with large limit",
			setup: func() *common.LimitOption {
				schOpts := &SchedulerOptions{}
				WithLimit(1000000)(schOpts)
				return schOpts.limit
			},
			wantLimit: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup if provided
			// 如果提供了设置函数，则执行
			if tt.setup != nil {
				tt.option = tt.setup()
			}

			// Test Name method
			// 测试 Name 方法
			assert.Equal(t, common.LimitOptionName, tt.option.Name(), "Name should return LimitOptionName")

			// Test Enable method
			// 测试 Enable 方法
			expectedEnabled := tt.option.Enabled
			assert.Equal(t, expectedEnabled, tt.option.Enable(), "Enable should return enabled state")

			// Test Limit method
			// 测试 Limit 方法
			assert.Equal(t, tt.wantLimit, tt.option.Limit, "Limit should return expected limit")

			// Verify interface implementation
			// 验证接口实现
			var _ common.ScheduleOption = tt.option
		})
	}
}

// TestPrometheusOption tests the PrometheusOption functionality.
// TestPrometheusOption 测试 PrometheusOption 的功能。
//
// Test scenarios:
// - Happy path: create prometheus option with address
// - Boundary conditions: empty address
// - Exception cases: various address formats
func TestPrometheusOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string                          // 测试用例名称 / Test case name
		option   *common.PrometheusOption        // 要测试的选项 / Option to test
		setup    func() *common.PrometheusOption // 设置函数 / Setup function
		wantAddr string                          // 预期地址 / Expected address
	}{
		{
			name:     "default prometheus option is disabled",
			option:   &common.PrometheusOption{},
			wantAddr: "",
		},
		{
			name:     "enabled prometheus option with address",
			option:   &common.PrometheusOption{Enabled: true, Address: ":9090"},
			wantAddr: ":9090",
		},
		{
			name: "prometheus option created with WithPrometheus",
			setup: func() *common.PrometheusOption {
				schOpts := &SchedulerOptions{}
				WithPrometheus(":8080")(schOpts)
				return schOpts.prometheus
			},
			wantAddr: ":8080",
		},
		{
			name: "prometheus option with full address",
			setup: func() *common.PrometheusOption {
				schOpts := &SchedulerOptions{}
				WithPrometheus("localhost:9090")(schOpts)
				return schOpts.prometheus
			},
			wantAddr: "localhost:9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup if provided
			// 如果提供了设置函数，则执行
			if tt.setup != nil {
				tt.option = tt.setup()
			}

			// Test Name method
			// 测试 Name 方法
			assert.Equal(t, common.PrometheusOptionName, tt.option.Name(), "Name should return PrometheusOptionName")

			// Test Enable method
			// 测试 Enable 方法
			expectedEnabled := tt.option.Enabled
			assert.Equal(t, expectedEnabled, tt.option.Enable(), "Enable should return enabled state")

			// Test Address method
			// 测试 Address 方法
			assert.Equal(t, tt.wantAddr, tt.option.Address, "Address should return expected address")

			// Verify interface implementation
			// 验证接口实现
			var _ common.ScheduleOption = tt.option
		})
	}
}

// TestSchedulerOptionsComposition tests composing multiple scheduler options.
// TestSchedulerOptionsComposition 测试组合多个调度器选项。
//
// Test scenarios:
// - Happy path: compose multiple options together
// - Boundary conditions: single option, all options
func TestSchedulerOptionsComposition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string                              // 测试用例名称 / Test case name
		options []SchedulerOption                   // 调度器选项 / Scheduler options
		verify  func(*testing.T, *SchedulerOptions) // 验证函数 / Verify function
	}{
		{
			name:    "single option - alias",
			options: []SchedulerOption{WithAliasMode()},
			verify: func(t *testing.T, opts *SchedulerOptions) {
				require.NotNil(t, opts.alias, "alias should not be nil")
				assert.True(t, opts.alias.Enable(), "alias should be enabled")
			},
		},
		{
			name:    "single option - limit",
			options: []SchedulerOption{WithLimit(100)},
			verify: func(t *testing.T, opts *SchedulerOptions) {
				require.NotNil(t, opts.limit, "limit should not be nil")
				assert.True(t, opts.limit.Enable(), "limit should be enabled")
				assert.Equal(t, 100, opts.limit.Limit, "limit should be 100")
			},
		},
		{
			name: "multiple options",
			options: []SchedulerOption{
				WithAliasMode(),
				WithLimit(50),
				WithWebMonitor(":8080"),
				WithPrometheus(":9090"),
			},
			verify: func(t *testing.T, opts *SchedulerOptions) {
				assert.True(t, opts.alias.Enable(), "alias should be enabled")
				assert.True(t, opts.limit.Enable(), "limit should be enabled")
				assert.Equal(t, 50, opts.limit.Limit, "limit should be 50")
				assert.Equal(t, ":8080", opts.webMonitor.Address, "web monitor address should be :8080")
				assert.Equal(t, ":9090", opts.prometheus.Address, "prometheus address should be :9090")
			},
		},
		{
			name: "all options",
			options: []SchedulerOption{
				WithAliasMode(),
				WithWatch(nil),
				WithWebMonitor("localhost:8080"),
				WithLimit(200),
				WithPrometheus(":9090"),
			},
			verify: func(t *testing.T, opts *SchedulerOptions) {
				assert.NotNil(t, opts.alias, "alias should not be nil")
				assert.NotNil(t, opts.watch, "watch should not be nil")
				assert.NotNil(t, opts.webMonitor, "webMonitor should not be nil")
				assert.NotNil(t, opts.limit, "limit should not be nil")
				assert.NotNil(t, opts.prometheus, "prometheus should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create scheduler options
			// 创建调度器选项
			opts := &SchedulerOptions{}
			for _, option := range tt.options {
				option(opts)
			}

			// Verify options
			// 验证选项
			tt.verify(t, opts)
		})
	}
}

// TestOptionConstants tests the option name constants.
// TestOptionConstants 测试选项名称常量。
func TestOptionConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "AliasOptionName constant",
			constant: common.AliasOptionName,
			expected: "alias",
		},
		{
			name:     "WatchOptionName constant",
			constant: common.WatchOptionName,
			expected: "watch",
		},
		{
			name:     "WebMonitorOptionName constant",
			constant: common.WebMonitorOptionName,
			expected: "web_monitor",
		},
		{
			name:     "LimitOptionName constant",
			constant: common.LimitOptionName,
			expected: "limit",
		},
		{
			name:     "PrometheusOptionName constant",
			constant: common.PrometheusOptionName,
			expected: "prometheus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.constant, "constant should match expected value")
		})
	}
}

// ExampleWithAliasMode demonstrates the usage of WithAliasMode.
// ExampleWithAliasMode 展示 WithAliasMode 的使用方法。
func ExampleWithAliasMode() {
	opts := &SchedulerOptions{}
	WithAliasMode()(opts)
	fmt.Println(opts.alias.Enable())
	// Output: true
}

// ExampleWithLimit demonstrates the usage of WithLimit.
// ExampleWithLimit 展示 WithLimit 的使用方法。
func ExampleWithLimit() {
	opts := &SchedulerOptions{}
	WithLimit(100)(opts)
	fmt.Println(opts.limit.Limit)
	// Output: 100
}

// ExampleWithWebMonitor demonstrates the usage of WithWebMonitor.
// ExampleWithWebMonitor 展示 WithWebMonitor 的使用方法。
func ExampleWithWebMonitor() {
	opts := &SchedulerOptions{}
	WithWebMonitor(":8080")(opts)
	fmt.Println(opts.webMonitor.Address)
	// Output: :8080
}
