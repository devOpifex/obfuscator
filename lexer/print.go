package lexer

import "fmt"

func (l *Lexer) Print() {
	fmt.Printf("Lexer with %v tokens\n", len(l.Files[0].Items))
	l.Files[0].Items.Print()
}
