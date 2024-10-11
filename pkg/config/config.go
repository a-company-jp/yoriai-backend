package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
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
		panic("No ENV_LOCATION found")
		return
	}
	if err := yaml.Unmarshal([]byte(loc), &Config); err != nil {
		panic(fmt.Errorf("failed to unmarshal yaml: %w", err))
		return
	}
}
