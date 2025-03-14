package environment

import (
	"strings"

	"github.com/sparkle-tech/obfuscator/lexer"
)

var KEY string = "DEFAULT"
var PROTECT []string

type Environment struct {
	variables []string
	functions []string
	arguments []string
	generics  []string
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

func Define(key string, protect string) {
	KEY = key
	PROTECT = strings.Split(protect, ",")
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

func isProtected(name string) bool {
	if len(PROTECT) == 0 {
		return false
	}

	for _, p := range PROTECT {
		if p == name {
			return true
		}
	}
	return false
}

func Mask(txt string) string {
	if isProtected(txt) {
		return txt
	}

	return cipher(txt, KEY)
}

func Unmask(ciphertext string) string {
	return decipher(ciphertext, KEY)
}

func cipher(plaintext, secret string) string {
	if len(plaintext) == 0 || len(secret) == 0 {
		return plaintext
	}

	// Derive a key sequence from the secret
	keySequence := make([]int, len(secret))
	for i, char := range secret {
		keySequence[i] = int(char) % 52
	}

	var result strings.Builder
	keyIndex := 0

	// For each character in the plaintext
	for _, char := range plaintext {
		// Get the current key value
		keyVal := keySequence[keyIndex%len(keySequence)]
		keyIndex++

		// Convert char to integer and apply the key (simple XOR-like operation)
		charCode := int(char)
		encoded := (charCode + keyVal) % 0x10000 // Keep within Unicode BMP

		// We need to encode this value using only a-z and A-Z (52 possible values)
		// This is essentially base-52 encoding
		// Each Unicode character requires up to 3 alphabet chars to represent (52^3 > 65536)

		// First, apply a deterministic scramble based on the secret
		scrambleFactor := 0
		for _, s := range secret {
			scrambleFactor = (scrambleFactor + int(s)) % 0x10000
		}
		encoded = (encoded + scrambleFactor) % 0x10000

		// Convert to base-52
		const alphabetSize = 52

		// Calculate how many digits we need (3 is enough for all Unicode BMP)
		digits := []int{
			encoded % alphabetSize,
			(encoded / alphabetSize) % alphabetSize,
			(encoded / (alphabetSize * alphabetSize)) % alphabetSize,
		}

		// Convert each digit to a letter (0-25 -> A-Z, 26-51 -> a-z)
		for _, digit := range digits {
			var encodedChar rune
			if digit < 26 {
				encodedChar = 'A' + rune(digit)
			} else {
				encodedChar = 'a' + rune(digit-26)
			}
			result.WriteRune(encodedChar)
		}
	}

	return result.String()
}

// Decipher reverses the encryption
func decipher(ciphertext, secret string) string {
	if len(ciphertext) == 0 || len(secret) == 0 {
		return ciphertext
	}

	// Check if the length is valid (must be multiple of 3)
	if len(ciphertext)%3 != 0 {
		return "Invalid ciphertext length"
	}

	// Derive the key sequence from the secret
	keySequence := make([]int, len(secret))
	for i, char := range secret {
		keySequence[i] = int(char) % 52
	}

	// Calculate the scramble factor
	scrambleFactor := 0
	for _, s := range secret {
		scrambleFactor = (scrambleFactor + int(s)) % 0x10000
	}

	var result strings.Builder
	keyIndex := 0

	// Process the ciphertext in groups of 3 characters
	for i := 0; i < len(ciphertext); i += 3 {
		// Convert each letter back to a digit
		digits := make([]int, 3)

		for j := 0; j < 3; j++ {
			char := rune(ciphertext[i+j])
			if char >= 'A' && char <= 'Z' {
				digits[j] = int(char - 'A')
			} else if char >= 'a' && char <= 'z' {
				digits[j] = int(char - 'a' + 26)
			} else {
				return "Invalid character in ciphertext"
			}
		}

		// Convert from base-52
		const alphabetSize = 52
		encoded := digits[0] + digits[1]*alphabetSize + digits[2]*alphabetSize*alphabetSize

		// Remove the scramble
		encoded = (encoded - scrambleFactor) % 0x10000
		if encoded < 0 {
			encoded += 0x10000
		}

		// Get the key value
		keyVal := keySequence[keyIndex%len(keySequence)]
		keyIndex++

		// Decode the character
		charCode := (encoded - keyVal) % 0x10000
		if charCode < 0 {
			charCode += 0x10000
		}

		// Convert back to a rune
		result.WriteRune(rune(charCode))
	}

	return result.String()
}
