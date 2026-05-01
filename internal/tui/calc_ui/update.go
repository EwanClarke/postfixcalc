package calcui

import (
	"fmt"

	"github.com/EwanClarke/postfixcalc/internal/engine"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	}

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
			case OnTopBar:
				m.handleTopBarPress()
			}
		case "/":
			m.mode = Edit
			m.inputField.Focus()
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			m.mode = Cursor
		}
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
		if msg.Type == tea.MouseLeft {
			m.mode = Cursor
			m.inputField.Blur()
		}
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
			m.topBarCol = 0
		}

	case OnTopBar:
		if rowDelta > 0 {
			// Move down into the screen area
			m.cursorLocation = OnScreen
		} else if colDelta != 0 {
			// Navigate left/right between top-bar items (0 = Deg/Rad, 1 = S/D)
			newCol := m.topBarCol + colDelta
			if newCol >= 0 && newCol <= 1 {
				m.topBarCol = newCol
			}
		}
	}
}

// handleTopBarPress handles Enter/Space on a focused top-bar toggle.
func (m *model) handleTopBarPress() {
	switch m.topBarCol {
	case 0:
		// Toggle AngleMode (Deg ↔ Rad)
		if m.angleMode == engine.Degrees {
			m.angleMode = engine.Radians
		} else {
			m.angleMode = engine.Degrees
		}
		// Recalculate immediately so any trig functions update
		m.executeCalculation()
	case 1:
		// Toggle S/D OutputMode — only available when last result is exact
		if !m.lastExact {
			return
		}
		if m.outputMode == Standard {
			m.outputMode = Decimal
		} else {
			m.outputMode = Standard
		}
		// Re-format the cached result immediately without recalculating
		m.reformatOutput()
	}
}

// reformatOutput re-renders outputField from the cached lastResult using the
// current outputMode. Called after toggling S/D.
func (m *model) reformatOutput() {
	if m.lastResult == nil {
		return
	}
	if m.outputMode == Decimal || !m.lastExact {
		m.outputField = engine.FormatDecimal(m.lastResult)
	} else {
		m.outputField = engine.FormatMixed(m.lastResult)
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
			m.lastResult = nil
			m.lastExact = true
		case "delete":
			v := m.inputField.Value()
			if len(v) > 0 {
				m.inputField.SetValue(v[:len(v)-1])
			}
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

func (m *model) handleMouseClick(msg tea.MouseMsg) {
	// No button click handling - only cursor mode switching
}

func (m *model) executeCalculation() {
	rawInput := m.inputField.Value()

	if len(rawInput) == 0 {
		m.outputField = ""
		m.lastResult = nil
		m.lastExact = true
		return
	}

	e := engine.NewEngine()
	e.AngleMode = m.angleMode
	result, exact, err := e.Calculate(rawInput)
	if err != nil {
		m.outputField = fmt.Sprintf("error %v", err)
		m.lastResult = nil
		m.lastExact = true
		return
	}

	m.lastResult = result
	m.lastExact = exact
	m.reformatOutput()
}
