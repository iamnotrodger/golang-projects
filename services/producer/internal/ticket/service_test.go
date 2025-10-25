package ticket

import (
	"errors"
	"testing"
	"time"

	"github.com/iamnotrodger/golang-kafka/services/producer/internal/config"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/writer"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestCreateTicket(t *testing.T) {
	type testDef struct {
		name          string
		ticket        *model.Ticket
		expectedError error
	}

	tests := []testDef{
		{
			name: "valid ticket",
			ticket: &model.Ticket{
				ID:        "ticket1",
				Title:     "Concert",
				Price:     50.0,
				CreatedAt: time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC),
			},
			expectedError: nil,
		},
		{
			name: "kafka write error",
			ticket: &model.Ticket{
				ID:        "ticket2",
				Title:     "Sports Event",
				Price:     100.0,
				CreatedAt: time.Date(2025, 10, 25, 12, 0, 0, 0, time.UTC),
			},
			expectedError: errors.New("kafka connection failed"),
		},
	}

	config.Global.KafkaTicketTopic = "test-topic"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := &writer.MockKafkaWriter{
				Messages:      []kafka.Message{},
				WriteErr:      tt.expectedError,
				ShouldCapture: tt.expectedError == nil,
			}

			service := &Service{
				kafkaWriter: mockWriter,
			}

			err := service.CreateTicket(tt.ticket)
			require.Equal(t, tt.expectedError, err)
		})
	}
}
