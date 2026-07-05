package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"pr3thiv-portfolio/internal/config"
)

var tabs = []string{"About", "Projects", "Blogs", "Links", "Playground"}

const linkTip = "Most terminals let you click (or Ctrl/Cmd+click) the URLs above."

 
const maxReaderWidth = 84

const headerBlinkInterval = 550 * time.Millisecond

type headerBlinkMsg struct{}

func headerBlinkCmd() tea.Cmd {
	return tea.Tick(headerBlinkInterval, func(time.Time) tea.Msg { return headerBlinkMsg{} })
}

type keyMap struct {
	Prev   key.Binding
	Next   key.Binding
	Up     key.Binding
	Down   key.Binding
	Top    key.Binding
	Bottom key.Binding
	Open   key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Prev:   key.NewBinding(key.WithKeys("h", "shift+tab", "[", "left")),
	Next:   key.NewBinding(key.WithKeys("l", "tab", "]", "right")),
	Up:     key.NewBinding(key.WithKeys("k", "up")),
	Down:   key.NewBinding(key.WithKeys("j", "down")),
	Top:    key.NewBinding(key.WithKeys("g", "home")),
	Bottom: key.NewBinding(key.WithKeys("G", "end")),
	Open:   key.NewBinding(key.WithKeys("enter")),
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c")),
}

type portfolio struct {
	content *config.Content
	styles  Styles

	vp     viewport.Model
	reader viewport.Model
	play   playground

	active       int
	blogCursor   int
	readerActive bool
	artIndex     int
	cursorOn     bool

	width  int
	height int
	ready  bool
}

func newPortfolio(c *config.Content, styles Styles) portfolio {
	artIdx := 0
	if len(c.ASCIIArts) > 0 {
		artIdx = rand.Intn(len(c.ASCIIArts))
	}
	return portfolio{
		content:  c,
		styles:   styles,
		play:     newPlayground(styles),
		artIndex: artIdx,
		cursorOn: true,
	}
}

func (m portfolio) Init() tea.Cmd {
	return headerBlinkCmd()
}

func (m *portfolio) setSize(w, h int) {
	m.width, m.height = w, h

	frameW := 6  
	frameH := 2  
	innerW := w - frameW
	if innerW < 10 {
		innerW = 10
	}

	header := lipgloss.Height(m.headerView())
	tabbar := lipgloss.Height(m.tabBarView(innerW))
	footer := 1
	chrome := header + tabbar + footer + 3

	vpHeight := h - frameH - chrome
	if vpHeight < 3 {
		vpHeight = 3
	}

	readerW := innerW - 4  
	if readerW > maxReaderWidth {
		readerW = maxReaderWidth
	}
	if readerW < 10 {
		readerW = 10
	}
	readerH := vpHeight - 1 
	if readerH < 3 {
		readerH = 3
	}

	if !m.ready {
		m.vp = viewport.New(innerW, vpHeight)
		m.reader = viewport.New(readerW, readerH)
		m.ready = true
	} else {
		m.vp.Width, m.vp.Height = innerW, vpHeight
		m.reader.Width, m.reader.Height = readerW, readerH
	}

	if tabs[m.active] == "Playground" {
		m.play.setSize(innerW, vpHeight-1)
	}
	m.refreshContent()
}

func (m *portfolio) refreshContent() {
	if tabs[m.active] == "Playground" {
		return
	}
	m.vp.SetContent(m.sectionContent(m.active, m.vp.Width))
	m.vp.GotoTop()
}

func (m *portfolio) setActive(i int) tea.Cmd {
	if tabs[m.active] == "Playground" {
		m.play.stop()
	}
	m.active = i
	m.readerActive = false

	if tabs[m.active] == "Playground" {
		m.play.setSize(m.vp.Width, m.vp.Height-1)
		return m.play.start()
	}
	m.refreshContent()
	return nil
}

func (m *portfolio) openReader() {
	if len(m.content.Blogs) == 0 {
		return
	}
	b := m.content.Blogs[m.blogCursor]

	var sb strings.Builder
	sb.WriteString(m.styles.ReaderTitle.Render(b.Title))
	sb.WriteString("\n")
	if b.Date != "" {
		sb.WriteString(m.styles.ReaderMeta.Render(b.Date))
		sb.WriteString("\n")
	}
	sb.WriteString(m.styles.TabGap.Render(strings.Repeat("─", m.reader.Width)))
	sb.WriteString("\n\n")

	body := strings.TrimSpace(b.Body)
	if body == "" {
		body = "This post doesn't have inline content yet. Read the full version at the link below."
	}
	sb.WriteString(m.styles.ReaderBody.Width(m.reader.Width).Render(body))

	if b.URL != "" {
		sb.WriteString("\n\n")
		sb.WriteString(m.styles.Muted.Render("link:  ") + m.styles.LinkURL.Render(b.URL))
	}

	m.reader.SetContent(sb.String())
	m.reader.GotoTop()
	m.readerActive = true
}

