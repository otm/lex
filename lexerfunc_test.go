package lex_test

import (
	"strings"

	"github.com/otm/lex"
)

// LexBegin is the initial lexer function
func LexBegin(lexer *lex.Lexer) lex.StateFn {
	lexer.AcceptWhitespace(lex.Ignore)

	if strings.HasPrefix(lexer.InputToEnd(), MinusSign) {
		return LexSingleLineComment
	}
	return LexCode
}

// LexSingleLineComment lexer
func LexSingleLineComment(lexer *lex.Lexer) lex.StateFn {
	lexer.Step(2, lex.Ignore)
	if strings.HasPrefix(lexer.InputToEnd(), LeftBracket) {
		return LexMultiLineComment(lexer)
	}

	for {
		if strings.HasPrefix(lexer.InputToEnd(), string(lex.NewLine)) {
			lexer.Emit(TComment)
			return LexBegin(lexer)
		}

		lexer.Next()

		if lexer.IsEOF() {
			lexer.Emit(TComment)
			lexer.Emit(lex.TEOF)
			lexer.Shutdown()
			return nil
		}
	}
}

// LexMultiLineComment lexer
func LexMultiLineComment(lexer *lex.Lexer) lex.StateFn {
	lexer.Step(2, lex.Ignore)

	for {
		if strings.HasPrefix(lexer.InputToEnd(), RightBracket) {
			lexer.Emit(TComment)
			return LexBegin(lexer)
		}

		lexer.Next()

		if lexer.IsEOF() {
			panic("EndOfFile")
		}
	}
}

// LexCode parsers random code
func LexCode(lexer *lex.Lexer) lex.StateFn {
	for {
		if strings.HasPrefix(lexer.InputToEnd(), MinusSign) {
			lexer.Start = lexer.Pos
			return LexSingleLineComment(lexer)
		}

		if strings.HasPrefix(lexer.InputToEnd(), string(lex.NewLine)) {
			lexer.Next()
			lexer.Start = lexer.Pos
			return LexCode(lexer)
		}

		lexer.Next()

		if lexer.IsEOF() {
			lexer.Emit(lex.TEOF)
			lexer.Shutdown()
			return nil
		}
	}
}
