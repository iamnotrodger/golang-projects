package ticket

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/iamnotrodger/golang-projects/pkg/proto/topics"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockTicketStore struct {
	mock.Mock
}

func (m *MockTicketStore) CreateTicket(ctx context.Context, ticket *topics.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func TestHandleMessage(t *testing.T) {
	type testDef struct {
		name          string
		message       kafka.Message
		store         *MockTicketStore
		expectedError error
	}

	validTicket := &topics.Ticket{
		Id:        "ticket-123",
		Title:     "Concert Ticket",
		Price:     99.99,
		CreatedAt: timestamppb.New(time.Date(2025, 10, 27, 12, 0, 0, 0, time.UTC)),
	}
	validTicketBytes, _ := proto.Marshal(validTicket)

	tests := []testDef{
		{
			name: "successfully handles valid message",
			message: kafka.Message{
				Key:   []byte("ticket-123"),
				Value: validTicketBytes,
			},
			store: func() *MockTicketStore {
				m := &MockTicketStore{}
				m.On("CreateTicket", mock.Anything, mock.MatchedBy(func(t *topics.Ticket) bool {
					return t.Id == "ticket-123" && t.Price == 99.99
				})).Return(nil)
				return m
			}(),
			expectedError: nil,
		},
		{
			name: "returns error on invalid protobuf",
			message: kafka.Message{
				Key:   []byte("ticket-456"),
				Value: []byte("invalid protobuf data"),
			},
			store:         &MockTicketStore{},
			expectedError: errors.New("cannot parse invalid wire-format data"),
		},
		{
			name: "returns error when store fails",
			message: kafka.Message{
				Key:   []byte("ticket-789"),
				Value: validTicketBytes,
			},
			store: func() *MockTicketStore {
				m := &MockTicketStore{}
				m.On("CreateTicket", mock.Anything, mock.Anything).Return(errors.New("database connection failed"))
				return m
			}(),
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "handles context cancellation",
			message: kafka.Message{
				Key:   []byte("ticket-999"),
				Value: validTicketBytes,
			},
			store: func() *MockTicketStore {
				m := &MockTicketStore{}
				m.On("CreateTicket", mock.Anything, mock.Anything).Return(context.Canceled)
				return m
			}(),
			expectedError: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := NewService(tt.store)
			ctx := context.Background()

			err := service.HandleMessage(ctx, tt.message)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			tt.store.AssertExpectations(t)
		})
	}
}
