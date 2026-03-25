package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(ctx context.Context, configPath string) error {
	m, err := New(ctx, configPath)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}
