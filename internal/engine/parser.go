package engine

import (
	"errors"
	"fmt"
)

type OperatorProps struct {
	Precedence int
	RightAssoc bool
}

var OperatorPropsMap = map[string]OperatorProps{
	"+":   {Precedence: 2, RightAssoc: false},
	"-":   {Precedence: 2, RightAssoc: false},
	"*":   {Precedence: 3, RightAssoc: false},
	"/":   {Precedence: 3, RightAssoc: false},
	"^":   {Precedence: 4, RightAssoc: true},
	"-u":  {Precedence: 5, RightAssoc: true},
	"sin": {Precedence: 5, RightAssoc: true},
	"cos": {Precedence: 5, RightAssoc: true},
	"tan": {Precedence: 5, RightAssoc: true},
}

type Parser struct {
	input         []Token
	previousToken Token
	hasPrev       bool
	operatorStack []Token
	outputQueue   []Token
	err           error
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) pushOperator(token Token) {
	p.operatorStack = append(p.operatorStack, token)
}

func (p *Parser) popOperator() Token {
	stackSize := len(p.operatorStack)
	if stackSize == 0 {
		return Token{Type: Error}
	}

	removedOperator := p.operatorStack[stackSize-1]
	p.operatorStack = p.operatorStack[:stackSize-1]
	return removedOperator
}

func (p *Parser) hasOperators() bool {
	return len(p.operatorStack) > 0
}

func (p *Parser) peekOperator() Token {
	if len(p.operatorStack) == 0 {
		return Token{Type: Error}
	}
	return p.operatorStack[len(p.operatorStack)-1]
}

func (p *Parser) enqueue(token Token) {
	p.outputQueue = append(p.outputQueue, token)
}

func (p *Parser) dequeue() Token {
	if len(p.outputQueue) == 0 {
		return Token{Type: Error}
	}

	removedToken := p.outputQueue[0]
	p.outputQueue = p.outputQueue[1:]
	return removedToken
}

func (p *Parser) setError(msg string) {
	p.err = errors.New(msg)
}

func (p *Parser) Convert(tokens []Token) ([]Token, error) {
	for _, t := range tokens {
		p.handleToken(t)

		if p.err != nil {
			return nil, p.err
		}
	}

	p.validateEnd()
	p.clearRemainingOperators()
	if p.err != nil {
		return nil, p.err
	}
	return p.outputQueue, nil
}

func (p *Parser) handleToken(token Token) {
	p.handleImplicitMult(token)
	p.validateSequence(token)

	switch token.Type {
	case Number, Variable:
		p.enqueue(token)
	case LeftBrace:
		p.pushOperator(token)
	case RightBrace:
		p.handleRightBrace()
	case Operator, Function, Negation:
		p.handleOperator(token)
	}

	p.previousToken = token
	p.hasPrev = true
}

func (p *Parser) handleRightBrace() {
	for len(p.operatorStack) > 0 && p.peekOperator().Type != LeftBrace {
		p.enqueue(p.popOperator())
	}
	p.popOperator()
}

func (p *Parser) handleOperator(token Token) {
	tokenProps := OperatorPropsMap[token.Value]
	for p.hasOperators() {
		topToken := p.peekOperator()
		if topToken.Type == LeftBrace {
			break
		}

		topProps := OperatorPropsMap[topToken.Value]
		if topProps.Precedence > tokenProps.Precedence ||
			(topProps.Precedence == tokenProps.Precedence && !tokenProps.RightAssoc) {

			p.enqueue(p.popOperator())
		} else {
			break
		}
	}
	p.pushOperator(token)
}

func (p *Parser) handleImplicitMult(currentToken Token) {
	if !p.hasPrev {
		return
	}

	prev := p.previousToken.Type
	curr := currentToken.Type

	insertMult := false
	if (prev == Number || prev == RightBrace) &&
		(curr == LeftBrace || curr == Function || curr == Variable) {

		insertMult = true
	}

	if prev == RightBrace && curr == Number {
		insertMult = true
	}

	if insertMult {
		multToken := Token{Type: Operator, Value: "*"}
		p.handleOperator(multToken)
		p.previousToken = multToken
	}
}

func (p *Parser) validateSequence(curr Token) {
	if !p.hasPrev {
		if curr.Type == Operator {
			p.setError("Syntax Error: expression cannot start with a binary operator")
		} else if curr.Type == RightBrace {
			p.setError("Syntax Error: expression cannot start with a closing bracket")
		}
		return
	}

	prev := p.previousToken

	switch curr.Type {
	case Operator, RightBrace:
		if prev.Type == Operator ||
			prev.Type == LeftBrace ||
			prev.Type == Negation ||
			prev.Type == Function {

			p.setError(fmt.Sprintf("Syntax Error: operator '%s' cannot follow '%s'", curr.Value, prev.Value))
		}
	case Number, Function, LeftBrace, Negation:
		if prev.Type == Number || prev.Type == RightBrace {
			p.setError(fmt.Sprintf("Syntax Error: %s value '%v' should not follow %s", curr.Type, curr.Value, prev.Type))
		}
	}
}

func (p *Parser) validateEnd() {
	last := p.previousToken
	if last.Type == Operator ||
		last.Type == Function ||
		last.Type == Negation ||
		last.Type == LeftBrace {
		p.setError(fmt.Sprintf("Syntax Error: expression ends prematurely after '%s'", last.Value))
	}
}

func (p *Parser) clearRemainingOperators() {
	for len(p.operatorStack) > 0 {
		operator := p.popOperator()

		if operator.Type == LeftBrace {
			p.setError("Mismatched Brackets: opening '(' was never closed")
			return
		} else if operator.Type == RightBrace {
			p.setError("Syntax Error: unexpected closing ')'")
			return
		}

		p.enqueue(operator)
	}
}
