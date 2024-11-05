package transpiler

import (
	"fmt"
	"testing"

	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
	"github.com/sparkle-tech/obfuscator/obfuscator"
	"github.com/sparkle-tech/obfuscator/parser"
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

	prog := p.Run()

	env := environment.New()
	o := obfuscator.New(env)
	o.Obfuscate(prog)
	o.Obfuscate(prog)

	trans := New(env)
	trans.Transpile(prog)
	fmt.Println(trans.GetCode())
}
