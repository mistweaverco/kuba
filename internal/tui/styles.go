package tui

import (
	"charm.land/bubbles/v2/table"
)

func secretsTableStyles() table.Styles {
	t := vhsEraTheme()
	s := table.DefaultStyles()

	s.Header = s.Header.
		Foreground(t.Fg).
		Background(t.Surface).
		Bold(true)

	s.Cell = s.Cell.Foreground(t.Fg)

	s.Selected = s.Selected.
		Foreground(t.Bg).
		Background(t.AccentMagenta).
		Bold(true)

	return s
}
