package parser

type OperatorProps struct {
	Precedence int
	RightAssoc bool
}

var OperatorPropsMap = map[string]OperatorProps {
	"+": {Precedence: 2, RightAssoc: false},
	"-": {Precedence: 2, RightAssoc: false},
	"*": {Precedence: 3, RightAssoc: false},
	"/": {Precedence: 3, RightAssoc: false},
	"^": {Precedence: 4, RightAssoc: true},
	"-u": {Precedence: 5, RightAssoc: true},
	"sin": {Precedence: 5, RightAssoc: true},
	"cos": {Precedence: 5, RightAssoc: true},
	"tan": {Precedence: 5, RightAssoc: true},
}
