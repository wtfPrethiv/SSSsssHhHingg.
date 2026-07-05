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

 
func NewStyles() Styles {
	return Styles{
		Logo:    lipgloss.NewStyle().Foreground(colorFg).Bold(true),
		Cursor:  lipgloss.NewStyle().Foreground(colorRed),
		Tagline: lipgloss.NewStyle().Foreground(colorMuted).Italic(true),
		Hint:    lipgloss.NewStyle().Foreground(colorFaint),

		App: lipgloss.NewStyle().Padding(1, 3),

		TabActive:   lipgloss.NewStyle().Foreground(colorRed).Bold(true).Underline(true).Padding(0, 1),
		TabInactive: lipgloss.NewStyle().Foreground(colorMuted).Padding(0, 1),
		TabGap:      lipgloss.NewStyle().Foreground(colorFaint),

		Heading:   lipgloss.NewStyle().Foreground(colorRed).Bold(true),
		Body:      lipgloss.NewStyle().Foreground(colorFg),
		Muted:     lipgloss.NewStyle().Foreground(colorMuted),
		Accent:    lipgloss.NewStyle().Foreground(colorRed),
		ItemTitle: lipgloss.NewStyle().Foreground(colorFg).Bold(true),
		ItemDesc:  lipgloss.NewStyle().Foreground(colorSubtle),
		ItemMeta:  lipgloss.NewStyle().Foreground(colorFaint),
		LinkLabel: lipgloss.NewStyle().Foreground(colorFg).Bold(true),
		LinkURL:   lipgloss.NewStyle().Foreground(colorCyan).Underline(true),
		Updated:   lipgloss.NewStyle().Foreground(colorMint),
		Ascii:     lipgloss.NewStyle().Foreground(colorRedDim),
		Selected:  lipgloss.NewStyle().Foreground(colorRed).Bold(true),

		ReaderTitle: lipgloss.NewStyle().Foreground(colorRed).Bold(true),
		ReaderMeta:  lipgloss.NewStyle().Foreground(colorAmber),
		ReaderBody:  lipgloss.NewStyle().Foreground(colorFg),

		PlayFrame: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorFaint),
		PlayInfo: lipgloss.NewStyle().Foreground(colorMuted),

		Footer:     lipgloss.NewStyle().Foreground(colorFaint),
		FooterKey:  lipgloss.NewStyle().Foreground(colorSubtle).Bold(true),
		FooterDesc: lipgloss.NewStyle().Foreground(colorFaint),
	}
}
