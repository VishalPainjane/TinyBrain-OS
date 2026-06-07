# Contract: Scheduler

## Owner

Package: `internal/scheduler/` (not yet implemented)

## Consumers

- API / Router — submit tasks for scheduling
- Telemetry — read-only queue and process metrics

## Scheduler Interface

```go
// Documentation only
type Scheduler interface {
    Enqueue(process *Process) error
    Schedule() error
    Preempt(pid string) error
    Boost() error
}
```

## Queue Contract (v1 FIFO → v0.7 MLFQ)

v1: single FIFO queue behind `Queue` interface.
v0.7: MLFQ with Q0–Q3, token quanta, boost/aging.

```go
// Documentation only
type Queue interface {
    Enqueue(process *Process) error
    Dequeue() (*Process, error)
    Peek() (*Process, error)
    Depth() int
}
```

## Responsibilities

- Admit processes to queues
- Select next process to run
- Preempt lower-priority work
- Request runtime load/unload (via interface — never direct)
- Apply swap heuristic (idle > 10s)

## Must NOT Own

- Direct model or inference access (INV-001)
- Process state storage (kernel owns table; scheduler requests transitions)
- Agent execution
- Registry writes

## Related

- [architecture/scheduler.md](../architecture/scheduler.md)
- [contracts/process.md](process.md)
- [contracts/runtime.md](runtime.md)
- [tasks/010-scheduler.md](../../tasks/010-scheduler.md)
- [planning/architecture-evolution/migration-paths.md](../../planning/architecture-evolution/migration-paths.md)

---
**Layer:** contract
