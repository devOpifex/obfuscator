package cli

import (
	"flag"
	"strings"
)

type CLI struct {
	In          *string
	Out         *string
	License     *string
	Key         *string
	Protect     *string
	Ignore      []string
	Deobfuscate *bool
}

func Cli() CLI {
	in := flag.String("in", "", "Directory of R files to obfuscate")
	out := flag.String("out", "", "Directory where to write the obfuscated files")
	key := flag.String("key", "", "Key to obfuscate")
	license := flag.String("license", "", "License to prepend to every obfuscated file, e.g.: license")
	protect := flag.String("protect", "", "Comma separated protected tokens, e.g.: foo,bar")
	ignore := flag.String("ignore", "", "Comma separated directories to ignore, e.g.: renv")
	deobfuscate := flag.Bool("deobfuscate", false, "Deobfuscate the obfuscated files")

	flag.Parse()

	return CLI{
		In:          in,
		Out:         out,
		Key:         key,
		License:     license,
		Protect:     protect,
		Ignore:      ignoreToSlice(*ignore),
		Deobfuscate: deobfuscate,
	}
}

func ignoreToSlice(ignore string) []string {
	if ignore == "" {
		return []string{}
	}

	return strings.Split(ignore, ",")
}
