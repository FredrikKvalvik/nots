//go:generate go tool golang.org/x/tools/cmd/stringer -type TokenType
package token

import "fmt"

type TokenType int

func (i Token) String() string {
	switch i.Type {
	case TokenTypeEOF:
		return "EOF"
	case TokenTypeError:
		return i.Val
	}
	if len(i.Val) > 10 {
		return fmt.Sprintf("[%s] %.10q...", i.Type, i.Val)
	}
	return fmt.Sprintf("[%s] %q", i.Type, i.Val)
}

const (
	TokenTypeError TokenType = iota // error occurred;
	// value is text of error
	TokenTypeDot // the cursor, spelled '.'
	TokenTypeEOF
	TokenTypeText // any type of text outside an expression

	TokenTypeIdentifier // an unqouted alphanumeric string starting with a letter, inside an expression
	TokenTypeNumber     // any valid integer/float value, inside an expression
	TokenTypeString     // any quoted string inside an expression

	TokenTypeLMeta // '{{'
	TokenTypeRMeta // '}}'

	TokenTypeLParen // '('
	TokenTypeRParen // ')'
	TokenTypePipe   // '|'
)

// item represents a token returned from the scanner.
type Token struct {
	Type       TokenType // Type, such as itemNumber.
	Val        string    // Value, such as "23.2".
	Start, End int
}
