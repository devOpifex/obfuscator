package transpiler

import (
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

	x <- foo(x, y = 23)`

	l := lexer.NewTest(code)

	l.Run()
	p := parser.New(l)

	p.Run()

	env := environment.New()
	o := obfuscator.New(env, p.Files())
	o.RunTwice()

	trans := New(env, o.Files())
	trans.Run()
	trans.Write("newPath")
}
