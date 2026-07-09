package runtime_test

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

type staticResolver struct {
	specs map[string]runtime.ModelSpec
	err   error
}

func (s staticResolver) Resolve(modelID string) (runtime.ModelSpec, error) {
	if s.err != nil {
		return runtime.ModelSpec{}, s.err
	}
	spec, ok := s.specs[modelID]
	if !ok {
		return runtime.ModelSpec{}, registry.ErrNotFound
	}
	return spec, nil
}

type loadFailProvider struct {
	loadErr error
}

func (p *loadFailProvider) LoadModel(string) error {
	return p.loadErr
}

func (p *loadFailProvider) UnloadModel(string) error {
	return runtime.ErrModelNotLoaded
}

func (p *loadFailProvider) Generate(runtime.GenerateRequest) (runtime.GenerateResponse, error) {
	return runtime.GenerateResponse{}, runtime.ErrModelNotLoaded
}

func (p *loadFailProvider) SaveContext(modelID, ctxID string) error {
	return nil
}

func (p *loadFailProvider) RestoreContext(modelID, ctxID string) error {
	return nil
}

func (p *loadFailProvider) FormatChat(modelID string, messages []runtime.ChatMessage, opts runtime.FormatChatOpts) (string, string, error) {
	return "", "", nil
}

func (p *loadFailProvider) GetMetadata(modelID string) (runtime.ModelCapabilities, error) {
	return runtime.ModelCapabilities{ModelID: modelID}, nil
}

func stubModelPath(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte("stub-gguf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func newIntegratedRuntime(t *testing.T, provider runtime.InferenceProvider, resolver staticResolver) (runtime.ModelRuntime, *loader.StubLoader) {
	t.Helper()

	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(8)
	rt := runtime.NewIntegratedModelRuntime(provider, ld, resolver, bus)
	return rt, ld
}

func TestIntegratedRuntime_Load_setsLoaderActive(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	provider := runtime.NewStubProvider()
	rt, ld := newIntegratedRuntime(t, provider, staticResolver{
		specs: map[string]runtime.ModelSpec{
			"model-a": {ID: "model-a", Path: path},
		},
	})

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}

	state, err := ld.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateActive {
		t.Fatalf("loader state = %s, want ACTIVE", state)
	}
}

func TestIntegratedRuntime_Unload_setsLoaderUnloaded(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	provider := runtime.NewStubProvider()
	rt, ld := newIntegratedRuntime(t, provider, staticResolver{
		specs: map[string]runtime.ModelSpec{
			"model-a": {ID: "model-a", Path: path},
		},
	})

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	if err := rt.UnloadModel("model-a"); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}

	state, err := ld.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateUnloaded {
		t.Fatalf("loader state = %s, want UNLOADED", state)
	}
}

func TestIntegratedRuntime_LoadProviderFail_rollbackLoader(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	loadErr := errors.New("provider load failed")
	provider := &loadFailProvider{loadErr: loadErr}
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(8)
	rt := runtime.NewIntegratedModelRuntime(provider, ld, staticResolver{
		specs: map[string]runtime.ModelSpec{
			"model-a": {ID: "model-a", Path: path},
		},
	}, bus)

	var eventCount int
	var mu sync.Mutex
	unsub := bus.Subscribe(events.TypeModelLoaded, func(events.Event) {
		mu.Lock()
		eventCount++
		mu.Unlock()
	})
	defer unsub()

	err := rt.LoadModel("model-a")
	if !errors.Is(err, loadErr) {
		t.Fatalf("LoadModel() error = %v, want %v", err, loadErr)
	}

	state, stateErr := ld.State("model-a")
	if stateErr != nil {
		t.Fatalf("State() error = %v", stateErr)
	}
	if state != loader.StateUnloaded {
		t.Fatalf("loader state after rollback = %s, want UNLOADED", state)
	}

	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if eventCount != 0 {
		t.Fatalf("ModelLoaded event count = %d, want 0", eventCount)
	}
}

