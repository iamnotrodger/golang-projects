package processes

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
	"github.com/redis/go-redis/v9"
)

type MessageHandler func(scores []model.Score)

type LeaderboardSubscriber struct {
	rdb     *redis.Client
	handler MessageHandler
}

func NewLeaderSubscriber(rdb *redis.Client, handler MessageHandler) *LeaderboardSubscriber {
	return &LeaderboardSubscriber{
		rdb:     rdb,
		handler: handler,
	}
}

func (l *LeaderboardSubscriber) Run(ctx context.Context, errChan chan error) {
	go l.start(ctx, errChan)
	<-ctx.Done()
}

func (h *LeaderboardSubscriber) start(ctx context.Context, errChan chan error) {
	pubsub := h.rdb.Subscribe(ctx, "leaderboard:top10")
	defer pubsub.Close()

	ch := pubsub.Channel()
	slog.Info("subscribed to leaderboard:top10")

	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
		case msg := <-ch:
			var scores []model.Score
			if err := json.Unmarshal([]byte(msg.Payload), &scores); err != nil {
				slog.Error("Error unmarshaling leaderboard data", "error", err)
				continue
			}

			h.handler(scores)
		}
	}
}
