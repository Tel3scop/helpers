package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metrics

// Init инициализация метрик. subsystem - например, grpc,  bucketsStart 0.0001 - 1 сек, bucketsFactor - шаг, bucketsCount - количество
func Init(namespace, appName, subsystem string, bucketsStart, bucketsFactor float64, bucketsCount int) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_requests_total",
				Help:      "Количество запросов к серверу",
			},
		),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_responses_total",
				Help:      "Количество ответов от сервера",
			},
			[]string{"status", "method"},
		),
		histogramResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_histogram_response_time_seconds",
				Help:      "Время ответа от сервера",
				Buckets:   prometheus.ExponentialBuckets(bucketsStart, bucketsFactor, bucketsCount),
			},
			[]string{"status"},
		),
	}

	return nil
}

// IncRequestCounter инкрементация количества запросов
func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

// IncResponseCounter инкрементация количества ответов
func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}

// HistogramResponseTimeObserve гистограмма для контроля времени ответа
func HistogramResponseTimeObserve(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
