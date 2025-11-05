package appctx

import (
	"context"

	"github.com/iamnotrodger/golang-projects/pkg/app"
	"github.com/iamnotrodger/golang-projects/pkg/health"
	"github.com/iamnotrodger/golang-projects/services/producer/internal/healthcheck"
	"github.com/iamnotrodger/golang-projects/services/producer/internal/metrics"
	"github.com/iamnotrodger/golang-projects/services/producer/internal/processes"
	"github.com/iamnotrodger/golang-projects/services/producer/internal/ticket"
)

type AppContext struct {
	ticketService *ticket.Service
	healthService *health.Service
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	return map[string]app.Runnable{
		"http": processes.NewHttpServer(processes.HttpServerServices{
			HealthService: appCtx.healthService,
			TicketService: appCtx.ticketService,
		}),
	}
}

func NewAppContext(ctx context.Context) *AppContext {
	appCtx := AppContext{}

	metrics.MustRegister()

	appCtx.ticketService = ticket.NewService()
	appCtx.healthService = health.NewService(map[string]health.HealthCheck{
		"kafka": healthcheck.NewKafkaCheck(),
	})

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	return nil
}
