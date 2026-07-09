package runtime_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func TestStubProvider_LoadUnload(t *testing.T) {
	t.Parallel()

	provider := runtime.NewStubProvider()

	if err := provider.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	if err := provider.LoadModel("model-a"); !errors.Is(err, runtime.ErrModelAlreadyLoaded) {
		t.Fatalf("LoadModel() duplicate error = %v, want ErrModelAlreadyLoaded", err)
	}
	if err := provider.UnloadModel("model-a"); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}
	if err := provider.UnloadModel("model-a"); !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("UnloadModel() missing error = %v, want ErrModelNotLoaded", err)
	}
}

func TestStubProvider_Generate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*runtime.StubProvider) error
		req     runtime.GenerateRequest
		wantOut string
		wantErr error
	}{
		{
			name: "loaded model returns stub response",
			setup: func(p *runtime.StubProvider) error {
				return p.LoadModel("model-a")
			},
			req: runtime.GenerateRequest{
				ModelID: "model-a",
				Prompt:  "hello",
			},
			wantOut: "stub response for model-a",
		},
		{
			name:  "unloaded model returns error",
			setup: func(*runtime.StubProvider) error { return nil },
			req: runtime.GenerateRequest{
				ModelID: "model-a",
				Prompt:  "hello",
			},
			wantErr: runtime.ErrModelNotLoaded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := runtime.NewStubProvider()
			if err := tt.setup(provider); err != nil {
				t.Fatalf("setup error = %v", err)
			}

			resp, err := provider.Generate(tt.req)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Generate() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if resp.ModelID != tt.req.ModelID {
				t.Errorf("ModelID = %q, want %q", resp.ModelID, tt.req.ModelID)
			}
			if resp.Output != tt.wantOut {
				t.Errorf("Output = %q, want %q", resp.Output, tt.wantOut)
			}
			if resp.TokensProduced != 1 {
				t.Errorf("TokensProduced = %d, want 1", resp.TokensProduced)
			}
		})
	}
}

func TestStubProvider_SaveRestoreContext(t *testing.T) {
	t.Parallel()

	provider := runtime.NewStubProvider()

	if err := provider.SaveContext("m1", "ctx-1"); err != nil {
		t.Fatalf("SaveContext() error = %v", err)
	}
	if err := provider.RestoreContext("m1", "ctx-1"); err != nil {
		t.Fatalf("RestoreContext() error = %v", err)
	}
	if err := provider.RestoreContext("m1", "missing"); !errors.Is(err, runtime.ErrContextNotFound) {
		t.Fatalf("RestoreContext() error = %v, want ErrContextNotFound", err)
	}
}

func TestModelRuntime_LoadUnload_PublishesEvents(t *testing.T) {
	t.Parallel()

	bus := events.NewChannelBus(8)
	rt := runtime.NewModelRuntime(runtime.NewStubProvider(), bus)

	var loaded []string
	var unloaded []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	unsubLoaded := bus.Subscribe(events.TypeModelLoaded, func(ev events.Event) {
		payload, ok := ev.Payload.(events.ModelLoadedPayload)
		if !ok {
			t.Errorf("ModelLoaded payload type = %T", ev.Payload)
			return
		}
		mu.Lock()
		loaded = append(loaded, payload.ModelID)
		mu.Unlock()
		wg.Done()
	})
	defer unsubLoaded()

	unsubUnloaded := bus.Subscribe(events.TypeModelUnloaded, func(ev events.Event) {
		payload, ok := ev.Payload.(events.ModelUnloadedPayload)
		if !ok {
			t.Errorf("ModelUnloaded payload type = %T", ev.Payload)
			return
		}
		mu.Lock()
		unloaded = append(unloaded, payload.ModelID)
		mu.Unlock()
		wg.Done()
	})
	defer unsubUnloaded()

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	if err := rt.UnloadModel("model-a"); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}

	waitFor(t, &wg, 2*time.Second)

	mu.Lock()
	defer mu.Unlock()
	if len(loaded) != 1 || loaded[0] != "model-a" {
		t.Errorf("loaded events = %v, want [model-a]", loaded)
	}
	if len(unloaded) != 1 || unloaded[0] != "model-a" {
		t.Errorf("unloaded events = %v, want [model-a]", unloaded)
	}
}

func TestModelRuntime_FailurePath_NoLifecycleEvents(t *testing.T) {
	t.Parallel()

	bus := events.NewChannelBus(8)
	rt := runtime.NewModelRuntime(runtime.NewStubProvider(), bus)

	var eventCount int
	var mu sync.Mutex

	record := func(events.Event) {
		mu.Lock()
		eventCount++
		mu.Unlock()
	}

	unsubLoaded := bus.Subscribe(events.TypeModelLoaded, record)
	defer unsubLoaded()
	unsubUnloaded := bus.Subscribe(events.TypeModelUnloaded, record)
	defer unsubUnloaded()

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	if err := rt.LoadModel("model-a"); !errors.Is(err, runtime.ErrModelAlreadyLoaded) {
		t.Fatalf("duplicate LoadModel() error = %v, want ErrModelAlreadyLoaded", err)
	}
	if err := rt.UnloadModel("model-a"); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}
	if err := rt.UnloadModel("model-a"); !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("duplicate UnloadModel() error = %v, want ErrModelNotLoaded", err)
	}

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if eventCount != 2 {
		t.Errorf("lifecycle event count = %d, want 2 (one load, one unload)", eventCount)
	}
}

func TestModelRuntime_Generate_DelegatesToStub(t *testing.T) {
	t.Parallel()

	rt := runtime.NewModelRuntime(runtime.NewStubProvider(), events.NewChannelBus(8))

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}

	resp, err := rt.Generate(runtime.GenerateRequest{
		ModelID: "model-a",
		Prompt:  "ping",
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Output != "stub response for model-a" {
		t.Errorf("Output = %q, want stub response for model-a", resp.Output)
	}
}

func TestModelRuntime_SwapStubModels(t *testing.T) {
	t.Parallel()

	rt := runtime.NewModelRuntime(runtime.NewStubProvider(), events.NewChannelBus(8))

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel(model-a) error = %v", err)
	}

	respA, err := rt.Generate(runtime.GenerateRequest{ModelID: "model-a", Prompt: "a"})
	if err != nil {
		t.Fatalf("Generate(model-a) error = %v", err)
	}
	if respA.ModelID != "model-a" {
		t.Errorf("model-a response ModelID = %q", respA.ModelID)
	}

	if err := rt.UnloadModel("model-a"); err != nil {
		t.Fatalf("UnloadModel(model-a) error = %v", err)
	}
	if err := rt.LoadModel("model-b"); err != nil {
		t.Fatalf("LoadModel(model-b) error = %v", err)
	}

	respB, err := rt.Generate(runtime.GenerateRequest{ModelID: "model-b", Prompt: "b"})
	if err != nil {
		t.Fatalf("Generate(model-b) error = %v", err)
	}
	if respB.Output != "stub response for model-b" {
		t.Errorf("Output = %q, want stub response for model-b", respB.Output)
	}
}

func TestModelRuntime_SaveRestoreContext(t *testing.T) {
	t.Parallel()

	rt := runtime.NewModelRuntime(runtime.NewStubProvider(), events.NewChannelBus(8))

	if err := rt.SaveContext("m1", "ctx-1"); err != nil {
		t.Fatalf("SaveContext() error = %v", err)
	}
	if err := rt.RestoreContext("m1", "ctx-1"); err != nil {
		t.Fatalf("RestoreContext() error = %v", err)
	}
}

func waitFor(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatal("timed out waiting for events")
	}
}
