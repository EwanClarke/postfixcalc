package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	screen := m.renderInputOutput()
	mainGrid := m.renderButtonGrid()
	footer := m.renderFooter()

	mainContent := lipgloss.JoinVertical(lipgloss.Left, screen, mainGrid)
	mainContent = m.styles.BorderAround.Render(mainContent)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, footer)
}

func (m model) renderInputOutput() string {
	screenContent := m.renderScreenContent()
	style := m.styles.Screen
	if m.mode == Edit {
		style = m.styles.ActiveScrn
	} else if m.cursorLocation == OnScreen {
		style = m.styles.ActiveScrn // Reuse active style for cursor
	}
	return style.Render(screenContent)
}

func (m model) renderButtonGrid() string {
	var rows []string
	for r, row := range m.primaryGrid {
		var currentRow []string
		for c, btn := range row {
			button := m.renderButton(btn, r, c)
			currentRow = append(currentRow, button)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
	}

	gridContents := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return m.styles.PrimaryGrid.Render(gridContents)
}

func (m model) renderButton(btn Button, row, col int) string {
	style := m.getButtonStyle(btn.Type)

	if m.cursorLocation == OnGrid && m.mode == Nav && m.cursor.row == row && m.cursor.col == col {
		style = m.getActiveButtonStyle(btn.Type)
	}

	return style.Render(btn.Label)
}

func (m model) getButtonStyle(btnType ButtonType) lipgloss.Style {
	switch btnType {
	case Operator:
		return m.styles.OpBtn
	case Action:
		return m.styles.ActnBtn
	case Control:
		return m.styles.CtrlBtn
	case Function:
		return m.styles.Button
	default:
		return m.styles.Button
	}
}

func (m model) getActiveButtonStyle(btnType ButtonType) lipgloss.Style {
	switch btnType {
	case Operator:
		return m.styles.ActiveOpBtn
	case Action:
		return m.styles.ActiveActnBtn
	case Control:
		return m.styles.ActiveCtrlBtn
	case Function:
		return m.styles.ActiveBtn
	default:
		return m.styles.ActiveBtn
	}
}

func (m model) renderFooter() string {
	return m.styles.Footer.Render("Press 'q' or 'ctrl+c' to exit.")
}

func (m model) renderScreenContent() string {
	inputStr := m.inputField.View()

	outputStr := m.styles.Output.Render("")
	if len(m.outputField) != 0 {
		outputStr = m.styles.Output.Render("= " + m.outputField)
	}

	return lipgloss.JoinVertical(lipgloss.Left, inputStr, outputStr)
}
