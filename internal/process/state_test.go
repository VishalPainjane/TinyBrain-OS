package process_test

import (
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
)

func TestProcessState_String(t *testing.T) {
	tests := []struct {
		state process.ProcessState
		want  string
	}{
		{process.New, "NEW"},
		{process.Ready, "READY"},
		{process.Running, "RUNNING"},
		{process.Waiting, "WAITING"},
		{process.Preempted, "PREEMPTED"},
		{process.Hibernated, "HIBERNATED"},
		{process.Terminated, "TERMINATED"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessState_String_Invalid(t *testing.T) {
	invalid := process.ProcessState(99)
	if got := invalid.String(); got != "UNKNOWN" {
		t.Errorf("String() = %q, want %q", got, "UNKNOWN")
	}
}

func TestProcessState_Valid(t *testing.T) {
	for _, state := range process.All() {
		if !state.Valid() {
			t.Errorf("state %s should be valid", state)
		}
	}

	invalid := process.ProcessState(-1)
	if invalid.Valid() {
		t.Error("negative state should be invalid")
	}

	invalid = process.ProcessState(99)
	if invalid.Valid() {
		t.Error("out-of-range state should be invalid")
	}
}

func TestAll_ContainsEveryState(t *testing.T) {
	states := process.All()
	if len(states) != 7 {
		t.Fatalf("All() len = %d, want 7", len(states))
	}

	seen := make(map[process.ProcessState]bool, len(states))
	for _, state := range states {
		if seen[state] {
			t.Errorf("duplicate state in All(): %s", state)
		}
		seen[state] = true
	}

	for _, state := range states {
		if !state.Valid() {
			t.Errorf("All() returned invalid state: %s", state)
		}
	}
}

func TestProcessState_StronglyTyped(t *testing.T) {
	var _ process.ProcessState = process.New

	// Distinct constants must not compare equal.
	distinct := []process.ProcessState{
		process.New,
		process.Ready,
		process.Running,
		process.Waiting,
		process.Preempted,
		process.Hibernated,
		process.Terminated,
	}

	for i := range distinct {
		for j := i + 1; j < len(distinct); j++ {
			if distinct[i] == distinct[j] {
				t.Errorf("states %s and %s must be distinct", distinct[i], distinct[j])
			}
		}
	}
}
