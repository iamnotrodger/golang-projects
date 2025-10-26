package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	go func() {
		runnable.ErrorCount = 0
		for {
			time.Sleep(runnable.Interval)
			if runnable.Terminated() {
				return
			}

			if runnable.IsThrowError {
				runnable.leaderMutex.Lock()
				errChan <- errors.New(fmt.Sprint(runnable.Name, ":error:", runnable.ErrorCount))
				runnable.ErrorCount++
				runnable.leaderMutex.Unlock()
			}
		}
	}()

	<-ctx.Done()
	time.Sleep(runnable.ShutdownDuration)
	runnable.leaderMutex.Lock()
	runnable.isTerminated = true
	runnable.leaderMutex.Unlock()
}

func (runnable *MockRunnable) Terminated() bool {
	runnable.leaderMutex.Lock()
	defer runnable.leaderMutex.Unlock()
	return runnable.isTerminated
}

func TestApplication_Start(t *testing.T) {
	type testDef struct {
		name         string
		processes    map[string]Runnable
		shutdownChan chan struct{}
		errChan      chan error
	}

	tests := []testDef{
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
			errChan:      make(chan error),
			shutdownChan: make(chan struct{}),
		},
	}

	for _, td := range tests {
		t.Run(td.name,
			func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				application := NewApplication(td.processes)
				errChan := application.Run(ctx, td.shutdownChan)

				errors := map[string]int{}

				for errors["p2"] < 2 && errors["p3"] < 2 {
					err := <-errChan
					msg := err.Error()
					processName := strings.Split(msg, ":")[0]
					errors[processName]++
				}

				cancel()

				isShutdownChannelOpen := true
				select {
				case _, isShutdownChannelOpen = <-td.shutdownChan:
				case <-time.After(time.Second * 2):
					t.Fatal("shutdown channel was not closed within timeout")
				}

				require.Equal(t, errors["p1"], 0)
				require.False(t, isShutdownChannelOpen)
				require.True(t, td.processes["p1"].(*MockRunnable).Terminated())
				require.True(t, td.processes["p2"].(*MockRunnable).Terminated())
				require.True(t, td.processes["p3"].(*MockRunnable).Terminated())
			})
	}
}
