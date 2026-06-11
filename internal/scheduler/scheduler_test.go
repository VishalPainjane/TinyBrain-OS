package scheduler_test

import (
	"errors"
	"fmt"
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

func newMLFQTestScheduler() (*process.ProcessTable, *scheduler.MLFQScheduler, *scheduler.MLFQQueue) {
	table := process.NewProcessTable()
	queue := scheduler.NewMLFQQueue()
	return table, scheduler.NewMLFQScheduler(table, queue), queue
}

func TestQueueLevelFromPriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		priority int
		want     int
	}{
		{priority: 10, want: 0},
		{priority: 8, want: 0},
		{priority: 7, want: 1},
		{priority: 5, want: 1},
		{priority: 4, want: 2},
		{priority: 1, want: 3},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("priority_%d", tt.priority), func(t *testing.T) {
			t.Parallel()
			if got := scheduler.QueueLevelFromPriority(tt.priority); got != tt.want {
				t.Errorf("QueueLevelFromPriority(%d) = %d, want %d", tt.priority, got, tt.want)
			}
		})
	}
}

func TestMLFQQueue_HigherLevelFirst(t *testing.T) {
	t.Parallel()

	q := scheduler.NewMLFQQueue()
	low := sampleProcess("pid-low")
	low.Priority = 1
	high := sampleProcess("pid-high")
	high.Priority = 10

	if err := q.Enqueue(low); err != nil {
		t.Fatalf("Enqueue(low) error = %v", err)
	}
	if err := q.Enqueue(high); err != nil {
		t.Fatalf("Enqueue(high) error = %v", err)
	}

	got, err := q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue() error = %v", err)
	}
	if got.PID != "pid-high" {
		t.Errorf("first PID = %q, want pid-high", got.PID)
	}
}

func TestMLFQScheduler_Preemption(t *testing.T) {
	t.Parallel()

	table, sched, _ := newMLFQTestScheduler()

	low := sampleProcess("pid-low")
	low.Priority = 1
	high := sampleProcess("pid-high")
	high.Priority = 10

	for _, p := range []process.Process{low, high} {
		if err := table.Create(p); err != nil {
			t.Fatalf("Create(%s) error = %v", p.PID, err)
		}
	}

	if err := sched.Enqueue(low); err != nil {
		t.Fatalf("Enqueue(low) error = %v", err)
	}
	if _, err := sched.Schedule(); err != nil {
		t.Fatalf("Schedule(low) error = %v", err)
	}

	if err := sched.Enqueue(high); err != nil {
		t.Fatalf("Enqueue(high) error = %v", err)
	}

	got, err := sched.Schedule()
	if err != nil {
		t.Fatalf("Schedule(high) error = %v", err)
	}
	if got.PID != "pid-high" {
		t.Errorf("scheduled PID = %q, want pid-high", got.PID)
	}

	lowStored, err := table.Get("pid-low")
	if err != nil {
		t.Fatalf("Get(pid-low) error = %v", err)
	}
	if lowStored.State != process.Ready {
		t.Errorf("pid-low State = %s, want READY (re-queued after preemption)", lowStored.State)
	}
	if got.State != process.Running {
		t.Errorf("pid-high State = %s, want RUNNING", got.State)
	}
}

func TestMLFQScheduler_TokenQuantumDemotion(t *testing.T) {
	t.Parallel()

	table, sched, _ := newMLFQTestScheduler()
	p := sampleProcess("pid-a")
	p.Priority = 10

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Enqueue(p); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}
	if _, err := sched.Schedule(); err != nil {
		t.Fatalf("Schedule() error = %v", err)
	}

	for i := 0; i < scheduler.TokenQuantum(0); i++ {
		if err := sched.RecordToken("pid-a"); err != nil {
			t.Fatalf("RecordToken(%d) error = %v", i, err)
		}
	}

	level, ok := sched.QueueLevel("pid-a")
	if !ok {
		t.Fatal("QueueLevel() missing pid-a")
	}
	if level != 1 {
		t.Errorf("QueueLevel after quantum = %d, want 1", level)
	}
}

