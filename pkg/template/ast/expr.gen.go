// THIS FILE IS GENERATED. DO NOT EDIT

package ast

import "github.com/fredrikkvalvik/nots/pkg/template/token"

type IdentifierExpr struct {
	Value string
}
func (n *IdentifierExpr) expressionNode() {}
func (n *IdentifierExpr) _node() {}

type NumberLiteralExpr struct {
	Value float64
}
func (n *NumberLiteralExpr) expressionNode() {}
func (n *NumberLiteralExpr) _node() {}

type StringLiteralExpr struct {
	Value string
}
func (n *StringLiteralExpr) expressionNode() {}
func (n *StringLiteralExpr) _node() {}

type ParenExpr struct {
	Expression Expr
}
func (n *ParenExpr) expressionNode() {}
func (n *ParenExpr) _node() {}

type PipeExpr struct {
	Left Expr
	Right Expr
}
func (n *PipeExpr) expressionNode() {}
func (n *PipeExpr) _node() {}

type BinaryExpr struct {
	Op token.TokenType
	Left Expr
	Right Expr
}
func (n *BinaryExpr) expressionNode() {}
func (n *BinaryExpr) _node() {}

type FunctionCallExpr struct {
	Callee Expr
	Arguments []Expr
}
func (n *FunctionCallExpr) expressionNode() {}
func (n *FunctionCallExpr) _node() {}

// this is gives us a compile time check to see of all the interafaces has ben properly implemented
func _() {
var _ Expr = &IdentifierExpr{}
var _ Expr = &NumberLiteralExpr{}
var _ Expr = &StringLiteralExpr{}
var _ Expr = &ParenExpr{}
var _ Expr = &PipeExpr{}
var _ Expr = &BinaryExpr{}
var _ Expr = &FunctionCallExpr{}
}