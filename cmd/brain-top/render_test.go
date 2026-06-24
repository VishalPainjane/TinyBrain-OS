package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
)

func TestRenderSnapshot_ProcessesAndQueues(t *testing.T) {
	t.Parallel()

	table := process.NewProcessTable()
	queue := scheduler.NewMLFQQueue()
	sched := scheduler.NewMLFQScheduler(table, queue)

	p1 := process.Process{
		PID: "pid-a", AgentRef: "agent-1", Priority: 10,
		KVCacheID: "kv-a", TaskID: "task-1",
	}
	p2 := process.Process{
		PID: "pid-b", AgentRef: "agent-2", Priority: 1,
		KVCacheID: "kv-b", TaskID: "task-2",
	}
	for _, p := range []process.Process{p1, p2} {
		if err := table.Create(p); err != nil {
			t.Fatalf("Create(%s) error = %v", p.PID, err)
		}
	}
	if err := table.UpdateState("pid-a", process.Running); err != nil {
		t.Fatalf("UpdateState(pid-a) error = %v", err)
	}
	if err := sched.Enqueue(p2); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}

	var buf bytes.Buffer
	colorEnabled = false
	defer func() { colorEnabled = true }()
	renderSnapshot(&buf, "test", processTableReader{table: table}, mlfqQueueReader{sched: sched}, nil)
	out := buf.String()

	for _, want := range []string{
		"brain-top test",
		"Processes",
		"pid-a", "RUNNING",
		"pid-b", "READY",
		"MLFQ Queues",
		"Resources",
		"Swap Monitor",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestRenderSnapshot_EmptyProcessTable(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	renderSnapshot(&buf, "test", processTableReader{table: process.NewProcessTable()}, nil, nil)
	if !strings.Contains(buf.String(), "(no processes)") {
		t.Fatalf("expected empty process panel, got:\n%s", buf.String())
	}
}

func TestRun_VersionAndSnapshot(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var stdout bytes.Buffer
	if code := run([]string{"version"}, &stdout, nil); code != 0 {
		t.Fatalf("run(version) = %d", code)
	}
	if !strings.Contains(stdout.String(), "brain-top") {
		t.Fatalf("version output = %q", stdout.String())
	}

	stdout.Reset()
	if code := run([]string{"snapshot"}, &stdout, nil); code != 0 {
		t.Fatalf("run(snapshot) = %d", code)
	}
	if !strings.Contains(stdout.String(), "Processes") {
		t.Fatalf("snapshot output = %q", stdout.String())
	}
}

func TestRun_HelpAndUnknown(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	if code := run([]string{"help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("run(help) = %d", code)
	}
	if !strings.Contains(stdout.String(), "Usage") {
		t.Fatalf("help output missing Usage:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"bogus"}, &stdout, &stderr); code != 2 {
		t.Fatalf("run(bogus) = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRun_InvalidInterval(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer

	// Invalid duration.
	if code := run([]string{"watch", "nope"}, &stdout, &stderr); code != 2 {
		t.Fatalf("run(watch nope) = %d, want 2", code)
	}

	stderr.Reset()

	// Too-short interval.
	if code := run([]string{"watch", "50ms"}, &stdout, &stderr); code != 2 {
		t.Fatalf("run(watch 50ms) = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "200ms") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRenderSnapshot_AllFourPanels(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	renderSnapshot(&buf, "test", processTableReader{}, mlfqQueueReader{}, nil)
	out := buf.String()

	panels := []string{"Processes", "Resources", "MLFQ Queues", "Swap Monitor"}
	for _, panel := range panels {
		if !strings.Contains(out, panel) {
			t.Errorf("output missing panel %q\n%s", panel, out)
		}
	}

	// Verify box-drawing characters are present.
	for _, ch := range []string{"┌", "┐", "└", "┘", "│"} {
		if !strings.Contains(out, ch) {
			t.Errorf("output missing box character %q", ch)
		}
	}
}

func TestRenderSnapshot_HeaderContainsVersion(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	renderSnapshot(&buf, "1.0.0", processTableReader{}, mlfqQueueReader{}, nil)
	if !strings.Contains(buf.String(), "brain-top 1.0.0") {
		t.Fatalf("header missing version:\n%s", buf.String())
	}
}
