package runtime

import (
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
)

type modelRuntime struct {
	provider InferenceProvider
	bus      events.EventBus
}

// NewModelRuntime returns a ModelRuntime that delegates to provider and publishes lifecycle events.
func NewModelRuntime(provider InferenceProvider, bus events.EventBus) ModelRuntime {
	return &modelRuntime{
		provider: provider,
		bus:      bus,
	}
}

// LoadModel loads modelID via the provider and emits ModelLoaded on success.
func (r *modelRuntime) LoadModel(modelID string) error {
	if err := r.provider.LoadModel(modelID); err != nil {
		return err
	}

	r.bus.Publish(events.NewEvent(events.TypeModelLoaded, events.ModelLoadedPayload{
		ModelID: modelID,
	}, time.Now()))
	return nil
}

// UnloadModel unloads modelID via the provider and emits ModelUnloaded on success.
func (r *modelRuntime) UnloadModel(modelID string) error {
	if err := r.provider.UnloadModel(modelID); err != nil {
		return err
	}

	r.bus.Publish(events.NewEvent(events.TypeModelUnloaded, events.ModelUnloadedPayload{
		ModelID: modelID,
	}, time.Now()))
	return nil
}

// Generate delegates inference to the provider.
func (r *modelRuntime) Generate(req GenerateRequest) (GenerateResponse, error) {
	return r.provider.Generate(req)
}

// SaveContext delegates context persistence to the provider.
func (r *modelRuntime) SaveContext(id string) error {
	return r.provider.SaveContext(id)
}

// RestoreContext delegates context restoration to the provider.
func (r *modelRuntime) RestoreContext(id string) error {
	return r.provider.RestoreContext(id)
}
