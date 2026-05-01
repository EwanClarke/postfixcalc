package graphui

import tea "github.com/charmbracelet/bubbletea"

func Start() error {
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
