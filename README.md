# ⚡ PortView

[![CI](https://github.com/Vinay-Madarkhandi/portview/actions/workflows/ci.yml/badge.svg)](https://github.com/Vinay-Madarkhandi/portview/actions/workflows/ci.yml)
[![Release](https://github.com/Vinay-Madarkhandi/portview/actions/workflows/release.yml/badge.svg)](https://github.com/Vinay-Madarkhandi/portview/releases/latest)

A sleek, real-time terminal UI for monitoring listening ports and their processes on Linux, macOS, and Windows.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for a polished TUI experience.

---

## Features

- 📊 **Structured table view** — protocol, port, address, process name, and PID
- 🔄 **Auto-refresh** — updates automatically, configurable via `config.json`
- ⌨️ **Manual refresh** — press `r` to refresh immediately
- 🔪 **Kill processes** — press `K` (shift+K), then confirm before terminating
- 🔀 **Sortable** — cycle through sort modes (port, process, protocol, PID) with `s`
- 🔎 **Protocol filter** — cycle through all, TCP, and UDP ports with `f`
- 🔍 **Search filter** — press `/` to filter by process name or port number
- 📋 **Clipboard copy** — copy the selected port or PID from the table
- 💾 **Export** — export visible rows to JSON or CSV
- 🎛️ **Config file** — customize refresh interval and color theme
- 🎨 **Styled UI** — color-highlighted header, selected row, and help bar
- 📐 **Responsive layout** — adapts to terminal width and height
- 🛡️ **Robust error handling** — graceful handling of missing PIDs, permissions, and command failures
- 🐧 **Linux support** — uses `ss` from `iproute2` for TCP/UDP socket discovery
- 🍎 **macOS support** — uses native `lsof` for TCP/UDP socket discovery
- 🪟 **Windows support** — uses PowerShell networking cmdlets for TCP/UDP socket discovery

---

## Installation

### Quick Install (recommended)

```bash
curl -sSL https://raw.githubusercontent.com/Vinay-Madarkhandi/portview/main/install.sh | bash
```

This will:
- Detect your OS (Linux/macOS/Windows) and architecture (amd64/arm64)
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
| Windows (x86_64)      | `portview-windows-amd64.exe` |
| Windows (ARM64)       | `portview-windows-arm64.exe` |

```bash
# Example: Linux x86_64
curl -Lo portview https://github.com/Vinay-Madarkhandi/portview/releases/latest/download/portview-linux-amd64
chmod +x portview
sudo mv portview /usr/local/bin/

# Example: Windows x86_64 (PowerShell)
iwr -OutFile portview.exe https://github.com/Vinay-Madarkhandi/portview/releases/latest/download/portview-windows-amd64.exe
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
# Basic usage (may not show all process names without elevated permissions)
portview

# Linux/macOS: recommended for full process info
sudo portview

# Windows: run your terminal as Administrator for full process info
```

---

## Key Bindings

| Key       | Action                          |
|-----------|---------------------------------|
| `↑` / `k` | Move selection up              |
| `↓` / `j` | Move selection down            |
| `r`       | Manual refresh                  |
| `K`       | Ask to terminate selected process |
| `y` / `Enter` | Confirm process termination |
| `n` / `Esc` | Cancel process termination |
| `c`       | Copy selected port              |
| `P`       | Copy selected PID               |
| `e`       | Export visible rows to CSV      |
| `E`       | Export visible rows to JSON     |
| `s`       | Cycle sort mode                 |
| `f`       | Cycle protocol filter           |
| `/`       | Search process name or port     |
| `Esc`     | Clear active search             |
| `q`       | Quit                            |

---

## Configuration

PortView reads configuration from your OS config directory:
- Linux: `~/.config/portview/config.json`
- macOS: `~/Library/Application Support/portview/config.json`
- Windows: `%AppData%\portview\config.json`

Set `PORTVIEW_CONFIG=/path/to/config.json` to use a custom location.

```json
{
  "refresh_interval_seconds": 5,
  "theme": "green",
  "accent_color": "#04B575",
  "header_foreground_color": "#FAFAFA",
  "muted_color": "#A0A0A0",
  "status_color": "#04B575",
  "error_color": "#FF4444"
}
```

Built-in themes: `purple`, `green`, `blue`, `amber`, `mono`.
Color fields are optional overrides and accept Lip Gloss-compatible color strings.

Exports are written to the current directory as `portview-ports-YYYYMMDD-HHMMSS.csv` or `.json`.

---

## Requirements

- **Linux** — uses the `ss` command from `iproute2`
- **macOS** — uses the built-in `lsof` command
- **Windows** — uses PowerShell `Get-NetTCPConnection`, `Get-NetUDPEndpoint`, and `Get-Process`
- **Clipboard on Linux** — uses `wl-copy`, `xclip`, or `xsel` when copying ports/PIDs
- Run with **`sudo`** on Linux/macOS or as **Administrator** on Windows for full process information

---

## Project Structure

```
portview/
├── cmd/portview/
│   └── main.go                  # Entry point
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # OS-specific port scanner command execution
│   │   ├── parser.go            # ss/lsof/PowerShell output parsing logic
│   │   └── parser_test.go       # Parser unit tests
│   ├── config/                  # Runtime config loading
│   ├── exporter/                # JSON/CSV export helpers
│   ├── clipboard/               # Clipboard command helpers
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

1. PortView executes `ss -tulpnH` on Linux, `lsof` on macOS, or PowerShell networking cmdlets on Windows to list listening TCP and UDP sockets
2. The output is parsed into structured `PortInfo` records
3. Results are displayed in a navigable table with automatic periodic refresh
4. Users can filter, search, export, copy, and terminate processes directly from the UI

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

## License

MIT

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style definitions
