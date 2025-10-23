package processes

import (
	"context"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/app"
)

type AppContext struct {
}

func NewAppContext(ctx context.Context) *AppContext {
	return &AppContext{}
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	return map[string]app.Runnable{
		"API": NewRouter(appCtx),
	}
}
