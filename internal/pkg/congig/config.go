package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
)

type Config struct {
	App     AppConfig
	GRPC    GRPCCConfig
	Clients ClientsConfig
}

type AppConfig struct {
	Env string `env:"ENV" env-default:"local"`
}

type GRPCCConfig struct {
	Port    int           `env:"PORT" env-required:"true"`
	Timeout time.Duration `env:"TIMEOUT" env-default:"5s"`
}

type Clients struct {
	Address      string        `env:"CLIENT_ADDRESS"`
	Timeout      time.Duration `env:"CLIENT_TIMEOUT"`
	RetriesCount int           `env:"CLIENT_RETRIES_COUNT"`
	//Insecure     bool          `env:"CLIENT_INSECURE"`
}

type ClientsConfig struct {
	SSO Clients `env:"CLIENT_CONFIG_SSO"`
}

var configInstance *Config
var configErr error

func GetConfig() (*Config, error) {
	if configInstance == nil {
		var readConfigOnce sync.Once

		readConfigOnce.Do(func() {
			configInstance = &Config{}
			configErr = cleanenv.ReadEnv(configInstance)
		})
	}

	return configInstance, configErr
}
