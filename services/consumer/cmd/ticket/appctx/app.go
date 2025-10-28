package appctx

import (
	"context"
	"log/slog"

	"github.com/iamnotrodger/golang-kafka/pkg/app"
	"github.com/iamnotrodger/golang-kafka/pkg/health"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/healthcheck"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/metrics"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/processes"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/ticket"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

type AppContext struct {
	kafkaReaderConfig *kafka.ReaderConfig
	dbClient          *pgx.Conn
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

	metrics.MustRegister()
	appCtx.initDBClient(ctx)

	ticketStore := ticket.NewStore(appCtx.dbClient)

	appCtx.healthService = health.NewService(map[string]health.HealthCheck{
		"kafka": healthcheck.NewKafkaCheck(),
	})
	appCtx.ticketService = ticket.NewService(ticketStore)

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	slog.Info("shutting down application context")

	err := a.dbClient.Close(ctx)
	if err != nil {
		slog.Error("failed to close database connection", "error", err.Error())
	}

	return err
}

func (a *AppContext) initDBClient(ctx context.Context) {
	var err error
	a.dbClient, err = pgx.Connect(ctx, config.Global.Secret.DatabaseURL)

	if err != nil {
		slog.Error("failed to connect to database", "error", err.Error())
		panic(err)
	}
}
