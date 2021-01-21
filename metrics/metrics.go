package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// LatencyHandler records basic stats about Shorty, which can be viewed at the /metrics handler
	LatencyHandler = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "shorty",
		Name:       "latency_of_handlers",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"handler_name"})
)

func init() {
	prometheus.MustRegister(LatencyHandler)
}
