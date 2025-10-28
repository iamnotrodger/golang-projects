package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobal(t *testing.T) {
	assert.Equal(t, Global.Env, defaultEnv)
	assert.Equal(t, Global.Port, defaultPort)
	assert.Equal(t, Global.LogLevel, defaultLogLevel)
	assert.Equal(t, Global.KafkaBroker, defaultKafkaBroker)
	assert.Equal(t, Global.KafkaTicketTopic, defaultKafkaTicketTopic)
	assert.Equal(t, Global.Secret.DatabaseURL, defaultDatabaseURL)
}

func TestLoadConfig(t *testing.T) {
	testConfig := New()

	originalPort := os.Getenv("PORT")
	originalLogLevel := os.Getenv("LOG_LEVEL")

	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("LOG_LEVEL", originalLogLevel)
	}()

	envLogLevel := "debug"
	os.Setenv("LOG_LEVEL", envLogLevel)
	os.Unsetenv("PORT")

	LoadConfig(testConfig)

	cases := []struct {
		name        string
		configValue any
		expected    any
	}{
		{
			name:        "uses env value when defined",
			configValue: testConfig.LogLevel,
			expected:    envLogLevel,
		},
		{
			name:        "uses default value when env not defined",
			configValue: testConfig.Port,
			expected:    defaultPort,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.configValue)
		})
	}
}
