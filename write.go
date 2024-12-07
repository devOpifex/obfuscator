package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func WriteWithDirCreate(path string, content string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
