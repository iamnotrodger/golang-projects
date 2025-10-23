package health

import "log/slog"

type HealthCheck interface {
	Ping() error
}

type Service struct {
	Checks map[string]HealthCheck
}

func NewService() *Service {
	return &Service{
		Checks: map[string]HealthCheck{
			"kafka": NewKafkaCheck(),
		},
	}
}

// TODO: make this async/concurrent
func (s *Service) Ping() error {
	for name, check := range s.Checks {
		if err := check.Ping(); err != nil {
			slog.Error("health check failed", "service", name, "error", err)
			return err
		}
	}
	return nil
}
