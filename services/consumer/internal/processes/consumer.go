package processes

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageHandler func(context.Context, kafka.Message) error

type Consumer struct {
	reader  *kafka.Reader
	handler KafkaMessageHandler
}

func NewConsumer(config *kafka.ReaderConfig, handler KafkaMessageHandler) *Consumer {
	return &Consumer{
		reader:  kafka.NewReader(*config),
		handler: handler,
	}
}

func (c *Consumer) Run(ctx context.Context, errChan chan error) {
	go c.start(ctx, errChan)
	c.stop(ctx)
}

func (c *Consumer) start(ctx context.Context, errChan chan error) {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			errChan <- err
			return
		}

		slog.Info(
			"received message",
			"topic", message.Topic,
			"partition", message.Partition,
			"offset", message.Offset,
			"key", string(message.Key),
			"value_length", len(message.Value),
		)

		if err := c.handler(ctx, message); err != nil {
			slog.Error("failed to handle kafka message", "error", err.Error(), "topic", message.Topic)
			continue
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			slog.Error("failed to commit kafka message", "error", err.Error(), "topic", message.Topic)
		}
	}
}

func (c *Consumer) stop(ctx context.Context) {
	<-ctx.Done()

	if err := c.reader.Close(); err != nil {
		slog.Error("consumer shutdown failed", "error", err)
	} else {
		slog.Info("consumer shutdown complete")
	}
}
