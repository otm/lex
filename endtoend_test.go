package lex_test

import (
	"fmt"
	"testing"

	"github.com/otm/lex"
)

var str = `
key=abcdefg

[User]
userName=otm
`

const (
	tLeftBracket lex.TokenType = iota + lex.TokenTypes
	tRightBracket
	tEqualSign

	tSection
	tKey
	tValue
)

var tokenMap = []string{
	"tLeftBracket",
	"tRightBracket",
	"tEqualSign",
	"tSection",
	"tKey",
	"tValue",
}

const (
	leftBracket  = '['
	rightBracket = ']'
	equalSign    = '='
)

var key = map[rune]lex.TokenType{
	leftBracket:  tEqualSign,
	rightBracket: tLeftBracket,
	equalSign:    tRightBracket,
}

func lexStart(l *lex.Lexer) lex.StateFn {
	l.AcceptWhitespace(lex.Ignore)

	r := l.Next()
	switch r {
	case lex.EOF:
		l.Emit(lex.TEOF)
		return nil
	case leftBracket:
		l.Backup()
		return lexLeftBracket
	default:
		if lex.IsAlphaNumeric(r) {
			return lexKey
		}
		return l.Errorf("Unknown charecter: %v", string(r))
	}
}

func lexLeftBracket(l *lex.Lexer) lex.StateFn {
	l.Accept("[", lex.Emit(tLeftBracket))
	return lexSection
}

func lexRightBracket(l *lex.Lexer) lex.StateFn {
	l.Accept("]", lex.Emit(tRightBracket))
	return lexStart
}

func lexKey(l *lex.Lexer) lex.StateFn {
	for {
		switch l.Next() {
		case lex.EOF:
			return l.Errorf("Unexpected EOF")
		case equalSign:
			if l.Pos > l.Start {
				l.Backup()
				l.Emit(tKey)
				return lexEqualSign
			}
			return l.Errorf("Unexpected equal sign")
		}
	}
}

func lexSection(l *lex.Lexer) lex.StateFn {
	for {
		switch l.Next() {
		case lex.EOF:
			return l.Errorf("Unexpected EOF in section")
		case '\n':
			return l.Errorf("Section: Missing right bracket")
		case rightBracket:
			l.Backup()
			l.Emit(tSection)
			return lexRightBracket
		}
	}
}

func lexEqualSign(l *lex.Lexer) lex.StateFn {
	l.Accept("=", lex.Emit(tEqualSign))
	return lexValue
}

func lexValue(l *lex.Lexer) lex.StateFn {
	for {
		switch l.Next() {
		case lex.EOF:
			if l.Pos > l.Start {
				l.Backup()
				l.Emit(tValue)
				l.Emit(lex.TEOF)
				return nil
			}
			return l.Errorf("Unexpected equal sign")
		case '\n':
			l.Backup()
			l.Emit(tValue)
			l.Next(lex.Ignore)
			return lexStart
		}
	}
}

func TestParseIni(t *testing.T) {
	var token lex.Token
	expected := []struct {
		typ   lex.TokenType
		line  int
		start int
		width int
		value string
	}{
		{tKey, 2, 1, 3, "key"},
		{tEqualSign, 2, 4, 1, "="},
		{tValue, 2, 5, 7, "abcdefg"},
		{tLeftBracket, 4, 14, 1, "["},
		{tSection, 4, 15, 4, "User"},
		{tRightBracket, 4, 19, 1, "]"},
		{tKey, 5, 21, 8, "userName"},
		{tEqualSign, 5, 29, 1, "="},
		{tValue, 5, 30, 3, "otm"},
	}

	i := 0
	l := lex.NewLexer("", str, lexStart)
	lex.UpdateTokenTypes(tokenMap)
	defer lex.RestoreTokenTypes()

	for {
		token = l.NextToken()
		switch token.Type {
		case lex.TError:
			fmt.Println(token)
			return
		case lex.TEOF:
			return
		default:
			if expected[i].typ != token.Type {
				t.Errorf("Incorrect toke.Type, expected: %v, got: %v, value: %v", expected[i].typ, token.Type, token.Value)
			}
			if expected[i].line != token.Line {
				t.Errorf("Incorrect toke.Line, expected: %v, got: %v, value: %v", expected[i].line, token.Line, token.Value)
			}
			if expected[i].start != token.Start {
				t.Errorf("Incorrect toke.Start, expected: %v, got: %v, value: %v", expected[i].start, token.Start, token.Value)
			}
			if expected[i].width != token.Width {
				t.Errorf("Incorrect toke.Width, expected: %v, got: %v, value: %v", expected[i].width, token.Width, token.Value)
			}
			if expected[i].value != token.Value {
				t.Errorf("Incorrect toke.Value, expected: %v, got: %v, value: %v", expected[i].value, token.Value, token.Value)
			}
		}
		i++
	}
}
