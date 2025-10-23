package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func assertCounterResults(t *testing.T, collector prometheus.Collector, name string, value float64) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricFamilies, err := reg.Gather()
	require.NoError(t, err)
	require.Len(t, metricFamilies, 1, "should have exactly one metric family")

	metricFamily := metricFamilies[0]
	require.NotNil(t, metricFamily)
	require.Greater(t, len(metricFamily.Metric), 0, "should have at least one metric")

	counterMetric := metricFamily.Metric[0]
	require.NotNil(t, counterMetric)

	require.Equal(t, name, metricFamily.GetName(), "metric name should match")
	require.Equal(t, value, counterMetric.GetCounter().GetValue(), "metric value should match")
}

func TestRecordTicketCreated(t *testing.T) {
	RecordTicketCreated()
	assertCounterResults(t, metric.TicketsCreatedCounter, "producer_tickets_created_total", 1)
}
