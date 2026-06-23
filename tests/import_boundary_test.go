package tests

import (
	"strings"
	"testing"
)

const modulePrefix = "github.com/VishalPainjane/TinyBrain-OS/"

// Architecture fitness tests enforce import boundaries from docs/architecture/invariants.md.
// They run in the fast CI job (CGO_ENABLED=0, no integration tags).

func TestImportBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pkg       string
		forbid    string
		invariant string
	}{
		{
			name:      "INV-001 scheduler does not import runtime",
			pkg:       modulePrefix + "internal/scheduler",
			forbid:    modulePrefix + "internal/runtime",
			invariant: "INV-001",
		},
		{
			name:      "INV-002 runtime does not import scheduler",
			pkg:       modulePrefix + "internal/runtime",
			forbid:    modulePrefix + "internal/scheduler",
			invariant: "INV-002",
		},
		{
			name:      "scheduler does not import inference",
			pkg:       modulePrefix + "internal/scheduler",
			forbid:    modulePrefix + "internal/inference",
			invariant: "INV-001",
		},
		{
			name:      "runtime does not import inference",
			pkg:       modulePrefix + "internal/runtime",
			forbid:    modulePrefix + "internal/inference",
			invariant: "INV-002",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertNoModuleDep(t, tt.pkg, tt.forbid, tt.invariant)
		})
	}
}

func TestCorePackagesDoNotImportInference(t *testing.T) {
	t.Parallel()

	// INV-008: inference engine only in adapter packages under internal/inference/.
	corePkgs := []string{
		modulePrefix + "internal/process",
		modulePrefix + "internal/events",
		modulePrefix + "internal/registry",
		modulePrefix + "internal/hardware",
		modulePrefix + "internal/loader",
		modulePrefix + "internal/scheduler",
		modulePrefix + "internal/kv",
		modulePrefix + "internal/swap",
		modulePrefix + "internal/agents",
		modulePrefix + "internal/runtime",
	}

	forbidden := modulePrefix + "internal/inference"
	for _, pkg := range corePkgs {
		pkg := pkg
		t.Run(shortPkg(pkg), func(t *testing.T) {
			t.Parallel()
			assertNoModuleDep(t, pkg, forbidden, "INV-008")
		})
	}
}

func shortPkg(pkg string) string {
	return strings.TrimPrefix(pkg, modulePrefix)
}

func assertNoModuleDep(t *testing.T, pkg, forbidden, invariant string) {
	t.Helper()

	deps := listDeps(t, pkg)
	for _, dep := range deps {
		if dep == forbidden || strings.HasPrefix(dep, forbidden+"/") {
			t.Fatalf("%s: %s imports forbidden %q: %s", invariant, pkg, forbidden, dep)
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
