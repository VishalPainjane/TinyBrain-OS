package main

import (
	"testing"
)

func TestColorize_Enabled(t *testing.T) {
	t.Parallel()

	colorEnabled = true
	defer func() { colorEnabled = true }()

	got := colorize("hello", ansiGreen)
	if got != ansiGreen+"hello"+ansiReset {
		t.Errorf("colorize(enabled) = %q", got)
	}
}

func TestColorize_Disabled(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	got := colorize("hello", ansiGreen)
	if got != "hello" {
		t.Errorf("colorize(disabled) = %q, want %q", got, "hello")
	}
}

func TestStateColor_AllStates(t *testing.T) {
	t.Parallel()

	cases := []struct {
		state string
		want  string
	}{
		{"RUNNING", ansiGreen},
		{"READY", ansiCyan},
		{"WAITING", ansiYellow},
		{"PREEMPTED", ansiRed},
		{"HIBERNATED", ansiBlue},
		{"NEW", ansiWhite},
		{"TERMINATED", ansiDim},
		{"UNKNOWN", ansiWhite},
	}

	for _, tc := range cases {
		if got := stateColor(tc.state); got != tc.want {
			t.Errorf("stateColor(%q) = %q, want %q", tc.state, got, tc.want)
		}
	}
}

func TestBar_Percentages(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	cases := []struct {
		used, total uint64
		wantSuffix  string
	}{
		{0, 100, "  0%"},
		{50, 100, " 50%"},
		{100, 100, "100%"},
		{0, 0, "  0%"},
	}

	for _, tc := range cases {
		got := bar(tc.used, tc.total, 20)
		if len(got) == 0 {
			t.Errorf("bar(%d, %d) returned empty string", tc.used, tc.total)
			continue
		}
		if got[0] != '[' {
			t.Errorf("bar(%d, %d) doesn't start with [: %q", tc.used, tc.total, got)
		}
		// Check the percentage suffix.
		if !containsSubstr(got, tc.wantSuffix) {
			t.Errorf("bar(%d, %d) = %q, want suffix %q", tc.used, tc.total, got, tc.wantSuffix)
		}
	}
}

func TestBar_ZeroWidth(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	got := bar(50, 100, 0)
	if got[0] != '[' {
		t.Errorf("bar with width=0 should default, got %q", got)
	}
}

func TestBold_And_Dim(t *testing.T) {
	t.Parallel()

	colorEnabled = true
	defer func() { colorEnabled = true }()

	b := bold("test")
	if b != ansiBold+"test"+ansiReset {
		t.Errorf("bold = %q", b)
	}

	d := dim("test")
	if d != ansiDim+"test"+ansiReset {
		t.Errorf("dim = %q", d)
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && findSubstring(s, sub)
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
