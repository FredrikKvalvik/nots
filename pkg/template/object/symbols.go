package object

type symbolType int

const (
	_ symbolType = iota

	SymbolTypeValue
	SymbolTypeFilter
	SymbolTypeFnValue
	SymbolTypeFunction
)

type Symbol interface {
	_symbol()

	Type() symbolType
	SymbolName() string
}

// === SYMBOL VALUE ===
var _ Symbol = &SymbolValue{}

type SymbolValue struct {
	Name string
	Val  Object
}

// Type implements Symbol.
func (s *SymbolValue) Type() symbolType {
	return SymbolTypeValue
}

// Type implements Symbol.
func (s *SymbolValue) SymbolName() string {
	return s.Name
}

// _symbol implements Symbol.
func (s *SymbolValue) _symbol() {}

// === SYMBOL FILTER ===
var _ Symbol = &SymbolFilter{}

// a filter takes a value (through a pipe), and does "something"
// with it, and returns a new value
type SymbolFilter struct {
	Fn   func(v Object) (Object, error)
	Name string
}

// Type implements Symbol.
func (s *SymbolFilter) Type() symbolType {
	return SymbolTypeFilter
}

// Type implements Symbol.
func (s *SymbolFilter) SymbolName() string {
	return s.Name
}

// _symbol implements Symbol.
func (s *SymbolFilter) _symbol() {}

// === SYMBOL FnValue ===
var _ Symbol = &SymbolFnValue{}

// a fnValue evaluates Fn to return a value.
type SymbolFnValue struct {
	Name string
	Fn   func() (Object, error)
}

// Type implements Symbol.
func (s *SymbolFnValue) Type() symbolType {
	return SymbolTypeFnValue
}

// Type implements Symbol.
func (s *SymbolFnValue) SymbolName() string {
	return s.Name
}

// _symbol implements Symbol.
func (s *SymbolFnValue) _symbol() {}
