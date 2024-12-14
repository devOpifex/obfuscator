package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sparkle-tech/obfuscator/ast"
	"github.com/sparkle-tech/obfuscator/lexer"
)

type obfs struct {
	files lexer.Files
}

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

	fl, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	rfl := lexer.File{
		Path:    path,
		Content: fl,
		Ast: &ast.Program{
			Statements: []ast.Statement{},
		},
	}

	o.files = append(o.files, rfl)

	return nil
}
