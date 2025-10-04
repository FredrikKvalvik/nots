//go:generate go run generator.go
package ast

import "fmt"

type Node interface {
	// add a private method to make the interface unusable outside package
	_node()
}

type Block interface {
	Node
	blockNode()

	String() string
}

type Expr interface {
	Node
	expressionNode()
	String() string
}

type Template struct {
	// a list of blocks is a mix of text and expressions
	//
	// Text a flat text with nothing more interesting to them
	//
	// Expr are parsed, and later evaluated to be replaced as text
	// in the finished template
	Blocks []Block
}

var _ Block = &BlockExpression{}

// var _ Expr = &BlockExpression{}

type BlockExpression struct {
	Expression Expr
}

// _node implements Expr.
func (m *BlockExpression) _node() {
	panic("unimplemented")
}

func (m *BlockExpression) blockNode() {}

func (b *BlockExpression) String() string {
	if b.Expression == nil {
		return "(block-expr nil)"
	}
	return fmt.Sprintf("(block-expr %s)", b.Expression.String())
}

type BlockText struct {
	Text string
}

// _node implements Expr.
func (m *BlockText) _node() {
	panic("unimplemented")
}

func (m *BlockText) blockNode() {}

func (b *BlockText) String() string {
	return fmt.Sprintf("(block-text \"%s\")", b.Text)
}
