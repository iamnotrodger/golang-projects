package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	defaultPort     = 8080
	defaultLogLevel = "info"
)

type Secret struct{}

type Spec struct {
	*Secret  `json:"-"`
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

func New() *Spec {
	secret := &Secret{}
	return &Spec{
		Secret:   secret,
		Port:     defaultPort,
		LogLevel: defaultLogLevel,
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
