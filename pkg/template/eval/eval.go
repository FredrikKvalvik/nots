package eval

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/fredrikkvalvik/nots/pkg/template/ast"
	"github.com/fredrikkvalvik/nots/pkg/template/object"
)

type Symbol = object.Symbol
type Object = object.Object

type Evaluator struct {
	// symbol table
	env      *Env
	template *ast.Template

	errors []error

	out bytes.Buffer
}

func New(template *ast.Template, env *Env) *Evaluator {
	e := &Evaluator{
		env:      env,
		template: template,
	}

	return e
}

func (e *Evaluator) Execute() (string, error) {
	// reset out if template is run multiple times
	e.reset()
	e.runEval()

	if len(e.errors) > 0 {
		return "", errors.Join(e.errors...)
	}

	return e.out.String(), nil
}

// ExecuteWriter does a buffers the output and writes the content to
// w. no writes are done when an error occurs
func (e *Evaluator) ExecuteWriter(w io.Writer) error {
	// reset out if template is run multiple times
	e.reset()
	e.runEval()

	if len(e.errors) > 0 {
		return errors.Join(e.errors...)
	}

	// write the buffer to w, only when no errors have occured
	_, err := io.Copy(w, &e.out)
	if err != nil {
		return fmt.Errorf("an error occured when writing to w: %w", err)
	}

	return nil
}

func (e *Evaluator) runEval() {
	for _, block := range e.template.Blocks {
		switch b := block.(type) {
		case *ast.BlockText:
			e.emit(b.Text)

		case *ast.BlockExpression:
			out, err := e.eval(b.Expression)
			if err != nil {
				e.errors = append(e.errors, err)
				continue
			}
			e.emit(out.ToString())
		}
	}
}

// emit add the text to the output
func (e *Evaluator) emit(text string) {
	e.out.Write([]byte(text))
}

// resets the output. useful if a template is evaluated multiple times
func (e *Evaluator) reset() {
	e.out = bytes.Buffer{}
	e.errors = nil
}

// func (e *Evaluator) errorf(msg string, v ...any) {
// 	e.errors = append(e.errors, fmt.Errorf(msg, v...))
// }

// evalauates an expression to a string
func (e *Evaluator) eval(expr ast.Expr) (Object, error) {
	switch ex := expr.(type) {
	case *ast.NumberLiteralExpr:
		return &object.ObjectNumber{Val: ex.Value}, nil

	case *ast.StringLiteralExpr:
		return &object.ObjectString{Val: ex.Value}, nil

	case *ast.IdentifierExpr:
		symbol := e.env.GetSymbol(ex.Value)
		if symbol == nil {
			return nil, fmt.Errorf("failed to resolve symbol=%s", ex.Value)
		}
		return &object.ObjectSymbol{Val: symbol}, nil

	case *ast.PipeExpr:
		return e.evalPipe(ex)

	case *ast.FunctionCallExpr:
		return e.evalFunctionCall(ex)

	default:
		panic("unexpected node: " + ex.String())
	}
}

func (e *Evaluator) evalPipe(n *ast.PipeExpr) (Object, error) {
	lo, err := e.eval(n.Left)
	if err != nil {
		return nil, err
	}
	// the right must be an indentifier to be a valid filter
	ident, ok := n.Right.(*ast.IdentifierExpr)
	if !ok {
		return nil, fmt.Errorf("expected Identifier, got %T", n.Right)
	}
	symbol := e.env.GetSymbol(ident.Value)
	if symbol == nil {
		return nil, fmt.Errorf("failed to resolve filter with name=%s", ident.Value)
	}

	filter, ok := symbol.(*object.SymbolFilter)
	if !ok {
		return nil, fmt.Errorf("symbol=%s is not a valid filter", ident.Value)
	}

	return filter.Fn(lo)
}

func (e *Evaluator) evalFunctionCall(n *ast.FunctionCallExpr) (Object, error) {
	callee, err := e.eval(n.Callee)
	if err != nil {
		return nil, err
	}

	// only type symbolFunction is callable
	if callee.ObjectType() != object.ObjectTypeSymbol {
		return nil, fmt.Errorf("can't call non-symbol: %s", callee.ObjectType())
	}
	callSymbolObj := callee.(*object.ObjectSymbol)

	// only type symbolFunction is callable
	if callSymbolObj.Val.Type() != object.SymbolTypeFunction {
		return nil, fmt.Errorf("can't call non-function symbol: %s", callSymbolObj)
	}

	funcSymbol := callSymbolObj.Val.(*object.SymbolFunction)
	args := []Object{}
	for _, arg := range n.Arguments {
		res, err := e.eval(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, res)
	}

	return funcSymbol.Call(args...)
}
