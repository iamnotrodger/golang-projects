package writer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type kafkaWriterAdapter struct {
	writer *kafka.Writer
}

func (k *kafkaWriterAdapter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return k.writer.WriteMessages(ctx, msgs...)
}

func (k *kafkaWriterAdapter) Close() error {
	return k.writer.Close()
}

func NewKafkaWriterAdapter(writer *kafka.Writer) KafkaWriter {
	return &kafkaWriterAdapter{writer: writer}
}
