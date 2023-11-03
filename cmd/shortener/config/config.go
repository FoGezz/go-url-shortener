package config

import (
	"errors"
	"flag"
	"strings"
)

var errInvalidFormat error = errors.New("address must be in format 'host:port'. i.e. localhost:8080")

const (
	defaultServerAddress   = ":8080"
	defaultResponseAddress = "http://localhost:8080"
)

type Config struct {
	ServerAddress   string
	ResponseAddress string
}

func (cfg *Config) ParseFlags() {
	flag.Func("a", "Example: -a localhost:8080", func(v string) error {
		if err := validateAddress(v); err != nil {
			return err
		}
		if v == "" {
			v = defaultServerAddress
		}
		cfg.ServerAddress = v
		return nil
	})
	flag.Func("b", "Example -b http://redirectdomain.com", func(v string) error {
		if err := validateAddress(v); err != nil {
			return err
		}
		if v == "" {
			v = defaultResponseAddress
		}
		cfg.ResponseAddress = v
		return nil
	})
	flag.Parse()
}

func validateAddress(v string) error {
	if parts := strings.Split(v, ":"); len(parts) != 2 {
		return errInvalidFormat
	}
	return nil
}
