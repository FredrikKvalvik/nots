package lexer

import (
	"testing"

	"github.com/fredrikkvalvik/nots/pkg/template/token"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		expect []token.Token
	}{
		{
			`hello world {{ ident "string" 123 1.1 }}`,
			[]token.Token{
				{Type: token.TokenTypeText, Val: "hello world "},
				{Type: token.TokenTypeLMeta, Val: "{{"},
				{Type: token.TokenTypeIdentifier, Val: "ident"},
				{Type: token.TokenTypeString, Val: `"string"`},
				{Type: token.TokenTypeNumber, Val: `123`},
				{Type: token.TokenTypeNumber, Val: `1.1`},
				{Type: token.TokenTypeRMeta, Val: "}}"},
				{Type: token.TokenTypeEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			r := require.New(t)
			l := NewLex("input", tt.input)

			toks := []token.Token{}
			for {
				tok := l.NextToken()
				if tok.Type == token.TokenTypeEOF {
					toks = append(toks, tok)
					break
				}

				toks = append(toks, tok)
			}

			r.Equal(len(tt.expect), len(toks), "expect input len and expect len to be equal")
			for idx, tok := range toks {
				r.Equal(tt.expect[idx].Type, tok.Type, "type must be equal")
				r.Equal(tt.expect[idx].Val, tok.Val, "lexeme must be equal")
			}
		})
	}
}
