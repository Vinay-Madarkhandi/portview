# ⚡ PortView

A sleek, real-time terminal UI for monitoring listening ports and their processes on Linux.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for a polished TUI experience.

## Features

- 📊 **Structured table view** — protocol, port, address, process name, and PID
- 🔄 **Auto-refresh** — updates every 3 seconds automatically
- ⌨️ **Manual refresh** — press `r` to refresh immediately
- 🔪 **Kill processes** — press `K` (shift+K) to send SIGTERM to the selected process
- 🔀 **Sortable** — cycle through sort modes (port, process, protocol, PID) with `s`
- 🎨 **Styled UI** — color-highlighted header, selected row, and help bar
- 📐 **Responsive layout** — adapts to terminal width and height
- 🛡️ **Robust error handling** — graceful handling of missing PIDs, permissions, and command failures
- 🐧 **TCP & UDP** — monitors both TCP (LISTEN) and UDP (UNCONN) sockets

## Requirements

- **Linux** (uses the `ss` command from `iproute2`)
- **Go 1.22+** (for building)
- Run with `sudo` for full process information (PID and process names)

## Installation

### From source

```bash
git clone https://github.com/yourusername/portview.git
cd portview
go build -o portview ./cmd/portview/
```

### Go install

```bash
go install github.com/yourusername/portview/cmd/portview@latest
```

## Usage

```bash
# Basic usage (may not show all process names without root)
./portview

# Recommended: run with sudo for full process info
sudo ./portview
```

## Key Bindings

| Key       | Action                          |
|-----------|---------------------------------|
| `↑` / `k` | Move selection up              |
| `↓` / `j` | Move selection down            |
| `r`       | Manual refresh                  |
| `K`       | Kill selected process (SIGTERM) |
| `s`       | Cycle sort mode                 |
| `q`       | Quit                            |

## Project Structure

```
portview/
├── cmd/
│   └── portview/
│       └── main.go              # Entry point
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # Port scanning (ss command execution)
│   │   ├── parser.go            # ss output parsing logic
│   │   └── parser_test.go       # Parser unit tests
│   ├── tui/
│   │   ├── model.go             # Bubble Tea model (Update/View/Init)
│   │   ├── styles.go            # Lip Gloss style definitions
│   │   └── keys.go              # Key binding definitions
│   └── types/
│       └── port.go              # PortInfo data type
├── go.mod
├── go.sum
└── README.md
```

## How It Works

1. PortView executes `ss -tulpnH` to list all listening TCP and UDP sockets
2. The output is parsed into structured `PortInfo` records
3. Results are displayed in a navigable table with automatic periodic refresh
4. Users can kill processes directly from the UI

## Future Improvements

- [ ] Filter ports by protocol (TCP/UDP toggle)
- [ ] Search/filter by process name or port number
- [ ] Copy port or PID to clipboard
- [ ] Export to JSON/CSV
- [ ] Configuration file for custom refresh intervals
- [ ] macOS support (via `lsof` or `netstat`)
- [ ] Color theme customization
- [ ] Confirmation dialog before killing a process


## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style definitions

