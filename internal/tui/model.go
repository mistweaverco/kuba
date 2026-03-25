package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistweaverco/kuba/internal/config"
	"github.com/mistweaverco/kuba/internal/lib/secrets"
)

type Screen int

const (
	screenEnvs Screen = iota
	screenSecrets
	screenView
	screenEdit
	screenCreate
	screenPickGCPLocation
	screenConfirmDelete
)

type secretRow struct {
	envVar   string
	value    string
	item     config.EnvItem
	provider string
	project  string
	refKind  string // secret-key | secret-path | value
	ref      string // secret-key or secret-path
}

type Model struct {
	ctx        context.Context
	configPath string

	cfg *config.KubaConfig

	screen Screen
	errMsg string

	winW int
	winH int

	envList list.Model

	selectedEnvName string
	selectedEnv     *config.Environment

	secretTable table.Model
	allRows     []secretRow
	maskValues  bool

	filter textinput.Model

	viewValue string

	editArea   textarea.Model
	editTarget *secretRow

	createEnvVar    textinput.Model
	createSecretKey textinput.Model
	createValue     textarea.Model
	createDesc      textinput.Model

	createGCPUseGlobal bool
	createGCPLocations []string
	gcpLocations       []string
	gcpLocList         list.Model

	confirmText string
}

type envItem struct{ name string }

func (e envItem) Title() string       { return e.name }
func (e envItem) Description() string { return "" }
func (e envItem) FilterValue() string { return e.name }

func New(ctx context.Context, configPath string) (Model, error) {
	cfg, err := config.LoadKubaConfig(configPath)
	if err != nil {
		return Model{}, err
	}

	envNames := make([]string, 0, len(cfg.Environments))
	for name := range cfg.Environments {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	items := make([]list.Item, 0, len(envNames))
	for _, n := range envNames {
		items = append(items, envItem{name: n})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Environments"
	l.SetShowHelp(true)

	filter := textinput.New()
	filter.Placeholder = "Filter secrets…"
	filter.CharLimit = 256
	filter.Prompt = "/ "

	editArea := textarea.New()
	editArea.Placeholder = "Secret value"
	editArea.ShowLineNumbers = false
	editArea.CharLimit = 0

	createEnvVar := textinput.New()
	createEnvVar.Placeholder = "ENV_VAR_NAME"
	createEnvVar.CharLimit = 256

	createSecretKey := textinput.New()
	createSecretKey.Placeholder = "provider secret key/id/path"
	createSecretKey.CharLimit = 512

	createValue := textarea.New()
	createValue.Placeholder = "Secret value"
	createValue.ShowLineNumbers = false
	createValue.CharLimit = 0

	createDesc := textinput.New()
	createDesc.Placeholder = "(optional) description"
	createDesc.CharLimit = 512

	gcpLocList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	gcpLocList.Title = "GCP Secret location"
	gcpLocList.SetShowHelp(true)
	gcpLocList.SetShowFilter(true)
	gcpLocList.SetFilteringEnabled(true)
	gcpLocList.SetShowStatusBar(true)

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Env Var", Width: 28},
			{Title: "Value", Width: 32},
			{Title: "Provider", Width: 10},
			{Title: "Ref", Width: 28},
		}),
		table.WithRows(nil),
		table.WithFocused(true),
	)

	return Model{
		ctx:                ctx,
		configPath:         configPath,
		cfg:                cfg,
		screen:             screenEnvs,
		envList:            l,
		secretTable:        t,
		maskValues:         true,
		filter:             filter,
		editArea:           editArea,
		createEnvVar:       createEnvVar,
		createSecretKey:    createSecretKey,
		createValue:        createValue,
		createDesc:         createDesc,
		createGCPUseGlobal: true,
		createGCPLocations: nil,
		gcpLocations:       nil,
		gcpLocList:         gcpLocList,
	}, nil
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Capture window size globally so we can size models
	// even when switching screens (Bubble Tea only sends this on resize/start).
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.winW = ws.Width
		m.winH = ws.Height
		innerW, innerH := panelInnerSize(ws.Width, ws.Height)
		// Keep the various lists/tables sized to panel inner area.
		m.envList.SetSize(innerW, innerH)
		m.gcpLocList.SetSize(innerW, innerH)
	}

	switch m.screen {
	case screenEnvs:
		return m.updateEnvs(msg)
	case screenSecrets:
		return m.updateSecrets(msg)
	case screenView:
		return m.updateView(msg)
	case screenEdit:
		return m.updateEdit(msg)
	case screenCreate:
		return m.updateCreate(msg)
	case screenPickGCPLocation:
		return m.updatePickGCPLocation(msg)
	case screenConfirmDelete:
		return m.updateConfirmDelete(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.screen {
	case screenEnvs:
		box := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2)
		if m.winW > 0 {
			box = box.Width(max(20, m.winW-4))
		}
		return box.Render(m.envList.View())
	case screenSecrets:
		return m.viewSecrets()
	case screenView:
		return m.viewModal("View secret", m.viewValue+"\n\n(esc to go back)")
	case screenEdit:
		return m.viewModal("Edit secret (ctrl+s to save)", m.editArea.View()+"\n\n(esc to cancel)")
	case screenCreate:
		body := strings.Builder{}
		body.WriteString("Env var:\n")
		body.WriteString(m.createEnvVar.View())
		body.WriteString("\n\nSecret key/id:\n")
		body.WriteString(m.createSecretKey.View())
		body.WriteString("\n\nValue:\n")
		body.WriteString(m.createValue.View())
		body.WriteString("\n\nDescription:\n")
		body.WriteString(m.createDesc.View())

		if m.selectedEnv != nil && m.selectedEnv.Provider == "gcp" {
			body.WriteString("\n\nReplication:\n")
			if m.createGCPUseGlobal {
				body.WriteString("  (x) Global (automatic replication)\n  ( ) User-managed location\n")
			} else {
				body.WriteString("  ( ) Global (automatic replication)\n  (x) User-managed location\n")
			}
			loc := "(none selected)"
			if len(m.createGCPLocations) == 1 {
				loc = m.createGCPLocations[0]
			} else if len(m.createGCPLocations) > 1 {
				loc = strings.Join(m.createGCPLocations, ", ")
			}
			body.WriteString("\nLocations: " + loc + "\n")
			body.WriteString("Tip: ctrl+g toggle global/user-managed, ctrl+l pick locations\n")
		}

		body.WriteString("\n(ctrl+s to create, esc to cancel)")
		return m.viewModal("Create secret & mapping", body.String())
	case screenPickGCPLocation:
		return m.gcpLocList.View()
	case screenConfirmDelete:
		return m.viewModal("Confirm delete", m.confirmText+"\n\n(y)es / (n)o")
	default:
		return ""
	}
}

