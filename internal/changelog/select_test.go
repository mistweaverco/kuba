package changelog

import (
	"strings"
	"testing"
)

const sample = `## 1.2.0 (2026-01-01)

* two

## <small>1.1.0 (2025-12-01)</small>

* one
`

func TestSelectLatest(t *testing.T) {
	out, err := Select(sample, "latest")
	if err != nil {
		t.Fatalf("Select(latest) error: %v", err)
	}
	if !strings.Contains(out, "## 1.2.0") {
		t.Fatalf("expected latest section to be 1.2.0, got: %s", out)
	}
	if strings.Contains(out, "1.1.0") {
		t.Fatalf("expected only latest section, got: %s", out)
	}
}

func TestSelectVersionWithAndWithoutVPrefix(t *testing.T) {
	out, err := Select(sample, "1.1.0")
	if err != nil {
		t.Fatalf("Select(1.1.0) error: %v", err)
	}
	if !strings.Contains(out, "1.1.0") {
		t.Fatalf("expected 1.1.0 section, got: %s", out)
	}

	out2, err := Select(sample, "v1.1.0")
	if err != nil {
		t.Fatalf("Select(v1.1.0) error: %v", err)
	}
	if out2 != out {
		t.Fatalf("expected same output for v-prefixed version")
	}
}

func TestSelectMissingVersion(t *testing.T) {
	if _, err := Select(sample, "9.9.9"); err == nil {
		t.Fatalf("expected error for missing version")
	}
}
