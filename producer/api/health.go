package api

import (
	"github.com/gin-gonic/gin"
)

type HealthDatabase interface {
	Ping() error
}

type HealthAPI struct {
	db HealthDatabase
}

func (h *HealthAPI) Health(ctx *gin.Context) {
	if err := h.db.Ping(); err != nil {
		ctx.JSON(500, gin.H{"status": "unhealthy"})
		return
	}
	ctx.JSON(200, gin.H{"status": "healthy"})
}
