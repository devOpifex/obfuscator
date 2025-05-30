package token

type ItemType int

type Item struct {
	Class ItemType
	Value string
	Line  int
	Pos   int
	Char  int
	File  string
}

type Items []Item

const (
	ItemError ItemType = iota

	// end of file
	ItemEOF

	// identifiers
	ItemIdent

	// quotes
	ItemDoubleQuote
	ItemSingleQuote

	// $
	ItemDollar

	// backtick
	ItemBacktick

	// infix %>%
	ItemInfix

	// comma,
	ItemComma

	// dot .
	ItemDot

	// dot dot ..
	ItemDoubleDot

	// question mark?
	ItemQuestion

	// boolean
	ItemBool

	// boolean
	ItemReturn

	// ...
	ItemThreeDot

	// native pipe
	ItemPipe

	// =
	ItemAssign

	// <<-
	ItemAssignParent

	// :=
	ItemWalrus

	// NULL
	ItemNULL

	// NA
	ItemNA
	ItemNan
	ItemNAString
	ItemNAReal
	ItemNAComplex
	ItemNAInteger

	// parens and brackets
	ItemLeftCurly
	ItemRightCurly
	ItemLeftParen
	ItemRightParen
	ItemLeftSquare
	ItemRightSquare
	ItemDoubleLeftSquare
	ItemDoubleRightSquare

	// "strings"
	ItemString

	// numbers
	ItemInteger
	ItemFloat

	// namespace::
	ItemNamespace
	// namespace:::
	ItemNamespaceInternal

	// colon
	ItemColon

	// + - / * ^
	ItemPlus
	ItemMinus
	ItemDivide
	ItemMultiply
	ItemPower
	ItemModulus

	// bang!
	ItemBang

	// comment
	ItemComment
	ItemExport

	// compare
	ItemDoubleEqual
	ItemLessThan
	ItemGreaterThan
	ItemNotEqual
	ItemLessOrEqual
	ItemGreaterOrEqual

	// if else
	ItemIf
	ItemElse
	ItemAnd
	ItemDoubleAnd
	ItemOr
	ItemDoubleOr
	ItemBreak

	// ^
	ItemCaret

	// ~
	ItemTilde

	// Infinite
	ItemInf

	// loop
	ItemFor
	ItemRepeat
	ItemWhile
	ItemNext
	ItemIn

	// function
	ItemFunction
	ItemBackslash

	// types
	ItemTypes
	ItemTypesPkg
	ItemTypesList
	ItemTypesDecl
)

const EOF = -1
