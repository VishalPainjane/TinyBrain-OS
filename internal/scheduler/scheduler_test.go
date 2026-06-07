package scheduler_test

import (
	"errors"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
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

func newTestScheduler() (*process.ProcessTable, scheduler.Scheduler, *scheduler.FIFOQueue) {
	table := process.NewProcessTable()
	queue := scheduler.NewFIFOQueue()
	return table, scheduler.NewFIFOScheduler(table, queue), queue
}

func TestFIFOQueue_EnqueueDequeue(t *testing.T) {
	t.Parallel()

	q := scheduler.NewFIFOQueue()
	a := sampleProcess("pid-a")
	b := sampleProcess("pid-b")

	if err := q.Enqueue(a); err != nil {
		t.Fatalf("Enqueue(a) error = %v", err)
	}
	if err := q.Enqueue(b); err != nil {
		t.Fatalf("Enqueue(b) error = %v", err)
	}
	if q.Depth() != 2 {
		t.Fatalf("Depth() = %d, want 2", q.Depth())
	}

	got, err := q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue() error = %v", err)
	}
	if got.PID != "pid-a" {
		t.Errorf("first PID = %q, want pid-a", got.PID)
	}

	got, err = q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue() error = %v", err)
	}
	if got.PID != "pid-b" {
		t.Errorf("second PID = %q, want pid-b", got.PID)
	}

	if _, err := q.Dequeue(); !errors.Is(err, scheduler.ErrQueueEmpty) {
		t.Fatalf("Dequeue() empty error = %v, want ErrQueueEmpty", err)
	}
}

func TestFIFOQueue_Peek(t *testing.T) {
	t.Parallel()

	q := scheduler.NewFIFOQueue()
	p := sampleProcess("pid-a")

	if err := q.Enqueue(p); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}

	got, err := q.Peek()
	if err != nil {
		t.Fatalf("Peek() error = %v", err)
	}
	if got.PID != "pid-a" {
		t.Errorf("Peek PID = %q, want pid-a", got.PID)
	}
	if q.Depth() != 1 {
		t.Errorf("Depth() after Peek = %d, want 1", q.Depth())
	}
}

func TestFIFOScheduler_Enqueue(t *testing.T) {
	t.Parallel()

	table, sched, queue := newTestScheduler()
	p := sampleProcess("pid-a")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Enqueue(p); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}
	if queue.Depth() != 1 {
		t.Fatalf("Depth() = %d, want 1", queue.Depth())
	}

	stored, err := table.Get("pid-a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.State != process.Ready {
		t.Errorf("State = %s, want READY", stored.State)
	}
}

func TestFIFOScheduler_ScheduleReturnsProcess(t *testing.T) {
	t.Parallel()

	table, sched, _ := newTestScheduler()
	p := sampleProcess("pid-a")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Enqueue(p); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}

	got, err := sched.Schedule()
	if err != nil {
		t.Fatalf("Schedule() error = %v", err)
	}
	if got.PID != "pid-a" {
		t.Errorf("PID = %q, want pid-a", got.PID)
	}
	if got.State != process.Running {
		t.Errorf("State = %s, want RUNNING", got.State)
	}

	stored, err := table.Get("pid-a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.State != process.Running {
		t.Errorf("table State = %s, want RUNNING", stored.State)
	}
}

func TestFIFOScheduler_TwoProcessOrdering(t *testing.T) {
	t.Parallel()

	table, sched, _ := newTestScheduler()
	a := sampleProcess("pid-a")
	b := sampleProcess("pid-b")

	for _, p := range []process.Process{a, b} {
		if err := table.Create(p); err != nil {
			t.Fatalf("Create(%s) error = %v", p.PID, err)
		}
		if err := sched.Enqueue(p); err != nil {
			t.Fatalf("Enqueue(%s) error = %v", p.PID, err)
		}
	}

	first, err := sched.Schedule()
	if err != nil {
		t.Fatalf("Schedule() first error = %v", err)
	}
	if first.PID != "pid-a" {
		t.Errorf("first PID = %q, want pid-a", first.PID)
	}

	second, err := sched.Schedule()
	if err != nil {
		t.Fatalf("Schedule() second error = %v", err)
	}
	if second.PID != "pid-b" {
		t.Errorf("second PID = %q, want pid-b", second.PID)
	}
}

func TestFIFOScheduler_ScheduleEmptyQueue(t *testing.T) {
	t.Parallel()

	_, sched, _ := newTestScheduler()

	_, err := sched.Schedule()
	if !errors.Is(err, scheduler.ErrQueueEmpty) {
		t.Fatalf("Schedule() error = %v, want ErrQueueEmpty", err)
	}
}

func TestFIFOScheduler_Preempt(t *testing.T) {
	t.Parallel()

	table, sched, _ := newTestScheduler()
	p := sampleProcess("pid-a")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Preempt("pid-a"); err != nil {
		t.Fatalf("Preempt() error = %v", err)
	}

	stored, err := table.Get("pid-a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.State != process.Preempted {
		t.Errorf("State = %s, want PREEMPTED", stored.State)
	}
}

func TestFIFOScheduler_BoostNoOp(t *testing.T) {
	t.Parallel()

	table, sched, queue := newTestScheduler()
	p := sampleProcess("pid-a")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Enqueue(p); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}
	if err := sched.Boost(); err != nil {
		t.Fatalf("Boost() error = %v", err)
	}
	if queue.Depth() != 1 {
		t.Errorf("Depth() = %d, want 1", queue.Depth())
	}

	stored, err := table.Get("pid-a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.State != process.Ready {
		t.Errorf("State = %s, want READY", stored.State)
	}
}

func TestFIFOScheduler_DuplicateEnqueue(t *testing.T) {
	t.Parallel()

	table, sched, _ := newTestScheduler()
	p := sampleProcess("pid-a")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Enqueue(p); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}
	if err := sched.Enqueue(p); !errors.Is(err, scheduler.ErrAlreadyQueued) {
		t.Fatalf("duplicate Enqueue() error = %v, want ErrAlreadyQueued", err)
	}
}

func TestFIFOScheduler_ScheduleRequiresReady(t *testing.T) {
	t.Parallel()

	table, sched, _ := newTestScheduler()
	p := sampleProcess("pid-a")

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := table.UpdateState("pid-a", process.Running); err != nil {
		t.Fatalf("UpdateState() error = %v", err)
	}

	queue := scheduler.NewFIFOQueue()
	if err := queue.Enqueue(p); err != nil {
		t.Fatalf("queue Enqueue() error = %v", err)
	}

	sched = scheduler.NewFIFOScheduler(table, queue)
	_, err := sched.Schedule()
	if !errors.Is(err, scheduler.ErrInvalidSchedule) {
		t.Fatalf("Schedule() error = %v, want ErrInvalidSchedule", err)
	}
}
