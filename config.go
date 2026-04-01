package template

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func Init() {
	const ConfigPath = "config.yaml"
	var Config struct {
		Version string `yaml:"version"`
	}
	bytes, err := os.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version: %s\n", Config.Version)
}
