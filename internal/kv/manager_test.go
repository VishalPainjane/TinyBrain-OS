package kv_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/kv"
)

func newTestManager(t *testing.T) (*kv.StubManager, events.EventBus) {
	t.Helper()
	bus := events.NewChannelBus(8)
	mgr := kv.NewStubManager(bus)
	return mgr, bus
}

func TestStubManager_AllocateGet(t *testing.T) {
	t.Parallel()

	mgr, _ := newTestManager(t)

	if err := mgr.Allocate("kv-1", "pid-a", 4096); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}

	got, err := mgr.Get("kv-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.KVCacheID != "kv-1" || got.PID != "pid-a" || got.SizeBytes != 4096 {
		t.Errorf("Get() = %+v, want kv-1 pid-a size 4096", got)
	}
	if got.Tier != kv.TierVRAM {
		t.Errorf("Tier = %s, want VRAM", got.Tier)
	}
}

func TestStubManager_SavePublishesKVStored(t *testing.T) {
	t.Parallel()

	mgr, bus := newTestManager(t)
	if err := mgr.Allocate("kv-1", "pid-a", 1024); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2) // Wait for both KVCompressed and KVStored
	
	unsub1 := bus.Subscribe(events.TypeKVCompressed, func(ev events.Event) {
		defer wg.Done()
		payload, ok := ev.Payload.(events.KVCompressedPayload)
		if !ok {
			t.Fatalf("Payload type = %T, want KVCompressedPayload", ev.Payload)
		}
		if payload.KVCacheID != "kv-1" {
			t.Errorf("Payload = %+v, want kv-1", payload)
		}
		if payload.CompressionRatio <= 1.0 {
			t.Errorf("Expected CompressionRatio > 1.0, got %f", payload.CompressionRatio)
		}
	})
	defer unsub1()
	
	unsub2 := bus.Subscribe(events.TypeKVStored, func(ev events.Event) {
		defer wg.Done()
		payload, ok := ev.Payload.(events.KVStoredPayload)
		if !ok {
			t.Fatalf("Payload type = %T, want KVStoredPayload", ev.Payload)
		}
		if payload.KVCacheID != "kv-1" || payload.PID != "pid-a" {
			t.Errorf("Payload = %+v, want kv-1 pid-a", payload)
		}
	})
	defer unsub2()

	if err := mgr.Save("kv-1"); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	waitForSubscriber(t, &wg)

	got, err := mgr.Get("kv-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Tier != kv.TierRAM {
		t.Errorf("Tier after Save = %s, want RAM", got.Tier)
	}
}

func TestStubManager_LoadPublishesKVLoaded(t *testing.T) {
	t.Parallel()

	mgr, bus := newTestManager(t)
	if err := mgr.Allocate("kv-1", "pid-a", 1024); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}
	if err := mgr.Save("kv-1"); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2) // Wait for both KVDecompressed and KVLoaded
	
	unsub1 := bus.Subscribe(events.TypeKVDecompressed, func(ev events.Event) {
		defer wg.Done()
		payload, ok := ev.Payload.(events.KVDecompressedPayload)
		if !ok {
			t.Fatalf("Payload type = %T, want KVDecompressedPayload", ev.Payload)
		}
		if payload.KVCacheID != "kv-1" {
			t.Errorf("Payload = %+v, want kv-1", payload)
		}
	})
	defer unsub1()
	
	unsub2 := bus.Subscribe(events.TypeKVLoaded, func(ev events.Event) {
		defer wg.Done()
		payload, ok := ev.Payload.(events.KVLoadedPayload)
		if !ok {
			t.Fatalf("Payload type = %T, want KVLoadedPayload", ev.Payload)
		}
		if payload.KVCacheID != "kv-1" || payload.PID != "pid-a" {
			t.Errorf("Payload = %+v, want kv-1 pid-a", payload)
		}
	})
	defer unsub2()

	if err := mgr.Load("kv-1"); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	waitForSubscriber(t, &wg)

	got, err := mgr.Get("kv-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Tier != kv.TierVRAM {
		t.Errorf("Tier after Load = %s, want VRAM", got.Tier)
	}
}

func TestStubManager_SaveLoadErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(m *kv.StubManager) error
		op      func(m *kv.StubManager) error
		wantErr error
	}{
		{
			name: "save missing block",
			setup: func(m *kv.StubManager) error {
				return nil
			},
			op: func(m *kv.StubManager) error {
				return m.Save("missing")
			},
			wantErr: kv.ErrBlockNotFound,
		},
		{
			name: "load from VRAM",
			setup: func(m *kv.StubManager) error {
				return m.Allocate("kv-1", "pid-a", 512)
			},
			op: func(m *kv.StubManager) error {
				return m.Load("kv-1")
			},
			wantErr: kv.ErrInvalidTier,
		},
		{
			name: "duplicate allocate",
			setup: func(m *kv.StubManager) error {
				return m.Allocate("kv-1", "pid-a", 512)
			},
			op: func(m *kv.StubManager) error {
				return m.Allocate("kv-1", "pid-b", 512)
			},
			wantErr: kv.ErrBlockExists,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mgr, _ := newTestManager(t)
			if err := tt.setup(mgr); err != nil {
				t.Fatalf("setup() error = %v", err)
			}
			if err := tt.op(mgr); !errors.Is(err, tt.wantErr) {
				t.Fatalf("op() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTier_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tier kv.Tier
		want string
	}{
		{tier: kv.TierVRAM, want: "VRAM"},
		{tier: kv.TierRAM, want: "RAM"},
		{tier: kv.TierNVMe, want: "NVMe"},
	}
	for _, tt := range tests {
		if got := tt.tier.String(); got != tt.want {
			t.Errorf("Tier(%d).String() = %q, want %q", tt.tier, got, tt.want)
		}
	}
}

func waitForSubscriber(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("subscriber did not receive event in time")
	}
}
