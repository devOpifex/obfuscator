package parser

import (
	"fmt"
	"testing"

	"github.com/sparkle-tech/obfuscator/lexer"
)

func TestBasic(t *testing.T) {
	code := `bar <- function(z = list(x = c(1)), c = "hello") {
	  z$u <- v
		return(z)
	}

	bar(list(x = 2), c = list(v = c(1)))`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}
