package tui

import "sort"

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
		for _, l := range m.gcpLocations {
			supported[l] = true
		}
		locs := make([]string, 0, len(pd.Regions))
		for _, r := range pd.Regions {
			if len(supported) == 0 || supported[r] {
				locs = append(locs, r)
			}
		}
		if len(locs) > 0 {
			m.createReplication = "user-managed"
			m.createLocations = locs
			sort.Strings(m.createLocations)
		}
	}
}
