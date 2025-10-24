package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "weather_api_calls_total",
		Help: "Total number of API calls to OpenWeatherMap",
	})

	apiLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "weather_api_latency_seconds",
		Help:    "API latency in seconds",
		Buckets: prometheus.DefBuckets,
	})

	kafkaPublishSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_publish_success_total",
		Help: "Successful Kafka publishes",
	})

	kafkaPublishFail = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_publish_fail_total",
		Help: "Failed Kafka publishes",
	})

	kafkaConsume = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_consume_total",
		Help: "Kafka consumer messages processed",
	})

	evaluatorRuns = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "evaluator_runs_total",
		Help: "Total evaluator runs",
	})

	evaluatorLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "evaluator_latency_seconds",
		Help:    "Evaluator run latency in seconds",
		Buckets: prometheus.DefBuckets,
	})
)

func InitMetrics() {
	prometheus.MustRegister(
		apiCalls,
		apiLatency,
		kafkaPublishSuccess,
		kafkaPublishFail,
		kafkaConsume,
		evaluatorRuns,
		evaluatorLatency,
	)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func IncAPICalls()                      { apiCalls.Inc() }
func ObserveAPILatency(v float64)       { apiLatency.Observe(v) }
func IncKafkaPublishSuccess()           { kafkaPublishSuccess.Inc() }
func IncKafkaPublishFailures()          { kafkaPublishFail.Inc() }
func IncKafkaConsume()                  { kafkaConsume.Inc() }
func IncEvaluatorRuns()                 { evaluatorRuns.Inc() }
func ObserveEvaluatorLatency(v float64) { evaluatorLatency.Observe(v) }
