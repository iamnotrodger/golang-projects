package ticket

import (
	"context"

	"github.com/iamnotrodger/golang-kafka/pkg/proto/topics"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	dbClient *pgx.Conn
}

func NewStore(dbClient *pgx.Conn) *Store {
	return &Store{
		dbClient: dbClient,
	}
}

func (s *Store) CreateTicket(ctx context.Context, ticket *topics.Ticket) error {
	_, err := s.dbClient.Exec(
		ctx,
		"INSERT INTO tickets (id, title, price, created_at) VALUES ($1, $2, $3, $4)",
		ticket.Id, ticket.Title, ticket.Price, ticket.CreatedAt,
	)
	return err
}
