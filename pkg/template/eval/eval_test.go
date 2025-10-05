package eval

import (
	"fmt"
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
		{
			name:     "fnValue evaluation",
			input:    `{{ fn_value }}`,
			expected: "from fn_value",
		},
		{
			name:     "function call, no args",
			input:    `{{ func0() }}`,
			expected: "from func0",
		},
		{
			name:     "function call, one arg",
			input:    `{{ echo("hello") }}`,
			expected: "hello",
		},
		{
			name:     "function call, function call result as arg",
			input:    `{{ echo(func0()) }}`,
			expected: "from func0",
		},
	}

	env := testEnv()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			l := lexer.NewLex(tt.name, tt.input)
			p := parser.New(l)
			template, err := p.Parse()
			r.NoError(err)

			ev := New(template, env)

			result, err := ev.Execute()
			r.NoError(err)

			r.Equal(tt.expected, result)
		})
	}
}

// test symbols for tests
func testEnv() *Env {
	return &Env{Symbols: map[string]Symbol{
		"capitalize": &object.SymbolFilter{
			Name: "capitalize",
			Fn: func(v Object) (Object, error) {
				return &object.ObjectString{
					Val: strings.ToUpper(v.ToString()),
				}, nil
			}},
		"lower": &object.SymbolFilter{
			Name: "lower",
			Fn: func(v Object) (Object, error) {
				return &object.ObjectString{Val: strings.ToLower(v.ToString())}, nil
			}},
		"hello_world": &object.SymbolValue{
			Name: "hello_world",
			Val:  &object.ObjectString{Val: "hello world"},
		},
		"fn_value": &object.SymbolFnValue{
			Name: "fn_value",
			Fn: func() (object.Object, error) {
				return &object.ObjectString{Val: "from fn_value"}, nil
			},
		},
		"func0": &object.SymbolFunction{
			Name: "func0",
			Fn: func(o ...object.Object) (object.Object, error) {
				if len(o) != 0 {
					return nil, fmt.Errorf("expect 0 args")
				}
				return &object.ObjectString{Val: "from func0"}, nil
			},
		},
		"echo": &object.SymbolFunction{
			Name: "echo",

			Fn: func(o ...object.Object) (object.Object, error) {
				return &object.ObjectString{Val: o[0].ToString()}, nil
			},
		},
	}}
}
