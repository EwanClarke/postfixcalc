package graphui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	graphContent := m.renderGraph()
	inputWithFooter := m.renderInputWithFooter()

	mainContent := lipgloss.JoinVertical(lipgloss.Left, graphContent, inputWithFooter)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		mainContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
}

func (m model) renderGraph() string {
	graph := m.canvas.Result()
	if graph == "" {
		width, height := m.canvas.Size()
		placeholder := lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			"Enter an expression to plot",
		)
		return m.styles.GraphBorder.Render(placeholder)
	}
	return m.styles.GraphBorder.Render(graph)
}

func (m model) renderInputWithFooter() string {
	inputView := m.inputField.View()
	if m.err != "" {
		inputView = m.styles.Error.Render(" " + m.err)
	} else {
		inputView = m.styles.InputBorder.Render(" " + inputView)
	}

	instructions := "Enter expression | q / ctrl+c to quit"
	footerWidth := m.termWidth/2 - 2
	footer := lipgloss.Place(footerWidth, 3, lipgloss.Center, lipgloss.Center, instructions, lipgloss.WithWhitespaceChars(" "))

	return lipgloss.JoinHorizontal(lipgloss.Left, inputView, footer)
}

func (m model) renderInput() string {
	inputView := m.inputField.View()
	return m.styles.InputBorder.Render(" " + inputView)
}

func (m model) renderFooter() string {
	instructions := "Enter expression | q / ctrl+c to quit"
	return m.styles.Footer.Render(instructions)
}
