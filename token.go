package lex

import "fmt"

// Token is a parsed token
type Token struct {
	Type  TokenType
	Value string

	// The start possition of the token
	Start int

	// The width of the token
	Width int

	// The line in the file
	Line int
}

func (t Token) String() string {
	return fmt.Sprintf("[TokenType(%v) %v:%v+%v] <%v>", t.Type, t.Line, t.Start, t.Width, t.Value)
}