func (m Model) updateEnvs(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		innerW, innerH := panelInnerSize(msg.Width, msg.Height)
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
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.envList, cmd = m.envList.Update(msg)
	return m, cmd
}

func (m *Model) reloadSecrets() error {
	if m.selectedEnv == nil {
		return fmt.Errorf("no environment selected")
	}

	factory := secrets.NewSecretManagerFactory()
	values, err := factory.GetSecretsForEnvironmentWithCache(m.ctx, m.selectedEnv, m.configPath, m.selectedEnvName)
	if err != nil {
		return err
	}

	items := m.selectedEnv.GetEnvItems()
	rows := make([]secretRow, 0, len(items))
	for _, it := range items {
		provider := it.Provider
		if provider == "" {
			provider = m.selectedEnv.Provider
		}
		project := it.Project
		if project == "" {
			project = m.selectedEnv.Project
		}

		refKind := "value"
		ref := ""
		if it.SecretKey != "" {
			refKind = "secret-key"
			ref = it.SecretKey
		} else if it.SecretPath != "" {
			refKind = "secret-path"
			ref = it.SecretPath
		}

		val := values[it.EnvironmentVariable]
		rows = append(rows, secretRow{
			envVar:   it.EnvironmentVariable,
			value:    val,
			item:     it,
			provider: provider,
			project:  project,
			refKind:  refKind,
			ref:      ref,
		})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].envVar < rows[j].envVar })
	m.allRows = rows
	m.applyFilterToTable()
	return nil
}

