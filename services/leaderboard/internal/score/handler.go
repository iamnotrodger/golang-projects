package score

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
)

type scoreService interface {
	SaveScore(ctx context.Context, score *model.Score) error
	PublishTopScores(ctx context.Context) error
}

type Handler struct {
	service scoreService
}

func NewHandler(service scoreService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	score := r.Group("/score")
	{
		score.POST("/", h.saveScore)
	}
}

func (h *Handler) saveScore(c *gin.Context) {
	var score model.Score
	if err := c.ShouldBindJSON(&score); err != nil {
		slog.Error("saveScore error parsing request", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SaveScore(c.Request.Context(), &score); err != nil {
		slog.Error("saveScore error saving score", "error", err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PublishTopScores(c.Request.Context()); err != nil {
		slog.Error("saveScore error publishing top scores", "error", err.Error())
	}

	c.Status(204)
}
