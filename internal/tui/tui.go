package tui

import tea "github.com/charmbracelet/bubbletea"

func Start() error {
	p := tea.NewProgram(InitialModel(), tea.WithMouseAllMotion())
	_, err := p.Run()
	return err
}
