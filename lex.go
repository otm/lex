/*
Package lex provieds a skeleton for writing lexers and tokenizers. The lexer is
a state machine which emits tokens. To get started with the package only token
types and state functions are needed to be defined.

Usage

The lexer emits tokens, and tokens are identified by its TokenType. Custom
TokenTypes are defined in the following way:
  const (
    tLeftBracket lex.TokenType = iota + lex.TokenTypes
    tRightBracket
    tEqualSign

    tSection
    tKey
    tValue
  )
One notable thing is the `iota + lex.TokenTypes` which is needed as there are
predefined constans allready in lex.

The lexer is a state machine, so the next step is defining state functions of the
type StateFn. Bellow is an example of a state function that would abort the lexing
process and return an error token. Return nil in the state function to abort the
lexer without an error.
  func lexStart(l *lex.Lexer) lex.StateFn {
    return l.Errorf("Unknown charecter: %v", string(r))
  }

Example

A simple lexer that only emits alpha-numeric or special charects could look like
the example bellow. Please note that the example is not optimized but to showcase
functionality, and that the power of a lexer is not displayed in such small example.
  package main

  const (
    tAlphaNumeric lex.TokenType = iota + lex.TokenTypes
    tSpcialChar
  )

  func lexAlphaNumeric(l *lex.Lexer) StateFn {
    l.AcceptWhitespace(lex.Ignore)

    r := l.Next()

    if r == lex.EOF {
      l.Emit(TEOF)
      return nil
    }

    if lex.IsAlphaNumeric(r) {
      l.Emit(tAlphaNumeric)
      return lexAlphaNumeric
    }

    l.Backup()
    return lexSpecialChar
  }

  func lexSpecialChar(l *lex.Lexer) StateFn {
    if l.Accept("-!#%&/()=?-", lex.Emit(tSpcialChar)) {
      return lexAlphaNumeric
    }

    r := l.Next()

    if r == lex.EOF {
      l.Emit(lex.TEOF)
      return nil
    }

    if lex.IsAlphaNumeric(r) {
      l.Backup()
      return lexAlphaNumeric
    }

    return l.Errorf("Syntax error: Unknown charecter `%v`", l.Next())
  }

  func main() {
    var str = `  lex! `
    var token lex.Token

    l := lex.NewLexer("", str, lexStart)
    for {
      token = l.NextToken()
      switch token.Type {
      case lex.TError:
        fmt.Printf("Error: %v\n", token)
        return
      case lex.TEOF:
        return
      default:
        fmt.Printf("%v\n", token)
    }

  }

Which would produce the following output:
  [TokenType(3) 1:2+1] <l>
  [TokenType(3) 1:3+1] <e>
  [TokenType(3) 1:4+1] <x>
  [TokenType(4) 1:5+1] <!>
The syntax above is [TokenType() <line>:<start byte>+<byte width>]

Further Reading

1. "Lexical Scanning in Go" by Rob Pike: https://www.youtube.com/watch?v=HxaD_trXwRE

2. "GoDoc for the Template Package": http://golang.org/src/text/template/parse/lex.go

*/
package lex
