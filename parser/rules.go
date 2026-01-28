package parser

type OperatorProps struct {
	Precedence int
	RightAssoc bool
}

var Operators = map[string]OperatorProps {
	"+": {Precedence: 2, RightAssoc: false},
	"-": {Precedence: 2, RightAssoc: false},
	"*": {Precedence: 3, RightAssoc: false},
	"/": {Precedence: 3, RightAssoc: false},
	"^": {Precedence: 4, RightAssoc: true},
}
