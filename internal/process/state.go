package process

// ProcessState is the lifecycle state of a TinyBrain agent process.
// Every process is always in exactly one of these states.
type ProcessState int

const (
	New ProcessState = iota
	Ready
	Running
	Waiting
	Preempted
	Hibernated
	Terminated
)

var stateNames = [...]string{
	New:        "NEW",
	Ready:      "READY",
	Running:    "RUNNING",
	Waiting:    "WAITING",
	Preempted:  "PREEMPTED",
	Hibernated: "HIBERNATED",
	Terminated: "TERMINATED",
}

// String returns the canonical state name from the process model.
func (s ProcessState) String() string {
	if !s.Valid() {
		return "UNKNOWN"
	}
	return stateNames[s]
}

// Valid reports whether s is one of the defined process states.
func (s ProcessState) Valid() bool {
	return s >= New && s <= Terminated
}

// All returns every defined process state in lifecycle order.
func All() []ProcessState {
	return []ProcessState{
		New,
		Ready,
		Running,
		Waiting,
		Preempted,
		Hibernated,
		Terminated,
	}
}
