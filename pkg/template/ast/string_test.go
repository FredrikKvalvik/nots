package ast

import (
	"testing"
)

func TestExpressionStringMethods(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name:     "identifier",
			expr:     &IdentifierExpr{Value: "myVar"},
			expected: "myVar",
		},
		{
			name:     "number literal",
			expr:     &NumberLiteralExpr{Value: 42.5},
			expected: "42.5",
		},
		{
			name:     "string literal",
			expr:     &StringLiteralExpr{Value: "hello"},
			expected: `"hello"`,
		},
		{
			name: "parenthesized expression",
			expr: &ParenExpr{
				Expression: &IdentifierExpr{Value: "test"},
			},
			expected: "(paren test)",
		},
		{
			name: "pipe expression",
			expr: &PipeExpr{
				Left:  &IdentifierExpr{Value: "input"},
				Right: &IdentifierExpr{Value: "filter"},
			},
			expected: "(pipe input filter)",
		},
		{
			name: "nested pipe expression",
			expr: &PipeExpr{
				Left: &PipeExpr{
					Left:  &IdentifierExpr{Value: "data"},
					Right: &IdentifierExpr{Value: "transform"},
				},
				Right: &IdentifierExpr{Value: "format"},
			},
			expected: "(pipe (pipe data transform) format)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBlockExpressionString(t *testing.T) {
	tests := []struct {
		name     string
		block    *BlockExpression
		expected string
	}{
		{
			name: "simple identifier",
			block: &BlockExpression{
				Expression: &IdentifierExpr{Value: "user"},
			},
			expected: "(block-expr user)",
		},
		{
			name: "pipe expression",
			block: &BlockExpression{
				Expression: &PipeExpr{
					Left:  &IdentifierExpr{Value: "name"},
					Right: &IdentifierExpr{Value: "upper"},
				},
			},
			expected: "(block-expr (pipe name upper))",
		},
		{
			name:     "nil expression",
			block:    &BlockExpression{Expression: nil},
			expected: "(block-expr nil)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.block.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBlockTextString(t *testing.T) {
	block := &BlockText{Text: "Hello, world!"}
	expected := `(block-text "Hello, world!")`
	result := block.String()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
