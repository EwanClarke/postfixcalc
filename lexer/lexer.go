package lexer
import (
	"strconv"
	"fmt"
)

type TokenType int
const (
	Number TokenType = iota
	Operator
	LeftBrace
	RightBrace
	Negation
	Error
)

func (t TokenType) String() string {
	return [...]string{"Number", "Operator", "LeftBrace", "RightBrace", "Negation", "Error"}[t]
}

var operatorMap = map[string]TokenType {
	"+": Operator,
	"-": Operator,
	"*": Operator,
	"/": Operator,
	"^": Operator,
	"(": LeftBrace,
	")": RightBrace,
}

type Token struct {
	Type TokenType
	Value string
}

type Lexer struct {
	input []rune
	pos int
	tokens []Token
}

func New(input string) *Lexer {
	return &Lexer {
		input: []rune(input),
	}
}

func (l *Lexer) Tokenise() []Token {
	for l.pos < len(l.input) {
		l.tokens = append(l.tokens, l.nextToken())
		fmt.Println(l.tokens)
	}
	return l.tokens
}

func (l *Lexer) nextToken() Token {
	var endPos int = l.pos
	for ; endPos<len(l.input); endPos++ {
		if l.input[endPos] == ' ' {
			break
		}
	}
	
	var newToken Token
	newToken.Value = string(l.input[l.pos:endPos])
	newToken.Type = l.Categorise(newToken.Value)

	l.pos = endPos+1
	return newToken
}

func (l *Lexer) Categorise(tokenValue string) TokenType {
	if tType, ok := operatorMap[tokenValue]; ok {
		return tType
	}

	if _, err := strconv.ParseFloat(tokenValue, 64); err == nil {
		return Number
	}

	return Error
}
