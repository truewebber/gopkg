package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsServer(hostPort string) *http.Server {
	const readHeaderTimeout = 3 * time.Second

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:              hostPort,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           mux,
	}
}
