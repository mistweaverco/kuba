package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

func formArrowNavCmd(form *huh.Form, dir string) (tea.Cmd, bool) {
	// Only intercept up/down when the focused field isn't one that uses arrow
	// keys for internal navigation (Select/MultiSelect/Text).
	f := form.GetFocusedField()
	switch f.(type) {
	case *huh.Text:
		return nil, false
	case *huh.Input, *huh.Confirm:
		// ok
	default:
		// Select/MultiSelect and others should keep arrow behavior.
		return nil, false
	}

	if dir == "up" {
		return form.PrevField(), true
	}
	if dir == "down" {
		return form.NextField(), true
	}
	return nil, false
}
