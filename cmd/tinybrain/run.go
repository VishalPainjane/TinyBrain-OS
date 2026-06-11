package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
	"github.com/VishalPainjane/TinyBrain-OS/internal/inference/llama"
	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

// runRun executes one-shot load → generate → unload and returns an exit code.
func runRun(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	modelID := fs.String("model", "", "model ID from registry (required)")
	prompt := fs.String("prompt", "", "prompt text (required)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *modelID == "" || *prompt == "" {
		fmt.Fprintln(stderr, "usage: tinybrain run --model ID --prompt TEXT")
		return 2
	}

	if !cgoEnabled() {
		fmt.Fprintln(stderr, "error: CGO is disabled — set CGO_ENABLED=1 and build llama.cpp for inference")
		return 1
	}

	profile, err := hardware.ProbeAndClassify()
	if err != nil {
		fmt.Fprintf(stderr, "probe: %v\n", err)
		return 1
	}

	reg, err := openRegistry()
	if err != nil {
		fmt.Fprintf(stderr, "open registry: %v\n", err)
		return 1
	}
	defer reg.Close()

	resolver := runtime.NewRegistryResolver(reg)
	cfg := llama.ConfigFromProbe(profile.Probe)
	if ng, ok := ngLayersFromEnv(); ok {
		cfg.NGLayers = ng
	}
	cfg.GreedySampler = true

	provider := llama.NewLlamaProvider(resolver, cfg)
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(16)
	subscribeRunEvents(bus, stderr)

	rt := runtime.NewIntegratedModelRuntime(provider, ld, resolver, bus)

	fmt.Fprintf(stderr, "TinyBrain %s | run model=%s profile=%s backend=%s\n",
		Version, *modelID, profile.Name, profile.Probe.Backend)

	loadStart := time.Now()
	if err := rt.LoadModel(*modelID); err != nil {
		fmt.Fprintf(stderr, "[load] failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stderr, "[load] %s … %.1fs\n", *modelID, time.Since(loadStart).Seconds())

	genStart := time.Now()
	resp, err := rt.Generate(runtime.GenerateRequest{
		ModelID: *modelID,
		Prompt:  *prompt,
	})
	genElapsed := time.Since(genStart)
	if err != nil {
		fmt.Fprintf(stderr, "[gen] failed: %v\n", err)
		_ = rt.UnloadModel(*modelID)
		return 1
	}
	tps := 0.0
	if genElapsed > 0 && resp.TokensProduced > 0 {
		tps = float64(resp.TokensProduced) / genElapsed.Seconds()
	}
	fmt.Fprintf(stderr, "[gen] tokens=%d elapsed=%.2fs tps=%.1f\n", resp.TokensProduced, genElapsed.Seconds(), tps)

	fmt.Fprintln(stdout, resp.Output)

	if err := rt.UnloadModel(*modelID); err != nil {
		fmt.Fprintf(stderr, "[unload] failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stderr, "[unload] %s\n", *modelID)
	return 0
}

func subscribeRunEvents(bus events.EventBus, stderr io.Writer) {
	bus.Subscribe(events.TypeModelLoaded, func(ev events.Event) {
		payload, ok := ev.Payload.(events.ModelLoadedPayload)
		if !ok {
			return
		}
		fmt.Fprintf(stderr, "[event] ModelLoaded %s\n", payload.ModelID)
	})
	bus.Subscribe(events.TypeModelUnloaded, func(ev events.Event) {
		payload, ok := ev.Payload.(events.ModelUnloadedPayload)
		if !ok {
			return
		}
		fmt.Fprintf(stderr, "[event] ModelUnloaded %s\n", payload.ModelID)
	})
}
