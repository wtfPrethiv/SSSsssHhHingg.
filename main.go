package main

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"pr3thiv-portfolio/internal/config"
	"pr3thiv-portfolio/internal/ui"
)

func main() {
	addr := envOr("PORT_ADDR", ":2222")
	hostKey := envOr("HOST_KEY_PATH", ".ssh/portfolio_ed25519")
	contentPath := envOr("CONTENT_PATH", "content.yaml")
	visitsPath := envOr("VISITS_LOG", "visits.log")

	logger := newLogger(visitsPath)

	// Validate content up front so a broken file fails loudly at boot.
	if _, err := config.Load(contentPath); err != nil {
		logger.Fatal("could not load content", "path", contentPath, "err", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKey),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler(contentPath, logger)),
			activeterm.Middleware(), // require an interactive terminal
			visitorLogging(logger),  // count and record every visitor
			logging.Middleware(),
		),
	)
	if err != nil {
		logger.Fatal("could not create server", "err", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("starting SSH portfolio server", "addr", addr)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			logger.Fatal("server error", "err", err)
		}
	}()

	<-done
	logger.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		logger.Error("shutdown error", "err", err)
	}
}

func teaHandler(contentPath string, logger *log.Logger) bubbletea.Handler {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "This program requires an interactive terminal.")
			return nil, nil
		}

		content, err := config.Load(contentPath)
		if err != nil {
			logger.Error("failed to load content for session", "err", err)
			wish.Fatalln(s, "Sorry, the portfolio is temporarily unavailable.")
			return nil, nil
		}

		// Build a renderer from the SSH session so Lip Gloss detects color
		// support from the connecting client, not the server's own stdout.
		renderer := bubbletea.MakeRenderer(s)
		return ui.NewRoot(content, renderer), []tea.ProgramOption{tea.WithAltScreen()}
	}
}

// visitorLogging records each connection so you can see how many people
// visited and when.

func visitorLogging(logger *log.Logger) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			host, _, err := net.SplitHostPort(s.RemoteAddr().String())
			if err != nil {
				host = s.RemoteAddr().String()
			}
			pty, _, _ := s.Pty()
			logger.Info("visitor connected",
				"time", time.Now().Format(time.RFC3339),
				"user", s.User(),
				"ip", host,
				"term", pty.Term,
				"client", s.Context().ClientVersion(),
			)
			start := time.Now()
			next(s)
			logger.Info("visitor disconnected",
				"user", s.User(),
				"ip", host,
				"duration", time.Since(start).Round(time.Second).String(),
			)
		}
	}
}

func newLogger(visitsPath string) *log.Logger {
	var w io.Writer = os.Stdout
	if f, err := os.OpenFile(visitsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
		w = io.MultiWriter(os.Stdout, f)
	}
	return log.NewWithOptions(w, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:         "portfolio",
	})
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
