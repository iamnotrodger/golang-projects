package score

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_SaveScore(t *testing.T) {
	type testCase struct {
		name          string
		score         *model.Score
		setupMock     func() (*redis.Client, redismock.ClientMock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "should save score successfully",
			score: &model.Score{
				ID:    "user1",
				Name:  "Alice",
				Value: 100,
			},
			setupMock: func() (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()

				mock.ExpectZAdd("leaderboard", redis.Z{
					Score:  100.0,
					Member: "Alice",
				}).SetVal(1)

				return rdb, mock
			},
			expectedError: nil,
		},
		{
			name: "should return error when ZAdd fails",
			score: &model.Score{
				ID:    "user1",
				Name:  "Alice",
				Value: 100,
			},
			setupMock: func() (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()

				mock.ExpectZAdd("leaderboard", redis.Z{
					Score:  100.0,
					Member: "Alice",
				}).SetErr(redis.Nil)

				return rdb, mock
			},
			expectedError: redis.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, _ := tt.setupMock()
			service := NewService(rdb)
			ctx := context.Background()
			err := service.SaveScore(ctx, tt.score)
			require.Equal(t, tt.expectedError, err)

		})
	}
}

func TestService_GetTopK(t *testing.T) {
	type result struct {
		scores []model.Score
		err    error
	}

	type testCase struct {
		name           string
		setupMock      func(k int) (*redis.Client, redismock.ClientMock)
		expectedResult result
	}

	tests := []testCase{
		{
			name: "should return top k scores",
			setupMock: func(k int) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()
				mock.ExpectZRevRangeWithScores("leaderboard", 0, int64(k-1)).SetVal([]redis.Z{
					{Score: 100, Member: "Alice"},
					{Score: 90, Member: "Bob"},
					{Score: 80, Member: "Charlie"},
				})

				return rdb, mock
			},
			expectedResult: result{
				scores: []model.Score{
					{Name: "Alice", Value: 100},
					{Name: "Bob", Value: 90},
					{Name: "Charlie", Value: 80},
				},
				err: nil,
			},
		},
		{
			name: "should return empty list when no scores",
			setupMock: func(k int) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()
				mock.ExpectZRevRangeWithScores("leaderboard", 0, int64(k-1)).SetVal([]redis.Z{})
				return rdb, mock
			},
			expectedResult: result{
				scores: []model.Score{},
				err:    nil,
			},
		},
		{
			name: "should return error when Redis returns an error",
			setupMock: func(k int) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()
				mock.ExpectZRevRangeWithScores("leaderboard", 0, int64(k-1)).SetErr(redis.Nil)
				return rdb, mock
			},
			expectedResult: result{
				scores: nil,
				err:    redis.Nil,
			},
		},
		{
			name: "should handle different k values",
			setupMock: func(k int) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()
				mock.ExpectZRevRangeWithScores("leaderboard", 0, int64(k-1)).SetVal([]redis.Z{
					{Score: 100, Member: "Alice"},
					{Score: 90, Member: "Bob"},
					{Score: 80, Member: "Charlie"},
					{Score: 70, Member: "David"},
					{Score: 60, Member: "Eve"},
				})
				return rdb, mock
			},
			expectedResult: result{
				scores: []model.Score{
					{Name: "Alice", Value: 100},
					{Name: "Bob", Value: 90},
					{Name: "Charlie", Value: 80},
					{Name: "David", Value: 70},
					{Name: "Eve", Value: 60},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := len(tt.expectedResult.scores)
			rdb, mock := tt.setupMock(k)
			service := NewService(rdb)
			ctx := context.Background()

			scores, err := service.GetTopK(ctx, k)

			if tt.expectedResult.err != nil {
				require.Equal(t, tt.expectedResult.err, err)
			}
			require.Equal(t, tt.expectedResult.scores, scores)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_PublishTopScores(t *testing.T) {
	type testCase struct {
		name          string
		scores        []model.Score
		setupMock     func(scores []model.Score) (*redis.Client, redismock.ClientMock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "should publish top scores successfully",
			scores: []model.Score{
				{Name: "Alice", Value: 100},
				{Name: "Bob", Value: 90},
				{Name: "Charlie", Value: 80},
			},
			setupMock: func(scores []model.Score) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()

				data, err := json.Marshal(scores)
				require.NoError(t, err)

				mock.ExpectPublish("leaderboard:top10", data).SetVal(1)

				return rdb, mock
			},
			expectedError: nil,
		},
		{
			name: "should return error when Publish fails",
			scores: []model.Score{
				{Name: "Alice", Value: 100},
			},
			setupMock: func(scores []model.Score) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()

				data, err := json.Marshal(scores)
				require.NoError(t, err)

				mock.ExpectPublish("leaderboard:top10", data).SetErr(redis.Nil)

				return rdb, mock
			},
			expectedError: redis.Nil,
		},
		{
			name:   "should publish empty list when no scores",
			scores: []model.Score{},
			setupMock: func(scores []model.Score) (*redis.Client, redismock.ClientMock) {
				rdb, mock := redismock.NewClientMock()

				data, err := json.Marshal(scores)
				require.NoError(t, err)

				mock.ExpectPublish("leaderboard:top10", data).SetVal(0)

				return rdb, mock
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb, mock := tt.setupMock(tt.scores)
			service := NewService(rdb)
			ctx := context.Background()

			err := service.PublishTopScores(ctx, tt.scores)
			require.Equal(t, tt.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
