package healthcheck

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.Match([]string{"GET", "HEAD"}, "/health", h.Ping)
}

func (h *Handler) Ping(ctx *gin.Context) {
	if err := h.service.Ping(ctx); err != nil {
		ctx.JSON(500, gin.H{"status": "unhealthy"})
		return
	}
	ctx.JSON(200, gin.H{"status": "healthy"})
}
