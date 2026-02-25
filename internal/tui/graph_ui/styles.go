package graphui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	GraphBorder lipgloss.Style
	InputBorder lipgloss.Style
	Footer      lipgloss.Style
	Error       lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		GraphBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("6")),
		InputBorder: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("5")).
			Padding(0),
		Footer: lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")),
	}
}
