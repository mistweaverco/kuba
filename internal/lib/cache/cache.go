package cache

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mistweaverco/kuba/internal/lib/log"
)

// Cache represents a SQLite-based cache for secrets
type Cache struct {
	db *sql.DB
}

// CacheEntry represents a cached secret entry
type CacheEntry struct {
	Path      string    `json:"path"`
	KubaEnv   string    `json:"kuba_env"`
	Env       string    `json:"env"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewCache creates a new cache instance
func NewCache() (*Cache, error) {
	logger := log.NewLogger()

	// Get cache directory
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	dbPath := filepath.Join(cacheDir, "db.sqlite")
	logger.Debug("Opening cache database", "path", dbPath)

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache database: %w", err)
	}

	cache := &Cache{db: db}

	// Initialize database schema
	if err := cache.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize cache schema: %w", err)
	}

	// Clean up expired entries
	if err := cache.cleanupExpired(); err != nil {
		logger.Debug("Failed to cleanup expired entries", "error", err)
		// Don't fail cache creation for cleanup errors
	}

	logger.Debug("Cache initialized successfully", "path", dbPath)
	return cache, nil
}

// Close closes the cache database connection
func (c *Cache) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// initSchema initializes the database schema
func (c *Cache) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS secrets (
		path TEXT NOT NULL,
		kuba_env TEXT NOT NULL,
		env TEXT NOT NULL,
		value TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL,
		PRIMARY KEY (path, kuba_env, env)
	);
	
	CREATE INDEX IF NOT EXISTS idx_expires_at ON secrets(expires_at);
	`

	_, err := c.db.Exec(query)
	return err
}

// cleanupExpired removes expired entries from the cache
func (c *Cache) cleanupExpired() error {
	query := `DELETE FROM secrets WHERE expires_at < datetime('now')`
	_, err := c.db.Exec(query)
	return err
}

// Set stores a secret in the cache
func (c *Cache) Set(path, kubaEnv, env, value string, ttl time.Duration) error {
	now := time.Now()
	expiresAt := now.Add(ttl)

	query := `
	INSERT OR REPLACE INTO secrets (path, kuba_env, env, value, created_at, expires_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := c.db.Exec(query, path, kubaEnv, env, value, now, expiresAt)
	return err
}

// Get retrieves a secret from the cache
func (c *Cache) Get(path, kubaEnv, env string) (string, bool, error) {
	query := `
	SELECT value FROM secrets 
	WHERE path = ? AND kuba_env = ? AND env = ? AND expires_at > datetime('now')
	`

	var value string
	err := c.db.QueryRow(query, path, kubaEnv, env).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return value, true, nil
}

// Delete removes a secret from the cache
func (c *Cache) Delete(path, kubaEnv, env string) error {
	query := `DELETE FROM secrets WHERE path = ? AND kuba_env = ? AND env = ?`
	_, err := c.db.Exec(query, path, kubaEnv, env)
	return err
}

// Clear removes all secrets from the cache
func (c *Cache) Clear() error {
	query := `DELETE FROM secrets`
	_, err := c.db.Exec(query)
	return err
}

// ClearByPath removes all secrets for a specific kuba.yaml path
func (c *Cache) ClearByPath(path string) error {
	query := `DELETE FROM secrets WHERE path = ?`
	_, err := c.db.Exec(query, path)
	return err
}

// ClearByEnvironment removes all secrets for a specific environment
func (c *Cache) ClearByEnvironment(path, kubaEnv string) error {
	query := `DELETE FROM secrets WHERE path = ? AND kuba_env = ?`
	_, err := c.db.Exec(query, path, kubaEnv)
	return err
}

// List returns all cached entries (for debugging/inspection)
func (c *Cache) List() ([]CacheEntry, error) {
	query := `
	SELECT path, kuba_env, env, value, created_at, expires_at
	FROM secrets
	ORDER BY path, kuba_env, env
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []CacheEntry
	for rows.Next() {
		var entry CacheEntry
		err := rows.Scan(&entry.Path, &entry.KubaEnv, &entry.Env, &entry.Value, &entry.CreatedAt, &entry.ExpiresAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// getCacheDir returns the cache directory path
func getCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Use platform-specific cache directories
	switch runtime.GOOS {
	case "darwin": // macOS
		return filepath.Join(homeDir, "Library", "Caches", "kuba"), nil
	case "windows":
		// Use LOCALAPPDATA if available, otherwise fall back to AppData\Local
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "kuba"), nil
		}
		return filepath.Join(homeDir, "AppData", "Local", "kuba"), nil
	default: // Linux and other Unix-like systems
		// Use XDG_CACHE_HOME if available, otherwise fall back to ~/.cache
		if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
			return filepath.Join(xdgCacheHome, "kuba"), nil
		}
		return filepath.Join(homeDir, ".cache", "kuba"), nil
	}
}
