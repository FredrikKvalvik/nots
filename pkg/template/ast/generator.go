//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"
)

// const stmt = "Stmt"
const expr = "Expr"

// const stmtMethod = "func (%s *%s) textNode() {}"
const exprMethod = "func (%s *%s) expressionNode() {}"

const packageName = "ast"
const tokenPkg = "github.com/fredrikkvalvik/nots/pkg/template/token"

type keyVal struct {
	key   string
	value string
}
type template struct {
	name  string
	props []keyVal
}

var exprs = []template{
	{
		name: "Identifier",
		props: []keyVal{
			{"Value", "string"},
		},
	},
	{
		name: "NumberLiteral",
		props: []keyVal{
			{"Value", "float64"},
		},
	},
	{
		name: "StringLiteral",
		props: []keyVal{
			{"Value", "string"},
		},
	},
	// {
	// 	name: "BooleanLiteral",
	// 	props: []keyVal{
	// 		{"Value", "bool"},
	// 	},
	// },
	// {
	// 	name: "Unary",
	// 	props: []keyVal{
	// 		{"Operand", "token.TokenType"},
	// 		{"Right", expr},
	// 	},
	// },
	// {
	// 	name: "Binary",
	// 	props: []keyVal{
	// 		{"Operand", "token.TokenType"},
	// 		{"Left", expr},
	// 		{"Right", expr},
	// 	},
	// },
	// {
	// 	name: "Logical",
	// 	props: []keyVal{
	// 		{"Operand", "token.TokenType"},
	// 		{"Left", expr},
	// 		{"Right", expr},
	// 	},
	// },
	{
		name: "Paren",
		props: []keyVal{
			{"Expression", expr},
		},
	},
	{
		name: "Pipe",
		props: []keyVal{
			{"Left", expr},
			{"Right", expr},
		},
	},
	// {
	// 	name: "FunctionLiteral",
	// 	props: []keyVal{
	// 		{"Arguments", "[]*Identifier" + expr},
	// 		{"Body", "*Block" + stmt},
	// 	},
	// },
}

// This will generate a file for statements and expressions
// the only unique part of the structs are the fields and the String method
func main() {
	expressionFile := generateNodes(expr, exprMethod, exprs)

	os.WriteFile("expr.gen.go", []byte(expressionFile), 0646)
}

func generateNodes(interfaceName, interfaceMethod string, tmpl []template) string {
	var f strings.Builder

	f.WriteString("// THIS FILE IS GENERATED. DO NOT EDIT\n\n")
	f.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	// f.WriteString(fmt.Sprintf(`import "%s"`+"\n\n", tokenPkg))

	for _, s := range tmpl {
		name := s.name + interfaceName

		f.WriteString(fmt.Sprintf("type %s struct {\n", name))

		// add token to all ast nodes
		kvs := append([]keyVal{ /*{"Token", "token.Token"}*/ }, s.props...)
		for _, kv := range kvs {
			f.WriteString(fmt.Sprintf("\t%s %s\n", kv.key, kv.value))
		}

		f.WriteString("}\n")
		f.WriteString(fmt.Sprintf(interfaceMethod, "n", name) + "\n")
		f.WriteString(fmt.Sprintf("func (n *%s) _node() {}", name) + "\n")
		// f.WriteString(fmt.Sprintf("func (n *%s) Lexeme() string { return n.Token.Val }\n", name))
		// f.WriteString(fmt.Sprintf("func (n *%s) Literal() any { return n.Token.Literal }\n", name))
		// f.WriteString(fmt.Sprintf("func (n *%s) GetToken() *token.Token { return &n.Token }\n", name))

		// create space for next struct
		f.WriteString("\n")
	}

	fmt.Fprint(&f, "// this is gives us a compile time check to see of all the interafaces has ben properly implemented\n")
	fmt.Fprintf(&f, "func _() {\n")
	for _, s := range tmpl {
		name := s.name + interfaceName

		fmt.Fprintf(&f, "var _ %s = &%s{}\n", interfaceName, name)
	}
	fmt.Fprint(&f, "}")

	return f.String()
}
