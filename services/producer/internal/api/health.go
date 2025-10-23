package api

import (
	"github.com/gin-gonic/gin"
)

type healthService interface {
	Ping() error
}

type HealthAPI struct {
	service healthService
}

func NewHealthAPI(service healthService) *HealthAPI {
	return &HealthAPI{
		service: service,
	}
}

func (h *HealthAPI) Health(ctx *gin.Context) {
	if err := h.service.Ping(); err != nil {
		ctx.JSON(500, gin.H{"status": "unhealthy"})
		return
	}
	ctx.JSON(200, gin.H{"status": "healthy"})
}
