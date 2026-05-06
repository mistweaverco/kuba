package kuba

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLookPathWithEnv_RespectsEnvPATH(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test uses unix executable bits")
	}

	dir := t.TempDir()
	exe := filepath.Join(dir, "turbo")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write temp executable: %v", err)
	}

	got, err := lookPathWithEnv("turbo", []string{"PATH=" + dir})
	if err != nil {
		t.Fatalf("expected to resolve turbo, got err: %v", err)
	}
	if got != exe {
		t.Fatalf("expected %q, got %q", exe, got)
	}
}

func TestLookPathWithEnv_ExplicitPathPassthrough(t *testing.T) {
	explicit := "/some/explicit/path"
	got, err := lookPathWithEnv(explicit, []string{"PATH=/nope"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got != explicit {
		t.Fatalf("expected %q, got %q", explicit, got)
	}
}

