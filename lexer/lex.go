package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/diagnostics"
	"github.com/sparkle-tech/obfuscator/token"
)

type File struct {
	Path    string
	Content []byte
	Items   token.Items
	Ast     *ast.Program
}

type Files []File

type Lexer struct {
	Files   Files
	filePos int
	input   string
	start   int
	pos     int
	width   int
	line    int // line number
	char    int // character number in line
	errors  diagnostics.Diagnostics
}

const stringNumber = "0123456789"
const stringAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const stringAlphaNum = stringAlpha + stringNumber
const stringMathOp = "+-*/^"

var exported regexp.Regexp = *regexp.MustCompile("\\@export")

func New(fl Files) *Lexer {
	return &Lexer{
		Files: fl,
	}
}

func NewCode(fl, code string) *Lexer {
	return New(
		Files{
			{
				Path:    fl,
				Content: []byte(code),
				Ast: &ast.Program{
					Statements: []ast.Statement{},
				},
			},
		},
	)
}

func NewTest(code string) *Lexer {
	return New(
		Files{
			{
				Path:    "test.vp",
				Content: []byte(code),
				Ast: &ast.Program{
					Statements: []ast.Statement{},
				},
			},
		},
	)
}

func (l *Lexer) HasError() bool {
	return len(l.errors) > 0
}

func (l *Lexer) Errors() diagnostics.Diagnostics {
	return l.errors
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	err := token.Item{
		Char:  l.char,
		Pos:   l.pos,
		Line:  l.line,
		Class: token.ItemError,
		Value: fmt.Sprintf(format, args...),
		File:  l.Files[l.filePos].Path,
	}
	l.errors = append(l.errors, diagnostics.NewError(err, err.Value))
	return nil
}

func (l *Lexer) emit(t token.ItemType) {
	// skip empty tokens
	if l.start == l.pos {
		return
	}

	l.Files[l.filePos].Items = append(l.Files[l.filePos].Items, token.Item{
		Char:  l.char,
		Line:  l.line,
		Pos:   l.pos,
		Class: t,
		Value: l.input[l.start:l.pos],
		File:  l.Files[l.filePos].Path,
	})
	l.start = l.pos
}

func (l *Lexer) emitEOF() {
	l.Files[l.filePos].Items = append(l.Files[l.filePos].Items, token.Item{Class: token.ItemEOF, Value: "EOF"})
}

// returns currently accepted token
func (l *Lexer) token() string {
	return l.input[l.start:l.pos]
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return token.EOF
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	l.char += l.width
	return r
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
	l.char -= l.width
}

func (l *Lexer) peek(n int) rune {
	var r rune
	for i := 0; i < n; i++ {
		r = l.next()
	}

	for i := 0; i < n; i++ {
		l.backup()
	}

	return r
}

type stateFn func(*Lexer) stateFn

func (l *Lexer) Run() {
	for i, f := range l.Files {
		l.filePos = i
		l.input = string(f.Content) + "\n"
		l.width = 0
		l.pos = 0
		l.start = 0
		l.line = 0
		l.char = 0
		l.Lex()

		// remove the EOF
		if i < len(l.Files)-1 {
			l.Files[l.filePos].Items = l.Files[l.filePos].Items[:len(l.Files[l.filePos].Items)-1]
		}

		if len(l.Files[l.filePos].Items) > 0 && l.Files[l.filePos].Items[len(l.Files[l.filePos].Items)-1].Class != token.ItemEOF {
			l.emitEOF()
		}
	}

	for i := range l.Files[l.filePos].Items {
		if l.Files[l.filePos].Items[i].Class == token.ItemDollar && l.Files[l.filePos].Items[i+2].Class == token.ItemLeftParen {
			l.Files[l.filePos].Items[i+1].Class = token.ItemMethod
		}
	}
}

func (l *Lexer) Lex() {
	for state := lexDefault; state != nil; {
		state = state(l)
	}
}

