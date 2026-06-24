package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteBox_BasicStructure(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	writeBox(&buf, "Test", 40, []string{"hello", "world"})
	out := buf.String()

	// Check top border.
	if !strings.Contains(out, "┌") {
		t.Error("missing top-left corner ┌")
	}
	if !strings.Contains(out, "┐") {
		t.Error("missing top-right corner ┐")
	}
	if !strings.Contains(out, "Test") {
		t.Error("missing title")
	}

	// Check content lines have side borders.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "hello") || strings.Contains(line, "world") {
			if !strings.HasPrefix(line, "│") {
				t.Errorf("content line missing left border: %q", line)
			}
			if !strings.HasSuffix(line, "│") {
				t.Errorf("content line missing right border: %q", line)
			}
		}
	}

	// Check bottom border.
	if !strings.Contains(out, "└") {
		t.Error("missing bottom-left corner └")
	}
	if !strings.Contains(out, "┘") {
		t.Error("missing bottom-right corner ┘")
	}
}

func TestWriteBox_EmptyContent(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	writeBox(&buf, "Empty", 30, nil)
	out := buf.String()

	if !strings.Contains(out, "┌") || !strings.Contains(out, "┘") {
		t.Fatalf("missing box borders:\n%s", out)
	}
}

func TestWriteBox_NoTitle(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	writeBox(&buf, "", 30, []string{"data"})
	out := buf.String()

	if !strings.Contains(out, "data") {
		t.Fatalf("missing content:\n%s", out)
	}
}

func TestPadRight(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		width int
		want  string
	}{
		{"hi", 5, "hi   "},
		{"hello", 5, "hello"},
		{"toolong", 3, "toolong"}, // wider than target — returned as-is
		{"", 3, "   "},
	}

	for _, tc := range cases {
		got := padRight(tc.input, tc.width)
		if got != tc.want {
			t.Errorf("padRight(%q, %d) = %q, want %q", tc.input, tc.width, got, tc.want)
		}
	}
}

func TestVisualLength(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"\033[31mhello\033[0m", 5}, // "hello" wrapped in ANSI red + reset
		{"", 0},
		{"\033[1m\033[32mAB\033[0m", 2}, // bold + green + "AB" + reset
	}

	for _, tc := range cases {
		got := visualLength(tc.input)
		if got != tc.want {
			t.Errorf("visualLength(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}
