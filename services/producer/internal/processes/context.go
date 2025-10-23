package processes

import (
	"context"
	"log/slog"

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
	appCtx.healthService = health.NewService()

	return &appCtx
}

func (a *AppContext) initKafkaWriter() {
	a.kafkaWriter = &kafka.Writer{
		Addr:                   kafka.TCP(config.Global.KafkaBroker),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	if config.Global.Env == "development" {
		conn, err := kafka.DialLeader(context.Background(), "tcp", config.Global.KafkaBroker, config.Global.KafkaTicketTopic, 0)
		if err != nil {
			slog.Error("failed to dial kafka leader and create topic", "error", err)
		}
		defer conn.Close()
	}
}
