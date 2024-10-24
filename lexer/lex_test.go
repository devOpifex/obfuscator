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

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	tokens :=
		[]token.ItemType{
			token.ItemIdent,
			token.ItemAssign,
			token.ItemInteger,
		}

	for i, token := range tokens {
		actual := l.Items[i].Class
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

	if len(l.Items) == 0 {
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
		actual := l.Items[i].Class
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
	code := `foo <- function(x, y = 2) {
  x + y
}

z <- foo(1, 2)

emit <- function(x, type) {
  UseMethod("emit")
}

#' @export
emit.lexer <- function(x, type) {
  if(!type %in% TYPES) {
    stop(sprintf("UNKNOWN type: %s", type))
  }

  x$tokens <- append(
    x$tokens,
    list(
      list(
        item = x$item,
        type = type
      )
    )
  )

  x$item <- ""

  invisible(x)
}

code <- function(x) {
  UseMethod("code")
}

#' @export
code.obfuscator <- function(x) {
  code <- ""

  has_code <- FALSE
  for(index in seq_along(x$tokens)) {
    token <- x$tokens[[index]]

    if(index > 1 && token$type == NEWLINE_T && x$tokens[[index - 1]]$type == NEWLINE_T) {
      next
    }

    # skip empty lines on top of script
    if(token$type == NEWLINE_T && !has_code){
      next
    }

    if(token$type == COMMENT_T){
      next
    }

    val <- token$item

    if(length(token$obfuscated)){
      val <- token$obfuscated
    }

    has_code <- TRUE
    code <- paste0(
      code,
      val
    )
  }

  invisible(code)
}
`

	l := NewTest(code)

	l.Run()

	if len(l.Items) == 0 {
		t.Fatal("No Items where lexed")
	}

	l.Print()
}
