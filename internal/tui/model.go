package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
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
	screenConfirmDelete
	screenError
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

	editTarget *secretRow
	editForm   *huh.Form
	editValue  string
	editSave   bool

	createForm        *huh.Form
	createEnvVar      string
	createSecretKey   string
	createValue       string
	createDesc        string
	createReplication string // "global" | "user-managed"
	createLocations   []string
	createAction      string // "create" | "cancel"
	createSummaryTick int
	createSummaryKey  string

	gcpLocations []string

	confirmText string
	deleteForm  *huh.Form
	deleteYes   bool

	errorForm   *huh.Form
	errorReturn Screen
	errorTitle  string
	errorText   string
}

type envItem struct{ name string }

func (e envItem) Title() string       { return e.name }
func (e envItem) Description() string { return "" }
func (e envItem) FilterValue() string { return e.name }

func New(ctx context.Context, configPath string) (*Model, error) {
	cfg, err := config.LoadKubaConfig(configPath)
	if err != nil {
		return nil, err
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

	return &Model{
		ctx:               ctx,
		configPath:        configPath,
		cfg:               cfg,
		screen:            screenEnvs,
		envList:           l,
		secretTable:       t,
		maskValues:        true,
		filter:            filter,
		createReplication: "global",
		createAction:      "create",
	}, nil
}

func (m *Model) Init() tea.Cmd {
	// Ensure we get an initial WindowSizeMsg even on terminals/platforms where
	// Bubble Tea may only deliver size updates after the first resize.
	return func() tea.Msg { return tea.RequestWindowSize() }
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Capture window size globally so we can size models
	// even when switching screens (Bubble Tea only sends this on resize/start).
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.winW = ws.Width
		m.winH = ws.Height
		innerW, innerH := panelInnerSize(ws.Width, ws.Height, panelStyle())
		// Keep the various lists/tables sized to panel inner area.
		m.envList.SetSize(innerW, innerH)
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
	case screenConfirmDelete:
		return m.updateConfirmDelete(msg)
	case screenError:
		return m.updateError(msg)
	default:
		return m, nil
	}
}

func (m *Model) View() tea.View {
	// Render "best effort" at any size. We avoid forcing minimum dimensions and
	// instead shrink content to whatever the terminal can actually show.
	switch m.screen {
	case screenEnvs:
		box := fitPanelToWindow(panelStyle(), m.winW, m.winH)
		v := tea.NewView(box.Render(m.envList.View()))
		v.AltScreen = true
		return v
	case screenSecrets:
		v := tea.NewView(m.viewSecrets())
		v.AltScreen = true
		return v
	case screenView:
		v := tea.NewView(m.viewModal("View secret", m.viewValue+"\n\n(esc to go back)"))
		v.AltScreen = true
		return v
	case screenEdit:
		if m.editForm == nil {
			v := tea.NewView(m.viewModal("Edit secret", "Loading…"))
			v.AltScreen = true
			return v
		}
		v := tea.NewView(m.viewModal("Edit secret", m.editForm.View()))
		v.AltScreen = true
		return v
	case screenCreate:
		if m.createForm == nil {
			v := tea.NewView(m.viewModal("Create secret & mapping", "Loading…"))
			v.AltScreen = true
			return v
		}
		v := tea.NewView(m.viewModal("Create secret & mapping", m.createForm.View()))
		v.AltScreen = true
		return v
	case screenConfirmDelete:
		if m.deleteForm == nil {
			v := tea.NewView(m.viewModal("Confirm delete", "Loading…"))
			v.AltScreen = true
			return v
		}
		v := tea.NewView(m.viewModal("Confirm delete", m.deleteForm.View()))
		v.AltScreen = true
		return v
	case screenError:
		if m.errorForm == nil {
			v := tea.NewView(m.viewModal("Error", "Loading…"))
			v.AltScreen = true
			return v
		}
		title := m.errorTitle
		if strings.TrimSpace(title) == "" {
			title = "Error"
		}
		v := tea.NewView(m.viewModal(title, m.errorForm.View()))
		v.AltScreen = true
		return v
	default:
		v := tea.NewView("")
		v.AltScreen = true
		return v
	}
}

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

