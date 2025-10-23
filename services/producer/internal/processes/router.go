package processes

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/api"
)

type Router struct {
	server *http.Server
}

func NewRouter(appContext *ApplicationContext) *Router {
	engine := gin.New()

	// g.Use(gin.LoggerWithFormatter(logFormatter), gin.Recovery(), gerror.Handler(), location.Default())
	engine.NoRoute(api.NotFound())

	healthHandler := &api.HealthAPI{}
	ticketHandler := &api.TicketAPI{}

	engine.Match([]string{"GET", "HEAD"}, "/health", healthHandler.Health)

	ticket := engine.Group("/ticket")
	{
		ticket.POST("/", ticketHandler.CreateTicket)
	}

	return &Router{server: &http.Server{Addr: "localhost:8080", Handler: engine}}
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
		slog.Error("api server shutdown failed", "error", err)
	} else {
		slog.Info("api server shutdown")
	}
}
