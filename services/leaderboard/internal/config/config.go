package config

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	defaultEnv           = "development"
	defaultPort          = 8080
	defaultLogLevel      = "info"
	defaultRedisAddr     = "localhost:6379"
	defaultRedisDb       = 0
	defaultRedisPassword = ""
)

type Secret struct {
	RedisPassword string `mapstructure:"redis_password"`
}

type Spec struct {
	*Secret   `json:"-"`
	Env       string `mapstructure:"env"`
	Port      int    `mapstructure:"port"`
	LogLevel  string `mapstructure:"log_level"`
	RedisAddr string `mapstructure:"redis_addr"`
	RedisDb   int    `mapstructure:"redis_db"`
}

func New() *Spec {
	secret := &Secret{
		RedisPassword: defaultRedisPassword,
	}
	return &Spec{
		Secret:    secret,
		Env:       defaultEnv,
		Port:      defaultPort,
		LogLevel:  defaultLogLevel,
		RedisAddr: defaultRedisAddr,
		RedisDb:   defaultRedisDb,
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
