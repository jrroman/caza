package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const Namespace = "caza"

var (
	InNetwork = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "in_network_tx",
			Help:      "Number of inter network transactions",
		},
		[]string{"zone"},
	)
	OutNetwork = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "out_network_tx",
			Help:      "Number of inter network transactions",
		},
		[]string{"src_zone", "dst_zone"},
	)
)

func init() {
	// Register prometheus metrics
	prometheus.MustRegister(InNetwork)
	prometheus.MustRegister(OutNetwork)
}

func Serve(addr string) error {
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	return http.ListenAndServe(addr, nil)
}
