package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// StateFn defines a lexer
type StateFn func(*Lexer) StateFn

// StateOpt are lexer options
type StateOpt func(*Lexer)

// Lexer defines our parser
type Lexer struct {
	Name   string
	Input  string
	Tokens chan Token
	State  StateFn

	Start int
	Pos   int
	Width int
	Line  int
}

// NewLexer returns a new lexer
func NewLexer(name, input string, inital StateFn) *Lexer {
	l := &Lexer{
		Name:   name,
		Input:  input,
		State:  inital,
		Tokens: make(chan Token, 3),
		Line:   1,
	}

	go l.Run()
	return l
}

// Run runs the state machine for the lexer.
func (l *Lexer) Run() {
	for l.State != nil {
		l.State = l.State(l)
	}
}

// Shutdown the parser
func (l *Lexer) Shutdown() {
	close(l.Tokens)
}

// Emit puts a token onto the token channel. The value of l token is
// read from the input based on the current lexer position.
func (l *Lexer) Emit(tokenType TokenType) {
	if tokenType != TIgnore {
		l.Tokens <- Token{Type: tokenType, Value: l.Input[l.Start:l.Pos], Line: l.Line, Start: l.Start, Width: l.Pos - l.Start}
	}
	l.Start = l.Pos
}

// Ignore ignors the currently parsed data
func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

// Next reads the next rune (character) from the input stream
// and advances the lexer position.
func (l *Lexer) Next(options ...StateOpt) rune {
	if l.Pos >= utf8.RuneCountInString(l.Input) {
		l.Width = 0
		return EOF
	}

	r, width := utf8.DecodeRuneInString(l.Input[l.Pos:])

	if r == '\n' {
		l.Line++
	}

	l.Width = width
	l.Pos += l.Width

	l.doOpts(options)
	return r
}

// Step the lexer n times, n can be possitive or -1
// options can be any StateOpt function
func (l *Lexer) Step(n int, options ...StateOpt) {
	action := l.Next
	if n < 0 {
		n = n * -1
		action = l.Backup
	}

	for ; n > 0; n-- {
		action()
	}

	l.doOpts(options)
}

// Peek returns next rune without advancing the parser
func (l *Lexer) Peek() rune {
	rune := l.Next()
	l.Backup()
	return rune
}

// Backup moves parser to last read token, if called more then once per call to
// Next it will not move the possition back
func (l *Lexer) Backup(options ...StateOpt) rune {
	if l.Width == 0 {
		result, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
		return result
	}

	l.Pos -= l.Width
	l.Width = 0

	result, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])

	if result == '\n' {
		l.Line--
	}

	l.doOpts(options)
	return result
}

// Accept consumes the next rune if it's from the valid set.
func (l *Lexer) Accept(valid string, options ...StateOpt) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		l.doOpts(options)
		return true
	}
	l.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set.
func (l *Lexer) AcceptRun(valid string, options ...StateOpt) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
	l.doOpts(options)
}

// LineNumber returns the current line number by calculating it
func (l *Lexer) LineNumber() int {
	return 1 + strings.Count(l.Input[:l.Pos], string(NewLine))
}

// Errorf returns an error token and terminates the scan
func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.Tokens <- Token{Type: TError, Start: l.Start, Line: l.Line, Value: fmt.Sprintf(format, args...)}
	return nil
}

// NextToken returns the next token from the channel
func (l *Lexer) NextToken() Token {
	return <-l.Tokens
}

// InputToEnd returns a slice of the input from the current lexer position
// to the end of the input string.
func (l *Lexer) InputToEnd() string {
	return l.Input[l.Pos:]
}

// IsEOF checks if we reached the end of file
func (l *Lexer) IsEOF() bool {
	return l.Pos >= len(l.Input)
}

// IsEndOfLine returns true if r is an end-of-line charecter
func IsEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// IsAlphaNumeric returns true if r alphabetic, numeric or underscore
func IsAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// AcceptWhitespace consumes a run of whitespace and newline charecters
func (l *Lexer) AcceptWhitespace(options ...StateOpt) {
	for {
		r := l.Next()

		if !unicode.IsSpace(r) && !IsEndOfLine(r) {
			l.Backup()
			l.doOpts(options)
			return
		}

		if r == EOF {
			l.doOpts(options)
			l.Emit(TEOF)
			return
		}
	}
}

// lexer functions

// LexString returns a StateFn for the given string
func (l *Lexer) LexString(s string, typ TokenType, fn StateFn) StateFn {
	return func(l *Lexer) StateFn {
		for _, runeValue := range s {
			r := l.Next()
			if runeValue != r {
				l.Errorf("LexString: expected: %v, got: %v", runeValue, r)
			}
		}
		l.Emit(typ)
		return fn
	}
}

// option functions

// Run options
func (l *Lexer) doOpts(options []StateOpt) {
	for _, opt := range options {
		opt(l)
	}
}

// Ignore output
func Ignore(l *Lexer) {
	l.Ignore()
}

// Emit token
func Emit(token TokenType) StateOpt {
	return func(l *Lexer) {
		l.Emit(token)
	}
}
