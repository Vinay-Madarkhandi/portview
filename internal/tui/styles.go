package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Vinay-Madarkhandi/portview/internal/config"
)

type palette struct {
	accent           string
	headerForeground string
	subtitle         string
	status           string
	error            string
	muted            string
}

var (
	HeaderStyle    lipgloss.Style
	FooterStyle    lipgloss.Style
	StatusStyle    lipgloss.Style
	ErrorStyle     lipgloss.Style
	HelpKeyStyle   lipgloss.Style
	HelpDescStyle  lipgloss.Style
	TitleText      string
	SubtitleText   string
	currentPalette palette
)

func init() {
	configureStyles(config.Default())
}

func configureStyles(cfg config.Config) {
	p := paletteForTheme(cfg.Theme)
	p.applyOverrides(cfg)
	currentPalette = p

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(p.headerForeground)).
		Background(lipgloss.Color(p.accent)).
		Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.muted))

	StatusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.status)).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.error)).
		Bold(true)

	HelpKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.accent)).
		Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.muted))

	TitleText = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(p.headerForeground)).
		Render("⚡ PortView")

	SubtitleText = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.subtitle)).
		Render("— listening ports monitor")
}

func paletteForTheme(name string) palette {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "green":
		return palette{accent: "#04B575", headerForeground: "#FAFAFA", subtitle: "#D6FFE8", status: "#04B575", error: "#FF4444", muted: "#A0A0A0"}
	case "blue":
		return palette{accent: "#2F80ED", headerForeground: "#FAFAFA", subtitle: "#DDEBFF", status: "#04B575", error: "#FF4444", muted: "#A0A0A0"}
	case "amber":
		return palette{accent: "#D97706", headerForeground: "#FAFAFA", subtitle: "#FFE8BD", status: "#04B575", error: "#FF4444", muted: "#A0A0A0"}
	case "mono":
		return palette{accent: "#FAFAFA", headerForeground: "#111111", subtitle: "#D0D0D0", status: "#FAFAFA", error: "#FF5555", muted: "#A0A0A0"}
	default:
		return palette{accent: "#7D56F4", headerForeground: "#FAFAFA", subtitle: "#D0D0FF", status: "#04B575", error: "#FF4444", muted: "#A0A0A0"}
	}
}

func (p *palette) applyOverrides(cfg config.Config) {
	if cfg.AccentColor != "" {
		p.accent = cfg.AccentColor
	}
	if cfg.HeaderForegroundColor != "" {
		p.headerForeground = cfg.HeaderForegroundColor
	}
	if cfg.MutedColor != "" {
		p.muted = cfg.MutedColor
	}
	if cfg.StatusColor != "" {
		p.status = cfg.StatusColor
	}
	if cfg.ErrorColor != "" {
		p.error = cfg.ErrorColor
	}
}