func TestIntegratedRuntime_Load_resolverFailure(t *testing.T) {
	t.Parallel()

	provider := runtime.NewStubProvider()
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(8)
	rt := runtime.NewIntegratedModelRuntime(provider, ld, staticResolver{err: registry.ErrNotFound}, bus)

	var eventCount int
	var mu sync.Mutex
	unsub := bus.Subscribe(events.TypeModelLoaded, func(events.Event) {
		mu.Lock()
		eventCount++
		mu.Unlock()
	})
	defer unsub()

	err := rt.LoadModel("model-a")
	if !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("LoadModel() error = %v, want ErrNotFound", err)
	}

	state, err := ld.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateNotLoaded {
		t.Fatalf("loader state = %s, want NOT_LOADED", state)
	}

	_, genErr := rt.Generate(runtime.GenerateRequest{ModelID: "model-a", Prompt: "hi"})
	if !errors.Is(genErr, runtime.ErrModelNotLoaded) {
		t.Fatalf("Generate() error = %v, want ErrModelNotLoaded", genErr)
	}

	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if eventCount != 0 {
		t.Fatalf("ModelLoaded event count = %d, want 0", eventCount)
	}
}

func TestIntegratedRuntime_Load_loaderFailure(t *testing.T) {
	t.Parallel()

	provider := runtime.NewStubProvider()
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(8)
	rt := runtime.NewIntegratedModelRuntime(provider, ld, staticResolver{
		specs: map[string]runtime.ModelSpec{
			"model-a": {ID: "model-a", Path: filepath.Join(t.TempDir(), "missing.gguf")},
		},
	}, bus)

	var eventCount int
	var mu sync.Mutex
	unsub := bus.Subscribe(events.TypeModelLoaded, func(events.Event) {
		mu.Lock()
		eventCount++
		mu.Unlock()
	})
	defer unsub()

	err := rt.LoadModel("model-a")
	if !errors.Is(err, loader.ErrStubPathInaccessible) {
		t.Fatalf("LoadModel() error = %v, want ErrStubPathInaccessible", err)
	}

	state, err := ld.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateNotLoaded {
		t.Fatalf("loader state = %s, want NOT_LOADED", state)
	}

	_, genErr := rt.Generate(runtime.GenerateRequest{ModelID: "model-a", Prompt: "hi"})
	if !errors.Is(genErr, runtime.ErrModelNotLoaded) {
		t.Fatalf("Generate() error = %v, want ErrModelNotLoaded", genErr)
	}

	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if eventCount != 0 {
		t.Fatalf("ModelLoaded event count = %d, want 0", eventCount)
	}
}

func TestIntegratedRuntime_LoadUnload_PublishesEvents(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	provider := runtime.NewStubProvider()
	bus := events.NewChannelBus(8)
	ld := loader.NewStubLoader()
	rt := runtime.NewIntegratedModelRuntime(provider, ld, staticResolver{
		specs: map[string]runtime.ModelSpec{
			"model-a": {ID: "model-a", Path: path},
		},
	}, bus)

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

func TestIntegratedRuntime_Generate_delegates(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	provider := runtime.NewStubProvider()
	rt, _ := newIntegratedRuntime(t, provider, staticResolver{
		specs: map[string]runtime.ModelSpec{
			"model-a": {ID: "model-a", Path: path},
		},
	})

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}

	resp, err := rt.Generate(runtime.GenerateRequest{ModelID: "model-a", Prompt: "ping"})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Output != "stub response for model-a" {
		t.Errorf("Output = %q, want stub response for model-a", resp.Output)
	}
}

func TestModelRuntime_LoaderLess_Unchanged(t *testing.T) {
	t.Parallel()

	bus := events.NewChannelBus(8)
	rt := runtime.NewModelRuntime(runtime.NewStubProvider(), bus)

	if err := rt.LoadModel("model-a"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}

	resp, err := rt.Generate(runtime.GenerateRequest{ModelID: "model-a", Prompt: "ping"})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Output != "stub response for model-a" {
		t.Errorf("Output = %q, want stub response for model-a", resp.Output)
	}

	if err := rt.UnloadModel("model-a"); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}
}
