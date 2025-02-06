package lexer

import "fmt"

func (l *Lexer) Print() {
	fmt.Printf("Lexer with %v tokens\n", len(l.Files[0].Items))
	for i, f := range l.Files {
		fmt.Printf("File %v\n", i)
		f.Items.Print()
	}
}
