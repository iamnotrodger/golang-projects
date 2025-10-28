package ticket

import (
	"context"
	"log/slog"

	"github.com/iamnotrodger/golang-kafka/pkg/proto/topics"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type ticketStore interface {
	CreateTicket(ctx context.Context, ticket *topics.Ticket) error
}

type Service struct {
	store ticketStore
}

func NewService(store ticketStore) *Service {
	return &Service{
		store: store,
	}
}

func (t *Service) HandleMessage(ctx context.Context, msg kafka.Message) error {
	var ticket topics.Ticket

	if err := proto.Unmarshal(msg.Value, &ticket); err != nil {
		slog.Error(
			"failed to unmarshal ticket from protobuf",
			"error", err.Error(),
			"raw_value", string(msg.Value),
		)
		return err
	}

	slog.Info("creating ticket", "id", ticket.Id)

	if err := t.store.CreateTicket(ctx, &ticket); err != nil {
		slog.Error("failed to create ticket", "error", err.Error())
		return err
	}

	return nil
}
