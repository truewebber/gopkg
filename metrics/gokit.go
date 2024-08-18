package metrics

import (
	"strconv"
	"time"

	gokitmetrics "github.com/go-kit/kit/metrics"
	gokitprometheus "github.com/go-kit/kit/metrics/prometheus"
	nativeprometheus "github.com/prometheus/client_golang/prometheus"
)

type recorder struct {
	histogram gokitmetrics.Histogram
	name      string
}

const (
	methodLabel   = "method"
	pathLabel     = "path"
	statusCode    = "status_code"
	recorderLabel = "recorder"
)

func NewLatencyRecorder(recorderName string) LatencyRecorder {
	return &recorder{
		name: recorderName,
		histogram: gokitprometheus.NewHistogramFrom(
			nativeprometheus.HistogramOpts{
				Namespace: "truewebber",
				Name:      "request_handling_seconds",
				Buckets:   []float64{.005, .01, .05, .1, .5, 1, 5, 10, 15, 30},
			},
			[]string{recorderLabel, methodLabel, pathLabel, statusCode},
		),
	}
}

func (r *recorder) RecordLatency(labels Labels, start time.Time) {
	r.histogram.With(
		recorderLabel, r.name,
		methodLabel, labels.Method,
		pathLabel, labels.Path,
		statusCode, strconv.Itoa(labels.StatusCode),
	).Observe(time.Since(start).Seconds())
}
