// template package is the main entry for working with templates.
// it abstracts all the lexing/parsing/env setup, and exports a simple to use interface
// for generating templates.
package template

import (
	"io"
	"os"

	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/eval"
	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
	"github.com/fredrikkvalvik/nots/pkg/template/parser"
)

type Object = eval.Object
type Template struct {
	t    *ast.Template
	eval *eval.Evaluator
	env  *eval.Env
}

func New(input string) (*Template, error) {
	l := lexer.NewLex("template", input)
	p := parser.New(l)

	t, err := p.Parse()
	if err != nil {
		return nil, err
	}

	env := newEnv()
	e := eval.New(t, env)

	template := Template{
		t:    t,
		eval: e,
		env:  env,
	}
	return &template, nil
}

func FromPath(path string) (*Template, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return New(string(b))
}

// return the string output of the evaluated template, or potential error
func (t *Template) Execute() (string, error) {
	return t.eval.Execute()
}

// execute the template and write the result to w. no write are done on error (unless writing to w errors)
func (t *Template) ExecuteWriter(w io.Writer) error {
	return t.eval.ExecuteWriter(w)
}

// register a static value with name
func (t *Template) RegisterStringValue(name, value string) {
	t.env.RegisterStringValue(name, value)
}

// register a static value with name
func (t *Template) RegisterNumberValue(name string, value float64) {
	t.env.RegisterNumberValue(name, value)
}

// register a 'filter'. A filter is used with the '|' operator to take a value obj
// on the left side of the pipe, and modify the value, outputing it its right side.
//
// ex: `value | filter` -> result
func (t *Template) RegisterFilter(name string, fn func(obj Object) (Object, error)) {
	t.env.RegisterFilter(name, fn)
}

// regisers a function that will be called when its name is evaluated.
//
// Useful for dynamic values like date or time. Can also do http requests and fetch external data
func (t *Template) RegisterFnValue(name string, fn func() (Object, error)) {
	t.env.RegisterFnValue(name, fn)
}
