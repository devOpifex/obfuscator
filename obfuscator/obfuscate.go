package obfuscator

import (
	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
)

type Obfuscator struct {
	env   *environment.Environment
	files lexer.Files
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

	case *ast.For:
		o.Obfuscate(node.Vector)
		o.env = environment.Enclose(o.env)
		o.Obfuscate(node.Value)
		o.env = environment.Open(o.env)

	case *ast.While:
		o.Obfuscate(node.Statement)
		o.env = environment.Enclose(o.env)
		o.Obfuscate(node.Value)
		o.env = environment.Open(o.env)

	case *ast.InfixExpression:
		o.Obfuscate(node.Left)

		if node.Operator == "=" || node.Operator == "<-" {
			switch n := node.Left.(type) {
			case *ast.Identifier:
				o.env.SetVariable(n.Value, environment.Variable{
					Name: n.Value,
				})
			}
		}

		o.Obfuscate(node.Right)
		return node.Right

	case *ast.IfExpression:
		o.Obfuscate(node.Condition)
		o.env = environment.Enclose(o.env)
		o.Obfuscate(node.Consequence)
		o.env = environment.Open(o.env)

		if node.Alternative != nil {
			o.env = environment.Enclose(o.env)
			o.Obfuscate(node.Alternative)
			o.env = environment.Open(o.env)
		}

	case *ast.FunctionLiteral:
		o.env = environment.Enclose(o.env)

		if node.Body != nil {
			o.Obfuscate(node.Body)
		}

		o.env = environment.Open(o.env)

	case *ast.CallExpression:
		o.obfuscateCallExpression(node)

	case *ast.Method:
		o.obfuscateMethod(node)
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

func (o *Obfuscator) obfuscateCallExpression(node *ast.CallExpression) {
	for _, a := range node.Arguments {
		o.Obfuscate(a.Value)
	}
}

func (o *Obfuscator) obfuscateMethod(node *ast.Method) {
	for _, a := range node.Arguments {
		o.Obfuscate(a.Value)
	}
}
