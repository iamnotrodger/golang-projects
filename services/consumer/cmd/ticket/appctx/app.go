package appctx

import (
	"context"

	"github.com/iamnotrodger/golang-kafka/pkg/app"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/health"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/processes"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/ticket"
	"github.com/segmentio/kafka-go"
)

type AppContext struct {
	kafkaReaderConfig *kafka.ReaderConfig
	healthService     *health.Service
	ticketService     *ticket.Service
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	return map[string]app.Runnable{
		"consumer": processes.NewConsumer(appCtx.kafkaReaderConfig, appCtx.ticketService.HandleMessage),
	}
}

func NewAppContext(ctx context.Context) *AppContext {
	appCtx := AppContext{}

	appCtx.kafkaReaderConfig = &kafka.ReaderConfig{
		Brokers: []string{config.Global.KafkaBroker},
		Topic:   config.Global.KafkaTicketTopic,
		GroupID: "ticket-consumer-group",
	}

	appCtx.healthService = health.NewService()
	appCtx.ticketService = ticket.NewService()

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	return nil
}
