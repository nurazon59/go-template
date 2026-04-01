package template

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Version string `yaml:"version"`
}

func Init(path string) *Config {
	cfg := newConfig()
	if path == "" {
		return cfg
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func newConfig() *Config {
	return &Config{
		Version: "0.1.0",
	}
}
