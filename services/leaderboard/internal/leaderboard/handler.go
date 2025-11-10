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

type Service interface {
	GetTopK(ctx context.Context, k int) ([]model.Score, error)
}

type leaderboardHub interface {
	RegisterClient(id string) chan []model.Score
	UnregisterClient(id string)
}

type Handler struct {
	service Service
	hub     leaderboardHub
}

func NewHandler(service Service, hub leaderboardHub) *Handler {
	return &Handler{service: service, hub: hub}
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
	clientChan := h.hub.RegisterClient(id)
	defer h.hub.UnregisterClient(id)

	initialScores, err := h.service.GetTopK(c.Request.Context(), 10)
	if err != nil {
		slog.Error("error getting initial leaderboard", "error", err)
	} else {
		h.sendSSEScores(c.Writer, initialScores)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			slog.Info("client disconnected", "id", id)
			return
		case <-ticker.C:
			h.writeSSE(c.Writer, "keepalive:\n\n")
		case scores, ok := <-clientChan:
			if !ok {
				slog.Info("hub closed connection", "id", id)
				return
			}
			h.sendSSEScores(c.Writer, scores)
		}
	}
}

func (h *Handler) sendSSEScores(w http.ResponseWriter, scores []model.Score) {
	data, err := json.Marshal(scores)
	if err != nil {
		slog.Error("error marshaling scores", "error", err)
		return
	}
	slog.Info("sending scores", "scores", scores)
	h.writeSSE(w, fmt.Sprintf("data: %s\n\n", data))
}

func (h *Handler) writeSSE(w http.ResponseWriter, message string) {
	fmt.Fprint(w, message)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
