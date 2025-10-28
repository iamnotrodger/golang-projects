package processes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-kafka/pkg/health"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/api"
	"github.com/iamnotrodger/golang-kafka/services/consumer/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpServer struct {
	server *http.Server
}

type HttpServerServices struct {
	HealthService *health.Service
}

func NewHttpServer(services HttpServerServices) *HttpServer {
	engine := gin.New()

	// g.Use(gin.LoggerWithFormatter(logFormatter), gin.Recovery(), gerror.Handler(), location.Default())
	engine.NoRoute(api.NotFound())

	healthHandler := api.NewHealthAPI(services.HealthService)

	engine.Match([]string{"GET", "HEAD"}, "/health", healthHandler.Health)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.Global.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      engine,
	}

	return &HttpServer{server}
}

func (h *HttpServer) Run(ctx context.Context, errChan chan error) {
	go h.start(errChan)
	h.stop(ctx)
}

func (h *HttpServer) start(errChan chan error) {
	errChan <- h.server.ListenAndServe()
}

func (h *HttpServer) stop(ctx context.Context) {
	<-ctx.Done()

	shutdownCtx := context.Background()
	if err := h.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("api server shutdown failed", "error", err.Error())
	} else {
		slog.Info("api server shutdown")
	}
}
