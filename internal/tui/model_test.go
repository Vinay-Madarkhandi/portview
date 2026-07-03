package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

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
