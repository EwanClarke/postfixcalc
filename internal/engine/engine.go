package engine

import (
	"errors"
	"maps"
	"math/big"
)

type AngleMode int

const (
	Degrees AngleMode = iota
	Radians
)

func (a AngleMode) String() string {
	return [...]string{"Deg", "Rad"}[a]
}

type Engine struct {
	lexer        *Lexer
	parser       *Parser
	evaluator    *Evaluator
	internalVars map[string]*big.Rat
	AngleMode    AngleMode
}

func NewEngine() *Engine {
	return &Engine{
		lexer:        NewLexer(),
		parser:       NewParser(),
		evaluator:    NewEvaluator(Degrees), // Default could be Degrees or Radians, but it gets overwritten when evaluating. We'll pass the mode explicitly during Evaluate.
		internalVars: float64MapToRat(InternalVars),
		AngleMode:    Radians, // default to Radians to match math package behaviour
	}
}

// float64MapToRat converts a map[string]float64 to map[string]*big.Rat.
// Constants like Pi and e will be represented as their float64 approximations,
// so any calculation using them will produce an inexact result.
func float64MapToRat(m map[string]float64) map[string]*big.Rat {
	out := make(map[string]*big.Rat, len(m))
	for k, v := range m {
		out[k] = new(big.Rat).SetFloat64(v)
	}
	return out
}

// ratToFloat64 converts a *big.Rat to float64 for use in graph/derivative calculations.
func ratToFloat64(r *big.Rat) float64 {
	f, _ := r.Float64()
	return f
}

// Calculate evaluates a mathematical expression string.
// Returns the result as a *big.Rat, a bool indicating whether the result is
// exact (true = pure rational arithmetic, false = float64 approximation was
// used e.g. for trig functions or non-integer powers), and any error.
func (e *Engine) Calculate(expression string) (*big.Rat, bool, error) {
	if expression == "" {
		return nil, false, errors.New("empty expression")
	}

	// Tokenize
	tokens, err := e.lexer.Tokenise(expression)
	if err != nil {
		return nil, false, err
	}

	// Parse to postfix
	e.parser = NewParser()
	postfixTokens, err := e.parser.Convert(tokens)
	if err != nil {
		return nil, false, err
	}

	// Evaluate
	e.evaluator = NewEvaluator(e.AngleMode)
	result, exact, err := e.evaluator.Evaluate(postfixTokens, e.internalVars)
	if err != nil {
		return nil, false, err
	}

	return result, exact, nil
}

func (e *Engine) Derive(postfixTokens []Token, xVal float64, hVal float64) (float64, error) {
	if hVal == 0 {
		return 0, errors.New("Math Error: division by zero")
	}
	vars := maps.Clone(e.internalVars)

	vars["x"] = new(big.Rat).SetFloat64(xVal)
	e.evaluator = NewEvaluator(e.AngleMode)
	xRat, _, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, err
	}

	vars["x"] = new(big.Rat).SetFloat64(xVal + hVal)
	e.evaluator = NewEvaluator(e.AngleMode)
	hRat, _, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, err
	}

	return (ratToFloat64(hRat) - ratToFloat64(xRat)) / hVal, nil
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
	vars["x"] = new(big.Rat).SetFloat64(x)
	e.evaluator = NewEvaluator(e.AngleMode)
	r, _, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, err
	}
	return ratToFloat64(r), nil
}

func (e *Engine) EvaluateWithDerivative(postfixTokens []Token, x, h float64) (y, dy float64, err error) {
	vars := maps.Clone(e.internalVars)
	vars["x"] = new(big.Rat).SetFloat64(x)
	e.evaluator = NewEvaluator(e.AngleMode)
	yRat, _, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, 0, err
	}
	y = ratToFloat64(yRat)

	vars["x"] = new(big.Rat).SetFloat64(x + h)
	e.evaluator = NewEvaluator(e.AngleMode)
	hRat, _, err := e.evaluator.Evaluate(postfixTokens, vars)
	if err != nil {
		return 0, 0, err
	}

	dy = (ratToFloat64(hRat) - y) / h
	return y, dy, nil
}
