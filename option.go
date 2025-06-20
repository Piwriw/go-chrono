package chrono

import "time"

const (
	AliasOptionName      = "alias"
	WatchOptionName      = "watch"
	TimoutOptionName     = "timeout"
	WebMonitorOptionName = "web_monitor"
)

// ChronoOption is the interface for options in chrono.
// ChronoOption 是 chrono 中选项的接口。
type ChronoOption interface {
	Name() string // Returns the name of the option
	// 返回选项名称
	Enable() bool // Returns whether the option is enabled
	// 返回选项是否启用
}

// WebMonitorOption represents the web monitor option.
// WebMonitorOption 表示 Web 监控选项。
type WebMonitorOption struct {
	// 是否启用 Web 监控选项
	// Whether the web monitor option is enabled
	enabled bool
	// Web 监控地址
	// Web monitor address
	address string
}

func (w *WebMonitorOption) Name() string {
	return WebMonitorOptionName
}

func (w *WebMonitorOption) Enable() bool {
	return w.enabled
}

func (w *WebMonitorOption) Address() string {
	return w.address
}

var _ ChronoOption = &WebMonitorOption{}

// AliasOption represents the alias option.
// AliasOption 表示别名选项。
type AliasOption struct {
	enabled bool // Whether the alias option is enabled
	// 是否启用别名选项
}

var _ ChronoOption = &AliasOption{}

// Name returns the name of the alias option.
// Name 返回别名选项的名称。
func (a *AliasOption) Name() string {
	return AliasOptionName
}

// Enable returns whether the alias option is enabled.
// Enable 返回别名选项是否启用。
func (a *AliasOption) Enable() bool {
	return a.enabled
}

// WatchOption represents the watch option.
// WatchOption 表示监听选项。
type WatchOption struct {
	// 是否启用监听选项
	enabled bool // Whether the watch option is enabled
	// 监听函数
	watchFunc func(event JobWatchInterface)
}

var _ ChronoOption = &WatchOption{}

// Name returns the name of the watch option.
// Name 返回监听选项的名称。
func (w *WatchOption) Name() string {
	return WatchOptionName
}

// Enable returns whether the watch option is enabled.
// Enable 返回监听选项是否启用。
func (w *WatchOption) Enable() bool {
	return w.enabled
}

func (w *WatchOption) WatchFunc() func(event JobWatchInterface) {
	return w.watchFunc
}

// TimeoutOption represents the timeout option.
// TimeoutOption 表示超时选项。
type TimeoutOption struct {
	// 是否启用超时选项
	enabled bool // Whether the timeout option is enabled
	// 超时时间
	timeout time.Duration
}

var _ ChronoOption = &TimeoutOption{}

func (t *TimeoutOption) Name() string {
	return TimoutOptionName
}

func (t *TimeoutOption) Enable() bool {
	return t.enabled
}

func (t *TimeoutOption) Timeout() time.Duration {
	return t.timeout
}
