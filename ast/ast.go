package ast

import (
	"bytes"

	"github.com/sparkle-tech/obfuscator/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	Item() token.Item
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Item() token.Item { return token.Item{} }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type CommentStatement struct {
	Token token.Item
	Value string
}

func (c *CommentStatement) Item() token.Item     { return c.Token }
func (c *CommentStatement) statementNode()       {}
func (c *CommentStatement) TokenLiteral() string { return c.Token.Value }
func (c *CommentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + "\n")

	return out.String()
}

type ExportStatement struct {
	Token token.Item
	Value string
}

func (e *ExportStatement) Item() token.Item     { return e.Token }
func (e *ExportStatement) statementNode()       {}
func (e *ExportStatement) TokenLiteral() string { return e.Token.Value }
func (e *ExportStatement) String() string {
	var out bytes.Buffer

	out.WriteString("#' @export")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Item // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) Item() token.Item     { return es.Token }
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Value }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Item // the { token
	Statements []Statement
}

func (bs *BlockStatement) Item() token.Item     { return bs.Token }
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString(";")
	}

	return out.String()
}

// Expressions
type Comma struct {
	Token token.Item
}

func (c *Comma) Item() token.Item     { return c.Token }
func (c *Comma) expressionNode()      {}
func (c *Comma) TokenLiteral() string { return c.Token.Value }
func (c *Comma) String() string {
	return ","
}

type Identifier struct {
	Token token.Item // the token.IDENT token
	Value string
}

func (i *Identifier) Item() token.Item     { return i.Token }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) String() string {
	return i.Value
}

type Attribute struct {
	Token token.Item
	Value string
}

func (a *Attribute) Item() token.Item     { return a.Token }
func (a *Attribute) expressionNode()      {}
func (a *Attribute) TokenLiteral() string { return a.Token.Value }
func (a *Attribute) String() string {
	return a.Value
}

type Square struct {
	Token token.Item
}

func (s *Square) Item() token.Item     { return s.Token }
func (s *Square) expressionNode()      {}
func (s *Square) TokenLiteral() string { return s.Token.Value }
func (s *Square) String() string {
	return s.Token.Value
}

type Boolean struct {
	Token token.Item
	Value bool
}

func (b *Boolean) Item() token.Item     { return b.Token }
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Value }
func (b *Boolean) String() string {
	if b.Value {
		return "TRUE"
	}

	return "FALSE"
}

type IntegerLiteral struct {
	Token token.Item
	Value string
}

func (il *IntegerLiteral) Item() token.Item     { return il.Token }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }

type FloatLiteral struct {
	Token token.Item
	Value string
}

func (fl *FloatLiteral) Item() token.Item     { return fl.Token }
func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FloatLiteral) String() string       { return fl.Token.Value }

type SquareRightLiteral struct {
	Token token.Item
	Value string
}

func (s *SquareRightLiteral) Item() token.Item     { return s.Token }
func (s *SquareRightLiteral) expressionNode()      {}
func (s *SquareRightLiteral) TokenLiteral() string { return s.Token.Value }
func (s *SquareRightLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(s.Value)
	out.WriteString("\n")

	return out.String()
}

type For struct {
	Token  token.Item
	Name   string
	Vector Expression
	Value  *BlockStatement
}

func (f *For) Item() token.Item     { return f.Token }
func (f *For) expressionNode()      {}
func (f *For) TokenLiteral() string { return f.Token.Value }
func (f *For) String() string {
	var out bytes.Buffer

	out.WriteString("for(")
	out.WriteString(f.Name)
	out.WriteString(" in ")
	out.WriteString(f.Vector.String())
	out.WriteString(")\n {")
	out.WriteString(f.Value.String())
	out.WriteString("}\n")

	return out.String()
}

type While struct {
	Token     token.Item
	Statement Statement
	Value     *BlockStatement
}

