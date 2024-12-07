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
		if err := t.write(); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *Transpiler) replaceRoot(out string) {
	t.file.Path = filepath.ToSlash(t.file.Path)
	path := strings.Split(t.file.Path, "/")
	path = path[1:len(path)]
	t.file.Path = filepath.Join(out, filepath.Join(path...))
}

func (t *Transpiler) write() error {
	dir := filepath.Dir(t.file.Path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(t.file.Path, []byte(t.GetCode()), 0644); err != nil {
		return err
	}

	return nil
}
