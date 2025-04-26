package chrono

const (
	AliasOptionName = "alias"
	WatchOptionName = "watch"
)

type ChronoOption interface {
	Name() string
	Enable() bool
}

type AliasOption struct {
	enabled bool
}

func (a *AliasOption) Name() string {
	return AliasOptionName
}

func (a *AliasOption) Enable() bool {
	return a.enabled
}

type WatchOption struct {
	enabled bool
}

func (w *WatchOption) Name() string {
	return WatchOptionName
}

func (w *WatchOption) Enable() bool {
	return w.enabled
}
