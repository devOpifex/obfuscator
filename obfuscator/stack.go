package obfuscator

type call struct {
	name   string
	method bool
}

type Stack []call

func (s Stack) Pop() Stack {
	return s[:len(s)-1]
}

func (s Stack) Push(el string, t bool) Stack {
	return append(s, call{el, t})
}

func (s Stack) Get() (bool, call) {
	if len(s) == 0 {
		return false, call{}
	}
	return true, s[len(s)-1]
}
