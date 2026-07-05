package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	blinkInterval = 450 * time.Millisecond
	typeInterval  = 130 * time.Millisecond
	holdDuration  = 900 * time.Millisecond
)
 

type preloaderDoneMsg struct{}

type blinkMsg struct{}
type typeMsg struct{}
type holdDoneMsg struct{}

type preloader struct {
	styles  Styles
	name    string
	runes   []rune
	typed   int
	cursor  bool
	width   int
	height  int
	holding bool
	done    bool
}

func newPreloader(name string, styles Styles) preloader {
	return preloader{
		styles: styles,
		name:   name,
		runes:  []rune(name),
		cursor: true,
	}
}

func blinkCmd() tea.Cmd {
	return tea.Tick(blinkInterval, func(time.Time) tea.Msg { return blinkMsg{} })
}

func typeCmd() tea.Cmd {
	return tea.Tick(typeInterval, func(time.Time) tea.Msg { return typeMsg{} })
}

func holdCmd() tea.Cmd {
	return tea.Tick(holdDuration, func(time.Time) tea.Msg { return holdDoneMsg{} })
}

func (m preloader) Init() tea.Cmd {
	return tea.Batch(blinkCmd(), typeCmd())
}

func (m preloader) Update(msg tea.Msg) (preloader, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
 
		m.done = true
		return m, func() tea.Msg { return preloaderDoneMsg{} }

	case blinkMsg:
		m.cursor = !m.cursor
		return m, blinkCmd()

	case typeMsg:
		if m.typed < len(m.runes) {
			m.typed++
			if m.typed == len(m.runes) && !m.holding {
				m.holding = true
				return m, holdCmd()
			}
			return m, typeCmd()
		}

	case holdDoneMsg:
		m.done = true
		return m, func() tea.Msg { return preloaderDoneMsg{} }
	}

	return m, nil
}

func (m preloader) View() string {
	typed := string(m.runes[:m.typed])

	cursor := " "
	if m.cursor {
		cursor = m.styles.Cursor.Render("█")
	}

	logo := m.styles.Logo.Render(typed) + cursor
	tagline := m.styles.Tagline.Render("loading portfolio…")
	hint := m.styles.Hint.Render("press any key to skip")

	block := lipgloss.JoinVertical(lipgloss.Center, logo, "", tagline, "", hint)

	if m.width == 0 || m.height == 0 {
		return block
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}
