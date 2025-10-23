package processes

import (
	"context"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/app"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/health"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/ticket"
	"github.com/segmentio/kafka-go"
)

type AppContext struct {
	kafkaWriter   *kafka.Writer
	ticketService *ticket.Service
	healthService *health.Service
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	return map[string]app.Runnable{
		"API": NewRouter(appCtx),
	}
}

func NewAppContext(ctx context.Context) *AppContext {
	appCtx := AppContext{}
	appCtx.initKafkaWriter()

	appCtx.ticketService = ticket.NewService(appCtx.kafkaWriter)
	appCtx.healthService = health.NewService(appCtx.kafkaWriter)

	return &appCtx
}

func (a *AppContext) initKafkaWriter() {
	a.kafkaWriter = &kafka.Writer{
		Addr:                   kafka.TCP(config.Global.KafkaBroker),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}
