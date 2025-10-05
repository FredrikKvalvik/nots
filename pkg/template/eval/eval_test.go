package eval

import (
	"strings"
	"testing"

	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
	"github.com/fredrikkvalvik/nots/pkg/template/object"
	"github.com/fredrikkvalvik/nots/pkg/template/parser"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple expression",
			input:    `{{ "hello world" }}`,
			expected: "hello world",
		},
		{
			name:     "expression with text",
			input:    `# hello {{ "world"|capitalize }}!`,
			expected: "# hello WORLD!",
		},
		{
			name:     "expression with filter",
			input:    `{{ "hello world"|capitalize }}`,
			expected: "HELLO WORLD",
		},
		{
			name:     "expression with variable",
			input:    `{{ hello_world }}`,
			expected: "hello world",
		},
		{
			name:     "expression with variable through filter",
			input:    `{{ hello_world|capitalize|lower }}`,
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			l := lexer.NewLex(tt.name, tt.input)
			p := parser.New(l)
			template, err := p.Parse()
			r.NoError(err)

			ev := New(template, &Env{Symbols: map[string]Symbol{
				"capitalize": &object.SymbolFilter{Fn: func(v Object) (Object, error) {
					return &object.ObjectString{strings.ToUpper(v.ToString())}, nil
				}},
				"lower": &object.SymbolFilter{Fn: func(v Object) (Object, error) {
					return &object.ObjectString{strings.ToLower(v.ToString())}, nil
				}},
				"hello_world": &object.SymbolValue{
					Name: "hello_world",
					Val:  &object.ObjectString{Val: "hello world"},
				},
			}})

			result, err := ev.Execute()
			r.NoError(err)

			r.Equal(tt.expected, result)
		})
	}
}
