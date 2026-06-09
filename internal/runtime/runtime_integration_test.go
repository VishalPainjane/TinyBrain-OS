//go:build cgo && integration

package runtime_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/inference/llama"
	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func newLlamaIntegratedRuntime(t *testing.T) (runtime.ModelRuntime, *loader.StubLoader, events.EventBus, string) {
	t.Helper()

	path := os.Getenv("TB_TEST_GGUF_PATH")
	if path == "" {
		t.Skip("TB_TEST_GGUF_PATH not set")
	}

	modelID := "test-model"
	reg := registry.NewModelRegistry()
	if err := reg.RegisterModel(registry.ModelDefinition{
		ID:   modelID,
		Path: path,
	}); err != nil {
		t.Fatalf("RegisterModel() error = %v", err)
	}

	resolver := runtime.NewRegistryResolver(reg)
	cfg := llama.DefaultConfig()
	cfg.GreedySampler = true
	provider := llama.NewLlamaProvider(resolver, cfg)
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(8)
	rt := runtime.NewIntegratedModelRuntime(provider, ld, resolver, bus)
	return rt, ld, bus, modelID
}

func TestIntegratedRuntime_Llama_LoadGenerateUnload(t *testing.T) {
	rt, ld, bus, modelID := newLlamaIntegratedRuntime(t)

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

	if err := rt.LoadModel(modelID); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}

	state, err := ld.State(modelID)
	if err != nil {
		t.Fatalf("State() after load error = %v", err)
	}
	if state != loader.StateActive {
		t.Fatalf("loader state after load = %s, want ACTIVE", state)
	}

	resp, err := rt.Generate(runtime.GenerateRequest{
		ModelID: modelID,
		Prompt:  "Hello",
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Output == "" {
		t.Fatal("Generate() Output is empty")
	}
	if resp.TokensProduced <= 0 {
		t.Fatalf("TokensProduced = %d, want > 0", resp.TokensProduced)
	}
	t.Logf("Generate output (%d tokens): %q", resp.TokensProduced, resp.Output)

	if err := rt.UnloadModel(modelID); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}

	state, err = ld.State(modelID)
	if err != nil {
		t.Fatalf("State() after unload error = %v", err)
	}
	if state != loader.StateUnloaded {
		t.Fatalf("loader state after unload = %s, want UNLOADED", state)
	}

	waitFor(t, &wg, 5*time.Second)

	mu.Lock()
	defer mu.Unlock()
	if len(loaded) != 1 || loaded[0] != modelID {
		t.Errorf("loaded events = %v, want [%s]", loaded, modelID)
	}
	if len(unloaded) != 1 || unloaded[0] != modelID {
		t.Errorf("unloaded events = %v, want [%s]", unloaded, modelID)
	}
}