func (m *Model) viewSecrets() string {
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

	box := fitPanelToWindow(panelStyle(), m.winW, m.winH)
	return box.Render(content)
}

func (m *Model) viewModal(title, body string) string {
	box := fitPanelToWindow(panelStyle(), m.winW, m.winH)
	return box.Render(lipgloss.NewStyle().Bold(true).Render(title) + "\n\n" + body)
}

func (m *Model) updateSecrets(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		innerW, innerH := panelInnerSize(msg.Width, msg.Height, panelStyle())

		// Size the table to exactly fill available inner height.
		// viewSecrets uses blocks separated by "\n\n", which contributes one empty
		// line between blocks. So total height is:
		// sum(blockHeights) + (numBlocks-1).
		headerH := lipgloss.Height(lipgloss.NewStyle().Bold(true).Render("Environment: X"))
		helpH := lipgloss.Height("enter:view  e:edit  n:new  d:delete  /:filter  m:mask  esc:back  q:quit")
		filterVisible := m.filter.Focused() || strings.TrimSpace(m.filter.Value()) != ""
		filterH := 0
		if filterVisible {
			filterH = 1 // single-line text input / display
		}
		numBlocks := 3
		if filterVisible {
			numBlocks = 4
		}
		nonTableH := headerH + filterH + helpH + (numBlocks - 1)
		m.secretTable.SetWidth(innerW)
		m.setSecretTableColumns(innerW)
		m.secretTable.SetHeight(clampMin(innerH-nonTableH, 1))
		m.filter.SetWidth(clamp(clampMin(innerW, 1), 1, 60))
		return m, nil
	case tea.KeyMsg:
		// Clear any previous error once the user interacts.
		if m.errMsg != "" {
			m.errMsg = ""
		}
		if m.filter.Focused() {
			var cmd tea.Cmd
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc", "enter":
				m.filter.Blur()
				return m, nil
			}
			m.filter, cmd = m.filter.Update(msg)
			m.applyFilterToTable()
			return m, cmd
		}
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
			m.editValue = r.value
			m.editSave = false
			m.editForm = m.newEditForm()
			m.screen = screenEdit
			return m, m.editForm.Init()
		case "d":
			r, ok := m.selectedRow()
			if !ok {
				return m, nil
			}
			if r.refKind != "secret-key" {
				m.errMsg = "delete is only supported for secret-key mappings"
				return m, nil
			}
			m.editTarget = &r
			m.confirmText = fmt.Sprintf("Delete provider secret '%s'?\n\nEnv var: %s\nProvider: %s", r.ref, r.envVar, r.provider)
			m.deleteYes = false
			m.deleteForm = m.newDeleteForm()
			m.screen = screenConfirmDelete
			return m, m.deleteForm.Init()
		case "n":
			m.createEnvVar = ""
			m.createSecretKey = ""
			m.createValue = ""
			m.createDesc = ""
			m.createReplication = "global"
			m.createLocations = nil
			m.createAction = "create"
			m.createSummaryTick = 0
			m.createSummaryKey = ""
			// Lazy-load GCP locations for region multiselect.
			if err := m.ensureGCPLocationsLoaded(); err != nil {
				m.errMsg = err.Error()
				return m, nil
			}
			m.createForm = m.newCreateForm()
			m.screen = screenCreate
			return m, m.createForm.Init()
		}
	}

	var cmd tea.Cmd
	m.secretTable, cmd = m.secretTable.Update(msg)
	return m, cmd
}

