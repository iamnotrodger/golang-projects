package healthcheck

import (
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaCheck struct{}

func NewKafkaCheck() *KafkaCheck {
	return &KafkaCheck{}
}

func (k *KafkaCheck) Ping() error {
	conn, err := kafka.Dial("tcp", config.Global.KafkaBroker)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
