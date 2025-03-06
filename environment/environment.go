package environment

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/sparkle-tech/obfuscator/lexer"
)

var KEY string = "DEFAULT"

type Environment struct {
	variables []string
	functions []string
	arguments []string
	paths     []string
	outer     *Environment
}

func Enclose(outer *Environment) *Environment {
	env := New()
	env.outer = outer
	return env
}

func Open(env *Environment) *Environment {
	return env.outer
}

func SetKey(key string) {
	KEY = key
}

func New() *Environment {
	return &Environment{
		outer: nil,
	}
}

func (e *Environment) GetVariable(name string, outer bool) bool {
	for _, v := range e.variables {
		if v == name {
			return true
		}
	}

	if e.outer != nil && outer {
		return e.outer.GetVariable(name, outer)
	}

	return false
}

func (e *Environment) SetVariable(name string) {
	if e.GetVariable(name, false) {
		return
	}

	e.variables = append(e.variables, name)
}

func (e *Environment) GetFunction(name string) bool {
	for _, f := range e.functions {
		if f == name {
			return true
		}
	}

	if e.outer != nil {
		return e.outer.GetFunction(name)
	}

	return false
}

func (e *Environment) SetFunction(name string) {
	if e.GetFunction(name) {
		return
	}

	e.functions = append(e.functions, name)
}

func (e *Environment) SetPaths(files lexer.Files) {
	for _, f := range files {
		spit := strings.Split(f.Path, "/")
		for i := range spit {
			// root path, we skip
			if i == 0 {
				continue
			}

			if i == len(spit)-1 {
				spit[i] = strings.ReplaceAll(spit[i], ".R", "")
			}
			e.setPath(spit[i])
		}
	}
}

func (e *Environment) setPath(name string) {
	e.paths = append(e.paths, name)
}

func (e *Environment) GetPath(name string) bool {
	for _, p := range e.paths {
		if p == name {
			return true
		}
	}
	return false
}

func Mask(txt string) string {
	hasher := sha1.New()
	hasher.Write([]byte(txt + KEY))
	sha := hex.EncodeToString(hasher.Sum(nil))
	hash := base64.StdEncoding.EncodeToString([]byte(sha))
	return fmt.Sprintf("%v", strings.TrimRight(hash, "=="))
}
