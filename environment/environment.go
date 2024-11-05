package environment

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

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
	hasher.Write([]byte(txt))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("`%v`", sha)
}
