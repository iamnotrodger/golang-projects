package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/iamnotrodger/golang-kafka/producer/router"
)

type Runnable interface {
	Run(context.Context, chan error)
}

type Application struct {
	processes map[string]Runnable
	wg        sync.WaitGroup
}

func BuildApplicationProcesses(appCtx *ApplicationContext) map[string]Runnable {
	return map[string]Runnable{
		"Router": router.NewRouter(),
	}
}

func NewApplication(processes map[string]Runnable) *Application {
	return &Application{
		processes: processes,
	}
}

func (app *Application) Run(ctx context.Context, shutdownChan chan struct{}) chan error {
	slog.Info("starting application")

	processesCtx, cancelProcesses := context.WithCancel(ctx)
	errChan := app.startProcesses(processesCtx)

	go app.waitForShutdown(shutdownContext{
		Context:         ctx,
		cancelProcesses: cancelProcesses,
		shutdownChan:    shutdownChan,
		errChan:         errChan,
	})

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

type shutdownContext struct {
	context.Context
	cancelProcesses context.CancelFunc
	shutdownChan    chan struct{}
	errChan         chan error
}

func (app *Application) waitForShutdown(shutdown shutdownContext) {
	<-shutdown.Done()
	slog.Info("application shutdown requested")

	shutdown.cancelProcesses()
	app.wg.Wait()
	close(shutdown.errChan)
	close(shutdown.shutdownChan)

	slog.Info("all processes stopped, application shutting down")
}
