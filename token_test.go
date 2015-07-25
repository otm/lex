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

	expected := fmt.Sprintf("[TokenType(%v) 1:5+6] <--->", TEOF)
	got := tok.String()
	if got != expected {
		t.Errorf("Incurrect stringer:\ngot: %v\nexpected: %v", got, expected)
	}
}
