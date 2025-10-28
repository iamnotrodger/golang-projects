package health

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockHealthCheck struct {
	mock.Mock
}

func (m *MockHealthCheck) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestService_Ping(t *testing.T) {
	type testDef struct {
		name          string
		checks        map[string]HealthCheck
		expectedError error
	}

	tests := []testDef{
		{
			name: "all checks pass",
			checks: func() map[string]HealthCheck {
				check1 := &MockHealthCheck{}
				check1.On("Ping").Return(nil)

				check2 := &MockHealthCheck{}
				check2.On("Ping").Return(nil)

				check3 := &MockHealthCheck{}
				check3.On("Ping").Return(nil)

				return map[string]HealthCheck{
					"check1": check1,
					"check2": check2,
					"check3": check3,
				}
			}(),
			expectedError: nil,
		},
		{
			name: "one check fails",
			checks: func() map[string]HealthCheck {
				check1 := &MockHealthCheck{}
				check1.On("Ping").Return(nil)

				check2 := &MockHealthCheck{}
				check2.On("Ping").Return(errors.New("check2 error"))

				check3 := &MockHealthCheck{}
				check3.On("Ping").Return(nil)

				return map[string]HealthCheck{
					"check1": check1,
					"check2": check2,
					"check3": check3,
				}
			}(),
			expectedError: errors.New("check2 error"),
		},
		{
			name: "multiple checks fail",
			checks: func() map[string]HealthCheck {
				check1 := &MockHealthCheck{}
				check1.On("Ping").Return(nil)

				check2 := &MockHealthCheck{}
				check2.On("Ping").Return(errors.New("check2 error"))

				check3 := &MockHealthCheck{}
				check3.On("Ping").Run(func(args mock.Arguments) {
					time.Sleep(200 * time.Millisecond)
				}).Return(errors.New("check3 error"))

				return map[string]HealthCheck{
					"check1": check1,
					"check2": check2,
					"check3": check3,
				}
			}(),
			expectedError: errors.New("check2 error"),
		},
		{
			name: "no checks configured",
			checks: func() map[string]HealthCheck {
				return map[string]HealthCheck{}
			}(),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.checks)

			err := service.Ping()
			require.Equal(t, tt.expectedError, err)

			for _, check := range tt.checks {
				if mockCheck, ok := check.(*MockHealthCheck); ok {
					mockCheck.AssertExpectations(t)
				}
			}
		})
	}
}
