package ui

import "github.com/charmbracelet/lipgloss"


var (
	colorRed    = lipgloss.Color("#FF3B30")
	colorRedDim = lipgloss.Color("#8C1D18")
	colorFg     = lipgloss.Color("#EAEAEA")
	colorMuted  = lipgloss.Color("#7A7A7A")
	colorFaint  = lipgloss.Color("#4A4A4A")
	colorSubtle = lipgloss.Color("#B5B5B5")
	colorMint   = lipgloss.Color("#3DDC97")
	colorAmber  = lipgloss.Color("#FFB454")
	colorCyan   = lipgloss.Color("#56C2E6")
)

 
var playgroundPalette = []lipgloss.Color{
	lipgloss.Color("#FF3B30"),
	lipgloss.Color("#FF9F0A"),
	lipgloss.Color("#FFD60A"),
	lipgloss.Color("#30D158"),
	lipgloss.Color("#40C8E0"),
	lipgloss.Color("#5E5CE6"),
	lipgloss.Color("#BF5AF2"),
	lipgloss.Color("#FF6482"),
}

 
type Styles struct {
	// Renderer is the per-session Lip Gloss renderer. Building styles from
	// this (instead of the package global) is what lets color work when the
	// server runs without a TTY of its own (e.g. under systemd) — the color
	// profile is detected from the connecting SSH client, not the server.
	Renderer *lipgloss.Renderer

	Logo    lipgloss.Style
	Cursor  lipgloss.Style
	Tagline lipgloss.Style
	Hint    lipgloss.Style

	App         lipgloss.Style
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style
	TabGap      lipgloss.Style

	Heading   lipgloss.Style
	Body      lipgloss.Style
	Muted     lipgloss.Style
	Accent    lipgloss.Style
	ItemTitle lipgloss.Style
	ItemDesc  lipgloss.Style
	ItemMeta  lipgloss.Style
	LinkLabel lipgloss.Style
	LinkURL   lipgloss.Style
	Updated   lipgloss.Style
	Ascii     lipgloss.Style
	Selected  lipgloss.Style

	ReaderTitle lipgloss.Style
	ReaderMeta  lipgloss.Style
	ReaderBody  lipgloss.Style

	PlayFrame lipgloss.Style
	PlayInfo  lipgloss.Style

	Footer     lipgloss.Style
	FooterKey  lipgloss.Style
	FooterDesc lipgloss.Style
}

 
func NewStyles(r *lipgloss.Renderer) Styles {
	return Styles{
		Renderer: r,

		Logo:    r.NewStyle().Foreground(colorFg).Bold(true),
		Cursor:  r.NewStyle().Foreground(colorRed),
		Tagline: r.NewStyle().Foreground(colorMuted).Italic(true),
		Hint:    r.NewStyle().Foreground(colorFaint),

		App: r.NewStyle().Padding(1, 3),

		TabActive:   r.NewStyle().Foreground(colorRed).Bold(true).Underline(true).Padding(0, 1),
		TabInactive: r.NewStyle().Foreground(colorMuted).Padding(0, 1),
		TabGap:      r.NewStyle().Foreground(colorFaint),

		Heading:   r.NewStyle().Foreground(colorRed).Bold(true),
		Body:      r.NewStyle().Foreground(colorFg),
		Muted:     r.NewStyle().Foreground(colorMuted),
		Accent:    r.NewStyle().Foreground(colorRed),
		ItemTitle: r.NewStyle().Foreground(colorFg).Bold(true),
		ItemDesc:  r.NewStyle().Foreground(colorSubtle),
		ItemMeta:  r.NewStyle().Foreground(colorFaint),
		LinkLabel: r.NewStyle().Foreground(colorFg).Bold(true),
		LinkURL:   r.NewStyle().Foreground(colorCyan).Underline(true),
		Updated:   r.NewStyle().Foreground(colorMint),
		Ascii:     r.NewStyle().Foreground(colorRedDim),
		Selected:  r.NewStyle().Foreground(colorRed).Bold(true),

		ReaderTitle: r.NewStyle().Foreground(colorRed).Bold(true),
		ReaderMeta:  r.NewStyle().Foreground(colorAmber),
		ReaderBody:  r.NewStyle().Foreground(colorFg),

		PlayFrame: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorFaint),
		PlayInfo: r.NewStyle().Foreground(colorMuted),

		Footer:     r.NewStyle().Foreground(colorFaint),
		FooterKey:  r.NewStyle().Foreground(colorSubtle).Bold(true),
		FooterDesc: r.NewStyle().Foreground(colorFaint),
	}
}
