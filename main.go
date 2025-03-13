package main

import (
	"log"

	"github.com/sparkle-tech/obfuscator/cli"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
	"github.com/sparkle-tech/obfuscator/obfuscator"
	"github.com/sparkle-tech/obfuscator/parser"
	"github.com/sparkle-tech/obfuscator/transpiler"
)

func main() {
	c := cli.Cli()

	if *c.Key == "" {
		log.Fatal("Must pass -key")
	}

	environment.Define(*c.Key, *c.Protect)

	if *c.In == "" || *c.Out == "" {
		log.Fatal("Must pass -in, and -out")
	}

	if *c.In == *c.Out {
		log.Fatal("Input == output")
	}

	license := readLicense(*c.License)

	obfs := &obfs{}
	err := obfs.readDir(*c.In)

	if err != nil {
		log.Fatal("Failed to read files")
	}

	l := lexer.New(obfs.files)
	l.Run()

	p := parser.New(l)
	p.Run()

	env := environment.New()
	env.SetPaths(l.Files)

	o := obfuscator.New(env, p.Files())
	o.Run()
	o.Run()

	t := transpiler.New(env, p.Files())
	t.Run()
	t.Write(*c.Out, license)
}
