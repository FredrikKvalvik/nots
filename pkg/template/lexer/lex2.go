package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fredrikkvalvik/nots/pkg/template/token"
)

const (
	leftMeta  = "{{"
	rightMeta = "}}"

	eof = rune(0)
)

type stateFn func(l *Lex) stateFn

type Lex struct {
	Name  string           // used only for error reports.
	input string           // the string being scanned.
	start int              // start position of this item.
	pos   int              // current position in the input.
	width int              // width of last rune read from input.
	items chan token.Token // channel of scanned items.
}

func NewLex(name, input string) *Lex {
	l := &Lex{
		Name:  name,
		input: input,
		items: make(chan token.Token, 2),
	}
	go l.run() // Concurrently run state machine.
	return l
}

func (l *Lex) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}

	close(l.items)
}

// return the next token
func (l *Lex) NextToken() (item token.Token) {
	return <-l.items
}

// emit a new token to ch
func (l *Lex) emit(t token.TokenType) {
	i := token.Token{
		Type: t,
		Val:  l.input[l.start:l.pos],
	}
	l.items <- i
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *Lex) next() (char rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	char, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return char
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lex) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lex) peek() rune {
	char := l.next()
	l.backup()
	return char
}

// emits an error and returns nil, stopping the lexer
func (l *Lex) errorf(str string, v ...any) stateFn {
	l.items <- token.Token{
		Type: token.TokenTypeError,
		Val:  fmt.Sprintf(str, v...),
	}
	return nil
}

// ignore current char
func (l *Lex) ignore() {
	l.start = l.pos
}

// accept consumes the next rune
// if it's from the valid set.
func (l *Lex) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lex) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func lexText(l *Lex) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
			if l.pos > l.start {
				l.emit(token.TokenTypeText)
			}
			return lexLeftMeta // Next state.

		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(token.TokenTypeText)
	}
	l.emit(token.TokenTypeEOF) // Useful to make EOF a token.
	return nil                 // Stop the run loop
}

func lexLeftMeta(l *Lex) stateFn {
	l.pos += len(leftMeta)
	l.emit(token.TokenTypeLMeta)
	return lexInsideAction // we are now iside an expression {{ HERE }}
}
func lexRightMeta(l *Lex) stateFn {
	l.pos += len(rightMeta)
	l.emit(token.TokenTypeRMeta)
	return lexText // we are now done with expression {{ ... }} HERE
}

func lexInsideAction(l *Lex) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) {
			return lexRightMeta
		}

		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")

		case isSpace(r):
			l.ignore()

		case r == '|':
			l.emit(token.TokenTypePipe)

		case r == '(':
			l.emit(token.TokenTypeLParen)

		case r == ')':
			l.emit(token.TokenTypeRParen)

		case r == '"':
			l.backup()
			return lexStringLiteral

		case isAlpha(r):
			l.backup()
			return lexIdentifier

		case isNumeric(r):
			l.backup()
			return lexNumber
		}
	}
}

func lexNumber(l *Lex) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)

	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(token.TokenTypeNumber)
	return lexInsideAction
}

func lexStringLiteral(l *Lex) stateFn {
	// there must be an opening '"' for a string to be valid
	if !l.accept(`"`) {
		return l.errorf("expected opening string rune")
	}

	for {
		ch := l.next()
		if ch == '\n' {
			return l.errorf("illegal newline in string literal")
		}
		if ch == eof {
			return l.errorf("unexpected end of file")
		}
		if ch == '"' {
			break
		}
	}

	l.emit(token.TokenTypeString)
	return lexInsideAction
}

func lexIdentifier(l *Lex) stateFn {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if !l.accept(valid) {
		return l.errorf("unexpected character=%s", string(l.peek()))
	}
	valid += "0123456789"
	l.acceptRun(valid)
	l.emit(token.TokenTypeIdentifier)

	return lexInsideAction
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isAlpha(r rune) bool {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return strings.ContainsRune(valid, r)
}

func isNumeric(r rune) bool {
	return unicode.IsNumber(r)
}
