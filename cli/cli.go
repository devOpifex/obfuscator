package cli

import (
	"flag"
)

type CLI struct {
	In     *string
	Out    *string
	Header *string
	Key    *string
}

func Cli() CLI {
	in := flag.String("in", "", "Directory of R files to obfuscate")
	out := flag.String("out", "", "Output name of file to write obfuscated code")
	key := flag.String("key", "", "Key to obfuscate")
	header := flag.String("header", "", "Header to append to obfuscated code, e.g.: license")

	flag.Parse()

	return CLI{
		In:     in,
		Out:    out,
		Key:    key,
		Header: header,
	}
}
