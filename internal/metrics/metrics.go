package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CreateOrderDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "create_order_duration_nanoseconds",
			Help:    "Time spent creating order",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(CreateOrderDuration)
}
