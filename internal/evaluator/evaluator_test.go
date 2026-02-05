package evaluator

import (
	"github.com/EwanClarke/postfixcalc/internal/lexer"
	"testing"
)

func TestNew(t *testing.T) {
	e := New()

	if len(e.resultStack) != 0 {
		t.Errorf("expected empty result stack, got %v", e.resultStack)
	}
}

func TestPushValue(t *testing.T) {
	e := New()
	e.pushValue(42.0)

	if len(e.resultStack) != 1 {
		t.Errorf("expected stack length 1 after push, got %d", len(e.resultStack))
	}

	if e.resultStack[0] != 42.0 {
		t.Errorf("expected value 42.0, got %f", e.resultStack[0])
	}
}

func TestPopValue(t *testing.T) {
	e := New()
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

func TestApplyBinaryOp(t *testing.T) {
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
		e := New()
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

func TestApplyUnaryOp(t *testing.T) {
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
		e := New()
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

func TestEvaluateBinary(t *testing.T) {
	e := New()
	e.pushValue(2.0)
	e.pushValue(3.0)

	token := lexer.Token{Type: lexer.Operator, Value: "+"}
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

	e2 := New()
	err = e2.evaluateBinary(token)
	if err == nil {
		t.Errorf("expected error for insufficient operands, got none")
	}
}

func TestEvaluateUnary(t *testing.T) {
	e := New()
	e.pushValue(5.0)

	token := lexer.Token{Type: lexer.Negation, Value: "-u"}
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

	e2 := New()
	err = e2.evaluateUnary(token)
	if err == nil {
		t.Errorf("expected error for insufficient operands, got none")
	}
}

func TestValidateResult(t *testing.T) {
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
			e := New()
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

func TestEvaluateSimple(t *testing.T) {
	e := New()
	tokens := []lexer.Token{
		{Type: lexer.Number, Value: "3"},
		{Type: lexer.Number, Value: "4"},
		{Type: lexer.Operator, Value: "+"},
	}

	result, err := e.Evaluate(tokens)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 7.0 {
		t.Errorf("expected result 7.0, got %f", result)
	}
}

func TestEvaluateComplex(t *testing.T) {
	e := New()
	tokens := []lexer.Token{
		{Type: lexer.Number, Value: "2"},
		{Type: lexer.Number, Value: "3"},
		{Type: lexer.Number, Value: "4"},
		{Type: lexer.Operator, Value: "*"},
		{Type: lexer.Operator, Value: "+"},
	}

	result, err := e.Evaluate(tokens)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 14.0 {
		t.Errorf("expected result 14.0, got %f", result)
	}
}

func TestEvaluateUnaryFunction(t *testing.T) {
	e := New()
	tokens := []lexer.Token{
		{Type: lexer.Number, Value: "0"},
		{Type: lexer.Function, Value: "sin"},
	}

	result, err := e.Evaluate(tokens)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 0.0 {
		t.Errorf("expected result 0.0, got %f", result)
	}
}
