package health

import "log/slog"

type HealthCheck interface {
	Ping() error
}

type Service struct {
	checks map[string]HealthCheck
}

func NewService(checks map[string]HealthCheck) *Service {
	return &Service{checks: checks}
}

// TODO: make this async/concurrent
func (s *Service) Ping() error {
	for name, check := range s.checks {
		if err := check.Ping(); err != nil {
			slog.Error("health check failed", "service", name, "error", err.Error())
			return err
		}
	}
	return nil
}
