package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Vinay-Madarkhandi/portview/internal/tui"
)

func main() {
	p := tea.NewProgram(
		tui.InitialModel(),
		tea.WithAltScreen(), // full-screen TUI
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
