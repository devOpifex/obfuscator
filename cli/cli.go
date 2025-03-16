package cli

import (
	"flag"
)

type CLI struct {
	In          *string
	Out         *string
	License     *string
	Key         *string
	Protect     *string
	Deobfuscate *bool
}

func Cli() CLI {
	in := flag.String("in", "", "Directory of R files to obfuscate")
	out := flag.String("out", "", "Directory where to write the obfuscated files")
	key := flag.String("key", "", "Key to obfuscate")
	license := flag.String("license", "", "License to prepend to every obfuscated file, e.g.: license")
	protect := flag.String("protect", "", "Comma separated protected tokens, e.g.: foo,bar")
	deobfuscate := flag.Bool("deobfuscate", false, "Deobfuscate the obfuscated files")

	flag.Parse()

	return CLI{
		In:          in,
		Out:         out,
		Key:         key,
		License:     license,
		Protect:     protect,
		Deobfuscate: deobfuscate,
	}
}
