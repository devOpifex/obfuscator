package parser

import (
	"fmt"

	"github.com/devOpifex/obfuscator/ast"
	"github.com/devOpifex/obfuscator/diagnostics"
	"github.com/devOpifex/obfuscator/lexer"
	"github.com/devOpifex/obfuscator/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN     // <- and = (lowest precedence)
	TILDE      // ~
	OR         // |
	AND        // &
	UNARY      // ! and unary - and +
	COMPARISON // == >= > < <= !=
	PLUS       // binary + and -
	STAR       // * and /
	PIPE       // %>% and |>
	COLON      // :
	CARET      // ^
	SUBSET     // [] [[]]
	DOLLAR     // $
	NAMESPACE  // :: and :::
	CALL       // ()
	INDEX      // highest precedence
)

var precedences = map[token.ItemType]int{
	// Assignment operators (lowest precedence)
	token.ItemAssign:       ASSIGN,
	token.ItemAssignParent: ASSIGN,
	token.ItemWalrus:       ASSIGN,

	// Tilde
	token.ItemTilde: TILDE,

	// Logical operators
	token.ItemOr:        OR,  // |
	token.ItemDoubleOr:  OR,  // ||
	token.ItemAnd:       AND, // &
	token.ItemDoubleAnd: AND, // &&

	// Unary operators
	token.ItemBang: UNARY, // !

	// Comparison operators
	token.ItemDoubleEqual:    COMPARISON, // ==
	token.ItemNotEqual:       COMPARISON, // !=
	token.ItemLessThan:       COMPARISON, // <
	token.ItemLessOrEqual:    COMPARISON, // <=
	token.ItemGreaterThan:    COMPARISON, // >
	token.ItemGreaterOrEqual: COMPARISON, // >=

	// Arithmetic operators
	token.ItemPlus:     PLUS,  // binary +
	token.ItemMinus:    PLUS,  // binary -
	token.ItemMultiply: STAR,  // *
	token.ItemDivide:   STAR,  // /
	token.ItemCaret:    CARET, // ^

	// Special operators
	token.ItemDollar:            DOLLAR,    // $
	token.ItemNamespace:         NAMESPACE, // ::
	token.ItemNamespaceInternal: NAMESPACE, // :::

	// Subsetting operators
	token.ItemLeftSquare:       SUBSET, // [
	token.ItemDoubleLeftSquare: SUBSET, // [[

	// Function call
	token.ItemLeftParen: CALL, // (

	// Other operators that need specific precedence
	token.ItemPipe:  PIPE, // |>
	token.ItemInfix: STAR, // %op%

	token.ItemColon: COLON, // :
}

