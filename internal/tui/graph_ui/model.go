package graphui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	styles     Styles
	inputField textinput.Model
	canvas     *Canvas
	termWidth  int
	termHeight int
	err        string
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter expression (e.g. x^2)"
	ti.Focus()

	c := NewCanvas()

	m := model{
		styles:     DefaultStyles(),
		inputField: ti,
		canvas:     c,
		termWidth:  80,
		termHeight: 24,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}
