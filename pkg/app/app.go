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

func (app *Application) Run(ctx context.Context, shutdownChan chan struct{}) chan error {
	slog.Info("starting application")

	errChan := app.startProcesses(ctx)
	go app.waitForShutdown(ctx, shutdownChan, errChan)

	return errChan
}

func (app *Application) startProcesses(ctx context.Context) chan error {
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

func (app *Application) waitForShutdown(ctx context.Context, shutdownChan chan struct{}, errChan chan error) {
	<-ctx.Done()

	slog.Info("application shutdown requested")
	app.wg.Wait()

	close(errChan)
	close(shutdownChan)

	slog.Info("all processes stopped, application shutting down")
}
