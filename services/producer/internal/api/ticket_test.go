package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
	"github.com/stretchr/testify/require"
)

// MockTicketService is a mock implementation of TicketService
type MockTicketService struct {
	CreateTicketFunc func(ticket *model.Ticket) error
}

func (m *MockTicketService) CreateTicket(ticket *model.Ticket) error {
	if m.CreateTicketFunc != nil {
		return m.CreateTicketFunc(ticket)
	}
	return nil
}

func TestTicketAPI_CreateTicket(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    any
		serviceError   error
		expectedStatus int
		checkBody      func(t *testing.T, body *bytes.Buffer)
	}{
		{
			name:           "successfully creates ticket",
			requestBody:    map[string]any{"id": "123", "title": "Concert Ticket", "price": 50.00},
			serviceError:   nil,
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body *bytes.Buffer) {
				var responseBody map[string]any
				err := json.Unmarshal(body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.Equal(t, "123", responseBody["id"])
				require.Equal(t, "Concert Ticket", responseBody["title"])
				require.Equal(t, 50.00, responseBody["price"])
				require.NotEmpty(t, responseBody["created_at"])
			},
		},
		{
			name:           "invalid JSON returns 400",
			requestBody:    []byte(`{"id":"123","title":}`),
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
			checkBody:      func(t *testing.T, body *bytes.Buffer) {},
		},
		{
			name:           "service error returns 500",
			requestBody:    map[string]any{"id": "456", "title": "Sports Ticket", "price": 75.50},
			serviceError:   errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			checkBody:      func(t *testing.T, body *bytes.Buffer) {},
		},
		{
			name:           "missing fields still processes",
			requestBody:    map[string]any{"title": "Movie Ticket"},
			serviceError:   nil,
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body *bytes.Buffer) {
				var responseBody map[string]any
				err := json.Unmarshal(body.Bytes(), &responseBody)
				require.NoError(t, err)
				require.Equal(t, "Movie Ticket", responseBody["title"])
				require.Equal(t, 0.0, responseBody["price"])
				require.NotEmpty(t, responseBody["created_at"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockTicketService{
				CreateTicketFunc: func(ticket *model.Ticket) error {
					return tt.serviceError
				},
			}

			ticketAPI := NewTicketAPI(mockService)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			ticketAPI.CreateTicket(ctx)
			require.Equal(t, tt.expectedStatus, w.Code)
			tt.checkBody(t, w.Body)
		})
	}
}
