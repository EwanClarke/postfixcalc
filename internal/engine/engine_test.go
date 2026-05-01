package engine

import (
	"math/big"
	"testing"
)

// ratFromString is a test helper to build a *big.Rat from a string like "3/2" or "7".
func ratFromString(s string) *big.Rat {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		panic("ratFromString: invalid rat string: " + s)
	}
	return r
}

func TestNewEngine(t *testing.T) {
	e := NewEngine()

	if e.lexer == nil || e.parser == nil || e.evaluator == nil {
		t.Errorf("expected engine components to be initialized")
	}
}

func TestNewLexer(t *testing.T) {
	l := NewLexer()

	if l.input != nil {
		t.Errorf("expected lexer input to be nil initially, got %v", l.input)
	}

	// Test that Tokenise sets the input
	_, err := l.Tokenise("2+3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if string(l.input) != "2+3" {
		t.Errorf("expected input '2+3' after tokenisation, got '%s'", string(l.input))
	}
}

func TestNewParser(t *testing.T) {
	p := NewParser()

	if p.operatorStack != nil && len(p.operatorStack) != 0 {
		t.Errorf("expected empty operator stack, got %v", p.operatorStack)
	}

	if p.outputQueue != nil && len(p.outputQueue) != 0 {
		t.Errorf("expected empty output queue, got %v", p.outputQueue)
	}
}

func TestNewEvaluator(t *testing.T) {
	e := NewEvaluator(Radians)

	if len(e.resultStack) != 0 {
		t.Errorf("expected empty result stack, got %v", e.resultStack)
	}
}

func TestEngineCalculate(t *testing.T) {
	tests := []struct {
		input     string
		expected  *big.Rat
		exact     bool
		expectErr bool
	}{
		{"1+2", ratFromString("3"), true, false},
		{"3*4", ratFromString("12"), true, false},
		{"10/2", ratFromString("5"), true, false},
		{"1/3", ratFromString("1/3"), true, false},
		{"1/3+1/6", ratFromString("1/2"), true, false},
		// sin(0) = 0 but is inexact (float64 round-trip)
		{"sin(0)", ratFromString("0"), false, false},
		{"-5+3", ratFromString("-2"), true, false},
		{"", nil, false, true},
		{"x + 1", nil, false, true},
	}

	for _, tt := range tests {
		e := NewEngine()
		result, exact, err := e.Calculate(tt.input)

		if tt.expectErr {
			if err == nil {
				t.Errorf("expected error for input '%s', got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input '%s': %v", tt.input, err)
				continue
			}
			if exact != tt.exact {
				t.Errorf("expected exact=%v for input '%s', got %v", tt.exact, tt.input, exact)
			}
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("expected result %v for input '%s', got %v", tt.expected, tt.input, result)
			}
		}
	}
}

func TestEvaluatorPushValue(t *testing.T) {
	e := NewEvaluator(Radians)
	e.pushValue(ratFromString("42"))

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after push, got %d", len(e.resultStack))
	}

	if e.resultStack[0].Cmp(ratFromString("42")) != 0 {
		t.Errorf("expected value 42, got %v", e.resultStack[0])
	}
}

func TestEvaluatorPopValue(t *testing.T) {
	e := NewEvaluator(Radians)
	e.pushValue(ratFromString("10"))
	e.pushValue(ratFromString("20"))

	val, err := e.popValue()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val.Cmp(ratFromString("20")) != 0 {
		t.Errorf("expected popped value 20, got %v", val)
	}

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after pop, got %d", len(e.resultStack))
	}

	if e.resultStack[0].Cmp(ratFromString("10")) != 0 {
		t.Errorf("expected remaining value 10, got %v", e.resultStack[0])
	}

	_, err = e.popValue()
	if err != nil {
		t.Errorf("unexpected error popping second value: %v", err)
	}

	_, err = e.popValue()
	if err == nil {
		t.Errorf("expected error when popping from empty stack, got none")
	}

	if err.Error() != "Stack Underflow: attempted removal from empty stack" {
		t.Errorf("expected stack underflow error, got %v", err)
	}
}

