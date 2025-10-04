package parser

import (
	"testing"

	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
	"github.com/stretchr/testify/require"
)

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.BlockExpression
	}{
		{
			name:  "number literal",
			input: "{{ 10 }}",
			expected: &ast.BlockExpression{
				Expression: &ast.NumberLiteralExpr{Value: 10},
			},
		},
		{
			name:  "number literal with trailing '.'",
			input: "{{ 10. }}",
			expected: &ast.BlockExpression{
				Expression: &ast.NumberLiteralExpr{Value: 10},
			},
		},
		{
			name:  "number literal with decimal point",
			input: "{{ 10.5 }}",
			expected: &ast.BlockExpression{
				Expression: &ast.NumberLiteralExpr{Value: 10.5},
			},
		},
		{
			name:  "string literal",
			input: `{{ "hello" }}`,
			expected: &ast.BlockExpression{
				Expression: &ast.StringLiteralExpr{Value: "hello"},
			},
		},
		{
			name:  "pipe with two identifiers",
			input: "{{ ident | ident2 }}",
			expected: &ast.BlockExpression{
				Expression: &ast.PipeExpr{
					Left:  &ast.IdentifierExpr{Value: "ident"},
					Right: &ast.IdentifierExpr{Value: "ident2"},
				},
			},
		},
		{
			name:  "multiple piped values",
			input: "{{ ident1 | ident2 | ident3 }}",
			expected: &ast.BlockExpression{
				Expression: &ast.PipeExpr{
					Left: &ast.PipeExpr{
						Left:  &ast.IdentifierExpr{Value: "ident1"},
						Right: &ast.IdentifierExpr{Value: "ident2"},
					},
					Right: &ast.IdentifierExpr{Value: "ident3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			l := lexer.NewLex(tt.name, tt.input)
			p := New(l)

			result, err := p.Parse()
			r.NoError(err)

			expect := &ast.Template{Blocks: []ast.Block{tt.expected}}
			r.Equal(expect.Blocks[0].String(), result.Blocks[0].String())
			// need to wrap the expression in a template for valid comparison
			r.Equal(expect.Blocks[0], result.Blocks[0])
		})
	}
}
