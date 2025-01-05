package lexer

import (
	"testing"

	"github.com/sparkle-tech/obfuscator/token"
)

func TestDeclare(t *testing.T) {
	code := `x = 1
`

	l := NewTest(code)

	l.Run()

	if len(l.Files[0].Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemIdent,
			token.ItemAssign,
			token.ItemInteger,
		}

	for i, token := range tokens {
		actual := l.Files[0].Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestFunction(t *testing.T) {
	code := `foo <- function(x, y = 2) {
  x + y
}

foo(1, 2)
`

	l := NewTest(code)

	l.Run()

	if len(l.Files[0].Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemIdent,
			token.ItemAssign,
			token.ItemFunction,
			token.ItemLeftParen,
			token.ItemIdent,
			token.ItemComma,
			token.ItemIdent,
			token.ItemAssign,
			token.ItemInteger,
			token.ItemRightParen,
			token.ItemLeftCurly,
			token.ItemIdent,
			token.ItemPlus,
			token.ItemIdent,
			token.ItemRightCurly,
			token.ItemIdent,
			token.ItemLeftParen,
			token.ItemInteger,
			token.ItemComma,
			token.ItemInteger,
			token.ItemRightParen,
		}

	for i, token := range tokens {
		actual := l.Files[0].Items[i].Class
		if actual != token {
			t.Fatalf(
				"token %v expected `%v`, got `%v`",
				i,
				token,
				actual,
			)
		}
	}
}

func TestReal(t *testing.T) {
	code := `box::use(
    ambiorix[Ambiorix],
    . / here[get_home, p_rint]
  )`

	l := NewTest(code)

	l.Run()

	if len(l.Files[0].Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
}
