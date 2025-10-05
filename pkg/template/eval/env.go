package eval

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
	e.Symbols[name] = &SymbolFilter{
		Name: name,
		Fn:   fn,
	}
}

func (e *Env) RegisterStringValue(name, value string) {
	e.Symbols[name] = &SymbolValue{Name: name, Val: &ObjectString{Val: value}}
}

func (e *Env) RegisterNumberValue(name string, value float64) {
	e.Symbols[name] = &SymbolValue{Name: name, Val: &ObjectNumber{Val: value}}
}
