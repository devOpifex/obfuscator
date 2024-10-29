package transpiler

import (
	"testing"

	"github.com/sparkle-tech/obfuscator/lexer"
	"github.com/sparkle-tech/obfuscator/parser"
)

func (trans *Obfuscator) testOutput(t *testing.T, expected string) {
	if trans.GetCode() == expected {
		return
	}
	t.Fatalf("expected:\n`%v`\ngot:\n`%v`", expected, trans.GetCode())
}

func TestBasic(t *testing.T) {
	code := `x <- 1
	y <- 2

	foo <- function(x, y = y){
	  total <- sum(x, y)
		return(total)
	}

	results <- foo(x = 2)`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	prog := p.Run()

	o := New()
	o.Obfuscate(prog)

	expectations := `x=1;y=2;foo=\(x){x+1};foo(x=2)`
	o.testOutput(t, expectations)
}
