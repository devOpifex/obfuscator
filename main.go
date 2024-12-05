package main

import (
	"log"
	"os"

	"github.com/sparkle-tech/obfuscator/cli"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
	"github.com/sparkle-tech/obfuscator/obfuscator"
	"github.com/sparkle-tech/obfuscator/parser"
	"github.com/sparkle-tech/obfuscator/transpiler"
)

func main() {
	c := cli.Cli()

	if *c.In == "" || *c.Out == "" || *c.Key == "" {
		log.Fatal("Must pass -in -out and -key")
	}

	if *c.In == *c.Out {
		log.Fatal("Input == output")
	}

	header := ""
	if *c.Header != "" {
		fl, err := os.ReadFile(*c.Header)

		if err != nil {
			log.Fatal("Failed to read -header")
		}

		header = string(fl)
	}

	obfs := &obfs{}
	err := obfs.readDir(*c.In)

	if err != nil {
		log.Fatal("Failed to read files")
	}

	l := lexer.New(obfs.files)
	l.Run()

	p := parser.New(l)
	p.Run()

	environment.SetKey(*c.Key)
	env := environment.New()
	o := obfuscator.New(env)
	o.RunTwice(programs)

	t := transpiler.New(env)
	t.Transpile(prog)
	err = writeString(*c.Out, t.GetCode(), header)

	if err != nil {
		log.Fatal("Failed to write obfuscated code")
	}
}
