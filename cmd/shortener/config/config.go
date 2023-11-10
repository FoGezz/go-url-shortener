package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
)

const (
	defaultServerAddress   = ":8080"
	defaultResponseAddress = "http://localhost:8080"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ResponseAddress string `env:"BASE_URL"`
	Alphabet        []rune
}

func (cfg *Config) Load() {
	cfg.ServerAddress = defaultServerAddress
	cfg.ResponseAddress = defaultResponseAddress
	cfg.parseFlags()
	cfg.parseEnv()
	cfg.Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

}

func (cfg *Config) parseEnv() {
	err := env.Parse(cfg)
	if err != nil {
		fmt.Println("Unable to load config:", err)
	}
}

func (cfg *Config) parseFlags() {
	flag.Func("a", "Example: -a localhost:8080", func(v string) error {
		cfg.ServerAddress = v
		return nil
	})
	flag.Func("b", "Example -b http://redirectdomain.com", func(v string) error {
		cfg.ResponseAddress = v
		return nil
	})
	flag.Parse()
}
