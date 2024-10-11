package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var Config config

type config struct {
	General      General      `yaml:"general"`
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
	body, err := os.ReadFile(loc)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
		return
	}
	if err := yaml.Unmarshal(body, &Config); err != nil {
		panic(fmt.Errorf("failed to unmarshal yaml: %w", err))
		return
	}
}
