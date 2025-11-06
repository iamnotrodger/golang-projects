package score

import (
	"context"
	"encoding/json"

	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

func NewService(rdb *redis.Client) *Service {
	return &Service{rdb: rdb}
}

func (s *Service) SaveScore(ctx context.Context, score *model.Score) error {
	err := s.rdb.ZAdd(ctx, "leaderboard", redis.Z{
		Score:  float64(score.Value),
		Member: score.Name,
	}).Err()
	return err
}

func (s *Service) PublishTopScores(ctx context.Context) error {
	topScores, err := s.getTopK(ctx, 10)
	if err != nil {
		return err
	}

	leaderboardData, err := json.Marshal(topScores)
	if err != nil {
		return err
	}

	return s.rdb.Publish(ctx, "leaderboard:top10", leaderboardData).Err()
}

func (s *Service) getTopK(ctx context.Context, k int) ([]model.Score, error) {
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