type (
	prefixParseFn  func() ast.Expression
	postfixParseFn func() ast.Expression
	infixParseFn   func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors diagnostics.Diagnostics

	pos int

	curToken  token.Item
	peekToken token.Item

	filePos int

	postfixParseFns map[token.ItemType]postfixParseFn
	prefixParseFns  map[token.ItemType]prefixParseFn
	infixParseFns   map[token.ItemType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: diagnostics.Diagnostics{},
	}

	p.prefixParseFns = make(map[token.ItemType]prefixParseFn)
	p.registerPrefix(token.ItemIdent, p.parseIdentifier)
	p.registerPrefix(token.ItemInteger, p.parseIntegerLiteral)
	p.registerPrefix(token.ItemFloat, p.parseFloatLiteral)
	p.registerPrefix(token.ItemBang, p.parsePrefixExpression)
	p.registerPrefix(token.ItemMinus, p.parsePrefixExpression)
	p.registerPrefix(token.ItemBool, p.parseBoolean)
	p.registerPrefix(token.ItemIf, p.parseIfExpression)
	p.registerPrefix(token.ItemFunction, p.parseFunctionLiteral)
	p.registerPrefix(token.ItemBackslash, p.parseFunctionLiteral)
	p.registerPrefix(token.ItemDoubleQuote, p.parseStringLiteral)
	p.registerPrefix(token.ItemSingleQuote, p.parseStringLiteral)
	p.registerPrefix(token.ItemBacktick, p.parseBacktickLiteral)
	p.registerPrefix(token.ItemNA, p.parseNA)
	p.registerPrefix(token.ItemDot, p.parseDot)
	p.registerPrefix(token.ItemDoubleDot, p.parseDoubleDot)
	p.registerPrefix(token.ItemNan, p.parseNan)
	p.registerPrefix(token.ItemNAComplex, p.parseNaComplex)
	p.registerPrefix(token.ItemNAReal, p.parseNaReal)
	p.registerPrefix(token.ItemNAInteger, p.parseNaInteger)
	p.registerPrefix(token.ItemInf, p.parseInf)
	p.registerPrefix(token.ItemNULL, p.parseNull)
	p.registerPrefix(token.ItemThreeDot, p.parseElipsis)
	p.registerPrefix(token.ItemFor, p.parseFor)
	p.registerPrefix(token.ItemWhile, p.parseWhile)
	p.registerPrefix(token.ItemComma, p.parseComma)
	p.registerPrefix(token.ItemRightSquare, p.parsePostfixSquare)
	p.registerPrefix(token.ItemDoubleRightSquare, p.parsePostfixSquare)
	p.registerPrefix(token.ItemLeftParen, p.parseLeftParen)
	p.registerPrefix(token.ItemRightParen, p.parseRightParen)
	p.registerPrefix(token.ItemLeftCurly, p.parseLeftCurly)

	p.infixParseFns = make(map[token.ItemType]infixParseFn)
	p.registerInfix(token.ItemInfix, p.parseInfixExpression)
	p.registerInfix(token.ItemOr, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleOr, p.parseInfixExpression)
	p.registerInfix(token.ItemAnd, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleAnd, p.parseInfixExpression)
	p.registerInfix(token.ItemPlus, p.parseInfixExpression)
	p.registerInfix(token.ItemMinus, p.parseInfixExpression)
	p.registerInfix(token.ItemDivide, p.parseInfixExpression)
	p.registerInfix(token.ItemMultiply, p.parseInfixExpression)
	p.registerInfix(token.ItemAssign, p.parseInfixExpression)
	p.registerInfix(token.ItemAssignParent, p.parseInfixExpression)
	p.registerInfix(token.ItemWalrus, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemNotEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemLessThan, p.parseInfixExpression)
	p.registerInfix(token.ItemGreaterThan, p.parseInfixExpression)
	p.registerInfix(token.ItemGreaterOrEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemPipe, p.parseInfixExpression)
	p.registerInfix(token.ItemDollar, p.parseInfixExpression)
	p.registerInfix(token.ItemColon, p.parseInfixExpression)
	p.registerInfix(token.ItemNamespace, p.parseInfixExpression)
	p.registerInfix(token.ItemNamespaceInternal, p.parseInfixExpression)
	p.registerInfix(token.ItemLeftSquare, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleLeftSquare, p.parseInfixExpression)
	p.registerInfix(token.ItemLeftParen, p.parseCallExpression)

	p.postfixParseFns = make(map[token.ItemType]postfixParseFn)
	p.registerPostfix(token.ItemRightSquare, p.parsePostfixSquare)
	p.registerPostfix(token.ItemDoubleRightSquare, p.parsePostfixSquare)

	return p
}

func (p *Parser) Run() {
	for i := range p.l.Files {
		p.filePos = i
		p.pos = 0

		p.nextToken()
		p.nextToken()

		for !p.curTokenIs(token.ItemEOF) && !p.curTokenIs(token.ItemError) {
			stmt := p.parseStatement()
			if stmt != nil {
				p.l.Files[i].Ast.Statements = append(p.l.Files[i].Ast.Statements, stmt)
			}
			p.nextToken()
		}
	}
}

func (p *Parser) Files() lexer.Files {
	return p.l.Files
}

