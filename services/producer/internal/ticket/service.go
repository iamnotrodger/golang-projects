package ticket

import (
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
	"github.com/segmentio/kafka-go"
)

type TicketService struct {
	kafkaWriter *kafka.Writer
}

func NewTicketService(kafkaWriter *kafka.Writer) *TicketService {
	return &TicketService{
		kafkaWriter: kafkaWriter,
	}
}

func (t *TicketService) CreateTicket(ticket *model.Ticket) error {
	// msg := kafka.Message{
	// 	Key: []byte(ticket.ID),
	// 	Value: []byte(),
	// }
	return nil
}
