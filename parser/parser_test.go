package parser

import (
	"testing"

	"github.com/devOpifex/obfuscator/lexer"
)

func TestBasic(t *testing.T) {
	code := `bar <- function(z = list(x = c(1)), c = "hello") {
	  z$u <- v
		return(z)
	}

	x <- 1

	bar(list(x = 2), c = list(v = c(1, 2)), x = 2)

	u <- sum(2, 3, 4)
	`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	p.Run()

	p.Print()
}
