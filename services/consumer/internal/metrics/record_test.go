package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func mapContainsAll(expected, actual prometheus.Labels) bool {
	for k, v := range expected {
		if val, ok := actual[k]; !ok || val != v {
			return false
		}
	}
	return true
}

func assertCounterResults(t *testing.T, collector prometheus.Collector, name string, value float64, labels prometheus.Labels) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	metricFamilies, err := reg.Gather()
	require.NoError(t, err)
	require.Len(t, metricFamilies, 1, "should have exactly one metric family")

	metricFamily := metricFamilies[0]
	require.NotNil(t, metricFamily)

	var counterMetric *dto.Metric
	for _, m := range metricFamily.Metric {
		metricLabels := make(map[string]string)
		for _, label := range m.GetLabel() {
			metricLabels[label.GetName()] = label.GetValue()
		}

		if mapContainsAll(labels, metricLabels) {
			counterMetric = m
			break
		}
	}

	require.NotNil(t, counterMetric, "no metric found with matching labels")
	require.Equal(t, name, metricFamily.GetName(), "metric name should match")
	require.Equal(t, value, counterMetric.GetCounter().GetValue(), "metric value should match")
}

func TestRecordTicketCreated(t *testing.T) {
	metric.TicketsCreatedCounter.Reset()
	RecordTicketCreated("ticketType")
	assertCounterResults(t, metric.TicketsCreatedCounter, "producer_tickets_created_total", 1, prometheus.Labels{"type": "ticketType"})
}

func TestRecordError(t *testing.T) {
	metric.ErrorCounter.Reset()
	RecordError("errorType")
	assertCounterResults(t, metric.ErrorCounter, "producer_error_total", 1, prometheus.Labels{"type": "errorType"})
}
