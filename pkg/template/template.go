package template

import (
	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/eval"
	"github.com/fredrikkvalvik/nots/pkg/template/lexer"
	"github.com/fredrikkvalvik/nots/pkg/template/parser"
)

type Template struct {
	t    *ast.Template
	eval *eval.Evaluator
	env  *eval.Env
}

func NewTemplate(input string) (*Template, error) {
	l := lexer.NewLex("", input)
	p := parser.New(l)

	t, err := p.Parse()
	if err != nil {
		return nil, err
	}

	env := newEnv
	e := eval.New(t, env())

	template := Template{
		t:    t,
		eval: e,
	}
	return &template, nil
}

func (t *Template) Execute() (string, error) {
	return t.eval.Execute()
}
