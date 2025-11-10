package leaderboard

import (
	"log/slog"
	"sync"

	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
)

type Hub struct {
	clients map[string]chan []model.Score
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]chan []model.Score),
	}
}

func (h *Hub) RegisterClient(id string) chan []model.Score {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientChan := make(chan []model.Score, 10)
	h.clients[id] = clientChan
	return clientChan
}

func (h *Hub) UnregisterClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientChan, exists := h.clients[id]
	if exists {
		close(clientChan)
		delete(h.clients, id)
	}
}

func (h *Hub) Broadcast(scores []model.Score) {
	h.mu.Lock()
	defer h.mu.Unlock()

	slog.Info("broadcasting scores", "count", len(scores))

	for id, clientChan := range h.clients {
		select {
		case clientChan <- scores:
		default:
			slog.Info("client channel full, skipping update", "client_id", id)
		}
	}
}

func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	slog.Info("Shutting down hub, closing all client connections", "count", len(h.clients))
	for id, clientChan := range h.clients {
		close(clientChan)
		delete(h.clients, id)
	}
}
