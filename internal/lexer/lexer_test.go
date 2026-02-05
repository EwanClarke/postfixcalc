package lexer

import "testing"

func TestCategorise(t *testing.T) {
	tests := []struct {
		tokenValue       string
		expectedCategory TokenType
	}{
		{"0", Number},
		{"0.1", Number},
		{"-1", Error},
		{"+", Operator},
		{"*", Operator},
		{"/", Operator},
		{"^", Operator},
		{"(", LeftBrace},
		{")", RightBrace},
		{" +", Error},
		{" ", Error},
	}

	l := Lexer{}
	for _, tt := range tests {
		actualCategory := l.Categorise(tt.tokenValue)

		if actualCategory != tt.expectedCategory {
			t.Errorf("token \"%s\" incorrectly categorised as %s, correct categorisation: %s", tt.tokenValue, actualCategory, tt.expectedCategory)
		}
	}
}

func TestContexualCategorisation(t *testing.T) {
	tests := []struct {
		previousTokens   []Token
		tokenValue       string
		expectedCategory TokenType
	}{
		{[]Token{}, "-", Negation},
		{[]Token{{Type: Operator}}, "-", Negation},
		{[]Token{{Type: LeftBrace}}, "-", Negation}, {[]Token{{Type: Negation}}, "-", Negation},
		{[]Token{{Type: Number}}, "-", Operator},
		{[]Token{{Type: RightBrace}}, "-", Operator},
	}

	for _, tt := range tests {
		l := Lexer{
			tokens: tt.previousTokens,
		}
		actualCategory := l.Categorise(tt.tokenValue)

		if actualCategory != tt.expectedCategory {
			t.Errorf("token \"%s\" incorrectly categorised as %s, correct categorisation: %s", tt.tokenValue, actualCategory, tt.expectedCategory)
		}
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		input         string
		pos           int
		expectedtoken string
	}{
		{"1 + 1", 0, "1"},
		{"1 + 1", 1, "+"},
		{"1 + 1", 2, "+"},
		{"3.14", 0, "3.14"},
		{"sin", 0, "sin"},
		{"()", 0, "("},
		{"x", 0, "x"},
	}

	for _, tt := range tests {
		l := Lexer{
			input: []rune(tt.input),
			pos:   tt.pos,
		}
		actualToken := l.ExtractToken()

		if actualToken != tt.expectedtoken {
			t.Errorf("token starting at pos: %v in \"%s\" incorrectly extracted as \"%s\", correct extraction \"%s\"", tt.pos, tt.input, actualToken, tt.expectedtoken)
		}
	}
}

func TestNew(t *testing.T) {
	input := "1 + 2 * 3"
	l := New(input)

	if string(l.input) != input {
		t.Errorf("expected input %s, got %s", input, string(l.input))
	}

	if l.pos != 0 {
		t.Errorf("expected position 0, got %d", l.pos)
	}

	if len(l.tokens) != 0 {
		t.Errorf("expected empty tokens slice, got %v", l.tokens)
	}
}

func TestTokenise(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
		expectedLen int
	}{
		{"1 + 2", false, 3},
		{"sin 90", false, 2},
		{"-3 + 5", false, 4},
		{"x + 2", true, 0},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tokens, err := l.Tokenise()

		if tt.expectError && err == nil {
			t.Errorf("expected error for input '%s', got none", tt.input)
		} else if !tt.expectError && err != nil {
			t.Errorf("unexpected error for input '%s': %v", tt.input, err)
		} else if !tt.expectError && len(tokens) != tt.expectedLen {
			t.Errorf("expected %d tokens for input '%s', got %d", tt.expectedLen, tt.input, len(tokens))
		}
	}
}

func TestIsNegation(t *testing.T) {
	tests := []struct {
		tokenValue     string
		previousTokens []Token
		expected       bool
	}{
		{"+", []Token{}, false},
		{"-", []Token{}, true},
		{"-", []Token{{Type: Operator}}, true},
		{"-", []Token{{Type: Number}}, false},
		{"-", []Token{{Type: LeftBrace}}, true},
		{"-", []Token{{Type: RightBrace}}, false},
		{"*", []Token{}, false},
	}

	for _, tt := range tests {
		l := Lexer{
			tokens: tt.previousTokens,
		}
		actual := l.isNegation(tt.tokenValue)

		if actual != tt.expected {
			t.Errorf("isNegation(\"%s\") with previous tokens %v = %v, expected %v", tt.tokenValue, tt.previousTokens, actual, tt.expected)
		}
	}
}
