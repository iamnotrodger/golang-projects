package leaderboard

import (
	"testing"

	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHub_NewHub(t *testing.T) {
	hub := NewHub()
	assert.NotNil(t, hub)
}

func TestHub_RegisterClient(t *testing.T) {
	hub := NewHub()
	clientChan := hub.RegisterClient("client1")
	require.NotNil(t, clientChan)
}

func TestHub_UnregisterClient(t *testing.T) {
	hub := NewHub()
	clientChan := hub.RegisterClient("client1")
	hub.UnregisterClient("client1")

	select {
	case _, ok := <-clientChan:
		require.False(t, ok)
	default:
		t.Fatal("expected client channel to be closed")
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	clientChan := hub.RegisterClient("client1")
	hub.Broadcast([]model.Score{{Name: "Alice", Value: 100}})
	select {
	case scores := <-clientChan:
		require.Equal(t, []model.Score{{Name: "Alice", Value: 100}}, scores)
	default:
		t.Fatal("expected client channel to receive scores")
	}
}

func TestHub_Shutdown(t *testing.T) {
	hub := NewHub()
	clientChan1 := hub.RegisterClient("client1")
	clientChan2 := hub.RegisterClient("client2")
	hub.Shutdown()

	select {
	case _, ok := <-clientChan1:
		require.False(t, ok)
	case _, ok := <-clientChan2:
		require.False(t, ok)
	default:
		t.Fatal("expected client channels to be closed")
	}
}
