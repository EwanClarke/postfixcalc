package graphui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const renderTimeout = 2 * time.Second

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.updateCanvasSize()
		m.inputField.Width = m.termWidth/2 - 4 // account for border and space
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inputField, cmd = m.inputField.Update(msg)
	if m.updateGraph() {
		return m, tea.Quit
	}

	return m, cmd
}

func (m *model) updateCanvasSize() {
	graphHeight := m.termHeight - 5 // input+footer (1) + graph border (2) + extra (2)
	if graphHeight < 1 {
		graphHeight = 1
	}
	graphWidth := m.termWidth - 4
	if graphWidth < 1 {
		graphWidth = 1
	}

	m.canvas.Resize(graphWidth, graphHeight)
	m.canvas.SetBounds(-10, 10, -5, 5)
}

func (m *model) updateGraph() bool {
	expr := m.inputField.Value()
	if len(expr) == 0 {
		m.err = ""
		m.canvas.SetResult("")
		return false
	}

	err := m.canvas.ChangeExpression(expr)
	if err != nil {
		m.canvas.SetResult("")
		return false
	}

	resultChan := make(chan string, 1)
	go func() {
		result, _ := m.canvas.Render()
		resultChan <- result
	}()

	select {
	case <-time.After(renderTimeout):
		return true
	case result := <-resultChan:
		m.canvas.SetResult(result)
		m.err = ""
		return false
	}
}
