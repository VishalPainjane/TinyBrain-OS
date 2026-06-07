package process_test

import (
	"errors"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
)

func sampleProcess(pid string) process.Process {
	return process.Process{
		PID:            pid,
		AgentRef:       "agent-1",
		Priority:       5,
		MemoryUsage:    1024,
		VRAMUsage:      512,
		KVCacheID:      "kv-1",
		LastExecution:  time.Unix(1_700_000_000, 0).UTC(),
		TokensProduced: 42,
		TaskID:         "task-1",
	}
}

func TestProcessTable_Create(t *testing.T) {
	table := process.NewProcessTable()
	p := sampleProcess("pid-1")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := table.Get("pid-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.State != process.New {
		t.Errorf("State = %s, want NEW", got.State)
	}
	if got.AgentRef != p.AgentRef {
		t.Errorf("AgentRef = %q, want %q", got.AgentRef, p.AgentRef)
	}
}

func TestProcessTable_Create_ForcesNewState(t *testing.T) {
	table := process.NewProcessTable()
	p := sampleProcess("pid-1")
	p.State = process.Running

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := table.Get("pid-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.State != process.New {
		t.Errorf("State = %s, want NEW (Create must force NEW)", got.State)
	}
}

func TestProcessTable_Create_DuplicatePID(t *testing.T) {
	table := process.NewProcessTable()
	p := sampleProcess("pid-1")

	if err := table.Create(p); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	if err := table.Create(p); !errors.Is(err, process.ErrDuplicatePID) {
		t.Errorf("second Create() error = %v, want ErrDuplicatePID", err)
	}
}

func TestProcessTable_Get_NotFound(t *testing.T) {
	table := process.NewProcessTable()

	_, err := table.Get("missing")
	if !errors.Is(err, process.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestProcessTable_List(t *testing.T) {
	table := process.NewProcessTable()

	if len(table.List()) != 0 {
		t.Fatal("List() on empty table should return empty slice")
	}

	for _, pid := range []string{"a", "b", "c"} {
		if err := table.Create(sampleProcess(pid)); err != nil {
			t.Fatalf("Create(%s) error = %v", pid, err)
		}
	}

	list := table.List()
	if len(list) != 3 {
		t.Fatalf("List() len = %d, want 3", len(list))
	}

	seen := make(map[string]bool, len(list))
	for _, p := range list {
		seen[p.PID] = true
	}
	for _, pid := range []string{"a", "b", "c"} {
		if !seen[pid] {
			t.Errorf("List() missing PID %q", pid)
		}
	}
}

func TestProcessTable_UpdateState(t *testing.T) {
	table := process.NewProcessTable()
	if err := table.Create(sampleProcess("pid-1")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := table.UpdateState("pid-1", process.Running); err != nil {
		t.Fatalf("UpdateState() error = %v", err)
	}

	got, err := table.Get("pid-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.State != process.Running {
		t.Errorf("State = %s, want RUNNING", got.State)
	}
}

func TestProcessTable_UpdateState_NotFound(t *testing.T) {
	table := process.NewProcessTable()

	err := table.UpdateState("missing", process.Ready)
	if !errors.Is(err, process.ErrNotFound) {
		t.Errorf("UpdateState() error = %v, want ErrNotFound", err)
	}
}

func TestProcessTable_UpdateState_InvalidState(t *testing.T) {
	table := process.NewProcessTable()
	if err := table.Create(sampleProcess("pid-1")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := table.UpdateState("pid-1", process.ProcessState(99))
	if err == nil {
		t.Fatal("UpdateState() with invalid state should error")
	}
}

func TestProcessTable_Delete_Terminated(t *testing.T) {
	table := process.NewProcessTable()
	if err := table.Create(sampleProcess("pid-1")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := table.UpdateState("pid-1", process.Terminated); err != nil {
		t.Fatalf("UpdateState() error = %v", err)
	}

	if err := table.Delete("pid-1"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := table.Get("pid-1")
	if !errors.Is(err, process.ErrNotFound) {
		t.Errorf("Get() after Delete error = %v, want ErrNotFound", err)
	}
}

func TestProcessTable_Delete_NotTerminated(t *testing.T) {
	table := process.NewProcessTable()
	if err := table.Create(sampleProcess("pid-1")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := table.Delete("pid-1")
	if !errors.Is(err, process.ErrNotTerminated) {
		t.Errorf("Delete() error = %v, want ErrNotTerminated", err)
	}
}

func TestProcessTable_Delete_NotFound(t *testing.T) {
	table := process.NewProcessTable()

	err := table.Delete("missing")
	if !errors.Is(err, process.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
