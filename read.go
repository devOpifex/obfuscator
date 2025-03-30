package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/environment"
	"github.com/sparkle-tech/obfuscator/lexer"
)

type obfs struct {
	files  lexer.Files
	ignore []string
}

var ignoreRegex = regexp.MustCompile("^__")

func (o *obfs) readDir(root string) error {
	err := filepath.WalkDir(root, o.walk)

	if err != nil {
		return err
	}

	return nil
}

func (o *obfs) walk(path string, directory fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if directory.IsDir() {
		return nil
	}

	ext := filepath.Ext(path)

	if ext != ".R" {
		return nil
	}

	if o.Ignore(path) {
		return nil
	}

	fl, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	pathSplit := strings.Split(filepath.ToSlash(path), "/")
	for i := range pathSplit {
		if i == len(pathSplit)-1 {
			pathSplit[i] = strings.ReplaceAll(pathSplit[i], ".R", "")
		}

		if ignoreRegex.MatchString(pathSplit[i]) {
			continue
		}

		pathSplit[i] = environment.Mask(pathSplit[i])
	}

	rfl := lexer.File{
		Path:       path,
		Obfuscated: filepath.Join(pathSplit...) + ".R",
		PathSlice:  strings.Split(path, "/"),
		Content:    fl,
		Ast: &ast.Program{
			Statements: []ast.Statement{},
		},
	}

	o.files = append(o.files, rfl)

	return nil
}

func (o *obfs) Ignore(path string) bool {
	for _, ignore := range o.ignore {
		if strings.Contains(path, ignore) {
			return true
		}
	}

	return false
}
