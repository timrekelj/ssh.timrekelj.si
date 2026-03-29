package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello, World!\n\nPress q to quit.\n"
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:2222"),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return model{}, []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting SSH server on port 2222")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
