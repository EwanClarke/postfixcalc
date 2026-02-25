package engine

import (
	"testing"
)

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
	e := NewEvaluator()

	if len(e.resultStack) != 0 {
		t.Errorf("expected empty result stack, got %v", e.resultStack)
	}
}

func TestEngineCalculate(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		expectErr bool
	}{
		{"1+2", 3.0, false},
		{"3*4", 12.0, false},
		{"10/2", 5.0, false},
		{"2^3", 8.0, false},
		{"sin(0)", 0.0, false},
		{"-5+3", -2.0, false},
		{"", 0.0, true},
		{"x + 1", 0.0, true},
	}

	for _, tt := range tests {
		e := NewEngine()
		result, err := e.Calculate(tt.input)

		if tt.expectErr {
			if err == nil {
				t.Errorf("expected error for input '%s', got none", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input '%s': %v", tt.input, err)
			}

			if result != tt.expected {
				t.Errorf("expected result %f for input '%s', got %f", tt.expected, tt.input, result)
			}
		}
	}
}

func TestEvaluatorPushValue(t *testing.T) {
	e := NewEvaluator()
	e.pushValue(42.0)

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after push, got %d", len(e.resultStack))
	}

	if e.resultStack[0] != 42.0 {
		t.Errorf("expected value 42.0, got %f", e.resultStack[0])
	}
}

func TestEvaluatorPopValue(t *testing.T) {
	e := NewEvaluator()
	e.pushValue(10.0)
	e.pushValue(20.0)

	val, err := e.popValue()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != 20.0 {
		t.Errorf("expected popped value 20.0, got %f", val)
	}

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after pop, got %d", len(e.resultStack))
	}

	if e.resultStack[0] != 10.0 {
		t.Errorf("expected remaining value 10.0, got %f", e.resultStack[0])
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
		a           float64
		b           float64
		expected    float64
		expectError bool
	}{
		{"+", 2.0, 3.0, 5.0, false},
		{"-", 5.0, 2.0, 3.0, false},
		{"*", 4.0, 6.0, 24.0, false},
		{"/", 10.0, 2.0, 5.0, false},
		{"/", 10.0, 0.0, 0.0, true},
		{"^", 2.0, 3.0, 8.0, false},
		{"^", 0.0, -1.0, 0.0, true},
		{"%", 2.0, 3.0, 0.0, true},
	}

	for _, tt := range tests {
		e := NewEvaluator()
		result, err := e.applyBinaryOp(tt.operation, tt.a, tt.b)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for operation %s %f %f, got none", tt.operation, tt.a, tt.b)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for operation %s %f %f: %v", tt.operation, tt.a, tt.b, err)
			}

			if result != tt.expected {
				t.Errorf("expected result %f for operation %s %f %f, got %f", tt.expected, tt.operation, tt.a, tt.b, result)
			}
		}
	}
}

func TestEvaluatorApplyUnaryOp(t *testing.T) {
	tests := []struct {
		operation   string
		a           float64
		expected    float64
		expectError bool
	}{
		{"-u", 5.0, -5.0, false},
		{"sin", 0.0, 0.0, false},
		{"cos", 0.0, 1.0, false},
		{"tan", 0.0, 0.0, false},
		{"log", 5.0, 0.0, true},
	}

	for _, tt := range tests {
		e := NewEvaluator()
		result, err := e.applyUnaryOp(tt.operation, tt.a)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for operation %s %f, got none", tt.operation, tt.a)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for operation %s %f: %v", tt.operation, tt.a, err)
			}

			if result != tt.expected {
				t.Errorf("expected result %f for operation %s %f, got %f", tt.expected, tt.operation, tt.a, result)
			}
		}
	}
}

func TestEvaluatorEvaluateBinary(t *testing.T) {
	e := NewEvaluator()
	e.pushValue(2.0)
	e.pushValue(3.0)

	token := Token{Type: Operator, Value: "+"}
	err := e.evaluateBinary(token)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after binary evaluation, got %d", len(e.resultStack))
	}

	if e.resultStack[0] != 5.0 {
		t.Errorf("expected result 5.0, got %f", e.resultStack[0])
	}

	e2 := NewEvaluator()
	err = e2.evaluateBinary(token)
	if err == nil {
		t.Errorf("expected error for insufficient operands, got none")
	}
}

func TestEvaluatorEvaluateUnary(t *testing.T) {
	e := NewEvaluator()
	e.pushValue(5.0)

	token := Token{Type: Negation, Value: "-u"}
	err := e.evaluateUnary(token)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after unary evaluation, got %d", len(e.resultStack))
	}

	if e.resultStack[0] != -5.0 {
		t.Errorf("expected result -5.0, got %f", e.resultStack[0])
	}

	e2 := NewEvaluator()
	err = e2.evaluateUnary(token)
	if err == nil {
		t.Errorf("expected error for insufficient operands, got none")
	}
}

func TestEvaluatorValidateResult(t *testing.T) {
	tests := []struct {
		name        string
		stackValues []float64
		expectError bool
	}{
		{"single result", []float64{42.0}, false},
		{"empty stack", []float64{}, true},
		{"multiple results", []float64{1.0, 2.0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvaluator()
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
	e := NewEvaluator()
	tokens := []Token{
		{Type: Number, Value: "3"},
		{Type: Number, Value: "4"},
		{Type: Operator, Value: "+"},
	}

	result, err := e.Evaluate(tokens, make(map[string]float64))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 7.0 {
		t.Errorf("expected result 7.0, got %f", result)
	}
}

func TestEvaluatorEvaluateComplex(t *testing.T) {
	e := NewEvaluator()
	tokens := []Token{
		{Type: Number, Value: "2"},
		{Type: Number, Value: "3"},
		{Type: Number, Value: "4"},
		{Type: Operator, Value: "*"},
		{Type: Operator, Value: "+"},
	}

	result, err := e.Evaluate(tokens, make(map[string]float64))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 14.0 {
		t.Errorf("expected result 14.0, got %f", result)
	}
}

func TestEvaluatorEvaluateUnaryFunction(t *testing.T) {
	e := NewEvaluator()
	tokens := []Token{
		{Type: Number, Value: "0"},
		{Type: Function, Value: "sin"},
	}

	result, err := e.Evaluate(tokens, make(map[string]float64))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 0.0 {
		t.Errorf("expected result 0.0, got %f", result)
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
		vars        map[string]float64
		expected    float64
		expectError bool
	}{
		{"defined variable", map[string]float64{"x": 3.0}, 5.0, false},
		{"undefined variable", map[string]float64{"y": 3.0}, 0.0, true},
		{"empty variables", map[string]float64{}, 0.0, true},
	}

	for _, tt := range tests {
		e := NewEvaluator()
		result, err := e.Evaluate(tokens, tt.vars)

		if tt.expectError {
			if err == nil {
				t.Errorf("expected error for test '%s', got none", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for test '%s': %v", tt.name, err)
			}

			if result != tt.expected {
				t.Errorf("expected result %f for test '%s', got %f", tt.expected, tt.name, result)
			}
		}
	}
}
