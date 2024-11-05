package main

import (
	"os"
)

func writeString(file, code, header string) error {
	f, err := os.Create(file)

	if err != nil {
		return err
	}

	defer f.Close()

	if header != "" {
		code = header + code
	}

	_, err = f.WriteString(code)

	return err
}
