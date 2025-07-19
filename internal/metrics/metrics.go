package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration)
}

// InstrumentHandler は指定したパスの HTTP リクエスト数と処理時間を計測するミドルウェアを返します
func InstrumentHandler(path string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		h.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		// メトリクスに記録
		httpRequests.WithLabelValues(r.Method, path, http.StatusText(rw.status)).Inc()
		httpDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

// responseWriter は WriteHeader 呼び出し時にステータスコードを保持するラッパーです
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader はステータスコードをキャプチャしてから実際に書き込みます
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
