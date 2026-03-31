package tui

import (
	"image/color"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// VHS Era terminal theme palette (from mistweaverco/vhs-era-theme.terminal).
// We keep this as the single source of truth for all TUI colors.
type theme struct {
	Fg               color.Color
	Bg               color.Color
	Muted            color.Color
	SelectionBg      color.Color
	Cursor           color.Color
	CursorText       color.Color
	BorderActive     color.Color
	BorderInactive   color.Color
	AccentPink       color.Color
	AccentMagenta    color.Color
	AccentBlue       color.Color
	AccentCyan       color.Color
	AccentGreen      color.Color
	Surface          color.Color
	SurfaceElevated  color.Color
	SurfaceDeep      color.Color
	SurfaceHighlight color.Color
}

func vhsEraTheme() theme {
	return theme{
		Fg:               lipgloss.Color("#dde1e6"),
		Bg:               lipgloss.Color("#161616"),
		Muted:            lipgloss.Color("#8d8d8d"),
		SelectionBg:      lipgloss.Color("#525252"),
		Cursor:           lipgloss.Color("#f2f4f8"),
		CursorText:       lipgloss.Color("#393939"),
		BorderActive:     lipgloss.Color("#ee5396"),
		BorderInactive:   lipgloss.Color("#ff7eb6"),
		AccentPink:       lipgloss.Color("#ff7eb6"),
		AccentMagenta:    lipgloss.Color("#ee5396"),
		AccentBlue:       lipgloss.Color("#33b1ff"),
		AccentCyan:       lipgloss.Color("#3ddbd9"),
		AccentGreen:      lipgloss.Color("#42be65"),
		Surface:          lipgloss.Color("#262626"),
		SurfaceElevated:  lipgloss.Color("#393939"),
		SurfaceDeep:      lipgloss.Color("#0d0d0d"),
		SurfaceHighlight: lipgloss.Color("#525252"),
	}
}

func titleStyle() lipgloss.Style {
	t := vhsEraTheme()
	return lipgloss.NewStyle().Bold(true).Foreground(t.AccentMagenta)
}

func sectionHeaderStyle() lipgloss.Style {
	t := vhsEraTheme()
	return lipgloss.NewStyle().Bold(true).Foreground(t.AccentPink)
}

func helpStyle() lipgloss.Style {
	t := vhsEraTheme()
	return lipgloss.NewStyle().Foreground(t.Muted)
}

func errorStyle() lipgloss.Style {
	t := vhsEraTheme()
	return lipgloss.NewStyle().Bold(true).Foreground(t.AccentPink)
}

func vhsListStyles() list.Styles {
	t := vhsEraTheme()

	// Use dark defaults as a baseline, then re-color with the VHS palette.
	s := list.DefaultStyles(true)

	s.Title = s.Title.
		Background(t.AccentMagenta).
		Foreground(t.Bg).
		Bold(true)

	// Keep status/help subdued.
	s.StatusBar = s.StatusBar.Foreground(t.Muted)
	s.StatusEmpty = s.StatusEmpty.Foreground(t.Muted)
	s.HelpStyle = s.HelpStyle.Foreground(t.Muted)

	// Filter prompt/cursor colors.
	fs := s.Filter
	fs.Cursor.Color = t.Cursor
	fs.Focused.Prompt = lipgloss.NewStyle().Foreground(t.AccentCyan)
	fs.Blurred.Prompt = lipgloss.NewStyle().Foreground(t.AccentCyan)
	fs.Focused.Text = lipgloss.NewStyle().Foreground(t.Fg)
	fs.Blurred.Text = lipgloss.NewStyle().Foreground(t.Fg)
	fs.Focused.Placeholder = lipgloss.NewStyle().Foreground(t.Muted)
	fs.Blurred.Placeholder = lipgloss.NewStyle().Foreground(t.Muted)
	s.Filter = fs

	s.ActivePaginationDot = s.ActivePaginationDot.Foreground(t.AccentPink)
	s.InactivePaginationDot = s.InactivePaginationDot.Foreground(t.SurfaceElevated)
	s.DividerDot = s.DividerDot.Foreground(t.SurfaceElevated)

	return s
}

func vhsEnvDelegate() list.DefaultDelegate {
	t := vhsEraTheme()

	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	// Selected item: VHS magenta marker + bright text.
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(t.AccentMagenta).
		Foreground(t.Fg).
		Bold(true)

	// Normal / dimmed / match.
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(t.Fg)
	d.Styles.DimmedTitle = d.Styles.DimmedTitle.Foreground(t.Muted)
	d.Styles.FilterMatch = lipgloss.NewStyle().Underline(true).Foreground(t.AccentCyan)

	return d
}

func vhsSecretsFilterStyles() textinput.Styles {
	t := vhsEraTheme()
	s := textinput.DefaultStyles(true)
	s.Cursor.Color = t.Cursor
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(t.AccentCyan)
	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(t.AccentCyan)
	s.Focused.Text = lipgloss.NewStyle().Foreground(t.Fg)
	s.Blurred.Text = lipgloss.NewStyle().Foreground(t.Fg)
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(t.Muted)
	s.Blurred.Placeholder = lipgloss.NewStyle().Foreground(t.Muted)
	return s
}

func themeVHSEra(isDark bool) *huh.Styles {
	t := vhsEraTheme()
	_ = isDark // Palette is intentionally fixed to the theme.

	st := huh.ThemeBase(true)

	// Form/group basics.
	st.Form.Base = lipgloss.NewStyle().Foreground(t.Fg)
	st.Group.Title = lipgloss.NewStyle().Bold(true).Foreground(t.AccentPink)
	st.Group.Description = lipgloss.NewStyle().Foreground(t.Muted)

	// Focus indicator.
	st.Focused.Base = st.Focused.Base.BorderForeground(t.BorderActive)
	st.Focused.Card = st.Focused.Base

	// Common field bits.
	st.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(t.AccentMagenta)
	st.Focused.Description = lipgloss.NewStyle().Foreground(t.Muted)
	st.Focused.ErrorIndicator = st.Focused.ErrorIndicator.Foreground(t.AccentPink)
	st.Focused.ErrorMessage = st.Focused.ErrorMessage.Foreground(t.AccentPink)

	// Select/multiselect.
	st.Focused.SelectSelector = st.Focused.SelectSelector.Foreground(t.AccentCyan)
	st.Focused.MultiSelectSelector = st.Focused.MultiSelectSelector.Foreground(t.AccentCyan)
	st.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(t.AccentGreen).SetString("[•] ")
	st.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(t.Muted).SetString("[ ] ")
	st.Focused.SelectedOption = lipgloss.NewStyle().Foreground(t.Fg)
	st.Focused.UnselectedOption = lipgloss.NewStyle().Foreground(t.Fg)

	// Text inputs.
	st.Focused.TextInput.Prompt = lipgloss.NewStyle().Foreground(t.AccentCyan)
	st.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(t.Fg)
	st.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(t.Muted)
	st.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(t.Cursor)
	st.Focused.TextInput.CursorText = lipgloss.NewStyle().Foreground(t.CursorText)

	// Buttons.
	st.Focused.FocusedButton = st.Focused.FocusedButton.Foreground(t.Bg).Background(t.AccentMagenta).Bold(true)
	st.Focused.BlurredButton = st.Focused.BlurredButton.Foreground(t.Fg).Background(t.Surface)

	// Notes.
	st.Focused.NoteTitle = lipgloss.NewStyle().Bold(true).Foreground(t.AccentPink).MarginBottom(1)

	// Blurred state mirrors focused but without the thick border and with more muted titles.
	st.Blurred = st.Focused
	st.Blurred.Base = st.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	st.Blurred.Card = st.Blurred.Base
	st.Blurred.Title = lipgloss.NewStyle().Bold(true).Foreground(t.AccentPink)
	st.Blurred.SelectSelector = lipgloss.NewStyle().SetString("  ")
	st.Blurred.MultiSelectSelector = lipgloss.NewStyle().SetString("  ")

	// Help.
	st.Help.ShortKey = st.Help.ShortKey.Foreground(t.Muted)
	st.Help.ShortDesc = st.Help.ShortDesc.Foreground(t.Muted)
	st.Help.FullKey = st.Help.FullKey.Foreground(t.Muted)
	st.Help.FullDesc = st.Help.FullDesc.Foreground(t.Muted)

	return st
}