func (m portfolio) Update(msg tea.Msg) (portfolio, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil

	case headerBlinkMsg:
		m.cursorOn = !m.cursorOn
		return m, headerBlinkCmd()

	case frameMsg:
		if tabs[m.active] == "Playground" {
			var cmd tea.Cmd
			m.play, cmd = m.play.update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.ready {
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m portfolio) handleKey(msg tea.KeyMsg) (portfolio, tea.Cmd) {
 
	if m.readerActive {
		switch msg.String() {
		case "esc", "q", "backspace":
			m.readerActive = false
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, keys.Top):
			m.reader.GotoTop()
			return m, nil
		case key.Matches(msg, keys.Bottom):
			m.reader.GotoBottom()
			return m, nil
		}
		var cmd tea.Cmd
		m.reader, cmd = m.reader.Update(msg)
		return m, cmd
	}

	if key.Matches(msg, keys.Quit) {
		return m, tea.Quit
	}

 
	if n := msg.String(); len(n) == 1 && n[0] >= '1' && n[0] <= '9' {
		idx := int(n[0] - '1')
		if idx < len(tabs) {
			return m, m.setActive(idx)
		}
	}

	switch {
	case key.Matches(msg, keys.Next):
		return m, m.setActive((m.active + 1) % len(tabs))
	case key.Matches(msg, keys.Prev):
		return m, m.setActive((m.active - 1 + len(tabs)) % len(tabs))
	}

	switch tabs[m.active] {
	case "Blogs":
		switch {
		case key.Matches(msg, keys.Down):
			if m.blogCursor < len(m.content.Blogs)-1 {
				m.blogCursor++
				m.refreshContent()
			}
			return m, nil
		case key.Matches(msg, keys.Up):
			if m.blogCursor > 0 {
				m.blogCursor--
				m.refreshContent()
			}
			return m, nil
		case key.Matches(msg, keys.Open):
			m.openReader()
			return m, nil
		}
	case "Playground":
		var cmd tea.Cmd
		m.play, cmd = m.play.update(msg)
		return m, cmd
	}

	switch {
	case key.Matches(msg, keys.Top):
		m.vp.GotoTop()
		return m, nil
	case key.Matches(msg, keys.Bottom):
		m.vp.GotoBottom()
		return m, nil
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m portfolio) View() string {
	if !m.ready {
		return ""
	}
	innerW := m.vp.Width

	var body string
	switch {
	case m.readerActive:
		body = m.readerView()
	case tabs[m.active] == "Playground":
		body = m.play.view()
	default:
		body = m.vp.View()
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		"",
		m.tabBarView(innerW),
		"",
		body,
		"",
		m.footerView(),
	)
	return m.styles.App.Render(view)
}

func (m portfolio) readerView() string {
	pct := int(m.reader.ScrollPercent()*100 + 0.5)
	if pct < 0 {
		pct = 0
	} else if pct > 100 {
		pct = 100
	}

	content := lipgloss.NewStyle().PaddingLeft(2).Render(m.reader.View())

	label := fmt.Sprintf(" %d%% read", pct)
	barW := m.reader.Width - lipgloss.Width(label)
	if barW < 4 {
		barW = 4
	}
	filled := barW * pct / 100
	bar := m.styles.Accent.Render(strings.Repeat("━", filled)) +
		m.styles.TabGap.Render(strings.Repeat("━", barW-filled))
	indicator := lipgloss.NewStyle().PaddingLeft(2).Render(bar + m.styles.ReaderMeta.Render(label))

	return lipgloss.JoinVertical(lipgloss.Left, content, indicator)
}

func (m portfolio) headerView() string {
	name := m.styles.Logo.Render(m.content.Name)
	cursor := " "
	if m.cursorOn {
		cursor = m.styles.Cursor.Render("█")
	}
	line := name + cursor
	if m.content.Tagline != "" {
		line += m.styles.Muted.Render("  " + m.content.Tagline)
	}
	return line
}

func (m portfolio) tabBarView(width int) string {
	rendered := make([]string, len(tabs))
	for i, t := range tabs {
		if i == m.active {
			rendered[i] = m.styles.TabActive.Render(t)
		} else {
			rendered[i] = m.styles.TabInactive.Render(t)
		}
	}
	row := strings.Join(rendered, m.styles.TabGap.Render("·"))

	ruleW := width
	if ruleW < 0 {
		ruleW = 0
	}
	rule := m.styles.TabGap.Render(strings.Repeat("─", ruleW))
	return lipgloss.JoinVertical(lipgloss.Left, row, rule)
}

func (m portfolio) footerView() string {
	fk := func(k, d string) string {
		return m.styles.FooterKey.Render(k) + " " + m.styles.FooterDesc.Render(d)
	}

	var hints []string
	switch {
	case m.readerActive:
		hints = []string{fk("j/k", "scroll"), fk("g/G", "top/bottom"), fk("q/esc", "back")}
	case tabs[m.active] == "Playground":
		hints = []string{fk("h/l", "tabs"), fk("space", "launch"), fk("f", "drop"), fk("r", "reset"), fk("q", "quit")}
	case tabs[m.active] == "Blogs":
		hints = []string{fk("h/l", "tabs"), fk("j/k", "select"), fk("enter", "read"), fk("q", "quit")}
	default:
		hints = []string{fk("h/l", "tabs"), fk("j/k", "scroll"), fk("g/G", "top/bottom"), fk("q", "quit")}
	}
	return m.styles.Footer.Render(strings.Join(hints, m.styles.FooterDesc.Render("  ·  ")))
}

func (m portfolio) sectionContent(idx, width int) string {
	switch tabs[idx] {
	case "About":
		return m.aboutSection(width)
	case "Projects":
		return m.projectsSection(width)
	case "Blogs":
		return m.blogsSection(width)
	case "Links":
		return m.linksSection(width)
	}
	return ""
}

func (m portfolio) aboutSection(width int) string {
	heading := m.styles.Heading.Render("# whoami")

	art := ""
	if len(m.content.ASCIIArts) > 0 {
		art = m.styles.Ascii.Render(m.content.ASCIIArts[m.artIndex])
	}
	artW := lipgloss.Width(art)

	textW := width
	twoCol := art != "" && width > artW+36
	if twoCol {
		textW = width - artW - 4
	}

	text := m.styles.Body.Width(textW).Render(strings.TrimSpace(m.content.About))

	var bodyBlock string
	if twoCol {
		gap := m.styles.Muted.Render("    ")
		bodyBlock = lipgloss.JoinHorizontal(lipgloss.Top, text, gap, art)
	} else if art != "" {
		bodyBlock = lipgloss.JoinVertical(lipgloss.Left, art, "", text)
	} else {
		bodyBlock = text
	}

	parts := []string{heading, "", bodyBlock}
	if m.content.LastUpdated != "" {
		updated := m.styles.Updated.Render("last updated: " + m.content.LastUpdated)
		parts = append(parts, "", updated)
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m portfolio) projectsSection(width int) string {
	parts := []string{m.styles.Heading.Render("# projects"), ""}
	if len(m.content.Projects) == 0 {
		parts = append(parts, m.styles.Muted.Render("Nothing here yet."))
	}
	for i, p := range m.content.Projects {
		title := m.styles.Accent.Render("▸ ") + m.styles.ItemTitle.Render(p.Name)
		desc := m.styles.ItemDesc.Width(width).Render(p.Description)
		block := lipgloss.JoinVertical(lipgloss.Left, title, desc)
		if p.URL != "" {
			block = lipgloss.JoinVertical(lipgloss.Left, block, "  "+m.renderLink("", p.URL))
		}
		parts = append(parts, block)
		if i < len(m.content.Projects)-1 {
			parts = append(parts, "")
		}
	}
	if len(m.content.Projects) > 0 {
		parts = append(parts, "", m.styles.Muted.Render(linkTip))
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m portfolio) blogsSection(width int) string {
	parts := []string{
		m.styles.Heading.Render("# blogs"),
		m.styles.Muted.Render("select with j/k, open with enter"),
		"",
	}
	if len(m.content.Blogs) == 0 {
		parts = append(parts, m.styles.Muted.Render("No posts yet."))
	}
	for i, b := range m.content.Blogs {
		var line string
		if i == m.blogCursor {
			line = m.styles.Selected.Render("▸ " + b.Title)
		} else {
			line = m.styles.ItemDesc.Render("  " + b.Title)
		}
		if b.Date != "" {
			line += m.styles.ItemMeta.Render("   " + b.Date)
		}
		parts = append(parts, line)
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m portfolio) linksSection(width int) string {
	parts := []string{m.styles.Heading.Render("# links"), ""}
	if len(m.content.Links) == 0 {
		parts = append(parts, m.styles.Muted.Render("No links yet."))
	}
	for _, l := range m.content.Links {
		label := m.styles.LinkLabel.Render(fmt.Sprintf("%-11s", l.Label))
		parts = append(parts, m.styles.Accent.Render("→ ")+m.renderLink(label, l.URL))
	}
	parts = append(parts, "", m.styles.Muted.Render(linkTip))
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

 
func (m portfolio) renderLink(label, url string) string {
	if url == "" {
		return m.styles.LinkURL.Render(label)
	}
	rawURL := m.styles.LinkURL.Render(url)
	if strings.TrimSpace(label) == "" {
		return rawURL
	}
	clickable := hyperlink(label, url, m.styles.LinkLabel)
	return clickable + m.styles.Muted.Render("  ") + rawURL
}

 
func hyperlink(text, url string, style lipgloss.Style) string {
	rendered := style.Render(text)
	if url == "" {
		return rendered
	}
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, rendered)
}
