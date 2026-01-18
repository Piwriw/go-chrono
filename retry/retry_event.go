// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"encoding/json"
	"errors"
	"time"
)

// RetryEvent contains detailed retry information.
// RetryEvent 重试事件详情
type RetryEvent struct {
	// EventID is the event ID.
	// EventID 事件 ID
	EventID string `json:"event_id"`

	// OriginalEventID is the original task event ID.
	// OriginalEventID 原始任务事件 ID
	OriginalEventID string `json:"original_event_id"`

	// Attempt is the retry attempt number (starts from 0).
	// Attempt 第几次尝试（从 0 开始）
	Attempt int `json:"attempt"`

	// StartTime is the start time.
	// StartTime 开始时间
	StartTime *time.Time `json:"start_time,omitempty"`

	// EndTime is the end time.
	// EndTime 结束时间
	EndTime *time.Time `json:"end_time,omitempty"`

	// Err is the error.
	// Err 错误
	Err error `json:"-"`

	// ErrorMessage is the error message for JSON serialization.
	// ErrorMessage 错误信息（用于 JSON 序列化）
	ErrorMessage string `json:"error,omitempty"`

	// NextRetryIn is the duration until next retry.
	// NextRetryIn 下次重试间隔
	NextRetryIn time.Duration `json:"next_retry_in"`
}

// MarshalJSON serializes to JSON.
// MarshalJSON 序列化为 JSON
//
// Returns:
//
//	[]byte - JSON bytes / JSON 字节数组
//	error - Error if serialization fails / 序列化失败时的错误
func (r *RetryEvent) MarshalJSON() ([]byte, error) {
	type Alias RetryEvent
	errMsg := ""
	if r.Err != nil {
		errMsg = r.Err.Error()
	}
	return json.Marshal(&struct {
		*Alias
		ErrorMessage string `json:"error,omitempty"`
	}{
		Alias:        (*Alias)(r),
		ErrorMessage: errMsg,
	})
}

// GetDuration returns the execution duration.
// GetDuration 获取执行耗时
//
// Returns:
//
//	time.Duration - Execution duration / 执行耗时
//	time.Duration - Zero if start or end time is nil / 如果开始或结束时间为 nil 则返回零值
func (r *RetryEvent) GetDuration() time.Duration {
	if r.StartTime == nil || r.EndTime == nil {
		return 0
	}
	return r.EndTime.Sub(*r.StartTime)
}

// NewRetryEvent creates a new retry event.
// NewRetryEvent 创建新的重试事件
//
// Parameters:
//
//	eventID - Event ID / 事件 ID
//	originalEventID - Original task event ID / 原始任务事件 ID
//	attempt - Retry attempt number (starts from 0) / 重试次数（从 0 开始）
//	startTime - Start time / 开始时间
//	endTime - End time / 结束时间
//	err - Error that occurred / 发生的错误
//	nextRetryIn - Duration until next retry / 下次重试间隔
//
// Returns:
//
//	*RetryEvent - New retry event / 新的重试事件
func NewRetryEvent(eventID, originalEventID string, attempt int, startTime, endTime *time.Time, err error, nextRetryIn time.Duration) *RetryEvent {
	errMsg := ""
	if err != nil && !errors.Is(err, nil) {
		errMsg = err.Error()
	}
	return &RetryEvent{
		EventID:         eventID,
		OriginalEventID: originalEventID,
		Attempt:         attempt,
		StartTime:       startTime,
		EndTime:         endTime,
		Err:             err,
		ErrorMessage:    errMsg,
		NextRetryIn:     nextRetryIn,
	}
}