func (p *Parser) Print() {
	for i := range p.l.Files {
		fmt.Println(p.l.Files[i].Ast.String())
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	if p.pos >= len(p.l.Files[p.filePos].Items) {
		return
	}
	p.peekToken = p.l.Files[p.filePos].Items[p.pos]
	p.pos++
}

func (p *Parser) debug() {
	fmt.Println("++++++++++++++++++++ Current ++++++++++++++++++++")
	fmt.Printf("line: %v - character: %v | ", p.curToken.Line+1, p.curToken.Char+1)
	p.curToken.Print()
	fmt.Println("++++++++++++++++++++ Peek")
	fmt.Printf("line: %v - character: %v | ", p.peekToken.Line+1, p.peekToken.Char+1)
	p.peekToken.Print()
}

func (p *Parser) curTokenIs(t token.ItemType) bool {
	return p.curToken.Class == t
}

func (p *Parser) peekTokenIs(t token.ItemType) bool {
	return p.peekToken.Class == t
}

func (p *Parser) expectPeek(t token.ItemType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectCurrent(t token.ItemType) bool {
	if p.curTokenIs(t) {
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) HasError() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() diagnostics.Diagnostics {
	return p.errors
}

func (p *Parser) peekError(t token.ItemType) {
	// we already got an error on the lexer: use it
	if p.peekToken.Class == token.ItemError {
		return
	}

	msg := fmt.Sprintf(
		"expected next token to be `%v`, got `%v` instead",
		t,
		p.peekToken.Class,
	)

	p.errors = append(
		p.errors,
		diagnostics.NewError(p.curToken, msg),
	)
}

func (p *Parser) noPrefixParseFnError(t token.ItemType) {
	msg := fmt.Sprintf(
		"no prefix parse function for `%v` found",
		t,
	)
	p.errors = append(
		p.errors,
		diagnostics.NewError(p.curToken, msg),
	)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Class {
	case token.ItemComment:
		return p.parseCommentStatement()
	case token.ItemExport:
		return p.parseExportStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFor() ast.Expression {
	lit := &ast.For{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	lit.Name = p.curToken.Value

	if !p.expectPeek(token.ItemIn) {
		return nil
	}

	p.nextToken()

	lit.Vector = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}

	if !p.curTokenIs(token.ItemRightParen) {
		return nil
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Value = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseWhile() ast.Expression {
	lit := &ast.While{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	p.nextToken()

	// Parse the condition as an expression with LOWEST precedence
	expr := p.parseExpression(LOWEST)

	// Create an expression statement to hold the expression
	lit.Statement = &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: expr,
	}

	// Explicitly check for right parenthesis
	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Value = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{
		Token: p.curToken,
		Value: "NULL",
	}
}

func (p *Parser) parseElipsis() ast.Expression {
	return &ast.Keyword{Token: p.curToken, Value: "..."}
}

func (p *Parser) parseDoubleDot() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "..",
	}
}

func (p *Parser) parseDot() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: ".",
	}
}

func (p *Parser) parseNA() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA",
	}
}

func (p *Parser) parseNan() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NaN",
	}
}

func (p *Parser) parseComma() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: ",",
	}
}

func (p *Parser) parseLeftParen() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "(",
	}
}

func (p *Parser) parseRightParen() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: ")",
	}
}

func (p *Parser) parseNaString() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_character_",
	}
}

func (p *Parser) parseNaReal() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_real_",
	}
}

func (p *Parser) parseNaComplex() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_complex_",
	}
}

func (p *Parser) parseNaInteger() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_integer_",
	}
}

func (p *Parser) parseInf() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "Inf",
	}
}

func (p *Parser) parseLeftCurly() ast.Expression {
	var exp ast.ExpressionBlock
	exp.Expression = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Class]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Class)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.ItemEOF) &&
		precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Class]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	postfix := p.postfixParseFns[p.peekToken.Class]

	if postfix == nil {
		return leftExp
	}

	p.nextToken()

	return &ast.PostfixExpression{
		Token:   p.curToken,
		Left:    leftExp,
		Postfix: p.curToken.Value,
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Class]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Class]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	return &ast.FloatLiteral{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curToken.Value == "true" || p.curToken.Value == "TRUE",
	}
}

