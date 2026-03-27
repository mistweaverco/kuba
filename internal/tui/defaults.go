package tui

import (
	"regexp"
	"sort"
	"strings"
)

func uniqueSorted(ss []string) []string {
	if len(ss) == 0 {
		return nil
	}
	sort.Strings(ss)
	out := ss[:0]
	var prev string
	for i, s := range ss {
		if i == 0 || s != prev {
			out = append(out, s)
			prev = s
		}
	}
	return out
}

func (m *Model) applyCreateDefaults() {
	if m.globalCfg == nil || m.globalCfg.Defaults == nil || m.globalCfg.Defaults.Providers == nil || m.selectedEnv == nil {
		return
	}
	pd, ok := m.globalCfg.Defaults.Providers[m.selectedEnv.Provider]
	if !ok {
		return
	}
	if len(pd.Regions) == 0 {
		return
	}

	// For GCP, regions map to Secret Manager locations.
	if m.selectedEnv.Provider == "gcp" {
		// Filter to supported locations if we have them loaded.
		supported := map[string]bool{}
		for _, l := range m.gcpLocationsAll {
			supported[l] = true
		}

		// Treat defaults as a "pre-filter" for the location list. Support:
		// - exact strings (e.g. "us-central1")
		// - regex patterns (e.g. "us-(east1|central1)")
		// If a pattern fails to compile, we fall back to exact matching.
		matchers := make([]func(string) bool, 0, len(pd.Regions))
		for _, raw := range pd.Regions {
			p := strings.TrimSpace(raw)
			if p == "" {
				continue
			}
			// Heuristic: only pay regex cost if it looks like a regex.
			looksRegex := strings.ContainsAny(p, `|.*+?()[]{}^$\`)
			if looksRegex {
				if re, err := regexp.Compile(p); err == nil {
					matchers = append(matchers, re.MatchString)
					continue
				}
			}
			// Exact match fallback.
			matchers = append(matchers, func(s string) bool { return s == p })
		}

		// If we have supported locations from the API, filter those; otherwise, we
		// can only fall back to the configured values (best effort).
		filtered := []string{}
		if len(m.gcpLocationsAll) > 0 {
			for _, loc := range m.gcpLocationsAll {
				for _, matches := range matchers {
					if matches(loc) {
						filtered = append(filtered, loc)
						break
					}
				}
			}
		} else {
			for _, r := range pd.Regions {
				r = strings.TrimSpace(r)
				if r == "" {
					continue
				}
				if len(supported) == 0 || supported[r] {
					filtered = append(filtered, r)
				}
			}
		}

		// Pre-filter the create form list (what you see) and preselect the same
		// set to preserve the existing "defaults imply selection" behavior.
		if len(filtered) > 0 {
			filtered = uniqueSorted(filtered)
			m.gcpLocations = filtered
			m.createReplication = "user-managed"
			m.createLocations = append([]string(nil), filtered...)
		} else if len(m.gcpLocationsAll) > 0 {
			m.gcpLocations = append([]string(nil), m.gcpLocationsAll...)
		}
	}
}
