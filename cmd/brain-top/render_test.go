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
	renderSnapshot(&buf, "test", processTableReader{table: table}, mlfqQueueReader{sched: sched})
	out := buf.String()

	for _, want := range []string{
		"brain-top test | snapshot",
		"[Processes]",
		"pid-a", "RUNNING",
		"pid-b", "READY",
		"[Queues]",
		"[Resources]",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestRenderSnapshot_EmptyProcessTable(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	renderSnapshot(&buf, "test", processTableReader{table: process.NewProcessTable()}, nil)
	if !strings.Contains(buf.String(), "(none)") {
		t.Fatalf("expected empty process panel, got:\n%s", buf.String())
	}
}

func TestRun_VersionAndSnapshot(t *testing.T) {
	t.Parallel()

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
	if !strings.Contains(stdout.String(), "[Processes]") {
		t.Fatalf("snapshot output = %q", stdout.String())
	}
}
