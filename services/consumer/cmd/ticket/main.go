package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamnotrodger/golang-kafka/pkg/app"
	"github.com/iamnotrodger/golang-kafka/services/consumer/cmd/ticket/appctx"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/config"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	shutdownChan := make(chan struct{})

	appCtx := appctx.NewAppContext(ctx)
	application := app.NewApplication(appctx.BuildAppProcesses(appCtx))
	errChan := application.Run(ctx, shutdownChan)

	exitCode := waitForTermination(terminationContext{
		Context:      ctx,
		cancel:       cancel,
		shutdownChan: shutdownChan,
		errChan:      errChan,
		appCtx:       appCtx,
	})
	return exitCode
}

type terminationContext struct {
	context.Context
	cancel       context.CancelFunc
	shutdownChan chan struct{}
	errChan      <-chan error
	appCtx       *appctx.AppContext
}

func waitForTermination(terminationCtx terminationContext) int {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	exitCode := 0

	select {
	case err, ok := <-terminationCtx.errChan:
		if !ok {
			slog.Error("error channel closed unexpectedly")
			exitCode = 1
		}
		if err != nil {
			slog.Error("application error", "error", err.Error())
			exitCode = 1
		}
	case sig := <-sigs:
		slog.Info("received signal, shutting down", "signal", sig)
	}

	terminationCtx.cancel()
	if err := waitForShutdown(terminationCtx); err != nil {
		exitCode = 1
	}

	return exitCode
}

func waitForShutdown(terminationCtx terminationContext) error {
	slog.Info("waiting for shutdown to complete")

	timer := time.NewTimer(30 * time.Second)
	var err error

	select {
	case <-terminationCtx.shutdownChan:
		slog.Info("application shutdown successful")
	case <-timer.C:
		errMsg := "ungraceful shutdown timeout reached"
		slog.Warn(errMsg)
		err = errors.New(errMsg)

	}

	timer.Stop()
	return err
}

func getLogLevel() slog.Level {
	switch config.Global.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
