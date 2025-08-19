package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("new config creation", func(t *testing.T) {
		flags := ConfigFlags{
			Version: true,
		}
		cfg := NewConfig(Config{Flags: flags})

		assert.Equal(t, true, cfg.Flags.Version)
	})

	t.Run("get config flags", func(t *testing.T) {
		flags := ConfigFlags{
			Version: false,
		}
		cfg := Config{Flags: flags}

		result := cfg.GetConfigFlags()
		assert.Equal(t, false, result.Version)
	})

	t.Run("config flags default values", func(t *testing.T) {
		var flags ConfigFlags
		cfg := Config{Flags: flags}

		result := cfg.GetConfigFlags()
		assert.Equal(t, false, result.Version)
	})
}
