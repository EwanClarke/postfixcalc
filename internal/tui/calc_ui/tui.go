package calcui

import tea "github.com/charmbracelet/bubbletea"

func Start() error {
	p := tea.NewProgram(InitialModel(), tea.WithMouseAllMotion(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
