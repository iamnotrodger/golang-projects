package metrics

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	TicketsCreatedCounter *prometheus.CounterVec
	ErrorCounter          *prometheus.CounterVec
}

var metric = metrics{
	TicketsCreatedCounter: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "producer_tickets_created_total",
			Help: "total number of created tickets",
		},
		[]string{"type"},
	),
	ErrorCounter: prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "producer_error_total",
			Help: "total number of errors",
		},
		[]string{"type"},
	),
}

func MustRegister() {
	prometheus.MustRegister(metric.TicketsCreatedCounter)
	prometheus.MustRegister(metric.ErrorCounter)
	initCounters()
}

func initCounters() {
	metric.TicketsCreatedCounter.WithLabelValues("REPLACE_ME").Add(0)
	metric.ErrorCounter.WithLabelValues("REPLACE_ME").Add(0)
}
