package ticket

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
	"github.com/segmentio/kafka-go"
)

type TicketService struct {
	kafkaWriter *kafka.Writer
}

func NewTicketService(kafkaWriter *kafka.Writer) *TicketService {
	return &TicketService{
		kafkaWriter: kafkaWriter,
	}
}

func (t *TicketService) CreateTicket(ticket *model.Ticket) error {
	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		slog.Error("failed to marshal ticket to json", "error", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(ticket.ID),
		Value: []byte(ticketJSON),
	}

	err = t.kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		slog.Error("failed to write ticket message to kafka", "error", err)
	}

	return err
}
