package environment

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

var KEY string = "DEFAULT"

type Environment struct {
	variables map[string]Variable
	functions map[string]Function
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
	v := make(map[string]Variable)
	f := make(map[string]Function)

	return &Environment{
		functions: f,
		variables: v,
		outer:     nil,
	}
}

func (e *Environment) GetVariable(name string, outer bool) (Variable, bool) {
	obj, ok := e.variables[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetVariable(name, outer)
	}
	return obj, ok
}

func (e *Environment) SetVariable(name string, val Variable) Variable {
	_, ok := e.GetVariable(name, false)
	if ok {
		return val
	}

	val.Obfuscated = mask(name)
	e.variables[name] = val
	return val
}

func (e *Environment) GetFunction(name string, outer bool) (Function, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetFunction(name, outer)
	}
	return obj, ok
}

func (e *Environment) SetFunction(name string, val Function) Function {
	_, ok := e.GetFunction(name, true)
	if ok {
		return val
	}

	val.Obfuscated = mask(name)
	e.functions[name] = val
	return val
}

func mask(txt string) string {
	hasher := sha1.New()
	hasher.Write([]byte(txt + KEY))
	sha := hex.EncodeToString(hasher.Sum(nil))
  hash := base64.StdEncoding.EncodeToString([]byte(sha))
	return fmt.Sprintf("`%v`", strings.TrimRight(hash, "=="))
}
