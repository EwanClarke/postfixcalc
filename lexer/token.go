package lexer

type TokenType int
const (
	Number TokenType = iota
	Operator
	Function
	LeftBrace
	RightBrace
	Negation
	Error
)

func (t TokenType) String() string {
	return [...]string{"Number", "Operator", "Function", "LeftBrace", "RightBrace", "Negation", "Error"}[t]
}

var operatorMap = map[string]TokenType {
	"+": Operator,
	"-": Operator,
	"*": Operator,
	"/": Operator,
	"^": Operator,
	"(": LeftBrace,
	")": RightBrace,
	"sin": Function,
	"cos": Function,
	"tan": Function,
}

type Token struct {
	Type TokenType
	Value string
}

