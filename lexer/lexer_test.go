package lexer

import "testing"

func TestCategorise(t *testing.T) {
	tests := []struct {
		tokenValue string
		expectedCategory TokenType
	} {
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
		previousTokens []Token
		tokenValue string
		expectedCategory TokenType
	} {
		{[]Token{}, "-", Negation,},
		{[]Token{{Type: Operator},}, "-", Negation,},
		{[]Token{{Type: LeftBrace},}, "-", Negation,}, {[]Token{{Type: Negation},}, "-", Negation,},
		{[]Token{{Type: Number},}, "-", Operator,},
		{[]Token{{Type: RightBrace},}, "-", Operator,},
	}

	for _, tt := range tests {
		l := Lexer {
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
		input string
		pos int
		expectedtoken string
	} {
		{"1 + 1", 0, "1",},
		{"1 + 1", 1, "+",},
	}

	for _, tt := range tests {
		l := Lexer {
			input: []rune(tt.input),
			pos: tt.pos,
		}
		actualToken := l.ExtractToken()

		if actualToken != tt.expectedtoken {
			t.Errorf("token starting at pos: %v in \"%s\" incorrectly extracted as \"%s\", correct extraction \"%s\"", tt.pos, tt.input, actualToken, tt.expectedtoken)
		}
	}
}
