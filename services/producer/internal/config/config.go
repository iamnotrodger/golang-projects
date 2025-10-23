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

type Spec struct {
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

var Global = Spec{
	Port:     defaultPort,
	LogLevel: defaultLogLevel,
}

func LoadConfig() {
	v := viper.New()
	v.SetConfigFile(".env")
	v.ReadInConfig()
	v.AutomaticEnv()

	setDefaults(v, Global)

	if err := v.Unmarshal(&Global); err != nil {
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
