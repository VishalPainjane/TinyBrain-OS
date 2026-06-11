package scheduler

import (
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
)

// MLFQ level and token-quantum constants per docs/architecture/scheduler.md.
const (
	NumLevels = 4

	// BoostEveryTokens triggers Boost when this many tokens were recorded since last boost.
	BoostEveryTokens = 500
	// BoostEveryDuration triggers Boost after this wall-clock interval.
	BoostEveryDuration = 30 // seconds; use with time.Second at call sites

	// SwapIdleThreshold is the minimum idle duration before swap is recommended.
	SwapIdleThreshold = 10 * time.Second
)

// TokenQuantum returns the per-level token quantum for queue level in [0, NumLevels).
func TokenQuantum(level int) int {
	quanta := [NumLevels]int{32, 64, 128, 256}
	if level < 0 {
		level = 0
	}
	if level >= NumLevels {
		level = NumLevels - 1
	}
	return quanta[level]
}

// QueueLevelFromPriority maps a process scheduling priority to an initial MLFQ level.
// Higher priority values start in higher queues (Q0 is highest).
func QueueLevelFromPriority(priority int) int {
	switch {
	case priority >= 8:
		return 0
	case priority >= 5:
		return 1
	case priority >= 2:
		return 2
	default:
		return 3
	}
}

// ShouldSwap reports whether p is eligible for swap per the idle heuristic.
// Running, new, and terminated processes are never swapped via this policy.
// See docs/architecture/scheduler.md.
func ShouldSwap(p process.Process, now time.Time) bool {
	switch p.State {
	case process.Running, process.New, process.Terminated:
		return false
	}
	if p.LastExecution.IsZero() {
		return false
	}
	return now.Sub(p.LastExecution) >= SwapIdleThreshold
}
