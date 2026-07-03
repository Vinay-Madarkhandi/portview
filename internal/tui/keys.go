package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Refresh  key.Binding
	Kill     key.Binding
	CopyPort key.Binding
	CopyPID  key.Binding
	Sort     key.Binding
	Filter   key.Binding
	Search   key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Kill: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "kill process"),
	),
	CopyPort: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy port"),
	),
	CopyPID: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "copy PID"),
	),
	Sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "cycle sort"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter protocol"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
