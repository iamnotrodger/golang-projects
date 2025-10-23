package metrics

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	TicketsCreatedCounter *prometheus.CounterVec
}

var metric = metrics{
	TicketsCreatedCounter: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "producer_tickets_created_total",
			Help: "total number of created tickets",
		},
		[]string{},
	),
}

func MustRegister() {
	prometheus.MustRegister(metric.TicketsCreatedCounter)
	initCounters()
}

func initCounters() {
	metric.TicketsCreatedCounter.WithLabelValues().Add(0)
}
