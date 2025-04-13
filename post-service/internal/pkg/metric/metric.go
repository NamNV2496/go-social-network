package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

var (
	totalHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "post_service_hits_total",
			Help: "Total number of requests served by the application.",
		},
	)
	hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "post_service_hits",
			Help: "Number of hits to the post service by status, method, and path",
		},
		[]string{"status", "method", "path"},
	)
	times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "post_service_response_times",
			Help:    "Response times of the post service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "method", "path"},
	)
	errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "post_service_errors_total",
			Help: "Total number of error requests processed by the MyApp web server.",
		},
		[]string{"path", "status"},
	)
	newPostCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "post_service_new_post_rate",
			Help: "Total number of new post requests",
		},
		[]string{"name"},
	)
	getPostCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "post_service_get_post_rate",
			Help: "Total number of get post request",
		},
		[]string{"name"},
	)
)

func InitPrometheus() {
	// Register metrics with Prometheus
	prometheus.MustRegister(totalHits)
	prometheus.MustRegister(hits)
	prometheus.MustRegister(times)
	prometheus.MustRegister(errorCount)
	prometheus.MustRegister(newPostCnt)
	prometheus.MustRegister(getPostCnt)
}

func MetricIncHits(status string, method, path string) {
	totalHits.Inc()
	hits.WithLabelValues(status, method, path).Inc()
}

func MetricIncError(status string, method, path string) {
	errorCount.WithLabelValues(status, method, path).Inc()
}

func MetricObserveTime(status string, method, path string, observeTime float64) {
	times.WithLabelValues(status, method, path).Observe(observeTime)
}

func MetricNewPostCnt(name string) {
	newPostCnt.WithLabelValues(name)

}

func MetricGetPostCnt(name string) {
	getPostCnt.WithLabelValues(name).Inc()
}
