package writer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type MockKafkaWriter struct {
	Messages      []kafka.Message
	WriteErr      error
	ShouldCapture bool
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.WriteErr != nil {
		return m.WriteErr
	}
	if m.ShouldCapture {
		m.Messages = append(m.Messages, msgs...)
	}
	return nil
}

func (m *MockKafkaWriter) Close() error {
	return nil
}
