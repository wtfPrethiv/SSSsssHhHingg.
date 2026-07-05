package ui

import (
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	playFPS   = 30
	maxBalls  = 60
	ballGlyph = "●"
)

// playGravity pulls particles down. Origin is the top-left corner, so a
// positive Y acceleration means "downward" on screen.
var playGravity = harmonica.Vector{X: 0, Y: 55, Z: 0}

type frameMsg time.Time

func frameCmd() tea.Cmd {
	return tea.Tick(time.Second/playFPS, func(t time.Time) tea.Msg { return frameMsg(t) })
}

type ball struct {
	proj   *harmonica.Projectile
	x, y   float64
	px, py float64
	color  lipgloss.Color
}

type playground struct {
	styles  Styles
	w, h    int // inner canvas size in cells
	balls   []ball
	running bool
}

func newPlayground(styles Styles) playground {
	return playground{styles: styles}
}

func (m *playground) setSize(w, h int) {
	// Leave room for the rounded border (1 cell each side).
	m.w = w - 2
	m.h = h - 2
	if m.w < 4 {
		m.w = 4
	}
	if m.h < 3 {
		m.h = 3
	}
}

func (m *playground) newBall() ball {
	dt := harmonica.FPS(playFPS)
	x := rand.Float64() * float64(m.w-1)
	y := float64(m.h - 1)
	vx := (rand.Float64()*2 - 1) * 22
	vy := -(18 + rand.Float64()*28) // launch upward
	return ball{
		proj:  harmonica.NewProjectile(dt, harmonica.Point{X: x, Y: y, Z: 0}, harmonica.Vector{X: vx, Y: vy, Z: 0}, playGravity),
		x:     x,
		y:     y,
		px:    x,
		py:    y,
		color: playgroundPalette[rand.Intn(len(playgroundPalette))],
	}
}

func (m *playground) spawn(n int) {
	for i := 0; i < n; i++ {
		if len(m.balls) >= maxBalls {
			m.balls = m.balls[1:]
		}
		m.balls = append(m.balls, m.newBall())
	}
}

// start makes the simulation live and seeds a few particles.
func (m *playground) start() tea.Cmd {
	m.running = true
	if len(m.balls) == 0 {
		m.spawn(8)
	}
	return frameCmd()
}

func (m *playground) stop() {
	m.running = false
}

func (m *playground) step() {
	dt := harmonica.FPS(playFPS)
	const damp = 0.78
	maxX := float64(m.w - 1)
	maxY := float64(m.h - 1)

	for i := range m.balls {
		b := &m.balls[i]
		p := b.proj.Update()
		b.x, b.y = p.X, p.Y

		vx := b.x - b.px
		vy := b.y - b.py
		bounced := false

		if b.x <= 0 {
			b.x, vx, bounced = 0, math.Abs(vx)*damp, true
		} else if b.x >= maxX {
			b.x, vx, bounced = maxX, -math.Abs(vx)*damp, true
		}
		if b.y >= maxY {
			b.y = maxY
			vy = -math.Abs(vy) * damp
			vx *= 0.96
			bounced = true
		} else if b.y < 0 {
			b.y, vy, bounced = 0, math.Abs(vy)*damp, true
		}

		if bounced {
			vel := harmonica.Vector{X: vx / dt, Y: vy / dt, Z: 0}
			b.proj = harmonica.NewProjectile(dt, harmonica.Point{X: b.x, Y: b.y, Z: 0}, vel, playGravity)
		}
		b.px, b.py = b.x, b.y
	}
}

func (m playground) update(msg tea.Msg) (playground, tea.Cmd) {
	switch msg := msg.(type) {
	case frameMsg:
		if !m.running {
			return m, nil
		}
		m.step()
		return m, frameCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case " ", "space":
			m.spawn(6)
		case "r", "R":
			m.balls = nil
			m.spawn(8)
		case "f", "F":
			m.spawn(1)
		}
	}
	return m, nil
}

func (m playground) view() string {
	grid := make([][]int, m.h)
	for y := range grid {
		grid[y] = make([]int, m.w)
	}
	for i := range m.balls {
		x := int(math.Round(m.balls[i].x))
		y := int(math.Round(m.balls[i].y))
		if x >= 0 && x < m.w && y >= 0 && y < m.h {
			grid[y][x] = i + 1
		}
	}

	var b strings.Builder
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			if idx := grid[y][x]; idx > 0 {
				b.WriteString(lipgloss.NewStyle().Foreground(m.balls[idx-1].color).Render(ballGlyph))
			} else {
				b.WriteByte(' ')
			}
		}
		if y < m.h-1 {
			b.WriteByte('\n')
		}
	}

	canvas := m.styles.PlayFrame.Render(b.String())
	info := m.styles.PlayInfo.Render("space launch · f drop one · r reset  —  harmonica physics")
	return lipgloss.JoinVertical(lipgloss.Left, canvas, info)
}
