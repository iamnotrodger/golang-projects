package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// MockHealthDatabase is a mock implementation of HealthDatabase
type MockHealthDatabase struct {
	PingFunc func() error
}

func (m *MockHealthDatabase) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
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
			mockDB := &MockHealthDatabase{
				PingFunc: func() error {
					return tt.dbPingError
				},
			}

			healthAPI := &HealthAPI{db: mockDB}

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
