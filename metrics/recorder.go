package metrics

import (
	"time"
)

type Labels struct {
	Method, Path string
	StatusCode   int
}

type LatencyRecorder interface {
	RecordLatency(labels Labels, start time.Time)
}
