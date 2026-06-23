package main

import (
	"testing"
)

func TestSchemeRegistration(t *testing.T) {
	// scheme is a package-level variable in main.go
	if !scheme.IsGroupRegistered("core.tinybrain.io") {
		t.Errorf("Expected core group to be registered")
	}
	if !scheme.IsGroupRegistered("memory.tinybrain.io") {
		t.Errorf("Expected memory group to be registered")
	}
}
