package metrics

func RecordTicketCreated() {
	metric.TicketsCreatedCounter.WithLabelValues().Inc()
}
