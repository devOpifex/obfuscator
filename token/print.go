package token

import (
	"fmt"
)

var ItemName = map[ItemType]string{
	ItemError:             "error",
	ItemExport:            "export",
	ItemIdent:             "identifier",
	ItemDoubleQuote:       "double quote",
	ItemSingleQuote:       "single quote",
	ItemAssign:            "assign",
	ItemWalrus:            "walrus",
	ItemLeftCurly:         "curly left",
	ItemRightCurly:        "curly right",
	ItemLeftParen:         "paren left",
	ItemRightParen:        "paren right",
	ItemLeftSquare:        "square left",
	ItemRightSquare:       "square right",
	ItemString:            "string",
	ItemInteger:           "integer",
	ItemFloat:             "float",
	ItemNamespace:         "namespace",
	ItemNamespaceInternal: "namespace internal",
	ItemComment:           "comment",
	ItemDoubleEqual:       "double equal",
	ItemLessThan:          "less than",
	ItemGreaterThan:       "greater than",
	ItemNotEqual:          "not equal",
	ItemLessOrEqual:       "less or equal",
	ItemGreaterOrEqual:    "greater or equal",
	ItemBool:              "boolean",
	ItemDollar:            "dollar sign",
	ItemComma:             "comma",
	ItemColon:             "colon",
	ItemQuestion:          "question mark",
	ItemBacktick:          "backtick",
	ItemInfix:             "infix",
	ItemIf:                "if",
	ItemBreak:             "break",
	ItemElse:              "else",
	ItemAnd:               "ampersand",
	ItemDoubleAnd:         "double ampersand",
	ItemOr:                "vertical bar",
	ItemDoubleOr:          "double or",
	ItemCaret:             "caret",
	ItemTilde:             "tilde",
	ItemReturn:            "return",
	ItemNULL:              "null",
	ItemNA:                "NA",
	ItemNan:               "NaN",
	ItemNAString:          "NA string",
	ItemNAReal:            "NA real",
	ItemNAComplex:         "NA complex",
	ItemNAInteger:         "NA integer",
	ItemPipe:              "native pipe",
	ItemModulus:           "modulus",
	ItemDoubleLeftSquare:  "double left square",
	ItemDoubleRightSquare: "double right square",
	ItemFor:               "for loop",
	ItemRepeat:            "repeat",
	ItemWhile:             "while loop",
	ItemNext:              "next",
	ItemIn:                "in",
	ItemFunction:          "function",
	ItemBackslash:         "backslash",
	ItemPlus:              "plus",
	ItemMinus:             "minus",
	ItemMultiply:          "multiply",
	ItemDivide:            "divide",
	ItemPower:             "power",
	ItemEOF:               "end of file",
	ItemTypes:             "type",
	ItemTypesPkg:          "type package",
	ItemTypesList:         "list type",
	ItemTypesDecl:         "type declaration",
	ItemBang:              "bang",
	ItemDot:               "dot",
	ItemDoubleDot:         "dot dot",
	ItemThreeDot:          "elipsis",
}

func (t ItemType) String() string {
	k := ItemName[t]
	return k
}

func (item Item) String() string {
	return ItemName[item.Class]
}

func pad(str string, min int) string {
	out := str
	l := len(str)

	var i int
	for l < min {
		pad := "-"

		if i == 0 || i == min {
			pad = " "
		}
		out = out + pad
		l = len(out)
		i++
	}

	return out
}

func (i Item) Print() {
	name := i.String()
	val := i.Value
	if val == "\n" {
		val = "\\_n"
	}

	name = pad(name, 30)
	fmt.Printf(
		"%s `%v` \t [file: %v line: %v, char: %v]\n",
		name, val, i.File, i.Line, i.Char,
	)
}

func (i Items) Print() {
	for _, v := range i {
		v.Print()
	}
}
