package obfuscator

import (
	"testing"

	"github.com/devOpifex/obfuscator/environment"
	"github.com/devOpifex/obfuscator/lexer"
	"github.com/devOpifex/obfuscator/parser"
)

func TestBasic(t *testing.T) {
	code := `x <- 1

	foo <- function(x, y = 1){
	  total <- sum(x, y)
		return(total)
	}

	results <- foo(x = x, y = \(z = 2L){
	 return(z + 2)
	})
	results <- foo(x, y = x)`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	p.Run()

	env := environment.New()
	o := New(env, p.Files())
	o.RunTwice()
}
