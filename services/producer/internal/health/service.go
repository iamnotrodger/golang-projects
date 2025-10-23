package health

import "log/slog"

type healthCheck interface {
	Ping() error
}

type Service struct {
	checks map[string]healthCheck
}

func NewService() *Service {
	return &Service{
		checks: map[string]healthCheck{
			"kafka": NewKafkaCheck(),
		},
	}
}

// TODO: make this async/concurrent
func (s *Service) Ping() error {
	for name, check := range s.checks {
		if err := check.Ping(); err != nil {
			slog.Error("health check failed", "service", name, "error", err)
			return err
		}
	}
	return nil
}
