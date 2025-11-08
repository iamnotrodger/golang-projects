package score

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockScoreService struct {
	mock.Mock
}

func (m *MockScoreService) SaveScore(ctx context.Context, score *model.Score) error {
	args := m.Called(ctx, score)
	return args.Error(0)
}

func (m *MockScoreService) GetTopK(ctx context.Context, k int) ([]model.Score, error) {
	args := m.Called(ctx, k)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Score), args.Error(1)
}

func (m *MockScoreService) PublishTopScores(ctx context.Context, topScores []model.Score) error {
	args := m.Called(ctx, topScores)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewHandler(t *testing.T) {
	mockService := new(MockScoreService)

	handler := NewHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestHandler_RegisterRoutes(t *testing.T) {
	t.Run("should register routes correctly", func(t *testing.T) {
		mockService := new(MockScoreService)
		handler := NewHandler(mockService)
		router := setupTestRouter()

		handler.RegisterRoutes(router)

		req := httptest.NewRequest(http.MethodPost, "/score/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_SaveScore(t *testing.T) {
	type testCase struct {
		name               string
		requestBody        any
		score              *model.Score
		expectedStatusCode int
		expectedBody       map[string]any
		setupMock          func(score *model.Score) *MockScoreService
	}

	topScores := []model.Score{
		{Name: "Alice", Value: 100},
		{Name: "Bob", Value: 90},
	}

	score := &model.Score{
		ID:    "user1",
		Name:  "Alice",
		Value: 100,
	}

	requestBody := map[string]any{
		"id":    "user1",
		"name":  "Alice",
		"value": 100,
	}

	tests := []testCase{
		{
			name:        "should save score successfully and return 204",
			requestBody: requestBody,
			score:       score,
			setupMock: func(score *model.Score) *MockScoreService {
				mockService := new(MockScoreService)
				mockService.On("SaveScore", mock.Anything, score).Return(nil)
				mockService.On("GetTopK", mock.Anything, mock.Anything).Return(topScores, nil)
				mockService.On("PublishTopScores", mock.Anything, topScores).Return(nil)
				return mockService
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:        "should return 500 when SaveScore fails",
			requestBody: requestBody,
			score:       score,
			setupMock: func(score *model.Score) *MockScoreService {
				mockService := new(MockScoreService)
				mockService.On("SaveScore", mock.Anything, score).Return(errors.New("redis save score error"))
				return mockService
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       map[string]any{"error": "redis save score error"},
		},
		{
			name:        "should return 204 when GetTopK fails",
			score:       score,
			requestBody: requestBody,
			setupMock: func(score *model.Score) *MockScoreService {
				mockService := new(MockScoreService)
				mockService.On("SaveScore", mock.Anything, score).Return(nil)
				mockService.On("GetTopK", mock.Anything, 10).Return(nil, errors.New("redis query error"))
				return mockService
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:        "should return 204 even when PublishTopScores fails",
			score:       score,
			requestBody: requestBody,
			setupMock: func(score *model.Score) *MockScoreService {
				mockService := new(MockScoreService)
				mockService.On("SaveScore", mock.Anything, score).Return(nil)
				mockService.On("GetTopK", mock.Anything, mock.Anything).Return(topScores, nil)
				mockService.On("PublishTopScores", mock.Anything, topScores).Return(errors.New("publish error"))
				return mockService
			},
			expectedStatusCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScoreService := tt.setupMock(tt.score)
			handler := NewHandler(mockScoreService)

			router := setupTestRouter()
			handler.RegisterRoutes(router)
			w := httptest.NewRecorder()

			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/score/", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedBody != nil {
				responseBody := map[string]any{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedBody, responseBody)
			}

			mockScoreService.AssertExpectations(t)
		})
	}
}
