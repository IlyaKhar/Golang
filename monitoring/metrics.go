package monitoring

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "day2",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "day2",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration seconds.",
			// Подправь под нагрузку
			Buckets: prometheus.DefBuckets, // [0.005..10]
		},
		[]string{"method", "path", "status"},
	)

	HTTPResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "day2",
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "HTTP response size in bytes.",
			Buckets:   []float64{200, 500, 1_000, 5_000, 10_000, 50_000, 100_000},
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	// Регистрация метрик единожды
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration, HTTPResponseSize)
}

// Handler для /metrics
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// Middleware для записи метрик
type respRec struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *respRec) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
func (r *respRec) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(p)
	r.size += n
	return n, err
}

func WithPrometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &respRec{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}

		labels := prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path, // ("/users/:id")
			"status": http.StatusText(status),
		}

		HTTPRequestsTotal.With(labels).Inc()
		HTTPRequestDuration.With(labels).Observe(time.Since(start).Seconds())
		HTTPResponseSize.With(labels).Observe(float64(rec.size))
	})
}