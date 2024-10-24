package parser

import (
	"fmt"
	"testing"

	"github.com/sparkle-tech/obfuscator/lexer"
)

func TestBasic(t *testing.T) {
	code := `foo <- function(x = 1, y = 2){
    return(x + y)
	}

	foo(2, 3)

	library(ambiorix)

	app <- Ambiorix$new()

	app$get("/about", \(req, res) {
	  res$send("hello")
	})`

	l := lexer.NewTest(code)

	l.Run()
	p := New(l)

	prog := p.Run()

	fmt.Println(prog.String())
}
