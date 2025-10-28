package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockHealthDatabase is a mock implementation of HealthDatabase
type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthAPI_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		dbPingError    error
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name:           "healthy database returns 200",
			dbPingError:    nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]any{"status": "healthy"},
		},
		{
			name:           "unhealthy database returns 500",
			dbPingError:    errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]any{"status": "unhealthy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockHealthService{}
			mockService.On("Ping", tt.dbPingError).Return(tt.dbPingError)

			healthAPI := NewHealthAPI(mockService)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

			healthAPI.Health(ctx)

			require.Equal(t, tt.expectedStatus, w.Code)

			expectedJSON, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)
			require.JSONEq(t, string(expectedJSON), w.Body.String())
		})
	}
}