func TestMLFQScheduler_Boost(t *testing.T) {
	t.Parallel()

	table, sched, queue := newMLFQTestScheduler()
	low := sampleProcess("pid-low")
	low.Priority = 1

	if err := table.Create(low); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := sched.Enqueue(low); err != nil {
		t.Fatalf("Enqueue() error = %v", err)
	}

	if err := sched.Boost(); err != nil {
		t.Fatalf("Boost() error = %v", err)
	}

	depths := sched.QueueDepths()
	if depths[0] != 1 {
		t.Errorf("Q0 depth = %d, want 1", depths[0])
	}
	if queue.Depth() != 1 {
		t.Errorf("total Depth() = %d, want 1", queue.Depth())
	}

	level, ok := sched.QueueLevel("pid-low")
	if !ok || level != 0 {
		t.Errorf("QueueLevel = %d, ok=%v, want 0 true", level, ok)
	}
}

func TestMLFQScheduler_QueueDepths(t *testing.T) {
	t.Parallel()

	table, sched, _ := newMLFQTestScheduler()
	for i, priority := range []int{10, 7, 4, 1} {
		p := sampleProcess(fmt.Sprintf("pid-%d", i))
		p.Priority = priority
		if err := table.Create(p); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if err := sched.Enqueue(p); err != nil {
			t.Fatalf("Enqueue() error = %v", err)
		}
	}

	depths := sched.QueueDepths()
	for level, want := range []int{1, 1, 1, 1} {
		if depths[level] != want {
			t.Errorf("Q%d depth = %d, want %d", level, depths[level], want)
		}
	}
}

func TestShouldSwap(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0).UTC()
	p := sampleProcess("pid-a")
	p.State = process.Waiting
	p.LastExecution = now.Add(-11 * time.Second)

	if !scheduler.ShouldSwap(p, now) {
		t.Error("ShouldSwap() = false, want true after 11s idle")
	}

	p.LastExecution = now.Add(-5 * time.Second)
	if scheduler.ShouldSwap(p, now) {
		t.Error("ShouldSwap() = true, want false after 5s idle")
	}

	p.State = process.Running
	p.LastExecution = now.Add(-1 * time.Hour)
	if scheduler.ShouldSwap(p, now) {
		t.Error("ShouldSwap() = true for RUNNING, want false")
	}
}

func TestMLFQScheduler_AutoBoostViaTokens(t *testing.T) {
	t.Parallel()

	table, sched, _ := newMLFQTestScheduler()

	runner := sampleProcess("pid-run")
	runner.Priority = 10
	waiter := sampleProcess("pid-wait")
	waiter.Priority = 1

	for _, p := range []process.Process{runner, waiter} {
		if err := table.Create(p); err != nil {
			t.Fatalf("Create(%s) error = %v", p.PID, err)
		}
	}
	if err := sched.Enqueue(runner); err != nil {
		t.Fatalf("Enqueue(runner) error = %v", err)
	}
	if err := sched.Enqueue(waiter); err != nil {
		t.Fatalf("Enqueue(waiter) error = %v", err)
	}
	if _, err := sched.Schedule(); err != nil {
		t.Fatalf("Schedule(runner) error = %v", err)
	}

	levelBefore, ok := sched.QueueLevel("pid-wait")
	if !ok || levelBefore != 3 {
		t.Fatalf("waiter level = %d, ok=%v, want 3 true", levelBefore, ok)
	}

	for i := 0; i < scheduler.BoostEveryTokens; i++ {
		if err := sched.RecordToken("pid-run"); err != nil {
			t.Fatalf("RecordToken(%d) error = %v", i, err)
		}
	}

	levelAfter, ok := sched.QueueLevel("pid-wait")
	if !ok || levelAfter != 0 {
		t.Errorf("waiter level after boost = %d, ok=%v, want 0 true", levelAfter, ok)
	}
	depths := sched.QueueDepths()
	if depths[0] != 1 {
		t.Errorf("Q0 depth after boost = %d, want 1", depths[0])
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
