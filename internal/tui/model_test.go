package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Vinay-Madarkhandi/portview/internal/config"
	"github.com/Vinay-Madarkhandi/portview/internal/exporter"
	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

func TestFilteredPorts(t *testing.T) {
	ports := []types.PortInfo{
		{Protocol: "tcp", Port: 80, Process: "nginx"},
		{Protocol: "udp", Port: 53, Process: "mDNSResponder"},
		{Protocol: "tcp", Port: 443, Process: "caddy"},
	}

	tests := []struct {
		name     string
		filter   protocolFilter
		search   string
		expected []int
	}{
		{name: "all", filter: filterAll, expected: []int{80, 53, 443}},
		{name: "tcp", filter: filterTCP, expected: []int{80, 443}},
		{name: "udp", filter: filterUDP, expected: []int{53}},
		{name: "process search", filter: filterAll, search: "nginx", expected: []int{80}},
		{name: "case insensitive process search", filter: filterAll, search: "mdns", expected: []int{53}},
		{name: "port search", filter: filterAll, search: "44", expected: []int{443}},
		{name: "protocol and search", filter: filterTCP, search: "53", expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{ports: ports, filter: tt.filter, search: tt.search}
			got := model.filteredPorts()

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d ports, got %d", len(tt.expected), len(got))
			}

			for i, port := range got {
				if port.Port != tt.expected[i] {
					t.Errorf("port %d = %d, want %d", i, port.Port, tt.expected[i])
				}
			}
		})
	}
}

func TestProtocolFilterString(t *testing.T) {
	tests := []struct {
		filter   protocolFilter
		expected string
	}{
		{filter: filterAll, expected: "all"},
		{filter: filterTCP, expected: "tcp"},
		{filter: filterUDP, expected: "udp"},
		{filter: protocolFilter(99), expected: "all"},
	}

	for _, tt := range tests {
		if got := tt.filter.String(); got != tt.expected {
			t.Errorf("String() = %q, want %q", got, tt.expected)
		}
	}
}

func TestUpdateSearch(t *testing.T) {
	model := InitialModel()
	model.ports = []types.PortInfo{{Protocol: "tcp", Port: 8080, Process: "nginx"}}
	model.searching = true

	updated, _ := model.updateSearch(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ng")})
	model = updated.(Model)
	if model.search != "ng" {
		t.Fatalf("expected search %q, got %q", "ng", model.search)
	}

	updated, _ = model.updateSearch(tea.KeyMsg{Type: tea.KeyBackspace})
	model = updated.(Model)
	if model.search != "n" {
		t.Fatalf("expected search %q after backspace, got %q", "n", model.search)
	}

	updated, _ = model.updateSearch(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(Model)
	if model.searching {
		t.Fatal("expected search mode to be inactive after escape")
	}
	if model.search != "" {
		t.Fatalf("expected search to be cleared, got %q", model.search)
	}
}

func TestCopySelectedValue(t *testing.T) {
	originalCopyText := copyText
	t.Cleanup(func() { copyText = originalCopyText })

	var copied string
	copyText = func(text string) error {
		copied = text
		return nil
	}

	status := copySelectedValue(table.Row{"TCP", "8080", "0.0.0.0", "nginx", "1234"}, 1, "port")
	if copied != "8080" {
		t.Fatalf("copied %q, want 8080", copied)
	}
	if status != "✓ Copied port 8080" {
		t.Fatalf("status = %q", status)
	}

	status = copySelectedValue(table.Row{"TCP", "8080", "0.0.0.0", "nginx", "1234"}, 4, "PID")
	if copied != "1234" {
		t.Fatalf("copied %q, want 1234", copied)
	}
	if status != "✓ Copied PID 1234" {
		t.Fatalf("status = %q", status)
	}
}

func TestCopySelectedValueErrors(t *testing.T) {
	originalCopyText := copyText
	t.Cleanup(func() { copyText = originalCopyText })

	status := copySelectedValue(nil, 1, "port")
	if status != "No row selected" {
		t.Fatalf("status = %q", status)
	}

	status = copySelectedValue(table.Row{"TCP", "8080", "0.0.0.0", "unknown", "-"}, 4, "PID")
	if status != "No PID available for this port" {
		t.Fatalf("status = %q", status)
	}

	copyText = func(text string) error {
		return errors.New("clipboard unavailable")
	}
	status = copySelectedValue(table.Row{"TCP", "8080", "0.0.0.0", "nginx", "1234"}, 1, "port")
	if status != "✗ Failed to copy port: clipboard unavailable" {
		t.Fatalf("status = %q", status)
	}
}

func TestExportVisible(t *testing.T) {
	originalExportPorts := exportPorts
	t.Cleanup(func() { exportPorts = originalExportPorts })

	var gotFormat exporter.Format
	exportPorts = func(format exporter.Format, ports []types.PortInfo) (string, error) {
		gotFormat = format
		if len(ports) != 1 {
			t.Fatalf("exported %d ports, want 1", len(ports))
		}
		return "ports.csv", nil
	}

	status := exportVisible(exporter.FormatCSV, []types.PortInfo{{Protocol: "tcp", Port: 8080}})
	if gotFormat != exporter.FormatCSV {
		t.Fatalf("format = %q, want csv", gotFormat)
	}
	if status != "✓ Exported 1 port(s) to ports.csv" {
		t.Fatalf("status = %q", status)
	}

	exportPorts = func(format exporter.Format, ports []types.PortInfo) (string, error) {
		return "", errors.New("disk full")
	}
	status = exportVisible(exporter.FormatJSON, []types.PortInfo{{Protocol: "tcp", Port: 8080}})
	if status != "✗ Export failed: disk full" {
		t.Fatalf("status = %q", status)
	}
}

func TestRequestKillConfirmation(t *testing.T) {
	status, pending := requestKillConfirmation(table.Row{"TCP", "8080", "127.0.0.1", "nginx", "1234"})
	if pending == nil {
		t.Fatal("expected pending kill row")
	}
	if status != "Kill PID 1234 (nginx)? y/Enter confirm, n/Esc cancel" {
		t.Fatalf("status = %q", status)
	}

	status, pending = requestKillConfirmation(table.Row{"TCP", "8080", "127.0.0.1", "unknown", "-"})
	if pending != nil {
		t.Fatal("expected no pending kill without PID")
	}
	if status != "No PID available for this port" {
		t.Fatalf("status = %q", status)
	}
}

func TestUpdateKillConfirmationCancel(t *testing.T) {
	model := Model{pendingKill: table.Row{"TCP", "8080", "127.0.0.1", "nginx", "1234"}}
	updated, _ := model.updateKillConfirmation(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(Model)
	if model.pendingKill != nil {
		t.Fatal("expected pending kill to be cleared")
	}
	if model.statusMsg != "Kill cancelled" {
		t.Fatalf("status = %q", model.statusMsg)
	}
}

func TestInitialModelWithConfig(t *testing.T) {
	model := initialModelWithConfig(config.Config{
		RefreshIntervalSeconds: 9,
		Theme:                  "green",
	})
	if model.refreshInterval() != 9*time.Second {
		t.Fatalf("refresh interval = %s, want 9s", model.refreshInterval())
	}
}
