package lex

import (
	"fmt"
	"testing"
)

func TestStringer(t *testing.T) {
	tok := Token{
		Type:  TEOF,
		Value: "---",
		Start: 5,
		Width: 6,
		Line:  1,
	}

	expected := fmt.Sprintf("[TokenType(TEOF) 1:5+6] <--->")
	got := tok.String()
	if got != expected {
		t.Errorf("Incurrect stringer:\ngot: %v\nexpected: %v", got, expected)
	}
}

func TestCustomStringer(t *testing.T) {
	const TTest = iota + TokenTypes
	UpdateTokenTypes([]string{"TTest"})
	defer RestoreTokenTypes()

	tok := Token{
		Type:  TTest,
		Value: "---",
		Start: 5,
		Width: 6,
		Line:  1,
	}

	expected := fmt.Sprintf("[TokenType(TTest) 1:5+6] <--->")
	got := tok.String()
	if got != expected {
		t.Errorf("Incurrect stringer:\ngot: %v\nexpected: %v", got, expected)
	}
}
