package environment

import (
	"fmt"
	"math/rand"
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

func (e *Environment) GetVariableParent(name string) (Variable, bool) {
	if e.outer == nil {
		return Variable{}, false
	}

	obj, ok := e.outer.GetVariable(name, true)

	return obj, ok
}

func (e *Environment) SetVariable(name string, val Variable) Variable {
	val.Obfuscated = mask()
	e.variables[name] = val
	return val
}

func (e *Environment) SetVariableUsed(name string) (Variable, bool) {
	obj, ok := e.variables[name]

	if !ok && e.outer != nil {
		return e.outer.SetVariableUsed(name)
	}

	e.variables[name] = obj

	return obj, ok
}

func (e *Environment) GetFunction(name string, outer bool) (Function, bool) {
	obj, ok := e.functions[name]
	if !ok && e.outer != nil && outer {
		obj, ok = e.outer.GetFunction(name, outer)
	}
	return obj, ok
}

func (e *Environment) SetFunction(name string, val Function) Function {
	val.Obfuscated = mask()
	e.functions[name] = val
	return val
}

func mask() string {
	var letterRunes = []rune("-!?_|./\\abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 22)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return fmt.Sprintf("`%s`", string(b))
}