func (w *While) Item() token.Item     { return w.Token }
func (w *While) expressionNode()      {}
func (w *While) TokenLiteral() string { return w.Token.Value }
func (w *While) String() string {
	var out bytes.Buffer

	out.WriteString("while(")
	out.WriteString(w.Statement.String())
	out.WriteString(")\n {")
	out.WriteString(w.Value.String())
	out.WriteString("}\n")

	return out.String()
}

type Null struct {
	Token token.Item
	Value string
}

func (n *Null) Item() token.Item     { return n.Token }
func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return n.Token.Value }
func (n *Null) String() string {
	var out bytes.Buffer

	out.WriteString(n.Value)

	return out.String()
}

type Keyword struct {
	Token token.Item
	Value string
}

func (kw *Keyword) Item() token.Item     { return kw.Token }
func (kw *Keyword) expressionNode()      {}
func (kw *Keyword) TokenLiteral() string { return kw.Token.Value }
func (kw *Keyword) String() string {
	var out bytes.Buffer

	out.WriteString(kw.Value)

	return out.String()
}

type StringLiteral struct {
	Token token.Item
	Str   string
}

func (sl *StringLiteral) Item() token.Item     { return sl.Token }
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(sl.Token.Value + sl.Str + sl.Token.Value)

	return out.String()
}

type BacktickLiteral struct {
	Token token.Item
	Value string
}

func (b *BacktickLiteral) Item() token.Item     { return b.Token }
func (b *BacktickLiteral) expressionNode()      {}
func (b *BacktickLiteral) TokenLiteral() string { return b.Token.Value }
func (b *BacktickLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("`" + b.Value + "`")

	return out.String()
}

type PrefixExpression struct {
	Token    token.Item // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Item() token.Item     { return pe.Token }
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Value }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type PostfixExpression struct {
	Token   token.Item // The prefix token, e.g. !
	Left    Expression
	Postfix string
}

func (pe *PostfixExpression) Item() token.Item     { return pe.Token }
func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Value }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Postfix)
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Item // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) Item() token.Item     { return ie.Token }
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String())
	if ie.Operator == "<-" {
		out.WriteString("=")
	}

	out.WriteString(ie.Operator)

	if ie.Right != nil {
		out.WriteString(ie.Right.String())
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Item // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) Item() token.Item     { return ie.Token }
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if(")
	out.WriteString(ie.Condition.String())
	out.WriteString("){")
	out.WriteString(ie.Consequence.String())
	out.WriteString("}")

	if ie.Alternative != nil {
		out.WriteString("else{")
		out.WriteString(ie.Alternative.String())
		out.WriteString("}")
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Item // The 'function' token
	Name       string
	Parameters []*Argument
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Item() token.Item     { return fl.Token }
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	if fl.Name != "" {
		out.WriteString(fl.Name + "=")
	}

	out.WriteString("\\")
	out.WriteString("(")
	for _, a := range fl.Parameters {
		out.WriteString(a.Name)
	}
	out.WriteString("){")
	out.WriteString(fl.Body.String())
	out.WriteString("}")

	return out.String()
}

type Parameter struct {
	Token    token.Item // The 'func' token
	Name     string
	Operator string
	Default  Expression
	Method   bool
}

func (p *Parameter) Item() token.Item     { return p.Token }
func (p *Parameter) expressionNode()      {}
func (p *Parameter) TokenLiteral() string { return p.Token.Value }
func (p *Parameter) String() string {
	var out bytes.Buffer

	out.WriteString(p.Name)
	if p.Operator != "" {
		out.WriteString(" " + p.Operator + " " + p.Default.String())
	}
	return out.String()
}

type Argument struct {
	Token token.Item
	Name  string
	Value Expression
}

type CallExpression struct {
	Token     token.Item // The '(' token
	Name      string
	Arguments []*Argument
}

func (ce *CallExpression) Item() token.Item     { return ce.Token }
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Name)
	out.WriteString("(")
	for _, a := range ce.Arguments {
		out.WriteString(a.Name)
	}
	out.WriteString(")")

	return out.String()
}
