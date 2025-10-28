package health

import (
	"log/slog"
)

type HealthCheck interface {
	Ping() error
}

type Service struct {
	checks map[string]HealthCheck
}

func NewService(checks map[string]HealthCheck) *Service {
	return &Service{checks: checks}
}

type healthCheckResult struct {
	name string
	err  error
}

func (s *Service) Ping() error {
	if len(s.checks) == 0 {
		return nil
	}

	results := make(chan healthCheckResult, len(s.checks))

	for name, check := range s.checks {
		go func(name string, check HealthCheck) {
			err := check.Ping()
			results <- healthCheckResult{name: name, err: err}
		}(name, check)
	}

	for range len(s.checks) {
		res := <-results
		if res.err != nil {
			slog.Error("health check failed", "service", res.name, "error", res.err.Error())
			return res.err
		}
	}

	return nil
}
