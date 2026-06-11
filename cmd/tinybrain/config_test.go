package main

import (
	"path/filepath"
	"testing"
)

func TestModelsDBPath_defaultUnderHome(t *testing.T) {
	t.Setenv("TB_MODELS_DB", "")
	path := modelsDBPath()
	if filepath.Base(path) != "models.db" {
		t.Fatalf("modelsDBPath() = %q, want models.db basename", path)
	}
}

func TestNgLayersFromEnv(t *testing.T) {
	t.Setenv("TB_NGLAYERS", "")
	if _, ok := ngLayersFromEnv(); ok {
		t.Fatal("expected unset TB_NGLAYERS")
	}
	t.Setenv("TB_NGLAYERS", "-1")
	v, ok := ngLayersFromEnv()
	if !ok || v != -1 {
		t.Fatalf("ngLayersFromEnv() = (%d, %v), want (-1, true)", v, ok)
	}
}

func TestNgLayersFromEnv_invalid(t *testing.T) {
	t.Setenv("TB_NGLAYERS", "abc")
	if _, ok := ngLayersFromEnv(); ok {
		t.Fatal("expected invalid TB_NGLAYERS to return ok=false")
	}
}

func TestModelsDBPath_envOverride(t *testing.T) {
	t.Setenv("TB_MODELS_DB", "/custom/path/models.db")
	if got := modelsDBPath(); got != "/custom/path/models.db" {
		t.Fatalf("modelsDBPath() = %q, want env override", got)
	}
}
