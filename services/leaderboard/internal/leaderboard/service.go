package leaderboard

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb       *redis.Client
	clients   map[string]chan []model.Score
	clientsMu sync.RWMutex
}

func NewService(rdb *redis.Client) *Service {
	return &Service{
		rdb:     rdb,
		clients: make(map[string]chan []model.Score),
	}
}

func (s *Service) Subscribe(ctx context.Context) error {
	pubsub := s.rdb.Subscribe(ctx, "leaderboard:top10")
	defer pubsub.Close()

	ch := pubsub.Channel()

	slog.Info("Subscribed to leaderboard:top10")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			var scores []model.Score
			if err := json.Unmarshal([]byte(msg.Payload), &scores); err != nil {
				slog.Error("Error unmarshaling leaderboard data", "error", err)
				continue
			}

			s.broadcast(scores)
		}
	}
}

func (s *Service) RegisterClient(id string) chan []model.Score {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	clientChan := make(chan []model.Score, 10)
	s.clients[id] = clientChan
	return clientChan
}

func (s *Service) UnregisterClient(id string) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	clientChan := s.clients[id]
	close(clientChan)
	delete(s.clients, id)
}

func (s *Service) broadcast(scores []model.Score) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	for id, clientChan := range s.clients {
		select {
		case clientChan <- scores:
		default:
			slog.Info("client channel full, skipping update", "client_id", id)
		}
	}
}

func (s *Service) GetTopK(ctx context.Context, k int) ([]model.Score, error) {
	results, err := s.rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, int64(k-1)).Result()
	if err != nil {
		return nil, err
	}

	scores := make([]model.Score, len(results))
	for i, result := range results {
		scores[i] = model.Score{
			Name:  result.Member.(string),
			Value: int(result.Score),
		}
	}

	return scores, nil
}
