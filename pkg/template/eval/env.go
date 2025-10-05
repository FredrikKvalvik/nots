package eval

import "github.com/fredrikkvalvik/nots/pkg/template/object"

type Env struct {
	Symbols map[string]Symbol
}

func NewEnv() *Env {
	return &Env{
		Symbols: map[string]Symbol{},
	}
}

func (env *Env) GetSymbol(name string) Symbol {
	return env.Symbols[name]
}

func (e *Env) RegisterFilter(name string, fn func(o Object) (Object, error)) {
	e.Symbols[name] = &object.SymbolFilter{
		Name: name,
		Fn:   fn,
	}
}

func (e *Env) RegisterStringValue(name, value string) {
	e.Symbols[name] = &object.SymbolValue{Name: name, Val: &object.ObjectString{Val: value}}
}

func (e *Env) RegisterNumberValue(name string, value float64) {
	e.Symbols[name] = &object.SymbolValue{Name: name, Val: &object.ObjectNumber{Val: value}}
}

func (e *Env) RegisterFnValue(name string, fn func() (Object, error)) {
	e.Symbols[name] = &object.SymbolFnValue{
		Name: name,
		Fn:   fn,
	}
}
func (e *Env) RegisterFunction(name string, fn func(objs ...Object) (Object, error), validArgs object.ValidArgs) {
	e.Symbols[name] = &object.SymbolFunction{
		ValidArgs: validArgs,
		Name:      name,
		Fn:        fn,
	}
}