func (m *Model) selectedRow() (secretRow, bool) {
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

func (m *Model) updateView(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *Model) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		km := msg.(tea.KeyMsg).String()
		if km == "esc" {
			m.screen = screenSecrets
			m.editTarget = nil
			m.editForm = nil
			return m, nil
		}
		if (km == "alt+up" || km == "alt+down") && m.editForm != nil {
			if km == "alt+up" {
				return m, m.editForm.PrevField()
			}
			return m, m.editForm.NextField()
		}
		if (km == "up" || km == "down") && m.editForm != nil {
			if cmd, ok := formArrowNavCmd(m.editForm, km); ok {
				return m, cmd
			}
		}
	}

	if m.editForm == nil {
		m.editForm = m.newEditForm()
		return m, m.editForm.Init()
	}

	var cmd tea.Cmd
	var mdl huh.Model
	mdl, cmd = m.editForm.Update(msg)
	if f, ok := mdl.(*huh.Form); ok {
		m.editForm = f
	}

	if m.editForm.State == huh.StateCompleted {
		m.screen = screenSecrets
		m.editForm = nil
		if m.editTarget != nil && m.editSave {
			if err := m.saveEdit(*m.editTarget, m.editValue); err != nil {
				m.errMsg = err.Error()
				m.editTarget = nil
				return m, nil
			}
			_ = m.reloadSecrets()
		}
		m.editTarget = nil
	}

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

func (m *Model) updateCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		// Constrain the form to the modal body area so it can scroll internally
		// and the action field remains reachable even on small terminals.
		if m.createForm != nil {
			innerW, innerH := panelInnerSize(ws.Width, ws.Height, panelStyle())
			// viewModal renders: title + "\n\n" + body
			bodyH := clampMin(innerH-2, 1)
			m.createForm = m.createForm.WithWidth(innerW).WithHeight(bodyH)
		}
		return m, nil
	}

	if _, ok := msg.(tea.KeyMsg); ok {
		km := msg.(tea.KeyMsg).String()
		if km == "esc" {
			m.screen = screenSecrets
			m.createForm = nil
			return m, nil
		}
		if (km == "alt+up" || km == "alt+down") && m.createForm != nil {
			if km == "alt+up" {
				return m, m.createForm.PrevField()
			}
			return m, m.createForm.NextField()
		}
		if (km == "up" || km == "down") && m.createForm != nil {
			if cmd, ok := formArrowNavCmd(m.createForm, km); ok {
				return m, cmd
			}
		}
	}

	// Ensure locations exist before rendering the form (so the multiselect has options).
	if err := m.ensureGCPLocationsLoaded(); err != nil {
		return m.openError(screenCreate, "Create failed", err.Error())
	}

	if m.createForm == nil {
		m.createForm = m.newCreateForm()
		if m.winW > 0 && m.winH > 0 {
			innerW, innerH := panelInnerSize(m.winW, m.winH, panelStyle())
			bodyH := clampMin(innerH-2, 1)
			m.createForm = m.createForm.WithWidth(innerW).WithHeight(bodyH)
		}
		return m, m.createForm.Init()
	}

	var cmd tea.Cmd
	var mdl huh.Model
	mdl, cmd = m.createForm.Update(msg)
	if f, ok := mdl.(*huh.Form); ok {
		m.createForm = f
	}
	// Recompute the Summary note only when replication/locations change.
	{
		locs := append([]string(nil), m.createLocations...)
		sort.Strings(locs)
		key := m.createReplication + "|" + strings.Join(locs, ",")
		if key != m.createSummaryKey {
			m.createSummaryKey = key
			m.createSummaryTick++
		}
	}

	if m.createForm.State == huh.StateCompleted {
		switch m.createAction {
		case "create":
			if err := m.doCreateFromForm(); err != nil {
				// Keep current field values, but reset form state so it's usable after closing the error.
				m.createForm = m.newCreateForm()
				m.screen = screenCreate
				return m.openError(screenCreate, "Create failed", err.Error())
			}
			// Reload config + secrets, then return to overview.
			if cfg, err := config.LoadKubaConfig(m.configPath); err == nil {
				m.cfg = cfg
				if env, err := m.cfg.GetEnvironment(m.selectedEnvName); err == nil {
					m.selectedEnv = env
				}
			}
			_ = m.reloadSecrets()
			m.screen = screenSecrets
			m.createForm = nil
			return m, nil
		default:
			// Explicit cancel: return to overview.
			m.screen = screenSecrets
			m.createForm = nil
			return m, nil
		}
	}

	return m, cmd
}

