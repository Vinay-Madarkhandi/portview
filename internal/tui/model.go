// Package tui implements the terminal user interface for PortView.
package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Vinay-Madarkhandi/portview/internal/clipboard"
	"github.com/Vinay-Madarkhandi/portview/internal/config"
	"github.com/Vinay-Madarkhandi/portview/internal/exporter"
	"github.com/Vinay-Madarkhandi/portview/internal/scanner"
	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

const (
	defaultRefreshInterval = 3 * time.Second
)

var copyText = clipboard.Write
var exportPorts = exporter.Export

type sortField int

const (
	sortByPort sortField = iota
	sortByProcess
	sortByProtocol
	sortByPID
	sortFieldCount
)

type protocolFilter int

const (
	filterAll protocolFilter = iota
	filterTCP
	filterUDP
	protocolFilterCount
)

func (s sortField) String() string {
	switch s {
	case sortByPort:
		return "port"
	case sortByProcess:
		return "process"
	case sortByProtocol:
		return "protocol"
	case sortByPID:
		return "pid"
	default:
		return "port"
	}
}

func (f protocolFilter) String() string {
	switch f {
	case filterTCP:
		return "tcp"
	case filterUDP:
		return "udp"
	default:
		return "all"
	}
}

type tickMsg time.Time

type scanResultMsg struct {
	ports []types.PortInfo
	err   error
}

type statusMsg string

type Model struct {
	table       table.Model
	ports       []types.PortInfo
	err         error
	width       int
	height      int
	statusMsg   string
	sortBy      sortField
	filter      protocolFilter
	search      string
	searching   bool
	pendingKill table.Row
	interval    time.Duration
	ready       bool
}

// InitialModel creates the initial application model.
func InitialModel() Model {
	cfg, err := config.Load()
	m := initialModelWithConfig(cfg)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Config: %v", err)
	}
	return m
}

func initialModelWithConfig(cfg config.Config) Model {
	cfg.ApplyDefaults()
	configureStyles(cfg)

	columns := defaultColumns(80)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(currentPalette.accent)).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(currentPalette.accent))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(currentPalette.headerForeground)).
		Background(lipgloss.Color(currentPalette.accent)).
		Bold(true)
	t.SetStyles(s)

	return Model{
		table:    t,
		sortBy:   sortByPort,
		interval: time.Duration(cfg.RefreshIntervalSeconds) * time.Second,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		doScan,
		tickCmd(m.refreshInterval()),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.table.SetColumns(defaultColumns(msg.Width))
		tableHeight := msg.Height - 7
		if tableHeight < 3 {
			tableHeight = 3
		}
		m.table.SetHeight(tableHeight)

	case tickMsg:
		cmds = append(cmds, doScan, tickCmd(m.refreshInterval()))

	case scanResultMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.ports = msg.ports
			m.sortPorts()
			m.refreshRows()
		}

	case statusMsg:
		m.statusMsg = string(msg)
		cmds = append(cmds, clearStatusAfter(3*time.Second))

	case clearStatusMsg:
		m.statusMsg = ""

	case tea.KeyMsg:
		if m.searching {
			return m.updateSearch(msg)
		}
		if m.pendingKill != nil {
			return m.updateKillConfirmation(msg)
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Refresh):
			m.statusMsg = "Refreshing..."
			cmds = append(cmds, doScan, clearStatusAfter(1*time.Second))

		case key.Matches(msg, keys.Kill):
			statusText, pending := requestKillConfirmation(m.table.SelectedRow())
			m.statusMsg = statusText
			m.pendingKill = pending
			if pending == nil {
				cmds = append(cmds, clearStatusAfter(2*time.Second))
			}

		case key.Matches(msg, keys.CopyPort):
			m.statusMsg = copySelectedValue(m.table.SelectedRow(), 1, "port")
			cmds = append(cmds, clearStatusAfter(2*time.Second))

		case key.Matches(msg, keys.CopyPID):
			m.statusMsg = copySelectedValue(m.table.SelectedRow(), 4, "PID")
			cmds = append(cmds, clearStatusAfter(2*time.Second))

		case key.Matches(msg, keys.ExportCSV):
			m.statusMsg = exportVisible(exporter.FormatCSV, m.filteredPorts())
			cmds = append(cmds, clearStatusAfter(3*time.Second))

		case key.Matches(msg, keys.ExportJSON):
			m.statusMsg = exportVisible(exporter.FormatJSON, m.filteredPorts())
			cmds = append(cmds, clearStatusAfter(3*time.Second))

		case key.Matches(msg, keys.Sort):
			m.sortBy = (m.sortBy + 1) % sortFieldCount
			m.sortPorts()
			m.refreshRows()
			m.statusMsg = fmt.Sprintf("Sorted by %s", m.sortBy)
			cmds = append(cmds, clearStatusAfter(2*time.Second))

		case key.Matches(msg, keys.Filter):
			m.filter = (m.filter + 1) % protocolFilterCount
			m.refreshRows()
			m.statusMsg = fmt.Sprintf("Showing %s ports", m.filter)
			cmds = append(cmds, clearStatusAfter(2*time.Second))

		case key.Matches(msg, keys.Search):
			m.searching = true
			m.statusMsg = ""

		default:
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading PortView..."
	}

	var b strings.Builder

	header := HeaderStyle.Width(m.width).Render(
		TitleText + " " + SubtitleText,
	)
	b.WriteString(header)
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("  ⚠ "+m.err.Error()) + "\n\n")
	}

	b.WriteString(m.table.View())
	b.WriteString("\n")

	if m.searching {
		search := m.search
		if search == "" {
			search = "type process or port"
		}
		b.WriteString("  " + StatusStyle.Render("Search: "+search) + "\n")
	} else if m.statusMsg != "" {
		b.WriteString("  " + StatusStyle.Render(m.statusMsg) + "\n")
	} else {
		visible := m.filteredPorts()
		visibleCount := len(visible)
		totalCount := len(m.ports)
		searchInfo := ""
		if m.search != "" {
			searchInfo = fmt.Sprintf(" · search %q", m.search)
		}
		info := HelpDescStyle.Render(fmt.Sprintf("  %d/%d port(s) · filter %s%s · sorted by %s · auto-refresh %s", visibleCount, totalCount, m.filter, searchInfo, m.sortBy, m.refreshInterval()))
		b.WriteString(info + "\n")
	}

	b.WriteString(renderHelpBar(m.width))

	return b.String()
}

