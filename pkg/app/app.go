package app

import (
	"context"
	"log/slog"
	"sync"
)

type Runnable interface {
	Run(context.Context, chan error)
}

type Application struct {
	processes map[string]Runnable
	wg        sync.WaitGroup
}

func NewApplication(processes map[string]Runnable) *Application {
	return &Application{
		processes: processes,
	}
}

func (app *Application) Run(ctx context.Context) chan error {
	slog.Info("starting application")
	errChan := make(chan error)

	for name, process := range app.processes {
		app.wg.Add(1)

		go func(name string, process Runnable) {
			slog.Info("starting process", "name", name)

			defer app.wg.Done()
			process.Run(ctx, errChan)

			slog.Info("process stopped", "name", name)
		}(name, process)
	}

	return errChan
}

func (app *Application) Shutdown() {
	slog.Info("application shutdown requested")
	app.wg.Wait()
	slog.Info("all processes stopped, application shutting down")
}
