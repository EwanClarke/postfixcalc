package tui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	NavMode lipgloss.Style
	EditMode lipgloss.Style

	Screen lipgloss.Style
	ActiveScrn lipgloss.Style
	Output lipgloss.Style

	Button lipgloss.Style
	ActiveBtn lipgloss.Style
	OpBtn lipgloss.Style
	ActiveOpBtn lipgloss.Style
	CtrlBtn lipgloss.Style
	ActiveCtrlBtn lipgloss.Style
	ActnBtn lipgloss.Style
	ActiveActnBtn lipgloss.Style

	Footer lipgloss.Style

	PrimaryGrid lipgloss.Style
}

func DefaultStyles() Styles {
	s := Styles{}

	width := 30

	s.Button = lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0,0, 0,0).
		Foreground(lipgloss.Color("7")).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("7"))

	s.ActiveBtn = s.Button.Copy().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("15"))

	s.OpBtn = s.Button.Copy().
		Foreground(lipgloss.Color("6")).
		BorderForeground(lipgloss.Color("6"))

	s.ActiveOpBtn = s.ActiveBtn.Copy().
		Foreground(lipgloss.Color("14")).
		BorderForeground(lipgloss.Color("14"))

	s.CtrlBtn = s.Button.Copy().
		Foreground(lipgloss.Color("1")).
		BorderForeground(lipgloss.Color("1"))

	s.ActiveCtrlBtn = s.ActiveBtn.Copy().
		Foreground(lipgloss.Color("9")).
		BorderForeground(lipgloss.Color("9"))

	s.ActnBtn = s.Button.Copy().
		Foreground(lipgloss.Color("2")).
		BorderForeground(lipgloss.Color("2"))

	s.ActiveActnBtn = s.ActiveBtn.Copy().
		Foreground(lipgloss.Color("10")).
		BorderForeground(lipgloss.Color("10"))

	s.Screen = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("5")).
		Width(width-2).
		Height(2)
	
	s.ActiveScrn = s.Screen.Copy().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("13"))

	s.Output = lipgloss.NewStyle().
		Width(width-4).
		Align(lipgloss.Right)
	
	s.Footer = lipgloss.NewStyle().
		Width(width-2).
		Border(lipgloss.RoundedBorder()).
		AlignHorizontal(lipgloss.Center)
	
	s.PrimaryGrid = lipgloss.NewStyle().
		Width(width).
		AlignHorizontal(lipgloss.Center)

	return s
}