func (m *Model) ensureGCPLocationsLoaded() error {
	if m.selectedEnv == nil || m.selectedEnv.Provider != "gcp" {
		return nil
	}
	if len(m.gcpLocations) > 0 {
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
	return nil
}

func (m *Model) doCreateFromForm() error {
	envVar := strings.TrimSpace(m.createEnvVar)
	secretKey := strings.TrimSpace(m.createSecretKey)
	val := m.createValue
	desc := strings.TrimSpace(m.createDesc)

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
			if m.createReplication == "user-managed" && len(m.createLocations) > 0 {
				gcpSM.SetCreateLocations(m.createLocations)
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

func (m *Model) updateConfirmDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		km := msg.(tea.KeyMsg).String()
		if km == "esc" {
			m.screen = screenSecrets
			m.editTarget = nil
			m.deleteForm = nil
			return m, nil
		}
		if (km == "alt+up" || km == "alt+down") && m.deleteForm != nil {
			if km == "alt+up" {
				return m, m.deleteForm.PrevField()
			}
			return m, m.deleteForm.NextField()
		}
		if (km == "up" || km == "down") && m.deleteForm != nil {
			if cmd, ok := formArrowNavCmd(m.deleteForm, km); ok {
				return m, cmd
			}
		}
	}

	if m.deleteForm == nil {
		m.deleteForm = m.newDeleteForm()
		return m, m.deleteForm.Init()
	}

	var cmd tea.Cmd
	var mdl huh.Model
	mdl, cmd = m.deleteForm.Update(msg)
	if f, ok := mdl.(*huh.Form); ok {
		m.deleteForm = f
	}

	if m.deleteForm.State == huh.StateCompleted {
		m.screen = screenSecrets
		m.deleteForm = nil
		if m.editTarget != nil && m.deleteYes {
			if err := m.doDelete(*m.editTarget); err != nil {
				m.errMsg = err.Error()
			} else {
				_ = m.reloadSecrets()
			}
		}
		m.editTarget = nil
	}

	return m, cmd
}

func (m *Model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		if km := msg.(tea.KeyMsg).String(); km == "esc" {
			m.screen = m.errorReturn
			m.errorForm = nil
			m.errorTitle = ""
			m.errorText = ""
			switch m.screen {
			case screenCreate:
				if m.createForm != nil {
					return m, m.createForm.Init()
				}
			case screenEdit:
				if m.editForm != nil {
					return m, m.editForm.Init()
				}
			case screenConfirmDelete:
				if m.deleteForm != nil {
					return m, m.deleteForm.Init()
				}
			}
			return m, nil
		}
	}

	if m.errorForm == nil {
		m.errorForm = m.newErrorForm(m.errorTitle, m.errorText)
		return m, m.errorForm.Init()
	}

	var cmd tea.Cmd
	var mdl huh.Model
	mdl, cmd = m.errorForm.Update(msg)
	if f, ok := mdl.(*huh.Form); ok {
		m.errorForm = f
	}

	if m.errorForm.State == huh.StateCompleted {
		m.screen = m.errorReturn
		m.errorForm = nil
		m.errorTitle = ""
		m.errorText = ""
		switch m.screen {
		case screenCreate:
			if m.createForm != nil {
				return m, m.createForm.Init()
			}
		case screenEdit:
			if m.editForm != nil {
				return m, m.editForm.Init()
			}
		case screenConfirmDelete:
			if m.deleteForm != nil {
				return m, m.deleteForm.Init()
			}
		}
	}

	return m, cmd
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

func (m *Model) newEditForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Secret value").
				Value(&m.editValue),
			huh.NewConfirm().
				Title("Save changes?").
				Affirmative("Save").
				Negative("Cancel").
				Value(&m.editSave),
		),
	)
}