func tickCmd(interval time.Duration) tea.Cmd {
	if interval <= 0 {
		interval = defaultRefreshInterval
	}
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) refreshInterval() time.Duration {
	if m.interval <= 0 {
		return defaultRefreshInterval
	}
	return m.interval
}

func doScan() tea.Msg {
	ports, err := scanner.ScanPorts()
	return scanResultMsg{ports: ports, err: err}
}

type clearStatusMsg struct{}

func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func killSelected(selected table.Row) (string, tea.Cmd) {
	if selected == nil {
		return "No row selected", clearStatusAfter(2 * time.Second)
	}

	pidStr := selected[4] // PID column
	if pidStr == "-" || pidStr == "" {
		return "No PID available for this port", clearStatusAfter(2 * time.Second)
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Sprintf("Invalid PID: %s", pidStr), clearStatusAfter(2 * time.Second)
	}

	if err := scanner.KillProcess(pid); err != nil {
		return fmt.Sprintf("✗ %s", err.Error()), clearStatusAfter(3 * time.Second)
	}

	status := fmt.Sprintf("✓ Requested termination for PID %d (%s)", pid, selected[3])

	return status, tea.Batch(
		clearStatusAfter(3*time.Second),
		delayedScan(500*time.Millisecond),
	)
}

func requestKillConfirmation(selected table.Row) (string, table.Row) {
	if selected == nil {
		return "No row selected", nil
	}

	pidStr := selected[4]
	if pidStr == "-" || pidStr == "" {
		return "No PID available for this port", nil
	}

	if _, err := strconv.Atoi(pidStr); err != nil {
		return fmt.Sprintf("Invalid PID: %s", pidStr), nil
	}

	return fmt.Sprintf("Kill PID %s (%s)? y/Enter confirm, n/Esc cancel", pidStr, selected[3]), selected
}

func (m Model) updateKillConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyCtrlC:
		return m, tea.Quit
	case msg.Type == tea.KeyEnter || strings.EqualFold(msg.String(), "y"):
		statusText, cmd := killSelected(m.pendingKill)
		m.pendingKill = nil
		m.statusMsg = statusText
		return m, cmd
	case msg.Type == tea.KeyEsc || strings.EqualFold(msg.String(), "n"):
		m.pendingKill = nil
		m.statusMsg = "Kill cancelled"
		return m, clearStatusAfter(2 * time.Second)
	default:
		return m, nil
	}
}

