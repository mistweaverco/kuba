package fileutils

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestJoinPathMatchesFilepathJoin(t *testing.T) {
	tests := [][]string{
		{"a", "b"},
		{"/tmp", "kuba.log"},
		{"", "kuba.log"},
	}

	for _, tt := range tests {
		got := JoinPath(tt...)
		want := filepath.Join(tt...)
		if got != want {
			t.Fatalf("JoinPath(%q) = %q, want %q", tt, got, want)
		}
	}
}

func TestJoinPathDoesNotPrefixSeparatorForRelativePath(t *testing.T) {
	got := JoinPath("temp", "kuba.log")
	if strings.HasPrefix(got, string(filepath.Separator)) {
		t.Fatalf("JoinPath unexpectedly prefixed separator: %q", got)
	}
}
