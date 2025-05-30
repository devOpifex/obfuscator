package diagnostics

import (
	"bytes"
	"fmt"

	"github.com/devOpifex/obfuscator/token"
)

// To match LSP specs
type Severity int

const (
	Fatal Severity = iota
	Warn
	Hint
	Info
)

type Diagnostic struct {
	Token    token.Item
	Message  string
	Severity Severity
}

type Diagnostics []Diagnostic

func New(token token.Item, message string, severity Severity) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: severity,
	}
}

func NewError(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Fatal,
	}
}

func NewWarning(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Warn,
	}
}

func NewInfo(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Info,
	}
}

func NewHint(token token.Item, message string) Diagnostic {
	return Diagnostic{
		Token:    token,
		Message:  message,
		Severity: Hint,
	}
}

func (d Diagnostics) String() string {
	var out bytes.Buffer

	for _, v := range d {
		out.WriteString(v.String())
	}

	return out.String()
}

func (v Diagnostic) String() string {
	var out bytes.Buffer
	out.WriteString("[" + v.Severity.String() + "]\t")
	out.WriteString(v.Token.File)
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%v", v.Token.Line))
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%v", v.Token.Char))
	out.WriteString(" " + v.Message + "\n")
	return out.String()
}

func (d Diagnostics) Print() {
	fmt.Printf("%v", d.String())
}

func (s Severity) String() string {
	if s == Fatal {
		return "ERROR"
	}

	if s == Warn {
		return "WARN"
	}

	if s == Info {
		return "INFO"
	}

	return "HINT"
}
