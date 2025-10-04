package cache

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a duration string in various formats
// Supports: "1d", "2w", "72h", "2y", "3600" (seconds), "true" (default 12h), "false" (disabled)
func ParseDuration(duration interface{}) (time.Duration, bool, error) {
	switch v := duration.(type) {
	case bool:
		if v {
			return 12 * time.Hour, true, nil // Default 12 hours
		}
		return 0, false, nil // Disabled
	case int:
		return time.Duration(v) * time.Second, true, nil
	case int64:
		return time.Duration(v) * time.Second, true, nil
	case float64:
		return time.Duration(int64(v)) * time.Second, true, nil
	case string:
		return parseDurationString(v)
	default:
		return 0, false, fmt.Errorf("unsupported duration type: %T", duration)
	}
}

// parseDurationString parses duration strings like "1d", "2w", "72h", "2y"
func parseDurationString(s string) (time.Duration, bool, error) {
	s = strings.TrimSpace(s)

	// Handle boolean strings
	if s == "true" {
		return 12 * time.Hour, true, nil
	}
	if s == "false" {
		return 0, false, nil
	}

	// First try parsing as Go duration string (e.g., "24h0m0s", "1h30m")
	if duration, err := time.ParseDuration(s); err == nil {
		return duration, true, nil
	}

	// Regular expression to match duration patterns
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)([smhdwy])$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		// Try parsing as plain number (seconds)
		if seconds, err := strconv.ParseFloat(s, 64); err == nil {
			return time.Duration(seconds) * time.Second, true, nil
		}
		return 0, false, fmt.Errorf("invalid duration format: %s", s)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid duration value: %s", matches[1])
	}

	unit := matches[2]
	var multiplier time.Duration

	switch unit {
	case "s":
		multiplier = time.Second
	case "m":
		multiplier = time.Minute
	case "h":
		multiplier = time.Hour
	case "d":
		multiplier = 24 * time.Hour
	case "w":
		multiplier = 7 * 24 * time.Hour
	case "y":
		multiplier = 365 * 24 * time.Hour
	default:
		return 0, false, fmt.Errorf("unsupported duration unit: %s", unit)
	}

	duration := time.Duration(value) * multiplier
	return duration, true, nil
}

// FormatDuration formats a duration into a human-readable string
func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "disabled"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 && days == 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 && days == 0 && hours == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	if len(parts) == 0 {
		return "0s"
	}

	return strings.Join(parts, "")
}
