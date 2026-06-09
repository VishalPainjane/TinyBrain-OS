package runtime

import (
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
)

type modelRuntime struct {
	provider InferenceProvider
	loader   loader.Loader
	resolver ModelResolver
	bus      events.EventBus
}

// NewModelRuntime returns a loader-less ModelRuntime that delegates to provider and publishes lifecycle events.
func NewModelRuntime(provider InferenceProvider, bus events.EventBus) ModelRuntime {
	return &modelRuntime{
		provider: provider,
		bus:      bus,
	}
}

// NewIntegratedModelRuntime returns a ModelRuntime that orchestrates resolver, loader, and provider.
func NewIntegratedModelRuntime(provider InferenceProvider, ld loader.Loader, resolver ModelResolver, bus events.EventBus) ModelRuntime {
	return &modelRuntime{
		provider: provider,
		loader:   ld,
		resolver: resolver,
		bus:      bus,
	}
}

// LoadModel loads modelID and emits ModelLoaded on success.
// Integrated path: resolve → loader.Load → provider.LoadModel, with loader rollback on provider failure.
func (r *modelRuntime) LoadModel(modelID string) error {
	if r.loader == nil {
		if err := r.provider.LoadModel(modelID); err != nil {
			return err
		}
		r.publishModelLoaded(modelID)
		return nil
	}

	spec, err := r.resolver.Resolve(modelID)
	if err != nil {
		return err
	}

	if err := r.loader.Load(modelID, spec.Path); err != nil {
		return err
	}

	if err := r.provider.LoadModel(modelID); err != nil {
		_ = r.loader.Unload(modelID)
		return err
	}

	r.publishModelLoaded(modelID)
	return nil
}

// UnloadModel unloads modelID and emits ModelUnloaded on success.
// Integrated path: provider.UnloadModel → loader.Unload.
func (r *modelRuntime) UnloadModel(modelID string) error {
	if r.loader == nil {
		if err := r.provider.UnloadModel(modelID); err != nil {
			return err
		}
		r.publishModelUnloaded(modelID)
		return nil
	}

	if err := r.provider.UnloadModel(modelID); err != nil {
		return err
	}

	if err := r.loader.Unload(modelID); err != nil {
		return err
	}

	r.publishModelUnloaded(modelID)
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

func (r *modelRuntime) publishModelLoaded(modelID string) {
	r.bus.Publish(events.NewEvent(events.TypeModelLoaded, events.ModelLoadedPayload{
		ModelID: modelID,
	}, time.Now()))
}

func (r *modelRuntime) publishModelUnloaded(modelID string) {
	r.bus.Publish(events.NewEvent(events.TypeModelUnloaded, events.ModelUnloadedPayload{
		ModelID: modelID,
	}, time.Now()))
}
