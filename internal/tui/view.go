package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {

	screenContent := m.renderScreenContent()

	screen := m.styles.Screen.Render(screenContent)
	if m.mode == Edit {
		screen = m.styles.ActiveScrn.Render(screenContent)

	}

	var rows []string
	for r, row := range m.primaryGrid {
		var currentRow []string

		for c, btn := range row {
			style := m.styles.Button
			if btn.Type == Operator {
				style = m.styles.OpBtn
			}

			switch btn.Type {
			case Operator:
				style = m.styles.OpBtn
			case Action:
				style = m.styles.ActnBtn
			case Control:
				style = m.styles.CtrlBtn
			case Function:
			}

			if m.mode == Nav && m.cursor.row == r && m.cursor.col == c {
				style = m.styles.ActiveBtn

				switch btn.Type {
				case Operator:
					style = m.styles.ActiveOpBtn
				case Action:
					style = m.styles.ActiveActnBtn
				case Control:
					style = m.styles.ActiveCtrlBtn
				case Function:
			}
			}
			currentRow = append(currentRow, style.Render(btn.Label))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
	}
	gridContents := lipgloss.JoinVertical(lipgloss.Left, rows...)
	mainGrid := m.styles.PrimaryGrid.Render(gridContents)
	mainContent := lipgloss.JoinVertical(lipgloss.Left, screen, mainGrid)
	footer := m.styles.Footer.Render("Press 'q' or 'ctrl+c' to exit.")
	return lipgloss.JoinVertical(lipgloss.Left, mainContent, footer)
}

func (m model) renderScreenContent() string {
	inputStr := m.inputField.View()

	outputStr := m.styles.Output.Render(m.outputField)

	return lipgloss.JoinVertical(lipgloss.Left, inputStr, outputStr)
}
