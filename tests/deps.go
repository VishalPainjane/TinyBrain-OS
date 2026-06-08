package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func moduleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

func goListDeps(pkg string) ([]string, error) {
	cmd := exec.Command("go", "list", "-deps", pkg)
	cmd.Dir = moduleRoot()
	out, err := cmd.Output()
	if err != nil {
		if len(out) > 0 {
			return nil, err
		}
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}
	return lines, nil
}
