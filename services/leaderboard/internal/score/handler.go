package score

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	score := r.Group("/score")
	{
		score.POST("/", h.CreateScore)
	}
}

func (h *Handler) CreateScore(c *gin.Context) {
}
