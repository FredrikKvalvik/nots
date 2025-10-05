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
		{
			name:  "parse function call, no arguments",
			input: "{{ func() }}",
			expected: &ast.BlockExpression{
				Expression: &ast.FunctionCallExpr{
					Callee:    &ast.IdentifierExpr{Value: "func"},
					Arguments: []ast.Expr{},
				},
			},
		},
		{
			name:  "parse function call, one arguments",
			input: "{{ func(1) }}",
			expected: &ast.BlockExpression{
				Expression: &ast.FunctionCallExpr{
					Callee: &ast.IdentifierExpr{Value: "func"},
					Arguments: []ast.Expr{
						&ast.NumberLiteralExpr{Value: 1},
					},
				},
			},
		},
		{
			name:  "parse function call, two arguments",
			input: "{{ func(1, 2) }}",
			expected: &ast.BlockExpression{
				Expression: &ast.FunctionCallExpr{
					Callee: &ast.IdentifierExpr{Value: "func"},
					Arguments: []ast.Expr{
						&ast.NumberLiteralExpr{Value: 1},
						&ast.NumberLiteralExpr{Value: 2},
					},
				},
			},
		},
		{
			name:  "parse function call, with function call argument",
			input: "{{ func(1, fn(a)) }}",
			expected: &ast.BlockExpression{
				Expression: &ast.FunctionCallExpr{
					Callee: &ast.IdentifierExpr{Value: "func"},
					Arguments: []ast.Expr{
						&ast.NumberLiteralExpr{Value: 1},
						&ast.FunctionCallExpr{
							Callee: &ast.IdentifierExpr{Value: "fn"},
							Arguments: []ast.Expr{
								&ast.IdentifierExpr{Value: "a"},
							},
						}},
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