func (m *Model) applyFilterToTable() {
	q := strings.TrimSpace(strings.ToLower(m.filter.Value()))

	filtered := make([]secretRow, 0, len(m.allRows))
	for _, r := range m.allRows {
		if q == "" || strings.Contains(strings.ToLower(r.envVar), q) || strings.Contains(strings.ToLower(r.ref), q) {
			filtered = append(filtered, r)
		}
	}

	trows := make([]table.Row, 0, len(filtered))
	for _, r := range filtered {
		val := r.value
		if m.maskValues {
			val = mask(val)
		}
		ref := r.ref
		if ref == "" {
			ref = r.refKind
		} else {
			ref = r.refKind + ":" + ref
		}
		trows = append(trows, table.Row{r.envVar, val, r.provider, ref})
	}
	m.secretTable.SetRows(trows)
}

func mask(v string) string {
	if v == "" {
		return ""
	}
	if len(v) <= 4 {
		return strings.Repeat("•", len(v))
	}
	return strings.Repeat("•", 8)
}

func (m Model) viewSecrets() string {
	header := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Environment: %s", m.selectedEnvName))
	help := "enter:view  e:edit  n:new  d:delete  /:filter  m:mask  esc:back  q:quit"
	if m.errMsg != "" {
		help = "Error: " + m.errMsg + "\n" + help
	}

	filterLine := ""
	if m.filter.Focused() {
		filterLine = m.filter.View()
	} else {
		if v := m.filter.Value(); strings.TrimSpace(v) != "" {
			filterLine = "/ " + v
		} else {
			filterLine = ""
		}
	}

	parts := []string{header}
	if filterLine != "" {
		parts = append(parts, filterLine)
	}
	parts = append(parts, m.secretTable.View(), help)
	content := strings.Join(parts, "\n\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2)
	if m.winW > 0 {
		box = box.Width(max(20, m.winW-4))
	}
	return box.Render(content)
}

func (m Model) viewModal(title, body string) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2)

	if m.winW > 0 {
		// Leave a small gutter so borders don't clip.
		box = box.Width(max(20, m.winW-4))
	}
	if m.winH > 0 {
		// Use most of the available height (best-effort; content may still scroll).
		box = box.Height(max(8, m.winH-4))
	}
	return box.Render(lipgloss.NewStyle().Bold(true).Render(title) + "\n\n" + body)
}

func (m Model) updateSecrets(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.errMsg = ""

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// reserve space for header/help
		innerW, innerH := panelInnerSize(msg.Width, msg.Height)
		m.secretTable.SetWidth(innerW)
		m.secretTable.SetHeight(max(4, innerH-6))
		m.filter.Width = min(60, msg.Width-6)
		m.editArea.SetWidth(min(90, msg.Width-10))
		m.editArea.SetHeight(max(6, msg.Height-12))
		m.createValue.SetWidth(min(90, msg.Width-10))
		m.createValue.SetHeight(max(6, msg.Height-16))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.filter.Focused() {
				m.filter.Blur()
				return m, nil
			}
			m.screen = screenEnvs
			return m, nil
		case "/":
			m.filter.Focus()
			return m, nil
		case "m":
			m.maskValues = !m.maskValues
			m.applyFilterToTable()
			return m, nil
		case "enter":
			r, ok := m.selectedRow()
			if !ok {
				return m, nil
			}
			m.viewValue = fmt.Sprintf("%s\n\n%s", r.envVar, r.value)
			m.screen = screenView
			return m, nil
		case "e":
			r, ok := m.selectedRow()
			if !ok {
				return m, nil
			}
			if r.refKind != "secret-key" {
				m.errMsg = "edit is only supported for secret-key mappings"
				return m, nil
			}
			m.editTarget = &r
			m.editArea.SetValue(r.value)
			m.editArea.Focus()
			m.screen = screenEdit
			return m, nil
		case "d":
			r, ok := m.selectedRow()
			if !ok {
				return m, nil
			}
			if r.refKind != "secret-key" {
				m.errMsg = "delete is only supported for secret-key mappings"
				return m, nil
			}
			m.confirmText = fmt.Sprintf("Delete provider secret '%s'?\n\nEnv var: %s\nProvider: %s", r.ref, r.envVar, r.provider)
			m.editTarget = &r
			m.screen = screenConfirmDelete
			return m, nil
		case "n":
			m.createEnvVar.SetValue("")
			m.createSecretKey.SetValue("")
			m.createValue.SetValue("")
			m.createDesc.SetValue("")
			m.createGCPUseGlobal = true
			m.createGCPLocations = nil
			m.createEnvVar.Focus()
			m.createSecretKey.Blur()
			m.createDesc.Blur()
			m.createValue.Blur()
			m.screen = screenCreate
			return m, nil
		}
	}

	if m.filter.Focused() {
		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		m.applyFilterToTable()
		return m, cmd
	}

	var cmd tea.Cmd
	m.secretTable, cmd = m.secretTable.Update(msg)
	return m, cmd
}

