package obfuscator

import (
	"github.com/devOpifex/obfuscator/ast"
	"github.com/devOpifex/obfuscator/environment"
	"github.com/devOpifex/obfuscator/lexer"
)

type Obfuscator struct {
	env    *environment.Environment
	ignore []string
	files  lexer.Files
}

func New(env *environment.Environment, files lexer.Files) *Obfuscator {
	return &Obfuscator{
		env:   env,
		files: files,
	}
}

func (o *Obfuscator) Run() {
	for _, p := range o.files {
		o.Obfuscate(p.Ast)
	}
}

func (o *Obfuscator) Files() lexer.Files {
	return o.files
}

func (o *Obfuscator) Obfuscate(node ast.Node) ast.Node {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return o.obfuscateProgram(node)

	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return o.Obfuscate(node.Expression)
		}

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			o.Obfuscate(s)
		}

	case *ast.Identifier:
		return node

	case *ast.PrefixExpression:
		o.Obfuscate(node.Right)

	case *ast.InfixExpression:
		o.Obfuscate(node.Left)
		if node.Operator == "<-" {
			node.Operator = "="
		}

		if _, ok := node.Left.(*ast.Identifier); ok && node.Operator == "=" {
			o.env.SetVariable(node.Left.String())
		}
		o.Obfuscate(node.Right)
		return node.Right

	case *ast.FunctionLiteral:
		if node.Name != "" {
			o.env.SetFunction(node.Name)
		}

	case *ast.CallExpression:
		return node
	}

	return node
}

func (o *Obfuscator) obfuscateProgram(program *ast.Program) ast.Node {
	var node ast.Node
	for _, statement := range program.Statements {
		o.Obfuscate(statement)
	}
	return node
}
