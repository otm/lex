package lex

import "strconv"

// Constants for common charecters
const (
	EOF     rune = 0 // EOF is given a rune for practical reasons
	NewLine      = '\n'
)

// TokenType defines a token
type TokenType int

func (t TokenType) String() string {
	if int(t) >= len(tokenMap) {
		return strconv.Itoa(int(t))
	}
	return tokenMap[t]
}

// UpdateTokenTypes adds token translations for the stringer function. It serves
// no other purpose then to make pretty output when printing tokens
func UpdateTokenTypes(tokens []string) {
	tokenMap = append(tokenMap, tokens...)
}

// RestoreTokenTypes replaces the current internal mapping between TokenType and
// string to the original state.
func RestoreTokenTypes() {
	tm := make([]string, TokenTypes)
	for i := 0; i < TokenTypes; i++ {
		tm[i] = tokenMap[i]
	}
	tokenMap = tm
}

// TokenTypes defines the number of TokenTypes that are defined in the package
const TokenTypes = 4

// Predefined tokens
const (
	TError TokenType = iota
	TIgnore
	TEOF
	TNewline
)

var tokenMap = []string{
	"TError",
	"TIgnore",
	"TEOF",
	"TNewline",
}
