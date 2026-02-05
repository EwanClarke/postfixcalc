package parser

import (
	"github.com/EwanClarke/postfixcalc/internal/lexer"
	"testing"
)

func TestPushOperator(t *testing.T) {
	p := New()
	token := lexer.Token{Type: lexer.Operator, Value: "+"}

	p.pushOperator(token)

	if len(p.operatorStack) != 1 {
		t.Errorf("expected operator stack length 1, got %d", len(p.operatorStack))
	}

	if p.operatorStack[0].Value != "+" {
		t.Errorf("expected operator '+', got %s", p.operatorStack[0].Value)
	}
}

func TestPopOperator(t *testing.T) {
	p := New()
	token := lexer.Token{Type: lexer.Operator, Value: "*"}
	p.pushOperator(token)

	popped := p.popOperator()

	if popped.Value != "*" {
		t.Errorf("expected popped operator '*', got %s", popped.Value)
	}

	if len(p.operatorStack) != 0 {
		t.Errorf("expected empty operator stack after pop, got %d items", len(p.operatorStack))
	}

	emptyPop := p.popOperator()
	if emptyPop.Type != lexer.Error {
		t.Errorf("expected error token when popping from empty stack, got %v", emptyPop)
	}
}

func TestHasOperators(t *testing.T) {
	p := New()

	if p.hasOperators() {
		t.Errorf("expected hasOperators() to be false for empty stack")
	}

	p.pushOperator(lexer.Token{Type: lexer.Operator, Value: "+"})

	if !p.hasOperators() {
		t.Errorf("expected hasOperators() to be true after pushing operator")
	}
}

func TestPeekOperator(t *testing.T) {
	p := New()

	emptyPeek := p.peekOperator()
	if emptyPeek.Type != lexer.Error {
		t.Errorf("expected error token when peeking empty stack, got %v", emptyPeek)
	}

	p.pushOperator(lexer.Token{Type: lexer.Operator, Value: "+"})
	p.pushOperator(lexer.Token{Type: lexer.Operator, Value: "*"})

	peeked := p.peekOperator()
	if peeked.Value != "*" {
		t.Errorf("expected peeked operator '*', got %s", peeked.Value)
	}

	if len(p.operatorStack) != 2 {
		t.Errorf("expected stack length 2 after peek, got %d", len(p.operatorStack))
	}
}

func TestEnqueueDequeue(t *testing.T) {
	p := New()
	token := lexer.Token{Type: lexer.Number, Value: "42"}

	p.enqueue(token)

	if len(p.outputQueue) != 1 {
		t.Errorf("expected output queue length 1 after enqueue, got %d", len(p.outputQueue))
	}

	dequeued := p.dequeue()
	if dequeued.Value != "42" {
		t.Errorf("expected dequeued token '42', got %s", dequeued.Value)
	}

	if len(p.outputQueue) != 0 {
		t.Errorf("expected empty queue after dequeue, got %d items", len(p.outputQueue))
	}

	emptyDequeue := p.dequeue()
	if emptyDequeue.Type != lexer.Error {
		t.Errorf("expected error token when dequeuing from empty queue, got %v", emptyDequeue)
	}
}

func TestSetError(t *testing.T) {
	p := New()

	if p.err != nil {
		t.Errorf("expected initial error to be nil, got %v", p.err)
	}

	p.setError("test error")

	if p.err == nil {
		t.Errorf("expected error to be set, got nil")
	}

	if p.err.Error() != "test error" {
		t.Errorf("expected error message 'test error', got %s", p.err.Error())
	}
}

func TestHandleRightBrace(t *testing.T) {
	p := New()

	p.pushOperator(lexer.Token{Type: lexer.Operator, Value: "+"})
	p.pushOperator(lexer.Token{Type: lexer.LeftBrace, Value: "("})
	p.pushOperator(lexer.Token{Type: lexer.Operator, Value: "*"})

	p.handleRightBrace()

	if len(p.operatorStack) != 1 {
		t.Errorf("expected operator stack length 1 after handleRightBrace, got %d", len(p.operatorStack))
	}

	if p.operatorStack[0].Value != "+" {
		t.Errorf("expected '+' at bottom of stack, got %v", p.operatorStack[0])
	}

	if len(p.outputQueue) != 1 {
		t.Errorf("expected output queue length 1, got %d", len(p.outputQueue))
	}

	if p.outputQueue[0].Value != "*" {
		t.Errorf("expected '*' in output queue, got %s", p.outputQueue[0].Value)
	}
}

