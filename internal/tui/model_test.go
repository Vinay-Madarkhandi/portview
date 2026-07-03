package tui

import (
	"testing"

	"github.com/Vinay-Madarkhandi/portview/internal/types"
)

func TestFilteredPorts(t *testing.T) {
	ports := []types.PortInfo{
		{Protocol: "tcp", Port: 80},
		{Protocol: "udp", Port: 53},
		{Protocol: "tcp", Port: 443},
	}

	tests := []struct {
		name     string
		filter   protocolFilter
		expected []int
	}{
		{name: "all", filter: filterAll, expected: []int{80, 53, 443}},
		{name: "tcp", filter: filterTCP, expected: []int{80, 443}},
		{name: "udp", filter: filterUDP, expected: []int{53}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{ports: ports, filter: tt.filter}
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
