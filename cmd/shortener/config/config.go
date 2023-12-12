package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

const (
	defaultServerAddress   = ":8080"
	defaultResponseAddress = "http://localhost:8080"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	ResponseAddress string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DbDSN           string `env:"DATABASE_DSN"`
	Alphabet        []rune
}

func (cfg *Config) String() string {
	return fmt.Sprintf(`
	ServerAddress: %s,
	ResponseAddress: %s,
	FileStoragePath: %s,
	Alphabet: "%s",
	DATABASE_DSN: "%s"
	`, cfg.ServerAddress, cfg.ResponseAddress, cfg.FileStoragePath, string(cfg.Alphabet), cfg.DbDSN)
}

func (cfg *Config) Load() {
	cfg.ServerAddress = defaultServerAddress
	cfg.ResponseAddress = defaultResponseAddress
	cfg.FileStoragePath = os.TempDir() + "short-url-db.json"
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
	flag.Func("f", "Example -f /tmp/testfile.json", func(v string) error {
		cfg.FileStoragePath = v
		return nil
	})
	flag.Func("d", "Example -d postgres://username:password@localhost:5432/database_name", func(v string) error {
		cfg.DbDSN = v
		return nil
	})
	flag.Parse()
}
