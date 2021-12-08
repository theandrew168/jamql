package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

var (
	defaultPort = "5000"
)

type Config struct {
	SpotifyClientID     string `toml:"spotify_client_id"`
	SpotifyClientSecret string `toml:"spotify_client_secret"`

	Port string `toml:"port"`
}

func Defaults() Config {
	cfg := Config{
		Port: defaultPort,
	}
	return cfg
}

func Read(data string) (Config, error) {
	var cfg Config
	_, err := toml.Decode(data, &cfg)
	if err != nil {
		return Config{}, nil
	}

	// apply defaults
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	return cfg, nil
}

func ReadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return Read(string(data))
}
