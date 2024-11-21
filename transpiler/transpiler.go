package transpiler

import (
	"strings"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/obfuscator"
	"github.com/sparkle-tech/obfuscator/token"
)

type Transpiler struct {
	code        []string
	env         *environment.Environment
	callStack   obfuscator.Stack
	methodStack obfuscator.Stack
}

func New(env *environment.Environment) *Transpiler {
	return &Transpiler{
		env: env,
	}
}

func (t *Transpiler) Transpile(node ast.Node) ast.Node {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return t.obfuscateProgram(node)

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
			t.addCode(";")
		}

	case *ast.Method:
		t.obfuscateMethod(node)

	case *ast.Identifier:
		v, ok := t.env.GetVariable(node.Value, true)

		if ok {
			t.addCode(v.Obfuscated)
			return node
		}

		fn, ok := t.env.GetFunction(node.Value, true)

		if ok {
			t.addCode(fn.Obfuscated)
			return node
		}

		t.addCode(node.Value)
		return node

	case *ast.Boolean:
		if node.Value {
			t.addCode("T")
		} else {
			t.addCode("F")
		}

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
		t.addCode("){")
		t.env = environment.Enclose(t.env)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = environment.Open(t.env)

	case *ast.While:
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode("){")
		t.env = environment.Enclose(t.env)
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = environment.Open(t.env)

	case *ast.InfixExpression:
		if node.Operator == "in" {
			t.addCode(" ")
		}

		if node.Operator == "<-" {
			node.Operator = "="
		}

		switch node.Right.(type) {
		case *ast.FunctionLiteral:
			t.env.SetFunction(node.Left.Item().Value, environment.Function{
				Name: node.Left.Item().Value,
			})
		}

		switch l := node.Left.(type) {
		case *ast.Identifier:
			if node.Operator != "=" {
				break
			}

			if t.inMethod() {
				break
			}

			t.env.SetVariable(l.Value, environment.Variable{
				Name: l.Value,
			})
		}

		t.Transpile(node.Left)

		t.addCode(node.Operator)

		if node.Operator == "in" {
			t.addCode(" ")
		}

		t.Transpile(node.Right)

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
			t.addCode("else{")
			t.env = environment.Enclose(t.env)
			t.Transpile(node.Alternative)
			t.env = environment.Open(t.env)
			t.addCode("}")
		}

	case *ast.FunctionLiteral:
		t.env = environment.Enclose(t.env)

		t.addCode("\\(")

		for i, p := range node.Parameters {
			t.Transpile(p.Expression)

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
		t.obfuscateCallExpression(node)
	}

	return node
}

func (t *Transpiler) obfuscateProgram(program *ast.Program) ast.Node {
	var node ast.Node

	for _, statement := range program.Statements {
		if statement.Item().Class == token.ItemComment {
			continue
		}
		t.Transpile(statement)
		t.addCode(";")
	}

	return node
}

func (t *Transpiler) obfuscateMethod(node *ast.Method) {
	t.methodStack = t.methodStack.Push(node.Name, false)
	name := node.Name
	fn, ok := t.env.GetFunction(name, true)

	if ok {
		name = fn.Obfuscated
	}

	t.addCode(name + "(")
	for i, a := range node.Arguments {
		t.Transpile(a)
		if i < len(node.Arguments)-1 {
			t.addCode(",")
		}
	}
	t.addCode(")")
	t.methodStack = t.methodStack.Pop()
}

func (t *Transpiler) obfuscateCallExpression(node *ast.CallExpression) {
	t.callStack = t.callStack.Push(node.Name, true)
	name := node.Name
	fn, ok := t.env.GetFunction(name, true)

	if ok {
		name = fn.Obfuscated
	}

	t.addCode(name + "(")
	for i, a := range node.Arguments {
		t.Transpile(a)
		if i < len(node.Arguments)-1 {
			t.addCode(",")
		}
	}
	t.addCode(")")
	t.callStack = t.callStack.Pop()
}

func (t *Transpiler) GetCode() string {
	return strings.Join(t.code, "")
}

func (t *Transpiler) addCode(code string) {
	t.code = append(t.code, code)
}

func (t *Transpiler) inMethod() bool {
	return len(t.methodStack) > 0
}
