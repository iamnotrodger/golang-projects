package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	defaultEnv              = "development"
	defaultPort             = 8080
	defaultLogLevel         = "info"
	defaultKafkaBroker      = "localhost:9092"
	defaultKafkaTicketTopic = "tickets"
	defaultDatabaseURL      = "postgres://user:password@localhost:5432/tickets"
)

type Secret struct {
	DatabaseURL string `mapstructure:"database_url"`
}

type Spec struct {
	*Secret          `json:"-"`
	Env              string `mapstructure:"env"`
	Port             int    `mapstructure:"port"`
	LogLevel         string `mapstructure:"log_level"`
	KafkaBroker      string `mapstructure:"kafka_broker"`
	KafkaTicketTopic string `mapstructure:"kafka_topic"`
}

func New() *Spec {
	secret := &Secret{
		DatabaseURL: defaultDatabaseURL,
	}
	return &Spec{
		Secret:           secret,
		Env:              defaultEnv,
		Port:             defaultPort,
		LogLevel:         defaultLogLevel,
		KafkaBroker:      defaultKafkaBroker,
		KafkaTicketTopic: defaultKafkaTicketTopic,
	}
}

var Global = New()

func LoadConfig(spec *Spec) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.ReadInConfig()
	v.AutomaticEnv()

	setDefaults(v, spec)

	if err := v.Unmarshal(spec); err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config %s", err))
	}
}

func setDefaults(v *viper.Viper, i any) {
	values := map[string]any{}
	if err := mapstructure.Decode(i, &values); err != nil {
		panic(err)
	}
	for key, defaultValue := range values {
		v.SetDefault(key, defaultValue)
	}
}

func init() {
	LoadConfig(Global)
}
