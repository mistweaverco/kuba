package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

var versionHeaderRE = regexp.MustCompile(`(?m)^##\s+(?:<small>)?v?([0-9]+\.[0-9]+\.[0-9]+)\s+\(`)

func Select(markdown, selector string) (string, error) {
	sel := strings.TrimSpace(strings.ToLower(selector))
	if sel == "" {
		return markdown, nil
	}

	versions := versionHeaderRE.FindAllStringSubmatch(markdown, -1)
	if len(versions) == 0 {
		return "", fmt.Errorf("no version sections found in changelog")
	}

	target := sel
	if sel == "latest" {
		target = versions[0][1]
	}
	target = strings.TrimPrefix(target, "v")

	lines := strings.Split(markdown, "\n")
	start := -1
	end := len(lines)

	for i, line := range lines {
		m := versionHeaderRE.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}
		v := m[1]
		if start == -1 && v == target {
			start = i
			continue
		}
		if start != -1 {
			end = i
			break
		}
	}

	if start == -1 {
		return "", fmt.Errorf("version '%s' not found in changelog", selector)
	}
	return strings.TrimSpace(strings.Join(lines[start:end], "\n")) + "\n", nil
}
