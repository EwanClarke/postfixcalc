package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)
type InputMode int
const (
	Nav InputMode = iota
	Edit
	Cursor
)

func (i InputMode) String() string {
	return [...]string{"Nav", "Edit", "Cursor"}[i]
}

type CursorPos struct {
	row int
	col int
}

type ButtonType int
const (
	Number ButtonType = iota
	Operator
	Function
	Control // clear, delete
	Action // evaluate
	Mode // Rad/Deg toggle
)

type Button struct {
	Label string
	Value string
	Type ButtonType
	Width int
	Action func()
}

type model struct{
	styles Styles
	inputField textinput.Model
	outputField string
	mode InputMode
	cursor CursorPos
	primaryGrid [][]Button
}

func InitialModel() model {
	ti := textinput.New()
	ti.CharLimit = 25
	ti.Width = 25

	m := model{
		mode: Nav,
		styles: DefaultStyles(),
		primaryGrid: [][]Button{
			{ {Label: "C", Type: Control}, {Label: "(", Type: Operator}, {Label: ")", Type: Operator}, {Label: "*", Type: Operator}, },
			{ {Label: "7"}, {Label: "8"}, {Label: "9"}, {Label: "/", Type: Operator}, },
			{ {Label: "4"}, {Label: "5"}, {Label: "6"}, {Label: "+", Type: Operator}, },
			{ {Label: "1"}, {Label: "2"}, {Label: "3"}, {Label: "-", Type: Operator}, },
			{ {Label: "."}, {Label: "0"}, {Label: "^", Type: Operator}, {Label: "=", Type: Action}, },
		},

		inputField: ti,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
