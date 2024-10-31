package transpiler

import (
	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
)

type Obfuscator struct {
	code []string
	env  *environment.Environment
}

func New(env *environment.Environment) *Obfuscator {
	return &Obfuscator{
		env: env,
	}
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
		o.addCode(node.Value)
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
		if node.Left != nil {
			o.Obfuscate(node.Left)
		}

		if node.Right != nil {
			o.Obfuscate(node.Right)
		}

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

		for _, p := range node.Parameters {
			o.Obfuscate(p.Expression)
		}

		if node.Body != nil {
			o.Obfuscate(node.Body)
		}

		o.env = environment.Open(o.env)

	case *ast.CallExpression:
		o.obfuscateCallExpression(node)
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
	name := node.Name
	fn, ok := o.env.GetFunction(name, true)

	if ok {
		name = fn.Obfuscated
	}

	for _, a := range node.Arguments {
		o.Obfuscate(a)
	}
}
