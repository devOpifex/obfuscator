package transpiler

import (
	"regexp"
	"strings"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
	"github.com/sparkle-tech/obfuscator/token"
)

type Transpiler struct {
	code          []string
	env           *environment.Environment
	file          lexer.File
	boxUse        bool
	lastNamespace string
}

var startWithDot = regexp.MustCompile(`^\.`)

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
			if _, ok := s.(*ast.CommentStatement); ok {
				continue
			}
			t.addCode(";")
		}

	case *ast.Identifier:
		if t.inBoxUse() {
			if t.env.GetPath(node.Value) {
				t.addCode(environment.Mask(node.Value))
				return node
			}

			if t.env.GetFunction(node.Value) {
				t.addCode(environment.Mask(node.Value))
				return node
			}

			if t.env.GetVariable(node.Value, true) {
				t.addCode(environment.Mask(node.Value))
				return node
			}

			t.addCode(node.Value)
			return node
		}

		if t.env.GetVariable(node.Value, true) {
			t.addCode(environment.Mask(node.Value))
			return node
		}

		if t.env.GetFunction(node.Value) {
			t.addCode(environment.Mask(node.Value))
			return node
		}

		t.addCode(node.Value)

		return node

	case *ast.Boolean:
		if node.Value {
			t.addCode("T")
			return node
		}

		t.addCode("F")

		return node

	case *ast.IntegerLiteral:
		t.addCode(node.Value)

	case *ast.FloatLiteral:
		t.addCode(node.Value)

	case *ast.StringLiteral:
		t.addCode(node.Token.Value + node.Str + node.Token.Value)

	case *ast.BacktickLiteral:
		t.addCode("`" + node.Value + "`")

	case *ast.PrefixExpression:
		t.addCode(node.Operator)
		t.Transpile(node.Right)

	case *ast.For:
		t.addCode("for(")
		t.env = environment.Enclose(t.env)
		t.env.SetVariable(node.Name)
		t.addCode(environment.Mask(node.Name))
		t.addCode(" in ")
		t.Transpile(node.Vector)
		t.addCode("){")
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = environment.Open(t.env)

	case *ast.While:
		t.env = environment.Enclose(t.env)
		t.addCode("while(")
		t.Transpile(node.Statement)
		t.addCode("){")
		t.Transpile(node.Value)
		t.addCode("}")
		t.env = environment.Open(t.env)

	case *ast.InfixExpression:
		if node.Operator == "<-" {
			node.Operator = "="
		}

		if node.Operator == "in" {
			node.Operator = " in "
		}

		// it's a pipe e.g.: %>%
		if strings.Contains(node.Operator, "%") {
			node.Operator = " " + node.Operator + " "
		}

		if _, ok := node.Left.(*ast.Identifier); ok && node.Operator == "=" {
			t.env.SetVariable(node.Left.String())
		}

		if _, ok := node.Left.(*ast.Identifier); ok && node.Operator == "::" {
			t.lastNamespace = node.Left.String()
		}

		t.Transpile(node.Left)

		t.addCode(node.Operator)
		t.Transpile(node.Right)
		return node.Right

	case *ast.Square:
		t.addCode(node.Token.Value)

	case *ast.PostfixExpression:
		t.Transpile(node.Left)
		t.addCode(node.Postfix)

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
		// these 2 branches are for fecking InfixExpression
		// `%>%` <- function(lhs, rhs) { #... }
		if node.Name != "" && !strings.Contains(node.Name, "%") {
			t.transpileFunctionName(node)
		}

		if strings.Contains(node.Name, "%") {
			t.addCode(node.Name + "=")
		}

		t.env = environment.Enclose(t.env)
		t.addCode("\\(")
		for i, p := range node.Parameters {
			if p.Name != "" && p.Name != "..." {
				t.env.SetVariable(p.Name)
				t.addCode(environment.Mask(p.Name))
			}
			if p.Name == "..." {
				t.addCode("...")
			}
			if p.Value != nil {
				t.addCode("=")
				t.Transpile(p.Value)
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

func (t *Transpiler) obfuscateCallExpression(node *ast.CallExpression) {
	// calls to box::use
	if t.lastNamespace == "box" && node.Name == "use" {
		t.setBoxUse()
	}

	ok := t.env.GetFunction(node.Name)
	if ok {
		t.addCode(environment.Mask(node.Name) + "(")
	}

	if !ok {
		t.addCode(node.Name + "(")
	}

	for i, a := range node.Arguments {
		if a.Name != "" && ok {
			t.addCode(environment.Mask(a.Name) + "=")
		}

		if a.Name != "" && !ok {
			t.addCode(a.Name + "=")
		}

		if a.Value != nil {
			t.Transpile(a.Value)
			if i < len(node.Arguments)-1 {
				t.addCode(",")
			}
		}
	}
	t.unsetBoxUse()
	t.addCode(")")
}

func (t *Transpiler) GetCode() string {
	return t.cleanCode()
}

// this is bad but I have no other fix right now
func (t *Transpiler) cleanCode() string {
	code := strings.Join(t.code, "")
	code = strings.ReplaceAll(code, "(;", "(")
	code = strings.ReplaceAll(code, ";)", ")")
	code = strings.ReplaceAll(code, ";,", ",")
	code = strings.ReplaceAll(code, ",;", ",")
	code = strings.ReplaceAll(code, "(,", "(")
	return code
}

func (t *Transpiler) addCode(code string) {
	t.code = append(t.code, code)
}

func (t *Transpiler) setBoxUse() {
	t.boxUse = true
}

func (t *Transpiler) unsetBoxUse() {
	t.boxUse = false
}

func (t *Transpiler) inBoxUse() bool {
	return t.boxUse
}

func (t *Transpiler) transpileFunctionName(node *ast.FunctionLiteral) {
	// we don't obfuscate function names that start with a dot, e.g.: .onLoad
	if startWithDot.MatchString(node.Name) {
		t.addCode(node.Name)
		return
	}

	// split on . to identify methods
	split := strings.Split(node.Name, ".")

	if len(split) < 2 {
		t.env.SetFunction(node.Name)
		t.addCode(environment.Mask(node.Name) + "=")
		return
	}

	// it's not a method (no dot)
	if t.env.GetGeneric(split[0]) {
		t.env.SetFunction(node.Name)
		t.addCode(environment.Mask(node.Name) + "=")
		return
	}

	// it's a method (we find an existing function before first dot)
	if ok := t.env.GetFunction(split[0]); ok {
		rest := strings.Join(split[1:], ".")
		t.addCode(environment.Mask(split[0]) + "." + rest)
		return
	}

	t.env.SetFunction(node.Name)
	t.addCode(environment.Mask(node.Name) + "=")
}
