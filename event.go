package statsy

import (
	"time"
)

type Event struct {
	Stream string
	Weight float32
	Time   time.Time
}

type JsonEvent struct {
	Stream string  `json:"stream"`
	Weight float32 `json:"weight,omitempty"`
	Time   int64   `json:"timestamp,omitempty"`
}

func (e *Event) JsonEvent() *JsonEvent {
	jsonEvent := &JsonEvent{
		Stream: e.Stream,
		Weight: e.Weight,
	}

	timestamp := e.Time.Unix()
	if timestamp > 0 {
		jsonEvent.Time = timestamp
	}

	return jsonEvent
}