func TestEvaluatorApplyBinaryOp(t *testing.T) {
	tests := []struct {
		operation   string
		a           *big.Rat
		b           *big.Rat
		expected    *big.Rat
		expectError bool
	}{
		{"+", ratFromString("2"), ratFromString("3"), ratFromString("5"), false},
		{"-", ratFromString("5"), ratFromString("2"), ratFromString("3"), false},
		{"*", ratFromString("4"), ratFromString("6"), ratFromString("24"), false},
		{"/", ratFromString("10"), ratFromString("2"), ratFromString("5"), false},
		{"/", ratFromString("10"), ratFromString("0"), nil, true},
		// Exact rational arithmetic: 1/3 + 1/6 = 1/2
		{"+", ratFromString("1/3"), ratFromString("1/6"), ratFromString("1/2"), false},
		// ^ error case
		{"^", ratFromString("0"), ratFromString("-1"), nil, true},
		{"%", ratFromString("2"), ratFromString("3"), nil, true},
	}

	for _, tt := range tests {
		e := NewEvaluator(Radians)
		result, err := e.applyBinaryOp(tt.operation, tt.a, tt.b)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for operation %s %v %v, got none", tt.operation, tt.a, tt.b)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for operation %s %v %v: %v", tt.operation, tt.a, tt.b, err)
				continue
			}
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("expected result %v for operation %s %v %v, got %v", tt.expected, tt.operation, tt.a, tt.b, result)
			}
		}
	}
}

func TestEvaluatorApplyUnaryOp(t *testing.T) {
	tests := []struct {
		operation   string
		a           *big.Rat
		expected    *big.Rat
		expectError bool
	}{
		{"-u", ratFromString("5"), ratFromString("-5"), false},
		{"sin", ratFromString("0"), ratFromString("0"), false},
		{"cos", ratFromString("0"), ratFromString("1"), false},
		{"tan", ratFromString("0"), ratFromString("0"), false},
		{"log", ratFromString("5"), nil, true},
	}

	for _, tt := range tests {
		e := NewEvaluator(Radians)
		result, err := e.applyUnaryOp(tt.operation, tt.a)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for operation %s %v, got none", tt.operation, tt.a)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for operation %s %v: %v", tt.operation, tt.a, err)
				continue
			}
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("expected result %v for operation %s %v, got %v", tt.expected, tt.operation, tt.a, result)
			}
		}
	}
}

func TestEvaluatorEvaluateBinary(t *testing.T) {
	e := NewEvaluator(Radians)
	e.pushValue(ratFromString("2"))
	e.pushValue(ratFromString("3"))

	token := Token{Type: Operator, Value: "+"}
	err := e.evaluateBinary(token)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after binary evaluation, got %d", len(e.resultStack))
	}

	if e.resultStack[0].Cmp(ratFromString("5")) != 0 {
		t.Errorf("expected result 5, got %v", e.resultStack[0])
	}

	e2 := NewEvaluator(Radians)
	err = e2.evaluateBinary(token)
	if err == nil {
		t.Errorf("expected error for insufficient operands, got none")
	}
}

func TestEvaluatorEvaluateUnary(t *testing.T) {
	e := NewEvaluator(Radians)
	e.pushValue(ratFromString("5"))

	token := Token{Type: Negation, Value: "-u"}
	err := e.evaluateUnary(token)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after unary evaluation, got %d", len(e.resultStack))
	}

	if e.resultStack[0].Cmp(ratFromString("-5")) != 0 {
		t.Errorf("expected result -5, got %v", e.resultStack[0])
	}

	e2 := NewEvaluator(Radians)
	err = e2.evaluateUnary(token)
	if err == nil {
		t.Errorf("expected error for insufficient operands, got none")
	}
}

func TestEvaluatorValidateResult(t *testing.T) {
	tests := []struct {
		name        string
		stackValues []*big.Rat
		expectError bool
	}{
		{"single result", []*big.Rat{ratFromString("42")}, false},
		{"empty stack", []*big.Rat{}, true},
		{"multiple results", []*big.Rat{ratFromString("1"), ratFromString("2")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvaluator(Radians)
			for _, val := range tt.stackValues {
				e.pushValue(val)
			}

			err := e.validateResult()

			if tt.expectError && err == nil {
				t.Errorf("expected error for case '%s', got none", tt.name)
			} else if !tt.expectError && err != nil {
				t.Errorf("unexpected error for case '%s': %v", tt.name, err)
			}
		})
	}
}

