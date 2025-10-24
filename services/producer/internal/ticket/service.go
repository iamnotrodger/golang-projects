package ticket

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/metrics"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	kafkaWriter *kafka.Writer
}

func NewService(kafkaWriter *kafka.Writer) *Service {
	return &Service{
		kafkaWriter: kafkaWriter,
	}
}

func (t *Service) CreateTicket(ticket *model.Ticket) error {
	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		slog.Error("failed to marshal ticket to json", "error", err)
		return err
	}

	msg := kafka.Message{
		Topic: config.Global.KafkaTicketTopic,
		Key:   []byte(ticket.ID),
		Value: []byte(ticketJSON),
	}

	err = t.kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		slog.Error("failed to write ticket message to kafka", "error", err)
	}

	metrics.RecordTicketCreated("ticket_type")

	return err
}
