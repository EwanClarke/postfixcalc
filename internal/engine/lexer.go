package engine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	input  []rune
	pos    int
	tokens []Token
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Tokenise(input string) ([]Token, error) {
	l.input = []rune(input)
	l.pos = 0
	l.tokens = l.tokens[:0] // Clear tokens slice

	for l.pos < len(l.input) {
		value := l.ExtractToken()
		if value == "" {
			break
		}

		tokenType := l.Categorise(value)

		if tokenType == Error {
			return nil, fmt.Errorf("unknown token '%s'", value)
		} else if tokenType == Function {
			value = strings.ToLower(value)
		} else if tokenType == Negation {
			value = "-u"
		}

		l.tokens = append(l.tokens, Token{Value: value, Type: tokenType})
	}
	return l.tokens, nil
}

func (l *Lexer) ExtractToken() string {
	for l.pos < len(l.input) && l.input[l.pos] == ' ' {
		l.pos++
	}
	if l.pos >= len(l.input) {
		return ""
	}

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

	lowercaseValue := strings.ToLower(tokenValue)
	if tType, ok := operatorMap[lowercaseValue]; ok {
		return tType
	}

	if value, err := strconv.ParseFloat(tokenValue, 64); err == nil {
		if value < 0 {
			return Error
		}
		return Number
	}

	if l.isVariable(tokenValue) {
		return Variable
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
	return previousTokenType == LeftBrace ||
		previousTokenType == Operator ||
		previousTokenType == Negation ||
		previousTokenType == Function
}

func (l *Lexer) isVariable(tokenValue string) bool {
	if len(tokenValue) == 0 {
		return false
	}

	for _, c := range tokenValue {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}
