package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	defaultSpotifyRedirectURI = "http://localhost:5000/callback"
	defaultPort = "5000"
)

type Config struct {
	// required
	SecretKey string `toml:"secret_key"`

	// optional - spotify
	SpotifyClientID     string `toml:"spotify_client_id"`
	SpotifyClientSecret string `toml:"spotify_client_secret"`
	SpotifyRedirectURI  string `toml:"spotify_redirect_uri"`

	// optional - general
	Port string `toml:"port"`
}

func Read(data string) (Config, error) {
	var cfg Config
	meta, err := toml.Decode(data, &cfg)
	if err != nil {
		return Config{}, nil
	}

	// gather extra values
	extra := []string{}
	for _, keys := range meta.Undecoded() {
		key := keys[0]
		extra = append(extra, key)
	}

	// error upon extra values
	if len(extra) > 0 {
		msg := strings.Join(extra, ", ")
		return Config{}, fmt.Errorf("extra config values: %s", msg)
	}

	// build set of present config keys
	present := make(map[string]bool)
	for _, keys := range meta.Keys() {
		key := keys[0]
		present[key] = true
	}

	// list of required config values
	required := []string{
		"secret_key",
	}

	// gather missing values
	missing := []string{}
	for _, key := range required {
		if _, ok := present[key]; !ok {
			missing = append(missing, key)
		}
	}

	// error upon missing values
	if len(missing) > 0 {
		msg := strings.Join(missing, ", ")
		return Config{}, fmt.Errorf("missing config values: %s", msg)
	}

	// apply defaults
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}
	if cfg.SpotifyRedirectURI == "" {
		cfg.SpotifyRedirectURI = defaultSpotifyRedirectURI
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
