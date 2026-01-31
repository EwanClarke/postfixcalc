package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/EwanClarke/postfixcalc/internal/lexer"
	"github.com/EwanClarke/postfixcalc/internal/parser"
	"github.com/EwanClarke/postfixcalc/internal/evaluator"
	"fmt"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case Nav:
		switch msg := msg.(type){
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				m.moveCursor(-1,0)
			case "down", "j":
				m.moveCursor(1, 0)
			case "left", "h":
				m.moveCursor(0, -1)
			case "right", "l":
				m.moveCursor(0, 1)
			case "enter", " ":
				// carry out button action
			case "/":
				m.mode = Edit
				m.inputField.Focus()
			}
		case tea.MouseMsg:
			m.mode = Cursor
		}
	case Edit:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.mode = Nav
				m.inputField.Blur()
			}
			before := m.inputField.View()
			var cmd tea.Cmd
			m.inputField, cmd = m.inputField.Update(msg)
			if m.inputField.View() != before {
				m.executeCalculation()
			}
			return m, cmd
		case tea.MouseMsg:
			m.mode = Cursor
		}
	case Cursor:
		switch msg.(type) {
		case tea.KeyMsg:
			m.mode = Nav
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
		}
	}
	
	return m, nil
}

func (m *model) moveCursor(rowDelta, colDelta int) {
	newRow := m.cursor.row + rowDelta
	if newRow >= 0 && newRow < len(m.primaryGrid) {
		m.cursor.row = newRow
	}

	newCol := m.cursor.col + colDelta
	if newCol >= 0 && newCol < len(m.primaryGrid[m.cursor.row]) {
		m.cursor.col = newCol
	}

	if m.cursor.col >= len(m.primaryGrid[m.cursor.row]) {
		m.cursor.col = len(m.primaryGrid[m.cursor.row])-1
	}
}

// func (m *model) updateInputs(msg tea.Msg) {

// }
//
func (m *model) executeCalculation() {
	rawInput := m.inputField.Value()

	if len(rawInput) == 0 {
		m.outputField = ""
	}
	
	l := lexer.New(rawInput)
	tokens, err := l.Tokenise()
	if err != nil {
		return
	}

	p := parser.New()
	postfix, err := p.Convert(tokens)
	if err != nil {
		return
	}

	e := evaluator.New()
	result, err := e.Evaluate(postfix)
	if err != nil {
		return
	} else {
		m.outputField = fmt.Sprintf("%g", result)
	}
	
}
