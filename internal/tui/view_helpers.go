package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) viewModal(title, body string) string {
	box := fitPanelToWindow(panelStyle(), m.winW, m.winH)
	return box.Render(lipgloss.NewStyle().Bold(true).Render(title) + "\n\n" + body)
}

func (m *Model) viewBusyModal(title, body string) string {
	return m.viewModal(title, strings.TrimSpace(body))
}
