package processes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/pkg/health"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/config"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/leaderboard"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/score"
)

type HttpServer struct {
	server *http.Server
}

type HttpServerServices struct {
	HealthService *health.Service
	ScoreService  *score.Service
}

func NewHttpServer(engine *gin.Engine, hub *leaderboard.Hub, services HttpServerServices) *HttpServer {
	scoreHandler := score.NewHandler(services.ScoreService)
	scoreHandler.RegisterRoutes(engine)

	leaderboardHandler := leaderboard.NewHandler(services.ScoreService, hub)
	leaderboardHandler.RegisterRoutes(engine)

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.Global.Port),
		WriteTimeout: 0,
		ReadTimeout:  15 * time.Second,
		Handler:      engine,
	}

	return &HttpServer{
		server: server,
	}
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
