package template

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var loadTests = map[string]struct {
		configFile string
		want       int
	}{
		"config file": {
			configFile: "testdata/config.yaml",
			want:       1,
		},
		"default when file missing": {
			configFile: filepath.Join(t.TempDir(), "missing.yaml"),
			want:       1,
		},
	}

	for name, test := range loadTests {
		t.Run(name, func(t *testing.T) {
			cfg, err := Load(test.configFile)
			assert.NoError(t, err)
			assert.Equal(t, test.want, cfg.Version)
		})
	}
}
