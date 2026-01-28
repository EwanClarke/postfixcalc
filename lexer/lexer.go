package lexer
import (
	"strconv"
	"fmt"
	"unicode"
)

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

func (l *Lexer) Tokenise() ([]Token, error) {
	for l.pos < len(l.input) {
		value := l.ExtractToken()
		if value == "" {
			break
		}
		
		tokenType := l.Categorise(value)
		if tokenType == Error {
			return nil, fmt.Errorf("unknown token '%s'", value)
		}

		l.tokens = append(l.tokens, Token{Value: value, Type: tokenType})
	}
	return l.tokens, nil
}

func (l *Lexer) ExtractToken() string {
	for l.pos < len(l.input) && l.input[l.pos] == ' ' {
		l.pos++
	}
	if l.pos >= len(l.input) {return ""}
	
	start := l.pos
	if _, ok := operatorMap[string(l.input[start])]; ok {
		l.pos++
	} else if unicode.IsDigit(l.input[start]) {
		for l.pos < len(l.input) && (unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
			l.pos++
		}
	} else if unicode.IsLetter(l.input[start]) {
		for l.pos < len(l.input) && unicode.IsLetter(l.input[l.pos]) {
			l.pos++
		}
	} else {
		l.pos++
	}
	
	return string(l.input[start:l.pos])
}

func (l *Lexer) Categorise(tokenValue string) TokenType {
	if l.isNegation(tokenValue) {
		return Negation
	}

	if tType, ok := operatorMap[tokenValue]; ok {
		return tType
	}

	if value, err := strconv.ParseFloat(tokenValue, 64); err == nil {
		if value < 0 {
			return Error
		}
		return Number
	}

	return Error
}

func (l *Lexer) isNegation(tokenValue string) bool {
	if tokenValue != "-" {
		return false
	}

	if len(l.tokens) == 0 {
		return true
	}
	
	previousTokenType := l.tokens[len(l.tokens)-1].Type
	if previousTokenType == LeftBrace || previousTokenType == Operator || previousTokenType == Negation {
		return true
	}
	return false
}
