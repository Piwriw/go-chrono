package common

import (
	"time"

	"github.com/piwriw/go-chrono/retry"
)

const (
	// AliasOptionName is the name constant for the alias option.
	// AliasOptionName 是别名选项的名称常量。
	AliasOptionName = "alias"
	// WatchOptionName is the name constant for the watch option.
	// WatchOptionName 是监听选项的名称常量。
	WatchOptionName = "watch"
	// WebMonitorOptionName is the name constant for the web monitor option.
	// WebMonitorOptionName 是 Web 监控选项的名称常量。
	WebMonitorOptionName = "web_monitor"
	// LimitOptionName is the name constant for the limit option.
	// LimitOptionName 是限制选项的名称常量。
	LimitOptionName = "limit"
	// PrometheusOptionName is the name constant for the prometheus option.
	// PrometheusOptionName 是 Prometheus 选项的名称常量。
	PrometheusOptionName = "prometheus"
	// RetryOptionName is the name constant for the retry option.
	// RetryOptionName 是重试选项的名称常量。
	RetryOptionName = "retry"
)

// ScheduleOption is the interface for options in chrono.
// ScheduleOption 是 chrono 中选项的接口。
type ScheduleOption interface {
	// Name Returns the name of the option.
	// 返回选项的名称。
	//
	// Returns:
	//	string - The option name / 选项名称
	Name() string
	// Enable Returns whether the option is Enabled.
	// 返回选项是否启用。
	//
	// Returns:
	//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
	Enable() bool
}

// WebMonitorOption represents the web monitor option.
// WebMonitorOption 表示 Web 监控选项。
type WebMonitorOption struct {
	// Enabled indicates whether the web monitor option is Enabled.
	// Enabled 表示是否启用 Web 监控选项。
	Enabled bool
	// Address is the web monitor Address.
	// Address 是 Web 监控地址。
	Address string
}

// Name returns the name of the web monitor option.
// Name 返回 Web 监控选项的名称。
//
// Returns:
//
//	string - The option name / 选项名称
func (w *WebMonitorOption) Name() string {
	return WebMonitorOptionName
}

// Enable returns whether the web monitor option is Enabled.
// Enable 返回 Web 监控选项是否启用。
//
// Returns:
//
//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
func (w *WebMonitorOption) Enable() bool {
	return w.Enabled
}

var _ ScheduleOption = &WebMonitorOption{}

// AliasOption represents the alias option.
// AliasOption 表示别名选项。
type AliasOption struct {
	// Enabled indicates whether the alias option is Enabled.
	// Enabled 表示是否启用别名选项。
	Enabled bool
}

var _ ScheduleOption = &AliasOption{}

// Name returns the name of the alias option.
// Name 返回别名选项的名称。
//
// Returns:
//
//	string - The option name / 选项名称
func (a *AliasOption) Name() string {
	return AliasOptionName
}

// Enable returns whether the alias option is Enabled.
// Enable 返回别名选项是否启用。
//
// Returns:
//
//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
func (a *AliasOption) Enable() bool {
	return a.Enabled
}

// WatchOption represents the watch option.
// WatchOption 表示监听选项。
type WatchOption struct {
	// Enabled indicates whether the watch option is Enabled.
	// Enabled 表示是否启用监听选项。
	Enabled bool
	// WatchFunc is the watch function for monitoring jobs events.
	// WatchFunc 是用于监控任务事件的监听函数。
	// Use interface{} to avoid circular dependency with monitor package
	// 使用 interface{} 避免与 monitor 包的循环依赖
	WatchFunc interface{}
}

var _ ScheduleOption = &WatchOption{}

// Name returns the name of the watch option.
// Name 返回监听选项的名称。
//
// Returns:
//
//	string - The option name / 选项名称
func (w *WatchOption) Name() string {
	return WatchOptionName
}

// Enable returns whether the watch option is Enabled.
// Enable 返回监听选项是否启用。
//
// Returns:
//
//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
func (w *WatchOption) Enable() bool {
	return w.Enabled
}

// TimeoutOption represents the timeout option.
// TimeoutOption 表示超时选项。
type TimeoutOption struct {
	// Enabled indicates whether the timeout option is Enabled.
	// Enabled 表示是否启用超时选项。
	Enabled bool
	// timeout is the timeout duration.
	// timeout 是超时时间。
	timeout time.Duration
}

// LimitOption represents the limit option.
// LimitOption 表示限制选项。
type LimitOption struct {
	// Enabled indicates whether the limit option is Enabled.
	// Enabled 表示是否启用限制选项。
	Enabled bool
	// Limit is the limit Limit.
	// Limit 是限制数量。
	Limit int
}

var _ ScheduleOption = &LimitOption{}

// Name returns the name of the limit option.
// Name 返回限制选项的名称。
//
// Returns:
//
//	string - The option name / 选项名称
func (l *LimitOption) Name() string {
	return LimitOptionName
}

// Enable returns whether the limit option is Enabled.
// Enable 返回限制选项是否启用。
//
// Returns:
//
//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
func (l *LimitOption) Enable() bool {
	return l.Enabled
}

// PrometheusOption represents the prometheus option.
// PrometheusOption 表示 Prometheus 选项。
type PrometheusOption struct {
	// Enabled indicates whether the prometheus option is Enabled.
	// Enabled 表示是否启用 Prometheus 选项。
	Enabled bool
	// Address is the prometheus Address.
	// Address 是 Prometheus 地址。
	Address string
}

var _ ScheduleOption = &PrometheusOption{}

// Name returns the name of the prometheus option.
// Name 返回 Prometheus 选项的名称。
//
// Returns:
//
//	string - The option name / 选项名称
func (p *PrometheusOption) Name() string {
	return PrometheusOptionName
}

// Enable returns whether the prometheus option is Enabled.
// Enable 返回 Prometheus 选项是否启用。
//
// Returns:
//
//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
func (p *PrometheusOption) Enable() bool {
	return p.Enabled
}

// RetryOption represents the retry option.
// RetryOption 表示重试选项。
type RetryOption struct {
	// Enabled indicates whether the retry option is Enabled.
	// Enabled 表示是否启用重试选项。
	Enabled bool
	// Config is the retry configuration.
	// Config 是重试配置。
	Config *retry.RetryConfig
}

var _ ScheduleOption = &RetryOption{}

// Name returns the name of the retry option.
// Name 返回重试选项的名称。
//
// Returns:
//
//	string - The option name / 选项名称
func (r *RetryOption) Name() string {
	return RetryOptionName
}

// Enable returns whether the retry option is Enabled.
// Enable 返回重试选项是否启用。
//
// Returns:
//
//	bool - True if Enabled, false otherwise / 如果启用返回 true，否则返回 false
func (r *RetryOption) Enable() bool {
	return r.Enabled
}
