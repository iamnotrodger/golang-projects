package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
)

type leaderboardService interface {
	RegisterClient(id string) chan []model.Score
	UnregisterClient(id string)
	GetTopK(ctx context.Context, k int) ([]model.Score, error)
}

type Handler struct {
	service leaderboardService
}

func NewHandler(service leaderboardService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/leaderboard/stream", h.HandleSSE)
}

func (h *Handler) HandleSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	id := c.Request.RemoteAddr
	clientChan := h.service.RegisterClient(id)
	defer h.service.UnregisterClient(id)

	initialScores, err := h.service.GetTopK(c.Request.Context(), 10)
	if err != nil {
		slog.Error("Error getting initial leaderboard", "error", err)
	} else {
		h.sendSSEScores(c.Writer, initialScores)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			slog.Info("Client disconnected", "id", id)
			return
		case <-ticker.C:
			h.sendSSEMessage(c.Writer, "keepalive:\n\n")
		case scores := <-clientChan:
			h.sendSSEScores(c.Writer, scores)
		}
	}
}

func (h *Handler) sendSSEScores(w http.ResponseWriter, scores []model.Score) {
	data, err := json.Marshal(scores)
	if err != nil {
		slog.Error("Error marshaling scores", "error", err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", data)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (h *Handler) sendSSEMessage(w http.ResponseWriter, message any) {
	fmt.Fprint(w, message)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
