package calcui

import (
	"math/big"

	"github.com/EwanClarke/postfixcalc/internal/engine"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

type CursorLocation int

const (
	OnGrid CursorLocation = iota
	OnScreen
	OnTopBar
)

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
	Action  // evaluate
)

type Button struct {
	Label string
	Value string
	Type  ButtonType
	Width int
}

type OutputMode int

const (
	Standard OutputMode = iota
	Decimal
)

func (o OutputMode) String() string {
	return [...]string{"S", "D"}[o]
}



type model struct {
	styles         Styles
	inputField     textinput.Model
	outputField    string
	mode           InputMode
	cursor         CursorPos
	cursorLocation CursorLocation
	topBarCol      int // 0 = AngleMode toggle, 1 = OutputMode (S/D) toggle
	outputMode     OutputMode
	angleMode      engine.AngleMode
	lastResult     *big.Rat // cached result for S/D re-formatting
	lastExact      bool     // false when result used float64 approximation (e.g. trig)
	primaryGrid    [][]Button
	termWidth      int
	termHeight     int
}

func InitialModel() model {
	ti := textinput.New()
	ti.CharLimit = 25
	ti.Width = 25

	m := model{
		mode:           Nav,
		styles:         DefaultStyles(),
		cursorLocation: OnGrid,
		topBarCol:      0,
		lastExact:      true,
		termWidth:      80, // default fallback
		termHeight:     24, // default fallback
		primaryGrid: [][]Button{
			{{Label: "C", Type: Control, Value: "clear"}, {Label: "(", Type: Operator, Value: "("}, {Label: ")", Type: Operator, Value: ")"}, {Label: "*", Type: Operator, Value: "*"}},
			{{Label: "7", Value: "7"}, {Label: "8", Value: "8"}, {Label: "9", Value: "9"}, {Label: "/", Type: Operator, Value: "/"}},
			{{Label: "4", Value: "4"}, {Label: "5", Value: "5"}, {Label: "6", Value: "6"}, {Label: "+", Type: Operator, Value: "+"}},
			{{Label: "1", Value: "1"}, {Label: "2", Value: "2"}, {Label: "3", Value: "3"}, {Label: "-", Type: Operator, Value: "-"}},
			{{Label: ".", Value: "."}, {Label: "0", Value: "0"}, {Label: "^", Type: Operator, Value: "^"}, {Label: "=", Type: Action, Value: "evaluate"}},
		},

		inputField: ti,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
