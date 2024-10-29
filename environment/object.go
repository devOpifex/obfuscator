package environment

import (
	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/token"
)

// this should be an interface but I haven't got the time right now
type Argument struct {
	Name       string
	Obfuscated string
}

type Function struct {
	Token      token.Item
	Value      *ast.FunctionLiteral
	Name       string
	Obfuscated string
	Arguments  []Argument
}

type Variable struct {
	Token      token.Item
	Name       string
	Obfuscated string
}
