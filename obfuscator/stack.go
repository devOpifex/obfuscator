package obfuscator

type stack []string

func (s stack) pop() stack {
	return s[:len(s)-1]
}

func (s stack) push(el string) stack {
	return append(s, el)
}
