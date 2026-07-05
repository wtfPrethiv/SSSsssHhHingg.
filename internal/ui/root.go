package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"pr3thiv-portfolio/internal/config"
)

type state int

const (
	statePreloader state = iota
	statePortfolio
)


type Root struct {
	state     state
	preloader preloader
	portfolio portfolio
	width     int
	height    int
}


func NewRoot(c *config.Content) Root {
	styles := NewStyles()
	return Root{
		state:     statePreloader,
		preloader: newPreloader(c.Name, styles),
		portfolio: newPortfolio(c, styles),
	}
}

func (m Root) Init() tea.Cmd {
	return m.preloader.Init()
}

func (m Root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		var cmd tea.Cmd
		m.preloader, _ = m.preloader.Update(msg)
		m.portfolio, cmd = m.portfolio.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case preloaderDoneMsg:
		m.state = statePortfolio

		if m.width > 0 && m.height > 0 {
			m.portfolio, _ = m.portfolio.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
		return m, m.portfolio.Init()
	}

	var cmd tea.Cmd
	switch m.state {
	case statePreloader:
		m.preloader, cmd = m.preloader.Update(msg)
	case statePortfolio:
		m.portfolio, cmd = m.portfolio.Update(msg)
	}
	return m, cmd
}

func (m Root) View() string {
	switch m.state {
	case statePortfolio:
		return m.portfolio.View()
	default:
		return m.preloader.View()
	}
}