func (m Model) selectedRow() (secretRow, bool) {
	i := m.secretTable.Cursor()
	if i < 0 || i >= len(m.secretTable.Rows()) {
		return secretRow{}, false
	}
	envVar := m.secretTable.Rows()[i][0]
	for _, r := range m.allRows {
		if r.envVar == envVar {
			return r, true
		}
	}
	return secretRow{}, false
}

func (m Model) updateView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			m.screen = screenSecrets
			return m, nil
		}
	}
	return m, nil
}

func (m Model) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.screen = screenSecrets
			m.editTarget = nil
			return m, nil
		case "ctrl+s":
			if m.editTarget == nil {
				m.screen = screenSecrets
				return m, nil
			}
			if err := m.saveEdit(*m.editTarget, m.editArea.Value()); err != nil {
				m.errMsg = err.Error()
				m.screen = screenSecrets
				return m, nil
			}
			_ = m.reloadSecrets()
			m.screen = screenSecrets
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.editArea, cmd = m.editArea.Update(msg)
	return m, cmd
}

func (m *Model) saveEdit(row secretRow, newValue string) error {
	factory := secrets.NewSecretManagerFactory()
	sm, err := factory.CreateSecretManager(m.ctx, row.provider, row.project)
	if err != nil {
		return err
	}
	defer sm.Close()

	mut, err := secrets.AsMutator(sm)
	if err != nil {
		return err
	}

	return mut.UpdateSecret(row.ref, newValue)
}

func (m Model) updateCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.createEnvVar.Width = min(80, msg.Width-10)
		m.createSecretKey.Width = min(80, msg.Width-10)
		m.createDesc.Width = min(80, msg.Width-10)
		m.createValue.SetWidth(min(90, msg.Width-10))
		m.createValue.SetHeight(max(6, msg.Height-16))
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.screen = screenSecrets
			return m, nil
		case "ctrl+g":
			if m.selectedEnv != nil && m.selectedEnv.Provider == "gcp" {
				m.createGCPUseGlobal = !m.createGCPUseGlobal
				if m.createGCPUseGlobal {
					m.createGCPLocations = nil
				}
				return m, nil
			}
		case "ctrl+l":
			if m.selectedEnv != nil && m.selectedEnv.Provider == "gcp" && !m.createGCPUseGlobal {
				if err := m.loadGCPLocationsForCreate(); err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
				m.screen = screenPickGCPLocation
				if m.winW > 0 && m.winH > 0 {
					m.gcpLocList.SetSize(m.winW, m.winH)
				}
				return m, nil
			}
		case "ctrl+s":
			if err := m.doCreate(); err != nil {
				m.errMsg = err.Error()
				m.screen = screenSecrets
				return m, nil
			}
			// Reload config + secrets
			cfg, err := config.LoadKubaConfig(m.configPath)
			if err == nil {
				m.cfg = cfg
				env, _ := m.cfg.GetEnvironment(m.selectedEnvName)
				m.selectedEnv = env
			}
			_ = m.reloadSecrets()
			m.screen = screenSecrets
			return m, nil
		case "enter":
			// advance focus
			if m.createEnvVar.Focused() {
				m.createEnvVar.Blur()
				m.createSecretKey.Focus()
				return m, nil
			}
			if m.createSecretKey.Focused() {
				m.createSecretKey.Blur()
				m.createValue.Focus()
				return m, nil
			}
			if m.createValue.Focused() {
				m.createValue.Blur()
				m.createDesc.Focus()
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	switch {
	case m.createEnvVar.Focused():
		m.createEnvVar, cmd = m.createEnvVar.Update(msg)
	case m.createSecretKey.Focused():
		m.createSecretKey, cmd = m.createSecretKey.Update(msg)
	case m.createValue.Focused():
		m.createValue, cmd = m.createValue.Update(msg)
	case m.createDesc.Focused():
		m.createDesc, cmd = m.createDesc.Update(msg)
	default:
		m.createEnvVar.Focus()
	}
	return m, cmd
}

