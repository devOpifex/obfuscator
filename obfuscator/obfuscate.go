package obfuscator

import (
	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
)

type Obfuscator struct {
	env       *environment.Environment
	callStack Stack
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
		o.Obfuscate(node.Right)

		if node.Operator == "::" || node.Operator == ":::" {
			return node
		}

		switch l := node.Left.(type) {
		case *ast.Identifier:
			switch n := node.Right.(type) {
			case *ast.FunctionLiteral:
				o.env.SetFunction(l.Value, environment.Function{
					Name:  l.Value,
					Value: n,
				})
			default:
				ok, c := o.callStack.Get()

				if ok && c.name != "" {
					break
				}

				if node.Operator == "$" {
					break
				}

				o.env.SetVariable(l.Value, environment.Variable{
					Name: l.Value,
				})
			}
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
	o.callStack = o.callStack.Push(node.Name, false)
	for _, a := range node.Arguments {
		o.Obfuscate(a)
	}
	o.callStack = o.callStack.Pop()
}

func (o *Obfuscator) obfuscateMethod(node *ast.Method) {
	o.callStack = o.callStack.Push(node.Name, true)
	for _, a := range node.Arguments {
		o.Obfuscate(a)
	}
	o.callStack = o.callStack.Pop()
}
