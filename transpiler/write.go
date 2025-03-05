package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (ts Transpilers) Write(out string, header string) {
	for _, t := range ts {
		t.replaceRoot(out)
		if err := t.write(header); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *Transpiler) replaceRoot(out string) {
	t.file.Obfuscated = filepath.ToSlash(t.file.Obfuscated)
	path := strings.Split(t.file.Obfuscated, "/")
	path = path[1:len(path)]
	t.file.Obfuscated = filepath.Join(out, filepath.Join(path...))
}

func (t *Transpiler) write(header string) error {
	dir := filepath.Dir(t.file.Obfuscated)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(t.file.Obfuscated, []byte(header+t.GetCode()), 0644); err != nil {
		return err
	}

	return nil
}
