package swap

import (
	"errors"
	"fmt"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/kv"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
)

var (
	// ErrNoKVCache is returned when a process has no KV cache ID.
	ErrNoKVCache = errors.New("process has no kv cache")
	// ErrCannotSwapRunning is returned when SwapOut is called on a running process.
	ErrCannotSwapRunning = errors.New("cannot swap running process")
	// ErrNotIdle is returned when the idle heuristic blocks swap-out.
	ErrNotIdle = errors.New("process not idle long enough to swap")
	// ErrNotHibernated is returned when SwapIn is called on a non-hibernated process.
	ErrNotHibernated = errors.New("process is not hibernated")
)

// Manager moves process KV state across memory tiers.
// See docs/architecture/memory.md and tasks/012-swap-manager.md.
type Manager interface {
	SwapOut(pid string) error
	SwapIn(pid string) error
}

// StubManager coordinates KV tier moves with process state and swap lifecycle events.
type StubManager struct {
	table *process.ProcessTable
	kv    kv.Manager
	bus   events.EventBus
	now   func() time.Time
}

// NewStubManager returns a swap manager using table, kv, and bus.
func NewStubManager(table *process.ProcessTable, kvm kv.Manager, bus events.EventBus) *StubManager {
	return &StubManager{
		table: table,
		kv:    kvm,
		bus:   bus,
		now:   time.Now,
	}
}

// SwapOut moves pid's KV from VRAM to RAM and transitions the process to HIBERNATED.
// The process must pass scheduler.ShouldSwap (idle > 10s) and must not be RUNNING.
func (m *StubManager) SwapOut(pid string) error {
	p, err := m.table.Get(pid)
	if err != nil {
		return err
	}
	if p.KVCacheID == "" {
		return ErrNoKVCache
	}
	if p.State == process.Running {
		return ErrCannotSwapRunning
	}
	if !scheduler.ShouldSwap(p, m.now()) {
		return ErrNotIdle
	}

	m.publishSwapStarted(kv.TierVRAM.String(), kv.TierRAM.String())

	if err := m.kv.Save(p.KVCacheID); err != nil {
		return err
	}
	if err := m.table.UpdateState(pid, process.Hibernated); err != nil {
		return err
	}

	m.publishSwapCompleted(kv.TierVRAM.String(), kv.TierRAM.String())
	return nil
}

// SwapIn restores pid's KV from RAM to VRAM and transitions HIBERNATED to READY.
func (m *StubManager) SwapIn(pid string) error {
	p, err := m.table.Get(pid)
	if err != nil {
		return err
	}
	if p.KVCacheID == "" {
		return ErrNoKVCache
	}
	if p.State != process.Hibernated {
		return fmt.Errorf("%w: %s in state %s", ErrNotHibernated, pid, p.State)
	}

	block, err := m.kv.Get(p.KVCacheID)
	if err != nil {
		return err
	}
	if block.Tier != kv.TierRAM {
		return fmt.Errorf("%w: kv %s in tier %s", kv.ErrInvalidTier, p.KVCacheID, block.Tier)
	}

	m.publishSwapStarted(kv.TierRAM.String(), kv.TierVRAM.String())

	if err := m.kv.Load(p.KVCacheID); err != nil {
		return err
	}
	if err := m.table.UpdateState(pid, process.Ready); err != nil {
		return err
	}

	m.publishSwapCompleted(kv.TierRAM.String(), kv.TierVRAM.String())
	return nil
}

func (m *StubManager) publishSwapStarted(fromTier, toTier string) {
	m.publish(events.TypeSwapStarted, events.SwapStartedPayload{
		FromModelID: fromTier,
		ToModelID:   toTier,
	})
}

func (m *StubManager) publishSwapCompleted(fromTier, toTier string) {
	m.publish(events.TypeSwapCompleted, events.SwapCompletedPayload{
		FromModelID: fromTier,
		ToModelID:   toTier,
	})
}

func (m *StubManager) publish(eventType events.Type, payload any) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(events.NewEvent(eventType, payload, m.now()))
}
