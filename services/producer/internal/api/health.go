package api

import (
	"github.com/gin-gonic/gin"
)

type HealthService interface {
	Ping() error
}

type HealthAPI struct {
	Service HealthService
}

func (h *HealthAPI) Health(ctx *gin.Context) {
	if err := h.Service.Ping(); err != nil {
		ctx.JSON(500, gin.H{"status": "unhealthy"})
		return
	}
	ctx.JSON(200, gin.H{"status": "healthy"})
}
