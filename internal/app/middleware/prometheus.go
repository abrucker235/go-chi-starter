package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	commonLabels = []string{"path"}
	requests     = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of requests",
		},
		commonLabels,
	)
	duration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			Help:    "HTTP Request Duration",
		},
		commonLabels,
	)
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(writer, request)

		path := request.URL.Path

		requests.WithLabelValues(path).Inc()

		duration.WithLabelValues(path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	})
}