type locItem struct {
	raw      string
	selected bool
}

func (l locItem) Title() string {
	if l.selected {
		return "[x] " + l.raw
	}
	return "[ ] " + l.raw
}
func (l locItem) Description() string { return "" }
func (l locItem) FilterValue() string { return l.raw }

func (m *Model) loadGCPLocationsForCreate() error {
	if m.selectedEnv == nil || m.selectedEnv.Provider != "gcp" {
		return nil
	}
	if m.gcpLocations != nil && len(m.gcpLocations) > 0 {
		return nil
	}

	factory := secrets.NewSecretManagerFactory()
	sm, err := factory.CreateSecretManager(m.ctx, "gcp", m.selectedEnv.Project)
	if err != nil {
		return err
	}
	defer sm.Close()

	gcpSM, ok := sm.(*secrets.GCPSecretManager)
	if !ok {
		return fmt.Errorf("unexpected gcp secret manager type")
	}

	locs, err := gcpSM.SupportedLocations(m.selectedEnv.Project)
	if err != nil {
		return err
	}
	m.gcpLocations = locs

	items := make([]list.Item, 0, len(locs))
	for _, l := range locs {
		items = append(items, locItem{raw: l, selected: false})
	}
	m.gcpLocList.SetItems(items)
	if m.winW > 0 && m.winH > 0 {
		m.gcpLocList.SetSize(m.winW, m.winH)
	}
	_ = m.refreshGcpLocListTitles()
	return nil
}

func (m Model) updatePickGCPLocation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.gcpLocList.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// pass-through for filter input
			if m.gcpLocList.FilterInput.Focused() {
				var cmd tea.Cmd
				m.gcpLocList, cmd = m.gcpLocList.Update(msg)
				return m, cmd
			}
			m.screen = screenCreate
			return m, nil
		case "ctrl+x":
			if it, ok := m.gcpLocList.SelectedItem().(locItem); ok {
				m.toggleGCPLocation(it.raw)
				cmd := m.refreshGcpLocListTitles()
				return m, cmd
			}
		case "enter":
			m.screen = screenCreate
			return m, nil
		case "ctrl+a":
			// pass-through for filter input
			if m.gcpLocList.FilterInput.Focused() {
				var cmd tea.Cmd
				m.gcpLocList, cmd = m.gcpLocList.Update(msg)
				return m, cmd
			}
			// select all
			m.createGCPLocations = append([]string(nil), m.gcpLocations...)
			cmd := m.refreshGcpLocListTitles()
			return m, cmd
		case "ctrl+c":
			// pass-through for filter input
			if m.gcpLocList.FilterInput.Focused() {
				var cmd tea.Cmd
				m.gcpLocList, cmd = m.gcpLocList.Update(msg)
				return m, cmd
			}
			// clear
			m.createGCPLocations = nil
			cmd := m.refreshGcpLocListTitles()
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.gcpLocList, cmd = m.gcpLocList.Update(msg)
	return m, cmd
}

func (m *Model) toggleGCPLocation(loc string) {
	for i, v := range m.createGCPLocations {
		if v == loc {
			m.createGCPLocations = append(m.createGCPLocations[:i], m.createGCPLocations[i+1:]...)
			return
		}
	}
	m.createGCPLocations = append(m.createGCPLocations, loc)
	sort.Strings(m.createGCPLocations)
}

