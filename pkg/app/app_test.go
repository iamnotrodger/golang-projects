package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockRunnable struct {
	Name             string
	Interval         time.Duration
	IsThrowError     bool
	isTerminated     bool
	ErrorCount       int
	ShutdownDuration time.Duration
	leaderMutex      sync.Mutex
}

func (runnable *MockRunnable) Run(ctx context.Context, errChan chan error) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		runnable.ErrorCount = 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(runnable.Interval):
				if runnable.IsThrowError {
					runnable.leaderMutex.Lock()
					errChan <- errors.New(fmt.Sprint(runnable.Name, ":error:", runnable.ErrorCount))
					runnable.ErrorCount++
					runnable.leaderMutex.Unlock()
				}
			}
		}
	}()

	<-ctx.Done()
	time.Sleep(runnable.ShutdownDuration)
	<-done
	runnable.leaderMutex.Lock()
	runnable.isTerminated = true
	runnable.leaderMutex.Unlock()
}

func (runnable *MockRunnable) Terminated() bool {
	runnable.leaderMutex.Lock()
	defer runnable.leaderMutex.Unlock()
	return runnable.isTerminated
}

func TestApplication(t *testing.T) {
	type testCase struct {
		name      string
		processes map[string]Runnable
	}

	tests := []testCase{
		{
			name: "should start all processes, and stop them gracefully",
			processes: map[string]Runnable{
				"p1": &MockRunnable{
					Name:         "p1",
					Interval:     time.Millisecond * 10,
					IsThrowError: false,
				},
				"p2": &MockRunnable{
					Name:         "p2",
					Interval:     time.Millisecond * 10,
					IsThrowError: true,
				},
				"p3": &MockRunnable{
					Name:         "p3",
					Interval:     time.Millisecond * 10,
					IsThrowError: true,
				},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.name,
			func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				application := NewApplication(td.processes)
				errChan := application.Run(ctx)
				defer close(errChan)

				cancel()
				application.Shutdown()

				require.True(t, td.processes["p1"].(*MockRunnable).Terminated())
				require.True(t, td.processes["p2"].(*MockRunnable).Terminated())
				require.True(t, td.processes["p3"].(*MockRunnable).Terminated())
			})
	}
}
