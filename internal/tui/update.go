package tui

import (
	"fmt"
	"github.com/EwanClarke/postfixcalc/internal/evaluator"
	"github.com/EwanClarke/postfixcalc/internal/lexer"
	"github.com/EwanClarke/postfixcalc/internal/parser"
	"github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if shouldQuit, cmd := m.handleGlobalKeys(msg); shouldQuit {
		return m, cmd
	}

	switch m.mode {
	case Nav:
		return m.handleNavMode(msg)
	case Edit:
		return m.handleEditMode(msg)
	case Cursor:
		return m.handleCursorMode(msg)
	}

	return m, nil
}

func (m model) handleGlobalKeys(msg tea.Msg) (bool, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return true, tea.Quit
		}
	}
	return false, nil
}

func (m model) handleNavMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1, 0)
		case "down", "j":
			m.moveCursor(1, 0)
		case "left", "h":
			m.moveCursor(0, -1)
		case "right", "l":
			m.moveCursor(0, 1)
		case "enter", " ":
			switch m.cursorLocation {
			case OnGrid:
				m.handleButtonPress()
			case OnScreen:
				m.mode = Edit
				m.inputField.Focus()
			}
		case "/":
			m.mode = Edit
			m.inputField.Focus()
		}
	case tea.MouseMsg:
		m.mode = Cursor
	}
	return m, nil
}

func (m model) handleEditMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = Nav
			m.inputField.Blur()
		default:
			before := m.inputField.View()
			var cmd tea.Cmd
			m.inputField, cmd = m.inputField.Update(msg)
			if m.inputField.View() != before {
				m.executeCalculation()
			}
			return m, cmd
		}
	case tea.MouseMsg:
		m.mode = Cursor
	}
	return m, nil
}

func (m model) handleCursorMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.mode = Nav
	}
	return m, nil
}

func (m *model) moveCursor(rowDelta, colDelta int) {
	switch m.cursorLocation {
	case OnGrid:
		newRow := m.cursor.row + rowDelta
		if newRow < 0 {
			m.cursorLocation = OnScreen
		} else if newRow < len(m.primaryGrid) {
			m.cursor.row = newRow
			newRowLastIndex := len(m.primaryGrid[newRow]) - 1
			if m.cursor.col > newRowLastIndex {
				m.cursor.col = newRowLastIndex
			}
		}

		if m.cursorLocation == OnGrid {
			newCol := m.cursor.col + colDelta
			if newCol >= 0 && newCol < len(m.primaryGrid[m.cursor.row]) {
				m.cursor.col = newCol
			}
		}
	case OnScreen:
		if rowDelta > 0 {
			m.cursorLocation = OnGrid
			m.cursor.row = 0
		} else if rowDelta < 0 {
			m.cursorLocation = OnTopBar
			m.cursor.row = 0
		}
	case OnTopBar:
		if rowDelta > 0 {
			m.cursorLocation = OnScreen
			m.cursor.row = 0
		}
	}
}

func (m *model) handleButtonPress() {
	selectedButton := m.getSelectedButton()

	switch selectedButton.Type {
	case Number, Operator, Function:
		m.inputField.SetValue(m.inputField.Value() + selectedButton.Value)
	case Control:
		switch selectedButton.Value {
		case "clear":
			m.inputField.SetValue("")
			m.outputField = ""
		case "delete":
			m.inputField.SetValue(m.inputField.Value()[:len(m.inputField.Value())-1])
		}
	case Action:
		switch selectedButton.Value {
		case "evaluate":
			m.executeCalculation()
		}
	}
}

func (m model) getSelectedButton() Button {
	return m.primaryGrid[m.cursor.row][m.cursor.col]
}

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
