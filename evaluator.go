<<<<<<< HEAD
package evaluator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/yourusername/weather-monitor/internal/api"
	"github.com/yourusername/weather-monitor/internal/metrics"
)

var (
	highTemp float64
	lowTemp  float64
)

func init() {
	// default thresholds, can be set through env
	highTemp = 35.0
	lowTemp = 0.0
	if v := os.Getenv("EVALUATOR_THRESHOLD_HIGH_TEMP"); v != "" {
		fmt.Sscanf(v, "%f", &highTemp)
	}
	if v := os.Getenv("EVALUATOR_THRESHOLD_LOW_TEMP"); v != "" {
		fmt.Sscanf(v, "%f", &lowTemp)
	}
}

// Evaluate inspects the Forecast payload and triggers alerts when thresholds met.
// This is a minimal example: it looks for "main.temp" fields in generic JSON.
func Evaluate(data *api.WeatherData) {
	// increment evaluate counter
	metrics.IncEvaluatorRuns()

	// naive: marshal forecast to bytes -> search for temps
	b, _ := json.Marshal(data.Forecast)
	// In production parse concrete fields; here we simply log and check rough strings
	log.Printf("Evaluator processing zip=%s", data.Zip)

	// For demo: write the whole forecast to log and (pretend) generate alerts
	// In real code, parse the forecast JSON into structs and check temps/winds.
	if string(b) != "" {
		// Example logic placeholder:
		// If you had parsed temps, you'd do:
		// if temp >= highTemp -> send alert, metrics.IncAlertHighTemp()
	}
	// Example: increment some metric for demonstration
	metrics.ObserveEvaluatorLatency(0.1)
}
=======
package evaluator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/yourusername/weather-monitor/internal/api"
	"github.com/yourusername/weather-monitor/internal/metrics"
)

var (
	highTemp float64
	lowTemp  float64
)

func init() {
	// default thresholds, can be set through env
	highTemp = 35.0
	lowTemp = 0.0
	if v := os.Getenv("EVALUATOR_THRESHOLD_HIGH_TEMP"); v != "" {
		fmt.Sscanf(v, "%f", &highTemp)
	}
	if v := os.Getenv("EVALUATOR_THRESHOLD_LOW_TEMP"); v != "" {
		fmt.Sscanf(v, "%f", &lowTemp)
	}
}

// Evaluate inspects the Forecast payload and triggers alerts when thresholds met.
// This is a minimal example: it looks for "main.temp" fields in generic JSON.
func Evaluate(data *api.WeatherData) {
	// increment evaluate counter
	metrics.IncEvaluatorRuns()

	// naive: marshal forecast to bytes -> search for temps
	b, _ := json.Marshal(data.Forecast)
	// In production parse concrete fields; here we simply log and check rough strings
	log.Printf("Evaluator processing zip=%s", data.Zip)

	// For demo: write the whole forecast to log and (pretend) generate alerts
	// In real code, parse the forecast JSON into structs and check temps/winds.
	if string(b) != "" {
		// Example logic placeholder:
		// If you had parsed temps, you'd do:
		// if temp >= highTemp -> send alert, metrics.IncAlertHighTemp()
	}
	// Example: increment some metric for demonstration
	metrics.ObserveEvaluatorLatency(0.1)
}
>>>>>>> 428f5cd2f762435819e1f11314a19742522374ff
