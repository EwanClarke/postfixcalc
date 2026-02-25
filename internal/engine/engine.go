package engine

import (
	"errors"
	"maps"
)

type Engine struct {
	lexer        *Lexer
	parser       *Parser
	evaluator    *Evaluator
	internalVars map[string]float64
}

func NewEngine() *Engine {
	return &Engine{
		lexer:        NewLexer(),
		parser:       NewParser(),
		evaluator:    NewEvaluator(),
		internalVars: InternalVars,
	}
}

func (e *Engine) Calculate(expression string) (float64, error) {
	if expression == "" {
		return 0, errors.New("empty expression")
	}

	// Tokenize
	tokens, err := e.lexer.Tokenise(expression)
	if err != nil {
		return 0, err
	}

	// Parse to postfix
	e.parser = NewParser()
	postfixTokens, err := e.parser.Convert(tokens)
	if err != nil {
		return 0, err
	}

	// Evaluate
	e.evaluator = NewEvaluator()
	result, err := e.evaluator.Evaluate(postfixTokens, e.internalVars)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (e *Engine) Derive(postfixTokens []Token, xVal float64, hVal float64) (float64, error) {
	if hVal == 0 {
		return 0, errors.New("Math Error: division by zero")
	}
	vars := maps.Clone(e.internalVars)

	vars["x"] = xVal
	e.evaluator = NewEvaluator()
	xResult, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, err
	}

	vars["x"] = xVal + hVal
	hResult, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, err
	}

	return (hResult - xResult) / hVal, nil
}

func (e *Engine) Parse(expression string) ([]Token, error) {
	tokens, err := e.lexer.Tokenise(expression)
	if err != nil {
		return nil, err
	}

	e.parser = NewParser()
	postfixTokens, err := e.parser.Convert(tokens)
	if err != nil {
		return nil, err
	}

	return postfixTokens, nil
}

func (e *Engine) EvaluateAt(postfixTokens []Token, x float64) (float64, error) {
	vars := maps.Clone(e.internalVars)
	vars["x"] = x
	e.evaluator = NewEvaluator()
	return e.evaluator.Evaluate(postfixTokens, vars)
}

func (e *Engine) EvaluateWithDerivative(postfixTokens []Token, x, h float64) (y, dy float64, err error) {
	vars := maps.Clone(e.internalVars)
	vars["x"] = x
	e.evaluator = NewEvaluator()
	y, err = e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, 0, err
	}

	vars["x"] = x + h
	hResult, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, 0, err
	}

	dy = (hResult - y) / h
	return y, dy, nil
}
