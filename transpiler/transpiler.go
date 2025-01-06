package transpiler

import (
	"strings"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
	"github.com/sparkle-tech/obfuscator/obfuscator"
	"github.com/sparkle-tech/obfuscator/token"
)

type Transpiler struct {
	code        []string
	env         *environment.Environment
	methodStack obfuscator.Stack
	file        lexer.File
}

type Transpilers []*Transpiler

func New(env *environment.Environment, files lexer.Files) Transpilers {
	var ts Transpilers

	for _, f := range files {
		ts = append(ts, &Transpiler{
			env:  env,
			file: f,
		})
	}

	return ts
}

func (ts Transpilers) Run() {
	for _, t := range ts {
		t.Transpile(t.file.Ast)
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

	case *ast.ExportStatement:
		t.addCode("\n#' @export\n")

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			t.Transpile(s)
			t.addCode(";")
		}

	case *ast.Method:
		t.obfuscateMethod(node)

	case *ast.Identifier:
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

		t.Transpile(node.Left)

		t.addCode(node.Operator)

		if node.Operator == "in" {
			t.addCode(" ")
		}

		t.Transpile(node.Right)

		if node.Operator == "[" {
			t.addCode("]")
		}

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
		for _, p := range node.Parameters {
			if p.Name != "" {
				t.addCode(p.Name + "=")
			}
			t.Transpile(p.Value)
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
		if statement.Item().Class == token.ItemExport {
			continue
		}
		t.addCode(";")
	}

	return node
}

func (t *Transpiler) obfuscateMethod(node *ast.Method) {
	t.addCode(node.Name + "(")
	for _, a := range node.Arguments {
		if a.Name != "" {
			t.addCode(a.Name + "=")
		}
		t.Transpile(a.Value)
	}
	t.addCode(")")
}

func (t *Transpiler) obfuscateCallExpression(node *ast.CallExpression) {
	t.addCode(node.Name + "(")
	for _, a := range node.Arguments {
		if a.Name != "" {
			t.addCode(a.Name + "=")
		}
		t.Transpile(a.Value)
	}
	t.addCode(")")
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
