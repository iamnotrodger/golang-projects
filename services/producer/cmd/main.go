package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/app"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/processes"
)

func main() {
	os.Exit(run())
}

func run() int {
	// config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
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

	appCtx := processes.NewApplicationContext(ctx)
	application := app.NewApplication(processes.BuildApplicationProcesses(appCtx))
	errChan := application.Run(ctx, shutdownChan)

	exitCode := waitForTermination(cancel, shutdownChan, errChan)
	return exitCode
}

func waitForTermination(cancel context.CancelFunc, shutdownChan chan struct{}, errChan <-chan error) int {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	exitCode := 0

	select {
	case err, ok := <-errChan:
		if !ok {
			slog.Error("error channel closed unexpectedly")
			exitCode = 1
		}
		if err != nil {
			slog.Error("application error", "error", err)
			exitCode = 1
		}
	case sig := <-sigs:
		slog.Info("received signal, shutting down", "signal", sig)
	}

	cancel()
	<-shutdownChan

	return exitCode
}

// func getLogLevel() slog.Level {
// 	switch config.Global.LogLevel {
// 	case "debug":
// 		return slog.LevelDebug
// 	case "info":
// 		return slog.LevelInfo
// 	case "warn":
// 		return slog.LevelWarn
// 	case "error":
// 		return slog.LevelError
// 	default:
// 		return slog.LevelInfo
// 	}
// }
