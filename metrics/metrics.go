package metrics

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	QuoteStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "route_master_quotes_status",
		},
		[]string{"fromToken", "toToken", "status"},
	)
)

func init() {
	prometheus.MustRegister(
		QuoteStatus,
	)
}

func NewMetricServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, nil)
}

func QuoteStatusCounter(fromToken, toToken, status string) {
	QuoteStatus.WithLabelValues(
		fromToken,
		toToken,
		status,
	).Add(1)
}