func lexDefault(l *Lexer) stateFn {
	r1 := l.peek(1)

	if r1 == token.EOF {
		l.emitEOF()
		return nil
	}

	if r1 == '"' {
		l.next()
		l.emit(token.ItemDoubleQuote)
		return l.lexString('"')
	}

	if r1 == '`' {
		return lexBacktick(l)
	}

	if r1 == '\'' {
		l.next()
		l.emit(token.ItemSingleQuote)
		return l.lexString('\'')
	}

	if r1 == '#' {
		return lexComment
	}

	if r1 == '\\' {
		l.next()
		l.emit(token.ItemBackslash)
		return lexDefault
	}

	// we parsed strings: we skip spaces and tabs
	if r1 == ' ' || r1 == '\t' {
		l.next()
		l.ignore()
		return lexDefault
	}

	if r1 == '\n' || r1 == '\r' {
		l.next()
		l.ignore()
		l.line++
		l.char = 0
		return lexDefault
	}

	// peek one more rune
	r2 := l.peek(2)

	if r1 == '[' && r2 == '[' {
		l.next()
		l.next()
		l.emit(token.ItemDoubleLeftSquare)
		return lexDefault
	}

	if r1 == ']' && r2 == ']' {
		l.next()
		l.next()
		l.emit(token.ItemDoubleRightSquare)
		return lexDefault
	}

	if r1 == '[' {
		l.next()
		l.emit(token.ItemLeftSquare)
		return lexDefault
	}

	if r1 == ']' {
		l.next()
		l.emit(token.ItemRightSquare)
		return lexDefault
	}

	if r1 == '.' && r2 == '.' && l.peek(3) == '.' {
		l.next()
		l.next()
		l.next()
		l.emit(token.ItemThreeDot)
		return lexDefault
	}

	if r1 == '.' {
		l.next()
		l.emit(token.ItemDot)
		return lexDefault
	}

	// if it's not %% it's an infix
	if r1 == '%' && r2 != '%' {
		return lexInfix
	}

	// it's a modulus
	if r1 == '%' && r2 == '%' {
		l.next()
		l.next()
		l.emit(token.ItemModulus)
		return lexDefault
	}

	if r1 == '=' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemDoubleEqual)
		return lexDefault
	}

	if r1 == '!' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemNotEqual)
		return lexDefault
	}

	if r1 == '!' {
		l.next()
		l.emit(token.ItemBang)
		return lexDefault
	}

	if r1 == '>' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemGreaterOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemLessOrEqual)
		return lexDefault
	}

	if r1 == '<' {
		l.next()
		l.emit(token.ItemLessThan)
		return lexDefault
	}

	if r1 == '>' {
		l.next()
		l.emit(token.ItemGreaterThan)
		return lexDefault
	}

	if r1 == '<' && r2 == '-' {
		l.next()
		l.next()
		l.emit(token.ItemAssign)
		return lexDefault
	}

	if r1 == ':' && r2 == ':' && l.peek(3) == ':' {
		l.next()
		l.next()
		l.next()
		l.emit(token.ItemNamespaceInternal)
		return lexIdentifier
	}

	if r1 == ':' && r2 == ':' {
		l.next()
		l.next()
		l.emit(token.ItemNamespace)
		return lexIdentifier
	}

	// we also emit namespace:: (above)
	// so we can assume this is not
	if r1 == ':' {
		l.next()
		l.emit(token.ItemColon)
		return lexDefault
	}

	if r1 == '&' {
		l.next()
		l.emit(token.ItemAnd)
		return lexDefault
	}

	if r1 == '|' && r2 == '>' {
		l.next()
		l.next()
		l.emit(token.ItemPipe)
		return lexDefault
	}

	if r1 == '|' {
		l.next()
		l.emit(token.ItemOr)
		return lexDefault
	}

	if r1 == '$' {
		l.next()
		l.emit(token.ItemDollar)
		return lexDefault
	}

	if r1 == '@' {
		l.next()
		l.emit(token.ItemAt)
		return lexDefault
	}

	if r1 == ',' {
		l.next()
		l.emit(token.ItemComma)
		return lexDefault
	}

	if r1 == '=' {
		l.next()
		l.emit(token.ItemAssign)
		return lexDefault
	}

	if r1 == '(' {
		l.next()
		l.emit(token.ItemLeftParen)
		return lexDefault
	}

	if r1 == ')' {
		l.next()
		l.emit(token.ItemRightParen)
		return lexIdentifier
	}

	if r1 == '{' {
		l.next()
		l.emit(token.ItemLeftCurly)
		return lexDefault
	}

	if r1 == '}' {
		l.next()
		l.emit(token.ItemRightCurly)
		return lexDefault
	}

	if r1 == '[' && r2 == '[' {
		l.next()
		l.emit(token.ItemDoubleLeftSquare)
		return lexDefault
	}

	if r1 == '[' {
		l.next()
		l.emit(token.ItemLeftSquare)
		return lexDefault
	}

	if r1 == ']' && r2 == ']' {
		l.next()
		l.emit(token.ItemDoubleRightSquare)
		return lexDefault
	}

	if r1 == ']' {
		l.next()
		l.emit(token.ItemRightSquare)
		return lexDefault
	}

	if r1 == '?' {
		l.next()
		l.emit(token.ItemQuestion)
		return lexDefault
	}

	if l.acceptNumber() {
		return lexNumber
	}

	if l.acceptMathOp() {
		return lexMathOp
	}

	if l.acceptAlphaNumeric() {
		return lexIdentifier
	}

	l.next()
	return lexDefault
}

func lexMathOp(l *Lexer) stateFn {
	l.acceptRun(stringMathOp)

	tk := l.token()

	if tk == "+" {
		l.emit(token.ItemPlus)
	}

	if tk == "-" {
		l.emit(token.ItemMinus)
	}

	if tk == "*" {
		l.emit(token.ItemMultiply)
	}

	if tk == "/" {
		l.emit(token.ItemDivide)
	}

	if tk == "^" {
		l.emit(token.ItemPower)
	}

	return lexDefault
}

