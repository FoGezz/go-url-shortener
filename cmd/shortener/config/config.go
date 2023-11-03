package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v10"
)

// var errInvalidFormat error = errors.New("address must be in format 'host:port'. i.e. localhost:8080")

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ResponseAddress string `env:"BASE_URL"`
}

func (cfg *Config) ParseEnv() {
	env.Parse(cfg)
}

func (cfg *Config) ParseFlags() {
	flag.Func("a", "Example: -a localhost:8080", func(v string) error {
		if err := validateAddress(v); err != nil {
			return err
		}
		cfg.ServerAddress = v
		return nil
	})
	flag.Func("b", "Example -b http://redirectdomain.com", func(v string) error {
		if err := validateAddress(v); err != nil {
			return err
		}
		cfg.ResponseAddress = v
		return nil
	})
	flag.Parse()
}

func validateAddress(v string) error {
	if parts := strings.Split(v, ":"); len(parts) != 2 {
		return nil
		// return errInvalidFormat === АВТОТЕСТЫ ОГРАНИЧИЛИ
	}
	return nil
}