func TestEvaluatorEvaluateSimple(t *testing.T) {
	e := NewEvaluator(Radians)
	tokens := []Token{
		{Type: Number, Value: "3"},
		{Type: Number, Value: "4"},
		{Type: Operator, Value: "+"},
	}

	result, exact, err := e.Evaluate(tokens, make(map[string]*big.Rat))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !exact {
		t.Errorf("expected exact result for integer addition")
	}
	if result.Cmp(ratFromString("7")) != 0 {
		t.Errorf("expected result 7, got %v", result)
	}
}

func TestEvaluatorEvaluateComplex(t *testing.T) {
	e := NewEvaluator(Radians)
	tokens := []Token{
		{Type: Number, Value: "2"},
		{Type: Number, Value: "3"},
		{Type: Number, Value: "4"},
		{Type: Operator, Value: "*"},
		{Type: Operator, Value: "+"},
	}

	result, exact, err := e.Evaluate(tokens, make(map[string]*big.Rat))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !exact {
		t.Errorf("expected exact result for integer arithmetic")
	}
	if result.Cmp(ratFromString("14")) != 0 {
		t.Errorf("expected result 14, got %v", result)
	}
}

func TestEvaluatorEvaluateUnaryFunction(t *testing.T) {
	e := NewEvaluator(Radians)
	tokens := []Token{
		{Type: Number, Value: "0"},
		{Type: Function, Value: "sin"},
	}

	result, exact, err := e.Evaluate(tokens, make(map[string]*big.Rat))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// sin is a float64 round-trip — result should be inexact
	if exact {
		t.Errorf("expected inexact result for sin function")
	}
	if result.Cmp(ratFromString("0")) != 0 {
		t.Errorf("expected result 0, got %v", result)
	}
}

func TestEvaluatorEvaluateWithVariable(t *testing.T) {
	tokens := []Token{
		{Type: Number, Value: "2"},
		{Type: Variable, Value: "x"},
		{Type: Operator, Value: "+"},
	}

	tests := []struct {
		name        string
		vars        map[string]*big.Rat
		expected    *big.Rat
		expectError bool
	}{
		{"defined variable", map[string]*big.Rat{"x": ratFromString("3")}, ratFromString("5"), false},
		{"undefined variable", map[string]*big.Rat{"y": ratFromString("3")}, nil, true},
		{"empty variables", map[string]*big.Rat{}, nil, true},
	}

	for _, tt := range tests {
		e := NewEvaluator(Radians)
		result, _, err := e.Evaluate(tokens, tt.vars)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for test '%s', got none", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for test '%s': %v", tt.name, err)
				continue
			}
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("expected result %v for test '%s', got %v", tt.expected, tt.name, result)
			}
		}
	}
}

// ---- FormatMixed tests ----

func TestFormatMixed(t *testing.T) {
	tests := []struct {
		input    *big.Rat
		expected string
	}{
		{ratFromString("0"), "0"},
		{ratFromString("3"), "3"},
		{ratFromString("-3"), "-3"},
		{ratFromString("7/2"), "3 1/2"},
		{ratFromString("-7/2"), "-3 1/2"},
		{ratFromString("1/3"), "1/3"},
		{ratFromString("-1/3"), "-1/3"},
		{ratFromString("6/3"), "2"},     // reduces to whole number
		{ratFromString("1/2"), "1/2"},   // proper fraction
		{ratFromString("10/4"), "2 1/2"}, // reduces: 5/2 = 2 1/2
	}

	for _, tt := range tests {
		result := FormatMixed(tt.input)
		if result != tt.expected {
			t.Errorf("FormatMixed(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDecimal(t *testing.T) {
	tests := []struct {
		input    *big.Rat
		expected string
	}{
		{ratFromString("0"), "0"},
		{ratFromString("3"), "3"},
		{ratFromString("1/2"), "0.5"},
		{ratFromString("7/2"), "3.5"},
		{ratFromString("-7/2"), "-3.5"},
	}

	for _, tt := range tests {
		result := FormatDecimal(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDecimal(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