func lexNumber(l *Lexer) stateFn {
	l.acceptRun(stringNumber)

	r1 := l.peek(1)

	if r1 == 'e' || r1 == 'E' {
		l.next()
		if r2 := l.peek(1); r2 == '+' || r2 == '-' {
			l.next()
		}
	}

	if l.accept(".") {
		l.acceptRun(stringNumber)
		l.emit(token.ItemFloat)
		return lexDefault
	}

	if l.peek(1) == 'L' {
		l.next()
		l.emit(token.ItemInteger)
		return lexDefault
	}

	l.emit(token.ItemInteger)
	return lexDefault
}

func lexComment(l *Lexer) stateFn {
	r := l.peek(1)
	for r != '\n' && r != token.EOF {
		l.next()
		r = l.peek(1)
	}

	if exported.Match([]byte(l.token())) {
		l.emit(token.ItemExport)
		return lexDefault
	}

	l.emit(token.ItemComment)

	return lexDefault
}

func lexBacktick(l *Lexer) func(l *Lexer) stateFn {
	l.next()
	r := l.peek(1)
	for r != '`' && r != token.EOF {
		r = l.next()
	}

	if r == token.EOF {
		l.next()
		return l.errorf("expecting closing backtick, got %v", l.token())
	}

	l.emit(token.ItemIdent)

	return lexDefault
}

func (l *Lexer) lexString(closing rune) func(l *Lexer) stateFn {
	return func(l *Lexer) stateFn {
		var c rune
		r := l.peek(1)
		for r != closing && r != token.EOF {
			c = l.next()
			r = l.peek(1)
		}

		// this means the closing is escaped so
		// it's not in fact closing:
		// we move the cursor and keep parsing string
		// e.g.: "hello \"world\""
		if c == '\\' && r == closing {
			l.next()
			return l.lexString(closing)
		}

		if r == token.EOF {
			l.next()
			return l.errorf("expecting closing quote, got %v", l.token())
		}

		l.emit(token.ItemString)

		r = l.next()

		if r == '"' {
			l.emit(token.ItemDoubleQuote)
		}

		if r == '\'' {
			l.emit(token.ItemSingleQuote)
		}

		return lexDefault
	}
}

func lexInfix(l *Lexer) stateFn {
	l.next()
	r := l.peek(1)
	for r != '%' && r != token.EOF {
		l.next()
		r = l.peek(1)
	}

	if r == token.EOF {
		l.next()
		return l.errorf("expecting closing %%, got %v", l.token())
	}

	l.next()

	l.emit(token.ItemInfix)

	return lexDefault
}

func lexIdentifier(l *Lexer) stateFn {
	l.acceptRun(stringAlphaNum + "_.")

	if l.peek(1) == '.' && l.peek(2) != '.' {
		l.acceptRun(stringAlphaNum + "_")
	}

	tk := l.token()

	if tk == "TRUE" || tk == "FALSE" {
		l.emit(token.ItemBool)
		return lexDefault
	}

	if tk == "if" {
		l.emit(token.ItemIf)
		return lexDefault
	}

	if tk == "else" {
		l.emit(token.ItemElse)
		return lexDefault
	}

	if tk == "NULL" {
		l.emit(token.ItemNULL)
		return lexDefault
	}

	if tk == "NA" {
		l.emit(token.ItemNA)
		return lexDefault
	}

	if tk == "Inf" {
		l.emit(token.ItemInf)
		return lexDefault
	}

	if tk == "while" {
		l.emit(token.ItemWhile)
		return lexDefault
	}

	if tk == "for" {
		l.emit(token.ItemFor)
		return lexFor
	}

	if tk == "function" {
		l.emit(token.ItemFunction)
		return lexDefault
	}

	if tk == "NaN" {
		l.emit(token.ItemNan)
		return lexDefault
	}

	if tk == "in" {
		l.emit(token.ItemIn)
		return lexDefault
	}

	if tk == "break" {
		l.emit(token.ItemBreak)
		return lexDefault
	}

	if tk == "next" {
		l.emit(token.ItemNext)
		return lexDefault
	}

	if tk == "repeat" {
		l.emit(token.ItemRepeat)
		return lexDefault
	}

	l.emit(token.ItemIdent)
	return lexDefault
}

func lexFor(l *Lexer) stateFn {
	r := l.peek(1)
	if r == ' ' {
		l.next()
		l.ignore()
	}

	if r == '\t' {
		l.next()
		l.ignore()
	}

	r = l.peek(1)

	if r != '(' {
		l.errorf("expecting `(`, got `%c`", r)
		return lexDefault
	}

	l.next()
	l.emit(token.ItemLeftParen)

	return lexIdentifier
}

func (l *Lexer) acceptNumber() bool {
	return l.accept(stringNumber)
}

func (l *Lexer) acceptMathOp() bool {
	return l.accept(stringMathOp)
}

func (l *Lexer) acceptAlphaNumeric() bool {
	return l.accept(stringAlphaNum)
}

func (l *Lexer) accept(rs string) bool {
	for strings.IndexRune(rs, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}
