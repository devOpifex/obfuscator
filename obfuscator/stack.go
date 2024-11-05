package obfuscator

type Stack []string

func (s Stack) Pop() Stack {
	return s[:len(s)-1]
}

func (s Stack) Push(el string) Stack {
	return append(s, el)
}

func (s Stack) Get() (bool, string) {
	if len(s) == 0 {
		return false, ""
	}
	return true, s[len(s)-1]
}
