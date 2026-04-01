package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var initTests = map[string]struct {
		configFile string
	}{
		"config file": {
			configFile: "testdata/config.yaml",
		},
	}

	for name, test := range initTests {
		t.Run(name, func(t *testing.T) {
			cfg := Init(test.configFile)
			assert.Equal(t, cfg.Version, "1")
		})
	}
}
