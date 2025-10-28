package metrics

func RecordTicketCreated(ticketType string) {
	metric.TicketsCreatedCounter.WithLabelValues(ticketType).Inc()
}

func RecordError(errorType string) {
	metric.ErrorCounter.WithLabelValues(errorType).Inc()
}
