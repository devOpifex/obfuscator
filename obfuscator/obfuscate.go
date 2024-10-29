package transpiler

import (
	"strings"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
)

type Obfuscator struct {
	code []string
	env  *environment.Environment
}

func (t *Obfuscator) Env() *environment.Environment {
	return t.env
}

func New() *Obfuscator {
	env := environment.New()

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

	case *ast.Comma:
		o.addCode(",")

	case *ast.Null:
		o.addCode("NULL")

	case *ast.Keyword:
		o.addCode(node.Value)

	case *ast.CommentStatement:
		o.addCode("")

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			o.Obfuscate(s)
			o.addCode(";")
		}

	case *ast.Identifier:
		v, ok := o.env.GetVariable(node.Value, true)

		if ok {
			o.addCode(v.Obfuscated)
			return node
		}

		fn, ok := o.env.GetFunction(node.Value, true)

		if ok {
			o.addCode(fn.Obfuscated)
			return node
		}

		o.addCode(node.Value)
		return node

	case *ast.Boolean:
		o.addCode(strings.ToUpper(node.String()))

	case *ast.IntegerLiteral:
		o.addCode(node.Value)

	case *ast.FloatLiteral:
		o.addCode(node.Value)

	case *ast.StringLiteral:
		o.addCode(node.Token.Value + node.Str + node.Token.Value)

	case *ast.PrefixExpression:
		o.addCode("(" + node.Operator)
		o.Obfuscate(node.Right)
		o.addCode(")")

	case *ast.For:
		o.addCode("for(")
		o.addCode(node.Name)
		o.addCode(" in ")
		o.Obfuscate(node.Vector)
		o.addCode("){")
		o.env = environment.Enclose(o.env)
		o.Obfuscate(node.Value)
		o.addCode("}")
		o.env = environment.Open(o.env)

	case *ast.While:
		o.addCode("while(")
		o.Obfuscate(node.Statement)
		o.addCode("){")
		o.env = environment.Enclose(o.env)
		o.Obfuscate(node.Value)
		o.addCode("}")
		o.env = environment.Open(o.env)

	case *ast.InfixExpression:
		if node.Operator == "in" {
			o.addCode(" ")
		}

		if node.Operator == "<-" {
			node.Operator = "="
		}

		if node.Left != nil {
			o.Obfuscate(node.Left)
		}

		o.addCode(node.Operator)

		if node.Operator == "in" {
			o.addCode(" ")
		}

		if node.Right != nil {
			o.Obfuscate(node.Right)
		}

	case *ast.Square:
		o.addCode(node.Token.Value)

	case *ast.IfExpression:
		o.addCode("if(")
		o.Obfuscate(node.Condition)
		o.addCode("){")
		o.env = environment.Enclose(o.env)
		o.Obfuscate(node.Consequence)
		o.env = environment.Open(o.env)
		o.addCode("}")

		if node.Alternative != nil {
			o.addCode("else{")
			o.env = environment.Enclose(o.env)
			o.Obfuscate(node.Alternative)
			o.env = environment.Open(o.env)
			o.addCode("}")
		}

	case *ast.FunctionLiteral:
		o.env = environment.Enclose(o.env)

		o.addCode("\\(")

		for i, p := range node.Parameters {
			o.Obfuscate(p.Expression)

			if i < len(node.Parameters)-1 {
				o.addCode(",")
			}
		}

		o.addCode("){")
		if node.Body != nil {
			o.Obfuscate(node.Body)
		}

		o.env = environment.Open(o.env)
		o.addCode("}")

	case *ast.CallExpression:
		o.obfuscateCallExpression(node)
	}

	return node
}

func (o *Obfuscator) obfuscateProgram(program *ast.Program) ast.Node {
	var node ast.Node

	for _, statement := range program.Statements {
		o.Obfuscate(statement)
		o.addCode(";")
	}

	return node
}

func (o *Obfuscator) obfuscateCallExpression(node *ast.CallExpression) {
	name := node.Name
	fn, ok := o.env.GetFunction(name, true)

	if ok {
		name = fn.Obfuscated
	}

	o.addCode(name + "(")
	for i, a := range node.Arguments {
		o.Obfuscate(a)
		if i < len(node.Arguments)-1 {
			o.addCode(",")
		}
	}
	o.addCode(")")
}

func (o *Obfuscator) GetCode() string {
	return strings.Join(o.code, "")
}

func (o *Obfuscator) addCode(code string) {
	o.code = append(o.code, code)
}
