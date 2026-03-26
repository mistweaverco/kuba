package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m *Model) updateEnvs(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		innerW, innerH := panelInnerSize(msg.Width, msg.Height, panelStyle())
		m.envList.SetSize(innerW, innerH)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if it, ok := m.envList.SelectedItem().(envItem); ok {
				m.selectedEnvName = it.name
				env, err := m.cfg.GetEnvironment(it.name)
				if err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
				m.selectedEnv = env
				if err := m.reloadSecrets(); err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
				m.screen = screenSecrets
				m.filter.SetValue("")
				m.filter.Blur()
				// Ensure the secrets table is sized immediately, even if we haven't
				// received a WindowSizeMsg on this screen yet.
				if m.winW > 0 && m.winH > 0 {
					return m.updateSecrets(tea.WindowSizeMsg{Width: m.winW, Height: m.winH})
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.envList, cmd = m.envList.Update(msg)
	return m, cmd
}
