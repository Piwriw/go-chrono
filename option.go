package chrono

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
	// Enable Returns whether the option is enabled.
	// 返回选项是否启用。
	//
	// Returns:
	//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
	Enable() bool
}

// WebMonitorOption represents the web monitor option.
// WebMonitorOption 表示 Web 监控选项。
type WebMonitorOption struct {
	// enabled indicates whether the web monitor option is enabled.
	// enabled 表示是否启用 Web 监控选项。
	enabled bool
	// address is the web monitor address.
	// address 是 Web 监控地址。
	address string
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

// Enable returns whether the web monitor option is enabled.
// Enable 返回 Web 监控选项是否启用。
//
// Returns:
//
//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
func (w *WebMonitorOption) Enable() bool {
	return w.enabled
}

// Address returns the web monitor address.
// Address 返回 Web 监控地址。
//
// Returns:
//
//	string - The web monitor address / Web 监控地址
func (w *WebMonitorOption) Address() string {
	return w.address
}

var _ ScheduleOption = &WebMonitorOption{}

// AliasOption represents the alias option.
// AliasOption 表示别名选项。
type AliasOption struct {
	// enabled indicates whether the alias option is enabled.
	// enabled 表示是否启用别名选项。
	enabled bool
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

// Enable returns whether the alias option is enabled.
// Enable 返回别名选项是否启用。
//
// Returns:
//
//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
func (a *AliasOption) Enable() bool {
	return a.enabled
}

// WatchOption represents the watch option.
// WatchOption 表示监听选项。
type WatchOption struct {
	// enabled indicates whether the watch option is enabled.
	// enabled 表示是否启用监听选项。
	enabled bool
	// watchFunc is the watch function for monitoring job events.
	// watchFunc 是用于监控任务事件的监听函数。
	watchFunc func(event JobWatchInterface)
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

// Enable returns whether the watch option is enabled.
// Enable 返回监听选项是否启用。
//
// Returns:
//
//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
func (w *WatchOption) Enable() bool {
	return w.enabled
}

// WatchFunc returns the watch function.
// WatchFunc 返回监听函数。
//
// Returns:
//
//	func(event JobWatchInterface) - The watch function / 监听函数
func (w *WatchOption) WatchFunc() func(event JobWatchInterface) {
	return w.watchFunc
}

// TimeoutOption represents the timeout option.
// TimeoutOption 表示超时选项。
type TimeoutOption struct {
	// enabled indicates whether the timeout option is enabled.
	// enabled 表示是否启用超时选项。
	enabled bool
	// timeout is the timeout duration.
	// timeout 是超时时间。
	timeout time.Duration
}

// LimitOption represents the limit option.
// LimitOption 表示限制选项。
type LimitOption struct {
	// enabled indicates whether the limit option is enabled.
	// enabled 表示是否启用限制选项。
	enabled bool
	// number is the limit number.
	// number 是限制数量。
	number int
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

// Enable returns whether the limit option is enabled.
// Enable 返回限制选项是否启用。
//
// Returns:
//
//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
func (l *LimitOption) Enable() bool {
	return l.enabled
}

// Limit returns the limit number.
// Limit 返回限制数量。
//
// Returns:
//
//	int - The limit number / 限制数量
func (l *LimitOption) Limit() int {
	return l.number
}

// PrometheusOption represents the prometheus option.
// PrometheusOption 表示 Prometheus 选项。
type PrometheusOption struct {
	// enabled indicates whether the prometheus option is enabled.
	// enabled 表示是否启用 Prometheus 选项。
	enabled bool
	// address is the prometheus address.
	// address 是 Prometheus 地址。
	address string
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

// Enable returns whether the prometheus option is enabled.
// Enable 返回 Prometheus 选项是否启用。
//
// Returns:
//
//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
func (p *PrometheusOption) Enable() bool {
	return p.enabled
}

// Address returns the prometheus address.
// Address 返回 Prometheus 地址。
//
// Returns:
//
//	string - The prometheus address / Prometheus 地址
func (p *PrometheusOption) Address() string {
	return p.address
}

// RetryOption represents the retry option.
// RetryOption 表示重试选项。
type RetryOption struct {
	// enabled indicates whether the retry option is enabled.
	// enabled 表示是否启用重试选项。
	enabled bool
	// config is the retry configuration.
	// config 是重试配置。
	config *retry.RetryConfig
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

// Enable returns whether the retry option is enabled.
// Enable 返回重试选项是否启用。
//
// Returns:
//
//	bool - True if enabled, false otherwise / 如果启用返回 true，否则返回 false
func (r *RetryOption) Enable() bool {
	return r.enabled
}

// Config returns the retry configuration.
// Config 返回重试配置。
//
// Returns:
//
//	*retry.RetryConfig - The retry configuration / 重试配置
func (r *RetryOption) Config() *retry.RetryConfig {
	return r.config
}
