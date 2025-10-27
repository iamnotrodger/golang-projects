package processes

import (
	"context"

	"github.com/iamnotrodger/golang-kafka/pkg/app"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/health"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/ticket"
)

type AppContext struct {
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

	appCtx.ticketService = ticket.NewService()
	appCtx.healthService = health.NewService()

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	return nil
}
