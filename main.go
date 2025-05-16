package main

import (
	"fmt"
	"log"

	"github.com/devOpifex/obfuscator/cli"
	"github.com/devOpifex/obfuscator/environment"
	"github.com/devOpifex/obfuscator/lexer"
	"github.com/devOpifex/obfuscator/obfuscator"
	"github.com/devOpifex/obfuscator/parser"
	"github.com/devOpifex/obfuscator/transpiler"
)

func main() {
	c := cli.Cli()

	if *c.Key == "" {
		log.Fatal("Must pass -key")
	}

	environment.Define(*c.Key, *c.Protect, *c.Deobfuscate)

	if *c.In == "" || *c.Out == "" {
		log.Fatal("Must pass -in and -out")
	}

	if *c.In == *c.Out {
		log.Fatal("Input == output")
	}

	license := readLicense(*c.License)

	if *c.Deobfuscate && *c.License != "" {
		fmt.Println("Deobfuscating, ignoring -license")
		license = ""
	}

	obfs := &obfs{ignore: c.Ignore}
	err := obfs.readDir(*c.In)

	if err != nil {
		log.Fatal("Failed to read files")
	}

	l := lexer.New(obfs.files)
	l.Run()

	if l.HasError() {
		log.Fatal(l.Errors())
		return
	}

	p := parser.New(l)
	p.Run()

	if p.HasError() {
		log.Fatal(p.Errors())
		return
	}

	env := environment.New()
	env.SetPaths(l.Files)

	o := obfuscator.New(env, p.Files())
	o.Run()
	o.Run()

	t := transpiler.New(env, p.Files())
	t.Run()
	t.Write(*c.Out, license)
}