func (m *Model) newDeleteForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Delete secret").
				Description(mdEscape(m.confirmText)),
			huh.NewConfirm().
				Title("Proceed?").
				Affirmative("Delete").
				Negative("Cancel").
				Value(&m.deleteYes),
		),
	)
}

func (m *Model) newErrorForm(title, text string) *huh.Form {
	if strings.TrimSpace(title) == "" {
		title = "Error"
	}
	if strings.TrimSpace(text) == "" {
		text = "Unknown error"
	}
	back := "back"
	return huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(title).
				Description(text),
			huh.NewSelect[string]().
				Title("").
				Options(huh.NewOption("Back", "back")).
				Value(&back),
		),
	)
}

func (m *Model) openError(returnTo Screen, title, text string) (tea.Model, tea.Cmd) {
	m.errorReturn = returnTo
	m.errorTitle = title
	m.errorText = text
	m.errorForm = m.newErrorForm(title, text)
	m.screen = screenError
	return m, m.errorForm.Init()
}

func (m *Model) newCreateForm() *huh.Form {
	mainFields := []huh.Field{
		huh.NewInput().
			Title("Env var").
			Placeholder("ENV_VAR_NAME").
			Value(&m.createEnvVar).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("env var is required")
				}
				return nil
			}),
		huh.NewInput().
			Title("Secret key/id").
			Placeholder("provider secret key/id/path").
			Value(&m.createSecretKey).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("secret key/id is required")
				}
				return nil
			}),
		huh.NewText().
			Title("Value").
			Value(&m.createValue),
		huh.NewInput().
			Title("Description (optional)").
			Value(&m.createDesc),
	}

	// Keep GCP region selection on the *same page* so it's hard to miss.
	if m.selectedEnv != nil && m.selectedEnv.Provider == "gcp" {
		mainFields = append(mainFields,
			huh.NewSelect[string]().
				Title("Replication").
				Options(
					huh.NewOption("Global (automatic replication)", "global"),
					huh.NewOption("User-managed (choose locations)", "user-managed"),
				).
				Value(&m.createReplication),
			huh.NewMultiSelect[string]().
				TitleFunc(func() string {
					if m.createReplication == "global" {
						return "Locations (set replication to user-managed to choose)"
					}
					return "Locations"
				}, &m.createReplication).
				OptionsFunc(func() []huh.Option[string] {
					if m.createReplication == "global" {
						return nil
					}
					opts := make([]huh.Option[string], 0, len(m.gcpLocations))
					for _, l := range m.gcpLocations {
						opts = append(opts, huh.NewOption(l, l))
					}
					return opts
				}, &m.createReplication).
				Value(&m.createLocations).
				Validate(func(v []string) error {
					if m.createReplication == "user-managed" && len(v) == 0 {
						return fmt.Errorf("select at least one location")
					}
					return nil
				}),
		)
	}

	summary := func() string {
		if m.selectedEnv == nil {
			return ""
		}

		envVar := strings.TrimSpace(m.createEnvVar)
		secretKey := strings.TrimSpace(m.createSecretKey)
		desc := strings.TrimSpace(m.createDesc)

		var b strings.Builder
		b.WriteString(fmt.Sprintf("Environment: %s\n", mdEscape(m.selectedEnvName)))
		b.WriteString(fmt.Sprintf("Provider: %s\n", mdEscape(m.selectedEnv.Provider)))
		b.WriteString(fmt.Sprintf("Project: %s\n", mdEscape(m.selectedEnv.Project)))
		if envVar != "" {
			b.WriteString(fmt.Sprintf("Env var: %s\n", mdEscape(envVar)))
		}
		if secretKey != "" {
			b.WriteString(fmt.Sprintf("Secret key/id: %s\n", mdEscape(secretKey)))
		}
		if desc != "" {
			b.WriteString(fmt.Sprintf("Description: %s\n", mdEscape(desc)))
		}

		if m.selectedEnv.Provider == "gcp" {
			b.WriteString("Replication: ")
			if m.createReplication == "user-managed" {
				b.WriteString("User-managed\n")
				if len(m.createLocations) > 0 {
					b.WriteString("Locations:\n")
					for _, l := range m.createLocations {
						b.WriteString("  - " + mdEscape(l) + "\n")
					}
				} else {
					b.WriteString("Locations: (none)\n")
				}
			} else {
				b.WriteString("Global\n")
			}
		}

		return strings.TrimSpace(b.String())
	}

	mainFields = append(mainFields,
		huh.NewNote().
			Title("Summary").
			DescriptionFunc(summary, &m.createSummaryTick),
		huh.NewSelect[string]().
			Title("Action").
			Options(
				huh.NewOption("Create secret & mapping", "create"),
				huh.NewOption("Cancel", "cancel"),
			).
			Value(&m.createAction),
	)

	return huh.NewForm(huh.NewGroup(mainFields...))
}

