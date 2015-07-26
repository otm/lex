package lex_test

import (
	"fmt"
	"testing"

	"github.com/otm/lex"
)

func TestParse(t *testing.T) {
	s := `
--comment

--comment2
function test

--comment3

--[[1
Multi2
3]]
var test  = 32

--comment4
test = 1
`

	expected := []struct {
		comment string
		line    int
	}{
		{comment: "comment", line: 2},
		{comment: "comment2", line: 4},
		{comment: "comment3", line: 7},
		{comment: `1
Multi2
3`, line: 11},
		{comment: "comment4", line: 14},
	}

	comments := Parse(s)

	if len(comments) != len(expected) {
		t.Errorf("Not asserting all. Comments: %v, expected: %v", len(comments), len(expected))
		return
	}

	for i, comment := range comments {
		if comment.Value != expected[i].comment {
			t.Errorf("Comment value incorrect:\ngot: %v\nexpexted: %v", expected[i].comment, comment.Value)
		}

		if comment.Row != expected[i].line {
			t.Errorf("Comment line incorrect:\ngot: %v\nexpexted: %v", expected[i].line, comment.Row)
		}

	}
}

func TestStringer(t *testing.T) {
	tok := lex.Token{
		Type:  TSLComment,
		Value: "---",
		Start: 5,
		Width: 6,
		Line:  1,
	}

	expected := "[TokenType(4) 1:5+6] <--->"
	got := tok.String()
	if got != expected {
		t.Errorf("Incurrect stringer:\ngot: %v\nexpected: %v", got, expected)
	}
}

const (
	MinusSign           = "--"
	LeftBracket  string = "[["
	RightBracket string = "]]"
)

const (
	TSLComment lex.TokenType = iota + lex.TokenTypes
	TRMLComment
	TLMLComment
	TComment
)

// Comment contains information about a comment
type Comment struct {
	Value string
	Row   int
}

// ErrNoMatch is returned when no matching comment can be found for a row
var ErrNoMatch = fmt.Errorf("No matching comment")

// Comments that has been parsed
type Comments []Comment

// Get returns a Comment or an error
func (c Comments) Get(row int) (*Comment, error) {
	result := &Comment{Row: row}
	for i := len(c) - 1; i >= 0; i-- {
		if c[i].Row != row {
			continue
		}
		result.Value = result.Value + c[i].Value
		row--
	}

	if result.Value == "" {
		return nil, ErrNoMatch
	}

	return result, nil
}

// Parse parses a string
func Parse(str string) Comments {
	var token lex.Token

	comments := make(Comments, 0)

	l := lex.NewLexer("", str, LexBegin)

	for {
		token = l.NextToken()

		if token.Type == TComment {
			comment := Comment{Value: token.Value, Row: token.Line}
			comments = append(comments, comment)
		}

		if token.Type == lex.TEOF {
			return comments
		}

	}
}
