package parser

import (
	"testing"

	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
	"github.com/stretchr/testify/require"
)

// full test to see that the parser spits out a
// mix of text and expressions
func TestParser(t *testing.T) {
	input := `# Hello {{ "world" }}!`
	expect := &ast.Template{
		Blocks: []ast.Block{
			&ast.BlockText{Text: "# Hello "},
			&ast.BlockExpression{
				Expression: &ast.StringLiteralExpr{Value: "world"},
			},
			&ast.BlockText{Text: "!"},
		},
	}

	r := require.New(t)

	l := lexer.NewLex("test input", input)
	p := New(l)

	tmp, err := p.Parse()
	r.NoError(err)
	r.Equal(expect, tmp)
}
