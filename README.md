# ⚡ PortView

[![CI](https://github.com/Vinay-Madarkhandi/portview/actions/workflows/ci.yml/badge.svg)](https://github.com/Vinay-Madarkhandi/portview/actions/workflows/ci.yml)
[![Release](https://github.com/Vinay-Madarkhandi/portview/actions/workflows/release.yml/badge.svg)](https://github.com/Vinay-Madarkhandi/portview/releases/latest)

A sleek, real-time terminal UI for monitoring listening ports and their processes on Linux and macOS.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for a polished TUI experience.

---

## Features

- 📊 **Structured table view** — protocol, port, address, process name, and PID
- 🔄 **Auto-refresh** — updates every 3 seconds automatically
- ⌨️ **Manual refresh** — press `r` to refresh immediately
- 🔪 **Kill processes** — press `K` (shift+K) to send SIGTERM to the selected process
- 🔀 **Sortable** — cycle through sort modes (port, process, protocol, PID) with `s`
- 🎨 **Styled UI** — color-highlighted header, selected row, and help bar
- 📐 **Responsive layout** — adapts to terminal width and height
- 🛡️ **Robust error handling** — graceful handling of missing PIDs, permissions, and command failures
- 🐧 **Linux support** — uses `ss` from `iproute2` for TCP/UDP socket discovery
- 🍎 **macOS support** — uses native `lsof` for TCP/UDP socket discovery

---

## Installation

### Quick Install (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/Vinay-Madarkhandi/portview/main/install.sh | bash
```

This will:
- Detect your OS (Linux/macOS) and architecture (amd64/arm64)
- Download the correct prebuilt binary from the latest release
- Install it to `/usr/local/bin` (or `~/.local/bin` if sudo is unavailable)

### Go Install

If you have Go installed:

```bash
go install github.com/Vinay-Madarkhandi/portview/cmd/portview@latest
```

### Manual Download

Download the binary for your platform from the [latest release](https://github.com/Vinay-Madarkhandi/portview/releases/latest):

| Platform              | Binary                  |
|-----------------------|-------------------------|
| Linux (x86_64)        | `portview-linux-amd64`  |
| Linux (ARM64)         | `portview-linux-arm64`  |
| macOS (Intel)         | `portview-darwin-amd64` |
| macOS (Apple Silicon) | `portview-darwin-arm64` |

```bash
# Example: Linux x86_64
curl -Lo portview https://github.com/Vinay-Madarkhandi/portview/releases/latest/download/portview-linux-amd64
chmod +x portview
sudo mv portview /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/Vinay-Madarkhandi/portview.git
cd portview
make build
sudo mv portview /usr/local/bin/
```

---

## Usage

```bash
# Basic usage (may not show all process names without root)
portview

# Recommended: run with sudo for full process info
sudo portview
```

---

## Key Bindings

| Key       | Action                          |
|-----------|---------------------------------|
| `↑` / `k` | Move selection up              |
| `↓` / `j` | Move selection down            |
| `r`       | Manual refresh                  |
| `K`       | Kill selected process (SIGTERM) |
| `s`       | Cycle sort mode                 |
| `q`       | Quit                            |

---

## Requirements

- **Linux** — uses the `ss` command from `iproute2`
- **macOS** — uses the built-in `lsof` command
- Run with **`sudo`** for full process information (PID and process names)

---

## Project Structure

```
portview/
├── cmd/portview/
│   └── main.go                  # Entry point
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # OS-specific port scanner command execution
│   │   ├── parser.go            # ss/lsof output parsing logic
│   │   └── parser_test.go       # Parser unit tests
│   ├── tui/
│   │   ├── model.go             # Bubble Tea model (Update/View/Init)
│   │   ├── styles.go            # Lip Gloss style definitions
│   │   └── keys.go              # Key binding definitions
│   └── types/
│       └── port.go              # PortInfo data type
├── .github/workflows/
│   ├── ci.yml                   # CI: test + vet on push/PR
│   └── release.yml              # Release: build + upload on tag
├── install.sh                   # One-line installer script
├── Makefile                     # Build, test, release targets
├── go.mod
├── go.sum
└── README.md
```

---

## How It Works

1. PortView executes `ss -tulpnH` on Linux or `lsof` on macOS to list listening TCP and UDP sockets
2. The output is parsed into structured `PortInfo` records
3. Results are displayed in a navigable table with automatic periodic refresh
4. Users can kill processes directly from the UI

---

## Development

```bash
# Run tests
make test

# Run go vet
make vet

# Format code
make fmt

# Build for current platform
make build

# Cross-compile release binaries (outputs to dist/)
make release
```

### Creating a Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

The [release workflow](.github/workflows/release.yml) will automatically build binaries for all platforms and create a GitHub Release.

---

## Future Improvements

- [ ] Filter ports by protocol (TCP/UDP toggle)
- [ ] Search/filter by process name or port number
- [ ] Copy port or PID to clipboard
- [ ] Export to JSON/CSV
- [ ] Configuration file for custom refresh intervals
- [ ] Color theme customization
- [ ] Confirmation dialog before killing a process

---

## License

MIT

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style definitions
