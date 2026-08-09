package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

// MetricsMiddleware records Prometheus counters/histograms for every request.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(time.Since(start).Seconds())
	}
}
