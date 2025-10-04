package ast

import "fmt"

func (n *IdentifierExpr) String() string {
	return n.Value
}

func (n *NumberLiteralExpr) String() string {
	return fmt.Sprintf("%g", n.Value)
}

func (n *StringLiteralExpr) String() string {
	return fmt.Sprintf("\"%s\"", n.Value)
}

func (n *ParenExpr) String() string {
	return fmt.Sprintf("(paren %s)", n.Expression.String())
}

func (n *PipeExpr) String() string {
	return fmt.Sprintf("(pipe %s %s)", n.Left.String(), n.Right.String())
}
