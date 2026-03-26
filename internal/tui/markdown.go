package tui

import "strings"

func mdEscape(s string) string {
	// huh Note fields render markdown; escape characters that would be
	// interpreted as formatting (notably "_" in ENV_VAR names).
	//
	// We keep this intentionally minimal: the goal is to preserve the exact
	// visible value rather than apply markdown styling.
	repl := strings.NewReplacer(
		"\\", "\\\\",
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
	)
	return repl.Replace(s)
}
