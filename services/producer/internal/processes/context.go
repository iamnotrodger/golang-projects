package processes

import (
	"context"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/app"
)

type ApplicationContext struct {
}

func NewApplicationContext(ctx context.Context) *ApplicationContext {
	return &ApplicationContext{}
}

func BuildApplicationProcesses(ctx *ApplicationContext) map[string]app.Runnable {
	return map[string]app.Runnable{
		"API": NewRouter(ctx),
	}
}
