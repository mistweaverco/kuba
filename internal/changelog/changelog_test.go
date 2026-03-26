package changelog

import (
	"strings"
	"testing"
)

func TestEmbeddedChangelogNotEmpty(t *testing.T) {
	if strings.TrimSpace(Markdown) == "" {
		t.Fatalf("embedded changelog markdown is empty")
	}
}
