package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_version(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"version"}, &out, &out)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(out.String(), Version) {
		t.Fatalf("output = %q, want version %s", out.String(), Version)
	}
}

func TestRun_help(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"help"}, &out, &out)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "tinybrain doctor") {
		t.Fatal("help missing doctor subcommand")
	}
}

func TestRun_unknownCommand(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"nope"}, &out, &out)
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
}

func TestRun_runMissingFlags(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"run"}, &out, &out)
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
}