func TestValidateSequence(t *testing.T) {
	tests := []struct {
		name          string
		previousToken lexer.Token
		hasPrev       bool
		currentToken  lexer.Token
		expectError   bool
	}{
		{"starts with operator", lexer.Token{}, false, lexer.Token{Type: lexer.Operator, Value: "+"}, true},
		{"starts with right brace", lexer.Token{}, false, lexer.Token{Type: lexer.RightBrace, Value: ")"}, true},
		{"operator after operator", lexer.Token{Type: lexer.Operator, Value: "+"}, true, lexer.Token{Type: lexer.Operator, Value: "*"}, true},
		{"number after operator", lexer.Token{Type: lexer.Operator, Value: "+"}, true, lexer.Token{Type: lexer.Number, Value: "5"}, false},
		{"number after number", lexer.Token{Type: lexer.Number, Value: "3"}, true, lexer.Token{Type: lexer.Number, Value: "5"}, true},
		{"right brace after number", lexer.Token{Type: lexer.Number, Value: "5"}, true, lexer.Token{Type: lexer.RightBrace, Value: ")"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{
				previousToken: tt.previousToken,
				hasPrev:       tt.hasPrev,
			}

			p.validateSequence(tt.currentToken)

			if tt.expectError && p.err == nil {
				t.Errorf("expected error for case '%s', got none", tt.name)
			} else if !tt.expectError && p.err != nil {
				t.Errorf("unexpected error for case '%s': %v", tt.name, p.err)
			}
		})
	}
}

func TestValidateEnd(t *testing.T) {
	tests := []struct {
		name        string
		lastToken   lexer.Token
		expectError bool
	}{
		{"ends with operator", lexer.Token{Type: lexer.Operator, Value: "+"}, true},
		{"ends with function", lexer.Token{Type: lexer.Function, Value: "sin"}, true},
		{"ends with number", lexer.Token{Type: lexer.Number, Value: "42"}, false},
		{"ends with right brace", lexer.Token{Type: lexer.RightBrace, Value: ")"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{
				previousToken: tt.lastToken,
			}

			p.validateEnd()

			if tt.expectError && p.err == nil {
				t.Errorf("expected error for case '%s', got none", tt.name)
			} else if !tt.expectError && p.err != nil {
				t.Errorf("unexpected error for case '%s': %v", tt.name, p.err)
			}
		})
	}
}

func TestHandleImplicitMult(t *testing.T) {
	tests := []struct {
		name                 string
		previousToken        lexer.Token
		hasPrev              bool
		currentToken         lexer.Token
		expectedMultInserted bool
	}{
		{"number followed by left brace", lexer.Token{Type: lexer.Number, Value: "2"}, true, lexer.Token{Type: lexer.LeftBrace, Value: "("}, true},
		{"right brace followed by number", lexer.Token{Type: lexer.RightBrace, Value: ")"}, true, lexer.Token{Type: lexer.Number, Value: "3"}, true},
		{"number followed by function", lexer.Token{Type: lexer.Number, Value: "2"}, true, lexer.Token{Type: lexer.Function, Value: "sin"}, true},
		{"operator followed by number", lexer.Token{Type: lexer.Operator, Value: "+"}, true, lexer.Token{Type: lexer.Number, Value: "3"}, false},
		{"no previous token", lexer.Token{}, false, lexer.Token{Type: lexer.Number, Value: "3"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{
				previousToken: tt.previousToken,
				hasPrev:       tt.hasPrev,
			}

			initialStackLen := len(p.operatorStack)
			p.handleImplicitMult(tt.currentToken)

			multInserted := len(p.operatorStack) > initialStackLen
			if multInserted != tt.expectedMultInserted {
				t.Errorf("expected mult insertion %v for case '%s', got %v", tt.expectedMultInserted, tt.name, multInserted)
			}
		})
	}
}

func TestNew(t *testing.T) {
	p := New()

	if len(p.input) != 0 {
		t.Errorf("expected empty input slice, got %v", p.input)
	}

	if len(p.operatorStack) != 0 {
		t.Errorf("expected empty operator stack, got %v", p.operatorStack)
	}

	if len(p.outputQueue) != 0 {
		t.Errorf("expected empty output queue, got %v", p.outputQueue)
	}

	if p.err != nil {
		t.Errorf("expected nil error, got %v", p.err)
	}
}
