// Package retry provides retry mechanism for scheduled jobs.
// 包 retry 提供调度任务的重试机制。
package retry

import (
	"encoding/json"
	"errors"
	"time"
)

// RetryEvent 重试事件详情
// RetryEvent contains detailed retry information
type RetryEvent struct {
	// EventID 事件 ID
	// EventID is the event ID
	EventID string `json:"event_id"`

	// OriginalEventID 原始任务事件 ID
	// OriginalEventID is the original task event ID
	OriginalEventID string `json:"original_event_id"`

	// Attempt 第几次尝试（从 0 开始）
	// Attempt is the retry attempt number (starts from 0)
	Attempt int `json:"attempt"`

	// StartTime 开始时间
	// StartTime is the start time
	StartTime time.Time `json:"start_time"`

	// EndTime 结束时间
	// EndTime is the end time
	EndTime time.Time `json:"end_time"`

	// Err 错误
	// Err is the error
	Err error `json:"-"`

	// ErrorMessage 错误信息（用于 JSON 序列化）
	// ErrorMessage is the error message for JSON serialization
	ErrorMessage string `json:"error,omitempty"`

	// NextRetryIn 下次重试间隔
	// NextRetryIn is the duration until next retry
	NextRetryIn time.Duration `json:"next_retry_in"`
}

// MarshalJSON 序列化为 JSON
// MarshalJSON serializes to JSON
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

// GetDuration 获取执行耗时
// GetDuration returns the execution duration
func (r *RetryEvent) GetDuration() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

// NewRetryEvent 创建新的重试事件
// NewRetryEvent creates a new retry event
func NewRetryEvent(eventID, originalEventID string, attempt int, startTime, endTime time.Time, err error, nextRetryIn time.Duration) *RetryEvent {
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
