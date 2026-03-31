package tui

import "charm.land/lipgloss/v2"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampMin(v, minV int) int {
	if v < minV {
		return minV
	}
	return v
}

func clamp(v, minV, maxV int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func panelStyle() lipgloss.Style {
	t := vhsEraTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.BorderActive).
		Padding(1, 2)
}

func fitPanelToWindow(s lipgloss.Style, winW, winH int) lipgloss.Style {
	if winW <= 0 || winH <= 0 {
		return s
	}
	// Width/Height represent the full block size (before margins), which
	// includes borders and padding. Using the full window size avoids
	// over-cropping (frame size is already accounted for by the style itself).
	w := clampMin(winW, 0)
	h := clampMin(winH, 0)
	return s.Width(w).Height(h).MaxWidth(w).MaxHeight(h)
}

func panelInnerSize(w, h int, panel lipgloss.Style) (int, int) {
	if w <= 0 || h <= 0 {
		return 0, 0
	}
	fw, fh := panel.GetFrameSize()
	return clampMin(w-fw, 0), clampMin(h-fh, 0)
}
