package ticket

import (
	"context"
	"log/slog"

	"github.com/iamnotrodger/golang-kafka/pkg/proto/topics"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/metrics"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/writer"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	kafkaWriter writer.KafkaWriter
}

func NewService() *Service {
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(config.Global.KafkaBroker),
		Topic:    config.Global.KafkaTicketTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Service{
		kafkaWriter: writer.NewKafkaWriterAdapter(kafkaWriter),
	}
}

func (t *Service) CreateTicket(ticket *model.Ticket) error {
	protoTicket := &topics.Ticket{
		Id:        ticket.ID,
		Title:     ticket.Title,
		Price:     float32(ticket.Price),
		CreatedAt: timestamppb.New(ticket.CreatedAt),
	}

	ticketBytes, err := proto.Marshal(protoTicket)
	if err != nil {
		slog.Error("failed to marshal ticket to protobuf", "error", err.Error())
		return err
	}

	msg := kafka.Message{
		Topic: config.Global.KafkaTicketTopic,
		Key:   []byte(ticket.ID),
		Value: ticketBytes,
	}

	err = t.kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		slog.Error("failed to write ticket message to kafka", "error", err.Error())
		return err
	}

	metrics.RecordTicketCreated("ticket_type")
	return nil
}
