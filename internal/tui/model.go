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

	"portview/internal/scanner"
	"portview/internal/types"
)

const (
	refreshInterval = 3 * time.Second
)

type sortField int

const (
	sortByPort sortField = iota
	sortByProcess
	sortByProtocol
	sortByPID
	sortFieldCount
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

type tickMsg time.Time

type scanResultMsg struct {
	ports []types.PortInfo
	err   error
}

type statusMsg string

type Model struct {
	table     table.Model
	ports     []types.PortInfo
	err       error
	width     int
	height    int
	statusMsg string
	sortBy    sortField
	ready     bool
}

// InitialModel creates the initial application model.
func InitialModel() Model {
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
		BorderForeground(lipgloss.Color("#7D56F4")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Bold(true)
	t.SetStyles(s)

	return Model{
		table:  t,
		sortBy: sortByPort,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		doScan,
		tickCmd(),
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
		cmds = append(cmds, doScan, tickCmd())

	case scanResultMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.ports = msg.ports
			m.sortPorts()
			m.table.SetRows(portsToRows(m.ports))
		}

	case statusMsg:
		m.statusMsg = string(msg)
		cmds = append(cmds, clearStatusAfter(3*time.Second))

	case clearStatusMsg:
		m.statusMsg = ""

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Refresh):
			m.statusMsg = "Refreshing..."
			cmds = append(cmds, doScan, clearStatusAfter(1*time.Second))

		case key.Matches(msg, keys.Kill):
			statusText, cmd := killSelected(m.table.SelectedRow())
			m.statusMsg = statusText
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, keys.Sort):
			m.sortBy = (m.sortBy + 1) % sortFieldCount
			m.sortPorts()
			m.table.SetRows(portsToRows(m.ports))
			m.statusMsg = fmt.Sprintf("Sorted by %s", m.sortBy)
			cmds = append(cmds, clearStatusAfter(2*time.Second))

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

	if m.statusMsg != "" {
		b.WriteString("  " + StatusStyle.Render(m.statusMsg) + "\n")
	} else {
		portCount := len(m.ports)
		info := HelpDescStyle.Render(fmt.Sprintf("  %d port(s) · sorted by %s · auto-refresh %s", portCount, m.sortBy, refreshInterval))
		b.WriteString(info + "\n")
	}

	b.WriteString(renderHelpBar(m.width))

	return b.String()
}

func tickCmd() tea.Cmd {
	return tea.Tick(refreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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

	status := fmt.Sprintf("✓ Sent SIGTERM to PID %d (%s)", pid, selected[3])

	return status, tea.Batch(
		clearStatusAfter(3*time.Second),
		delayedScan(500*time.Millisecond),
	)
}

func delayedScan(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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
		{"s", "sort"},
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
