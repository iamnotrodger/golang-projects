package ticket

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (t *Service) HandleMessage(ctx context.Context, msg kafka.Message) error {
	return nil
}
