package eval

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
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
			input:    `{{ hello_world|capitalize }}`,
			expected: "HELLO WORLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			l := lexer.NewLex(tt.name, tt.input)
			p := parser.New(l)
			template, err := p.Parse()
			r.NoError(err)

			ev := New(template, &Env{symbols: map[string]Symbol{
				"capitalize": &SymbolFilter{Fn: func(v Object) (Object, error) {
					if v.ObjectType() != ObjectTypeString {
						return nil, fmt.Errorf("expected string, got=%s", v.ObjectType())
					}

					return &ObjectString{Val: strings.ToUpper(v.ToString())}, nil
				}},
				"hello_world": &SymbolValue{
					Name: "hello_world",
					Val:  &ObjectString{Val: "hello world"},
				},
			},
			})

			result, err := ev.Execute()
			r.NoError(err)

			r.Equal(tt.expected, result)
		})
	}
}
