//go:generate go tool golang.org/x/tools/cmd/stringer -type ObjectType
package eval

import "fmt"

type ObjectType int

const (
	_ ObjectType = iota

	ObjectTypeNumber
	ObjectTypeString
	ObjectTypeSymbol
)

type Object interface {
	ObjectType() ObjectType
	ToString() string
}

type ObjectString struct {
	Val string
}

func (*ObjectString) ObjectType() ObjectType {
	return ObjectTypeString
}
func (o *ObjectString) ToString() string {
	return o.Val
}

type ObjectNumber struct {
	Val float64
}

func (*ObjectNumber) ObjectType() ObjectType {
	return ObjectTypeNumber
}
func (o *ObjectNumber) ToString() string {
	return fmt.Sprint(o.Val)
}

type ObjectSymbol struct {
	Name string
	Val  Symbol
}

func (*ObjectSymbol) ObjectType() ObjectType {
	return ObjectTypeSymbol
}
func (o *ObjectSymbol) ToString() string {
	switch symbol := o.Val.(type) {
	case *SymbolValue:
		return symbol.Val.ToString()

	case *SymbolFilter:
		return fmt.Sprintf(`[filter "%s"]`, symbol.Name)

	default:
		panic("unexpected symbol type: " + symbol.SymbolName())
	}
}
