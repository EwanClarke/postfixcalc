package parser

import (
	"errors"
	"fmt"
	"github.com/EwanClarke/postfixcalc/internal/lexer"
)

func New() *Parser {
	return &Parser{}
}

type Parser struct {
	input         []lexer.Token
	previousToken lexer.Token
	hasPrev       bool
	operatorStack []lexer.Token
	outputQueue   []lexer.Token
	err           error
}

func (p *Parser) pushOperator(token lexer.Token) {
	p.operatorStack = append(p.operatorStack, token)
}

func (p *Parser) popOperator() lexer.Token {
	stackSize := len(p.operatorStack)
	if stackSize == 0 {
		return lexer.Token{Type: lexer.Error}
	}

	removedOperator := p.operatorStack[stackSize-1]
	p.operatorStack = p.operatorStack[:stackSize-1]
	return removedOperator
}

func (p *Parser) hasOperators() bool {
	return len(p.operatorStack) > 0
}

func (p *Parser) peekOperator() lexer.Token {
	if len(p.operatorStack) == 0 {
		return lexer.Token{Type: lexer.Error}
	}
	return p.operatorStack[len(p.operatorStack)-1]
}

func (p *Parser) enqueue(token lexer.Token) {
	p.outputQueue = append(p.outputQueue, token)
}

func (p *Parser) dequeue() lexer.Token {
	if len(p.outputQueue) == 0 {
		return lexer.Token{Type: lexer.Error}
	}

	removedToken := p.outputQueue[0]
	p.outputQueue = p.outputQueue[1:]
	return removedToken
}

func (p *Parser) setError(msg string) {
	p.err = errors.New(msg)
}

func (p *Parser) Convert(tokens []lexer.Token) ([]lexer.Token, error) {
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

func (p *Parser) handleToken(token lexer.Token) {
	p.handleImplicitMult(token)
	p.validateSequence(token)

	switch token.Type {
	case lexer.Number:
		p.enqueue(token)
	case lexer.LeftBrace:
		p.pushOperator(token)
	case lexer.RightBrace:
		p.handleRightBrace()
	case lexer.Operator, lexer.Function, lexer.Negation:
		p.handleOperator(token)
	}

	p.previousToken = token
	p.hasPrev = true
}

func (p *Parser) handleRightBrace() {
	for len(p.operatorStack) > 0 && p.peekOperator().Type != lexer.LeftBrace {
		p.enqueue(p.popOperator())
	}
	p.popOperator()
}

func (p *Parser) handleOperator(token lexer.Token) {
	tokenProps := OperatorPropsMap[token.Value]
	for p.hasOperators() {
		topToken := p.peekOperator()
		if topToken.Type == lexer.LeftBrace {
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

func (p *Parser) handleImplicitMult(currentToken lexer.Token) {
	if !p.hasPrev {
		return
	}

	prev := p.previousToken.Type
	curr := currentToken.Type

	insertMult := false
	if (prev == lexer.Number || prev == lexer.RightBrace) &&
		(curr == lexer.LeftBrace || curr == lexer.Function) {

		insertMult = true
	}

	if prev == lexer.RightBrace && curr == lexer.Number {
		insertMult = true
	}

	if insertMult {
		multToken := lexer.Token{Type: lexer.Operator, Value: "*"}
		p.handleOperator(multToken)
		p.previousToken = multToken
	}
}

func (p *Parser) validateSequence(curr lexer.Token) {
	if !p.hasPrev {
		if curr.Type == lexer.Operator {
			p.setError("Syntax Error: expression cannot start with a binary operator")
		} else if curr.Type == lexer.RightBrace {
			p.setError("Syntax Error: expression cannot start with a closing bracket")
		}
		return
	}

	prev := p.previousToken

	switch curr.Type {
	case lexer.Operator, lexer.RightBrace:
		if prev.Type == lexer.Operator ||
			prev.Type == lexer.LeftBrace ||
			prev.Type == lexer.Negation ||
			prev.Type == lexer.Function {

			p.setError(fmt.Sprintf("Syntax Error: operator '%s' cannot follow '%s'", curr.Value, prev.Value))
		}
	case lexer.Number, lexer.Function, lexer.LeftBrace, lexer.Negation:
		if prev.Type == lexer.Number || prev.Type == lexer.RightBrace {
			p.setError(fmt.Sprintf("Syntax Error: %s value '%v' should not follow %s", curr.Type, curr.Value, prev.Type))
		}
	}
}

func (p *Parser) validateEnd() {
	last := p.previousToken
	if last.Type == lexer.Operator ||
		last.Type == lexer.Function ||
		last.Type == lexer.Negation ||
		last.Type == lexer.LeftBrace {
		p.setError(fmt.Sprintf("Syntax Error: expression ends prematurely after '%s'", last.Value))
	}
}

func (p *Parser) clearRemainingOperators() {
	for len(p.operatorStack) > 0 {
		operator := p.popOperator()

		if operator.Type == lexer.LeftBrace {
			p.setError("Mismatched Brackets: opening '(' was never closed")
			return
		} else if operator.Type == lexer.RightBrace {
			p.setError("Syntax Error: unexpected closing ')'")
			return
		}

		p.enqueue(operator)
	}
}
