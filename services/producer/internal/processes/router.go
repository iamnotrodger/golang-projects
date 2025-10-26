package processes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/api"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	server *http.Server
}

func NewRouter(appCtx *AppContext) *Router {
	engine := gin.New()
	metrics.MustRegister()

	// g.Use(gin.LoggerWithFormatter(logFormatter), gin.Recovery(), gerror.Handler(), location.Default())
	engine.NoRoute(api.NotFound())

	healthHandler := api.NewHealthAPI(appCtx.healthService)
	ticketHandler := api.NewTicketAPI(appCtx.ticketService)

	engine.Match([]string{"GET", "HEAD"}, "/health", healthHandler.Health)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	ticket := engine.Group("/ticket")
	{
		ticket.POST("/", ticketHandler.CreateTicket)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.Global.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      engine,
	}

	return &Router{server}
}

func (r *Router) Run(ctx context.Context, errChan chan error) {
	go r.start(errChan)
	r.stop(ctx)
}

func (r *Router) start(errChan chan error) {
	errChan <- r.server.ListenAndServe()
}

func (r *Router) stop(ctx context.Context) {
	<-ctx.Done()

	shutdownCtx := context.Background()
	if err := r.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("api server shutdown failed", "error", err.Error())
	} else {
		slog.Info("api server shutdown")
	}
}
