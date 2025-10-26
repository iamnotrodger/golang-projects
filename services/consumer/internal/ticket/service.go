package ticket

import (
	"context"
	"log/slog"

	"github.com/iamnotrodger/golang-kafka/pkg/proto/topics"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (t *Service) HandleMessage(ctx context.Context, msg kafka.Message) error {
	var ticket topics.Ticket

	if err := proto.Unmarshal(msg.Value, &ticket); err != nil {
		slog.Error("failed to unmarshal ticket from protobuf",
			"error", err.Error(),
			"raw_value", string(msg.Value))
		return err
	}

	slog.Info(
		"successfully parsed ticket",
		"id", ticket.Id,
		"title", ticket.Title,
		"price", ticket.Price,
	)

	return nil
}
