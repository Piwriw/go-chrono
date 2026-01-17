// retry/retry_event_test.go
package retry

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestRetryEvent(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(5 * time.Second)
	event := &RetryEvent{
		EventID:         "test-event-1",
		OriginalEventID: "original-1",
		Attempt:         2,
		StartTime:       startTime,
		EndTime:         endTime,
		NextRetryIn:     10 * time.Second,
	}

	t.Run("RetryEvent fields are accessible", func(t *testing.T) {
		if event.EventID != "test-event-1" {
			t.Errorf("EventID = %s, want test-event-1", event.EventID)
		}
		if event.Attempt != 2 {
			t.Errorf("Attempt = %d, want 2", event.Attempt)
		}
	})

	t.Run("GetDuration returns correct duration", func(t *testing.T) {
		duration := event.GetDuration()
		if duration != 5*time.Second {
			t.Errorf("GetDuration() = %v, want 5s", duration)
		}
	})
}

func TestNewRetryEvent(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(3 * time.Second)
	testErr := errors.New("test error")

	event := NewRetryEvent("event-123", "original-456", 1, startTime, endTime, testErr, 15*time.Second)

	if event.EventID != "event-123" {
		t.Errorf("EventID = %s, want event-123", event.EventID)
	}
	if event.OriginalEventID != "original-456" {
		t.Errorf("OriginalEventID = %s, want original-456", event.OriginalEventID)
	}
	if event.Attempt != 1 {
		t.Errorf("Attempt = %d, want 1", event.Attempt)
	}
	if event.Err != testErr {
		t.Errorf("Err = %v, want %v", event.Err, testErr)
	}
	if event.ErrorMessage != "test error" {
		t.Errorf("ErrorMessage = %s, want 'test error'", event.ErrorMessage)
	}
	if event.NextRetryIn != 15*time.Second {
		t.Errorf("NextRetryIn = %v, want 15s", event.NextRetryIn)
	}
}

func TestRetryEventMarshalJSON(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	endTime := startTime.Add(5 * time.Second)
	testErr := errors.New("something went wrong")

	event := NewRetryEvent("event-1", "original-1", 2, startTime, endTime, testErr, 10*time.Second)

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if result["event_id"] != "event-1" {
		t.Errorf("event_id = %v, want event-1", result["event_id"])
	}
	if result["original_event_id"] != "original-1" {
		t.Errorf("original_event_id = %v, want original-1", result["original_event_id"])
	}
	if int(result["attempt"].(float64)) != 2 {
		t.Errorf("attempt = %v, want 2", result["attempt"])
	}
	if result["error"] != "something went wrong" {
		t.Errorf("error = %v, want 'something went wrong'", result["error"])
	}
}