func copySelectedValue(selected table.Row, column int, label string) string {
	if selected == nil {
		return "No row selected"
	}
	if column < 0 || column >= len(selected) {
		return fmt.Sprintf("No %s available", label)
	}

	value := strings.TrimSpace(selected[column])
	if value == "" || value == "-" {
		return fmt.Sprintf("No %s available for this port", label)
	}

	if err := copyText(value); err != nil {
		return fmt.Sprintf("✗ Failed to copy %s: %v", label, err)
	}
	return fmt.Sprintf("✓ Copied %s %s", label, value)
}

func exportVisible(format exporter.Format, ports []types.PortInfo) string {
	path, err := exportPorts(format, ports)
	if err != nil {
		return fmt.Sprintf("✗ Export failed: %v", err)
	}
	return fmt.Sprintf("✓ Exported %d port(s) to %s", len(ports), path)
}

func delayedScan(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		m.searching = false
		if m.search == "" {
			m.statusMsg = "Search cleared"
		} else {
			m.statusMsg = fmt.Sprintf("Searching %q", m.search)
		}
		return m, clearStatusAfter(2 * time.Second)
	case tea.KeyEsc:
		m.searching = false
		m.search = ""
		m.refreshRows()
		m.statusMsg = "Search cleared"
		return m, clearStatusAfter(2 * time.Second)
	case tea.KeyBackspace, tea.KeyCtrlH:
		if m.search != "" {
			runes := []rune(m.search)
			m.search = string(runes[:len(runes)-1])
			m.refreshRows()
		}
	default:
		if msg.String() == "/" && m.search == "" {
			return m, nil
		}
		if len(msg.Runes) > 0 {
			m.search += string(msg.Runes)
			m.refreshRows()
		}
	}

	return m, nil
}

func (m *Model) sortPorts() {
	sort.SliceStable(m.ports, func(i, j int) bool {
		switch m.sortBy {
		case sortByProcess:
			return strings.ToLower(m.ports[i].Process) < strings.ToLower(m.ports[j].Process)
		case sortByProtocol:
			return m.ports[i].Protocol < m.ports[j].Protocol
		case sortByPID:
			return m.ports[i].PID < m.ports[j].PID
		default: // sortByPort
			return m.ports[i].Port < m.ports[j].Port
		}
	})
}

func (m *Model) refreshRows() {
	m.table.SetRows(portsToRows(m.filteredPorts()))
}

func (m Model) filteredPorts() []types.PortInfo {
	search := strings.ToLower(strings.TrimSpace(m.search))
	if m.filter == filterAll && search == "" {
		return m.ports
	}

	protocol := m.filter.String()
	filtered := make([]types.PortInfo, 0, len(m.ports))
	for _, port := range m.ports {
		if m.filter != filterAll && port.Protocol != protocol {
			continue
		}
		if search != "" && !portMatchesSearch(port, search) {
			continue
		}
		filtered = append(filtered, port)
	}
	return filtered
}

func portMatchesSearch(port types.PortInfo, search string) bool {
	return strings.Contains(strings.ToLower(port.Process), search) ||
		strings.Contains(strconv.Itoa(port.Port), search)
}

func defaultColumns(width int) []table.Column {
	usable := width - 10
	if usable < 50 {
		usable = 50
	}

	return []table.Column{
		{Title: "PROTO", Width: max(6, usable*8/100)},
		{Title: "PORT", Width: max(6, usable*10/100)},
		{Title: "ADDRESS", Width: max(10, usable*25/100)},
		{Title: "PROCESS", Width: max(10, usable*35/100)},
		{Title: "PID", Width: max(6, usable*10/100)},
	}
}

func portsToRows(ports []types.PortInfo) []table.Row {
	rows := make([]table.Row, len(ports))
	for i, p := range ports {
		rows[i] = table.Row{
			strings.ToUpper(p.Protocol),
			p.PortString(),
			p.Address,
			p.Process,
			p.PIDString(),
		}
	}
	return rows
}

func renderHelpBar(width int) string {
	pairs := []struct{ key, desc string }{
		{"↑/↓", "navigate"},
		{"r", "refresh"},
		{"K", "kill"},
		{"c", "copy port"},
		{"P", "copy PID"},
		{"e", "CSV"},
		{"E", "JSON"},
		{"s", "sort"},
		{"f", "filter"},
		{"/", "search"},
		{"q", "quit"},
	}

	var parts []string
	for _, p := range pairs {
		parts = append(parts,
			HelpKeyStyle.Render(p.key)+" "+HelpDescStyle.Render(p.desc),
		)
	}

	helpLine := "  " + strings.Join(parts, "  │  ")
	return FooterStyle.Width(width).Render(helpLine)
}