func (p *Parser) parseCommentStatement() ast.Statement {
	return &ast.CommentStatement{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseExportStatement() ast.Statement {
	return &ast.ExportStatement{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	str := &ast.StringLiteral{
		Token: p.curToken,
	}

	// it's an empty string ""
	if p.peekTokenIs(p.curToken.Class) {
		p.nextToken()
		return str
	}

	p.expectPeek(token.ItemString)

	str.Str = p.curToken.Value

	p.nextToken()

	return str
}

func (p *Parser) parseBacktickLiteral() ast.Expression {
	bt := &ast.BacktickLiteral{
		Token: p.curToken,
	}

	// it's an empty backtick ""
	if p.peekTokenIs(token.ItemBacktick) {
		p.nextToken()
		return bt
	}

	p.expectPeek(token.ItemIdent)

	bt.Value = p.curToken.Value

	p.nextToken()

	return bt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Value,
	}

	p.nextToken()

	expression.Right = p.parseExpression(UNARY)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// it's a function declaration (bit hacky)
	if p.curTokenIs(token.ItemAssign) && (p.peekTokenIs(token.ItemFunction) || p.peekTokenIs(token.ItemBackslash)) {
		return p.parseNamedFunctionLiteral(left)
	}

	operator := p.curToken.Value

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: operator,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// skip paren left (
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ItemElse) {
		p.nextToken()

		if !p.expectPeek(token.ItemLeftCurly) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.ItemRightCurly) && !p.curTokenIs(token.ItemEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionParameters() []*ast.Argument {
	var params []*ast.Argument

	// Handle empty parameter list
	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		return params
	}

	p.nextToken() // move past opening paren

	if p.curTokenIs(token.ItemLeftCurly) {
		params = append(params, &ast.Argument{
			Token: p.curToken,
			Value: p.parseLeftCurly(),
		})
		p.nextToken()
	}

	for !p.curTokenIs(token.ItemRightParen) && !p.curTokenIs(token.ItemEOF) {
		if p.curTokenIs(token.ItemComma) {
			p.nextToken()
			continue
		}

		if p.curTokenIs(token.ItemComment) {
			p.nextToken()
			continue
		}

		arg := &ast.Argument{
			Token: p.curToken,
		}

		// Check for named parameter (identifier followed by =)
		if p.curTokenIs(token.ItemIdent) && p.peekTokenIs(token.ItemAssign) {
			name := p.curToken.Value
			p.nextToken() // move past identifier
			p.nextToken() // move past =
			arg.Name = name
			arg.Value = p.parseExpression(LOWEST)
		} else {
			// Unnamed parameter
			arg.Value = p.parseExpression(LOWEST)
		}

		params = append(params, arg)

		// Break if we're at the end
		if p.peekTokenIs(token.ItemRightParen) {
			p.nextToken()
			break
		}

		p.nextToken() // move past comma
	}

	return params
}

func (p *Parser) parseFunctionArguments() []*ast.Argument {
	var params []*ast.Argument

	// Handle empty parameter list
	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		return params
	}

	p.nextToken() // move past opening paren

	for !p.curTokenIs(token.ItemRightParen) && !p.curTokenIs(token.ItemEOF) {
		if p.curTokenIs(token.ItemComma) {
			p.nextToken()
			continue
		}

		arg := &ast.Argument{
			Token: p.curToken,
		}

		// Check for named parameter (identifier followed by =)
		if p.curTokenIs(token.ItemIdent) && p.peekTokenIs(token.ItemAssign) {
			name := p.curToken.Value
			p.nextToken() // move past identifier
			p.nextToken() // move past =
			arg.Name = name
			arg.Value = p.parseExpression(LOWEST)
		} else {
			// Unnamed parameter
			arg.Name = p.curToken.Value
		}

		params = append(params, arg)

		// Break if we're at the end
		if p.peekTokenIs(token.ItemRightParen) {
			p.nextToken()
			break
		}

		p.nextToken() // move past comma
	}

	return params
}

func (p *Parser) parseNamedFunctionLiteral(name ast.Expression) ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken, Name: name.String()}

	if !(p.peekTokenIs(token.ItemFunction) || p.peekTokenIs(token.ItemBackslash)) {
		p.peekError(token.ItemFunction)
		return nil
	}
	p.nextToken()

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionArguments()

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionArguments()

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseSquare() ast.Expression {
	return &ast.Square{
		Token: p.curToken,
	}
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Name: function.Item().Value}

	exp.Arguments = p.parseFunctionParameters()

	if !p.curTokenIs(token.ItemRightParen) {
		return nil
	}

	return exp
}

func (p *Parser) parsePostfixSquare() ast.Expression {
	return &ast.Square{
		Token: p.curToken,
	}
}

func (p *Parser) registerPrefix(tokenType token.ItemType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.ItemType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.ItemType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}
