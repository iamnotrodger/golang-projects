package health

import "github.com/segmentio/kafka-go"

type Service struct {
	kafkaWriter *kafka.Writer
}

func NewService(kafkaWriter *kafka.Writer) *Service {
	return &Service{
		kafkaWriter: kafkaWriter,
	}
}

func (h *Service) PingKafka() error {
	return nil
}