func (m *Model) refreshGcpLocListTitles() tea.Cmd {
	// Preserve current filter text/state so toggling doesn't wipe it.
	filterText := m.gcpLocList.FilterInput.Value()

	// Update list items so titles show selection state.
	items := make([]list.Item, 0, len(m.gcpLocations))
	selected := make(map[string]bool, len(m.createGCPLocations))
	for _, l := range m.createGCPLocations {
		selected[l] = true
	}
	for _, l := range m.gcpLocations {
		items = append(items, locItem{raw: l, selected: selected[l]})
	}
	cmd := m.gcpLocList.SetItems(items)

	// Restore filter text (SetItems can cause internal filter recalcs).
	m.gcpLocList.SetFilterText(filterText)
	if strings.TrimSpace(filterText) == "" {
		m.gcpLocList.SetFilterState(list.Unfiltered)
	} else {
		m.gcpLocList.SetFilterState(list.FilterApplied)
	}

	// If filtering is active, clamp cursor against visible (filtered) items.
	// Otherwise clamp against full list.
	visible := m.gcpLocList.VisibleItems()
	if len(visible) > 0 {
		cursor := m.gcpLocList.Cursor()
		if cursor < 0 {
			cursor = 0
		}
		if cursor >= len(visible) {
			cursor = len(visible) - 1
		}
		m.gcpLocList.Select(cursor)
	} else {
		m.gcpLocList.ResetSelected()
	}

	if m.winW > 0 && m.winH > 0 {
		m.gcpLocList.SetSize(m.winW, m.winH)
	}
	// Keep cursor stable-ish (Bubble list will clamp).
	m.gcpLocList.Title = fmt.Sprintf("GCP Secret locations (%d selected) — toggle: ctrl+x, done: enter, all: ctrl+a, clear: ctrl+c", len(m.createGCPLocations))

	return cmd
}

func (m *Model) doCreate() error {
	envVar := strings.TrimSpace(m.createEnvVar.Value())
	secretKey := strings.TrimSpace(m.createSecretKey.Value())
	val := m.createValue.Value()
	desc := strings.TrimSpace(m.createDesc.Value())

	if envVar == "" || secretKey == "" {
		return fmt.Errorf("env var and secret key are required")
	}

	// Create secret in provider for this environment.
	provider := m.selectedEnv.Provider
	project := m.selectedEnv.Project

	factory := secrets.NewSecretManagerFactory()
	sm, err := factory.CreateSecretManager(m.ctx, provider, project)
	if err != nil {
		return err
	}
	defer sm.Close()

	mut, err := secrets.AsMutator(sm)
	if err != nil {
		return err
	}

	// If GCP, optionally set per-create locations for replication on the manager instance.
	if provider == "gcp" {
		if gcpSM, ok := sm.(*secrets.GCPSecretManager); ok {
			if !m.createGCPUseGlobal && len(m.createGCPLocations) > 0 {
				gcpSM.SetCreateLocations(m.createGCPLocations)
			} else {
				gcpSM.SetCreateLocations(nil)
			}
		}
	}

	if err := mut.CreateSecret(secretKey, val, desc); err != nil {
		return err
	}

	// Add mapping to kuba.yaml.
	if err := config.AddOrUpdateEnvSecretKeyMapping(m.configPath, m.selectedEnvName, envVar, secretKey); err != nil {
		return err
	}

	return nil
}

func (m Model) updateConfirmDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "esc":
			m.screen = screenSecrets
			m.editTarget = nil
			return m, nil
		case "y":
			if m.editTarget != nil {
				if err := m.doDelete(*m.editTarget); err != nil {
					m.errMsg = err.Error()
				} else {
					_ = m.reloadSecrets()
				}
			}
			m.screen = screenSecrets
			m.editTarget = nil
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) doDelete(row secretRow) error {
	factory := secrets.NewSecretManagerFactory()
	sm, err := factory.CreateSecretManager(m.ctx, row.provider, row.project)
	if err != nil {
		return err
	}
	defer sm.Close()

	mut, err := secrets.AsMutator(sm)
	if err != nil {
		return err
	}

	if err := mut.DeleteSecret(row.ref, true); err != nil {
		return err
	}

	// Also remove the mapping from kuba.yaml so it doesn't reappear on refresh.
	if err := config.RemoveEnvMapping(m.configPath, m.selectedEnvName, row.envVar); err != nil {
		return err
	}

	// Reload config/env so subsequent actions use updated inheritance/mappings.
	if cfg, err := config.LoadKubaConfig(m.configPath); err == nil {
		m.cfg = cfg
		if env, err := m.cfg.GetEnvironment(m.selectedEnvName); err == nil {
			m.selectedEnv = env
		}
	}

	return nil
}

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

func panelInnerSize(w, h int) (int, int) {
	// Border takes 2 cols/rows, padding(1,2) takes 4 cols and 2 rows.
	innerW := w - 8
	innerH := h - 4
	if innerW < 10 {
		innerW = 10
	}
	if innerH < 4 {
		innerH = 4
	}
	return innerW, innerH
}
