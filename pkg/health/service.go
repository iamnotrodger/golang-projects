package health

import (
	"context"
	"log/slog"
)

type HealthCheck interface {
	Ping(ctx context.Context) error
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

func (s *Service) Ping(ctx context.Context) error {
	if len(s.checks) == 0 {
		return nil
	}

	checkCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan healthCheckResult, len(s.checks))

	for name, check := range s.checks {
		go func(name string, check HealthCheck) {
			err := check.Ping(checkCtx)
			results <- healthCheckResult{name: name, err: err}
		}(name, check)
	}

	for range len(s.checks) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res := <-results:
			if res.err != nil {
				slog.Error("health check failed", "service", res.name, "error", res.err.Error())
				return res.err
			}
		}
	}

	return nil
}
