package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMustRegister(t *testing.T) {
	MustRegister()

	assert.True(t, prometheus.Unregister(metric.TicketsCreatedCounter))
}
