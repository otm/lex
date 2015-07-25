package lex

// Constants for common charecters
const (
	EOF     rune = 0 // EOF is given a rune for practical reasons
	NewLine      = '\n'
)

// TokenType defines a token
type TokenType int

// TokenTypes defines the number of TokenTypes that are defined in the package
const TokenTypes = 3

// Predefined tokens
const (
	TError TokenType = iota
	TIgnore
	TEOF
	TNewline
)
