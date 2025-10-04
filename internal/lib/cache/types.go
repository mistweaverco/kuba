package cache

import "time"

// CacheConfig represents the caching configuration
type CacheConfig struct {
	Enabled bool          `yaml:"enabled"`
	TTL     time.Duration `yaml:"ttl"`
}
