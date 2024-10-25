package transpiler

import (
	"strings"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
)

type Transpiler struct {
	code []string
	env  *environment.Environment
}

func (t *Transpiler) Env() *environment.Environment {
	return t.env
}

func New() *Transpiler {
	env := environment.New()

	return &Transpiler{
		env: env,
	}
}

func (t *Transpiler) Transpile(node ast.Node) ast.Node {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return t.transpileProgram(node)

	case *ast.ExpressionStatement:
		if node.Expression != nil {
			return t.Transpile(node.Expression)
		}

	case *ast.Comma:
		t.addCode(",")

	case *ast.Null:
		t.addCode("NULL")

	case *ast.Keyword:
		t.addCode(node.Value)

	case *ast.CommentStatement:
		t.addCode("")

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			t.Transpile(s)
		}

	case *ast.Attribute:
		t.addCode(node.Value)
		return node

	case *ast.Identifier:
		t.addCode(node.Value)
		return node

	case *ast.Boolean:
		t.addCode(strings.ToUpper(node.String()))

	case *ast.IntegerLiteral:
		t.addCode(node.Value)

	case *ast.FloatLiteral:
		t.addCode(node.Value)

	case *ast.StringLiteral:
		t.addCode(node.Token.Value + node.Str + node.Token.Value)

	case *ast.PrefixExpression:
		t.addCode("(" + node.Operator)
		t.Transpile(node.Right)
		t.addCode(")")

	case *ast.For:
		t.addCode("for(")
		t.addCode(node.Name)
		t.addCode(" in ")
		t.Transpile(node.Vector)
		t.addCode(") {")
		t.env = environment.Enclose(t.env)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = environment.Open(t.env)

	case *ast.While:
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode(") {")
		t.env = environment.Enclose(t.env)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = environment.Open(t.env)

	case *ast.InfixExpression:
		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Left != nil {
			t.Transpile(node.Left)
		}

		if node.Operator == "in" {
			t.addCode(" ")
		} else {
			t.addCode(node.Operator)
		}

		if node.Right != nil {
			t.Transpile(node.Right)
		}

		if node.Operator == "<-" {
			t.addCode(";")
		}

	case *ast.Square:
		t.addCode(node.Token.Value)

	case *ast.IfExpression:
		t.addCode("if(")
		t.Transpile(node.Condition)
		t.addCode("){")
		t.env = environment.Enclose(t.env)
		t.Transpile(node.Consequence)
		t.env = environment.Open(t.env)
		t.addCode("}")

		if node.Alternative != nil {
			t.addCode(" else {")
			t.env = environment.Enclose(t.env)
			t.Transpile(node.Alternative)
			t.env = environment.Open(t.env)
			t.addCode("}")
		}

	case *ast.FunctionLiteral:
		t.env = environment.Enclose(t.env)

		t.addCode("\\(")

		for i, p := range node.Parameters {
			t.env.SetVariable(
				p.Token.Value,
				environment.Variable{
					Token: p.Token,
				},
			)

			t.addCode(p.Name)

			if p.Operator == "=" {
				t.addCode("=")
				t.Transpile(p.Default)
			}

			if i < len(node.Parameters)-1 {
				t.addCode(",")
			}
		}

		t.addCode("){")
		if node.Body != nil {
			t.Transpile(node.Body)
		}

		t.env = environment.Open(t.env)
		t.addCode("}")

	case *ast.CallExpression:
		t.transpileCallExpression(node)
		t.addCode(";")
	}

	return node
}

func (t *Transpiler) transpileProgram(program *ast.Program) ast.Node {
	var node ast.Node

	for _, statement := range program.Statements {
		t.Transpile(statement)
	}

	return node
}

func (t *Transpiler) transpileCallExpression(node *ast.CallExpression) {
	t.addCode(node.Name + "(")
	for i, a := range node.Arguments {
		t.Transpile(a.Value)
		if i < len(node.Arguments)-1 {
			t.addCode(",")
		}
	}
	t.addCode(")")
}

func (t *Transpiler) GetCode() string {
	return strings.Join(t.code, "")
}

func (t *Transpiler) addCode(code string) {
	t.code = append(t.code, code)
}
