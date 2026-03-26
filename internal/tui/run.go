package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"
)

func Run(ctx context.Context, configPath string) error {
	m, err := New(ctx, configPath)
	if err != nil {
		return err
	}

	_, err = tea.NewProgram(m).Run()
	return err
}
