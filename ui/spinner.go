package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spinner spinner.Model
	quit    chan struct{}
	done    bool
}

func NewSpinner() *model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return &model{
		spinner: s,
		quit:    make(chan struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.done {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.QuitMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		return ""
	}
	return m.spinner.View() + " Thinking..."
}

// This is what you'll call from main()
func RunSpinnerUntil(done <-chan struct{}) {
	m := NewSpinner()
	program := tea.NewProgram(m)

	go func() {
		<-done
		program.Send(tea.Quit())
	}()

	_, _ = program.Run()
}
