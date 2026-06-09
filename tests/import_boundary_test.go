package tests

import (
	"strings"
	"testing"
)

func TestSchedulerDoesNotImportInference(t *testing.T) {
	t.Parallel()

	deps := listDeps(t, "github.com/VishalPainjane/TinyBrain-OS/internal/scheduler")
	for _, dep := range deps {
		if strings.Contains(dep, "internal/inference") {
			t.Fatalf("scheduler imports inference: %s", dep)
		}
	}
}

func TestRuntimeDoesNotImportInference(t *testing.T) {
	t.Parallel()

	deps := listDeps(t, "github.com/VishalPainjane/TinyBrain-OS/internal/runtime")
	for _, dep := range deps {
		if strings.Contains(dep, "internal/inference") {
			t.Fatalf("runtime imports inference: %s", dep)
		}
	}
}

func listDeps(t *testing.T, pkg string) []string {
	t.Helper()

	out, err := goListDeps(pkg)
	if err != nil {
		t.Fatalf("go list -deps %s: %v", pkg, err)
	}
	return out
}