func mdEscape(s string) string {
	// huh Note fields render markdown; escape characters that would be
	// interpreted as formatting (notably "_" in ENV_VAR names).
	//
	// We keep this intentionally minimal: the goal is to preserve the exact
	// visible value rather than apply markdown styling.
	repl := strings.NewReplacer(
		"\\", "\\\\",
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
	)
	return repl.Replace(s)
}

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
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
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

func (m *Model) setSecretTableColumns(innerW int) {
	// bubbles/table does not automatically reflow column widths when SetWidth is
	// called, so we recompute widths to avoid horizontal overflow on small terminals.
	//
	// Note: bubbles/table renders each cell with padding and inserts spacing between
	// columns. Column.Width applies to the cell content area, but the final rendered
	// width includes extra "chrome". We therefore reserve a conservative overhead
	// budget so the table never exceeds innerW.
	const cols = 4
	// Default table styles include left+right cell padding (commonly 1 each),
	// plus at least 1 char gap between columns.
	cellPaddingLR := 2 // 1 left + 1 right (conservative)
	colGaps := cols - 1
	overhead := cols*cellPaddingLR + colGaps
	avail := innerW - overhead
	if avail < cols { // at least 1 char per column
		avail = cols
	}

	// Minimums for usability.
	minEnv, minVal, minProv, minRef := 8, 10, 8, 10
	minTotal := minEnv + minVal + minProv + minRef
	if avail < minTotal {
		// If extremely narrow, shrink everything but keep provider readable-ish.
		minEnv, minVal, minRef = 6, 6, 6
		minTotal = minEnv + minVal + minProv + minRef
	}

	envW, valW, provW, refW := minEnv, minVal, minProv, minRef
	remaining := avail - (envW + valW + provW + refW)
	if remaining < 0 {
		remaining = 0
	}

	// Distribute extra width with a bias towards Value and Ref.
	// Order: Value, Ref, Env Var, Provider.
	add := func(w *int, n int) {
		if n <= 0 {
			return
		}
		*w += n
		remaining -= n
	}
	for remaining > 0 {
		add(&valW, min(remaining, 2))
		if remaining == 0 {
			break
		}
		add(&refW, min(remaining, 2))
		if remaining == 0 {
			break
		}
		add(&envW, 1)
		if remaining == 0 {
			break
		}
		add(&provW, 1)
	}

	m.secretTable.SetColumns([]table.Column{
		{Title: "Env Var", Width: envW},
		{Title: "Value", Width: valW},
		{Title: "Provider", Width: provW},
		{Title: "Ref", Width: refW},
	})
}
