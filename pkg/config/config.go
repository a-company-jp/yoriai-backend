package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

var Config config

type config struct {
	LineConfig   LineConfig   `yaml:"line"`
	VonageConfig VonageConfig `yaml:"vonage"`
}

func init() {
	Config = config{}
	loc := os.Getenv("ENV_LOCATION")
	if loc == "" {
		slog.Error("No ENV_LOCATION found")
		return
	}
	if err := yaml.Unmarshal([]byte(loc), &Config); err != nil {
		slog.Error("Failed to unmarshal yaml", err)
		return
	}
}
