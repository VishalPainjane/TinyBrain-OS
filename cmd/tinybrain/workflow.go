package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/agents"
	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
	"github.com/VishalPainjane/TinyBrain-OS/internal/inference/llama"
	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/router"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
)

// runWorkflow executes a sequential 2-agent task pipeline.
func runWorkflow(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("workflow", flag.ContinueOnError)
	fs.SetOutput(stderr)
	prompt := fs.String("prompt", "", "initial prompt text for the planner (required)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	
	if *prompt == "" {
		fmt.Fprintln(stderr, "usage: tinybrain workflow --prompt TEXT")
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

	if err := registry.LoadModelsYAML(filepath.Join("testdata", "models.yaml"), reg); err != nil {
		if !strings.Contains(err.Error(), "duplicate id") && !strings.Contains(err.Error(), "duplicate ID") {
			fmt.Fprintf(stderr, "load models: %v\n", err)
			return 1
		}
	}

	areg := registry.NewAgentRegistry()
	if err := registry.LoadAgentsYAML(filepath.Join("testdata", "fleet.yaml"), areg); err != nil {
		fmt.Fprintf(stderr, "load fleet: %v\n", err)
		return 1
	}

	resolver := runtime.NewRegistryResolver(reg)
	cfg := llama.ConfigFromProbe(profile.Probe)
	if ng, ok := ngLayersFromEnv(); ok {
		cfg.NGLayers = ng
	}
	cfg.GreedySampler = true

	provider := llama.NewLlamaProvider(resolver, cfg)
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(16)

	rt := runtime.NewIntegratedModelRuntime(provider, ld, resolver, bus)

	if err := rt.LoadModel("tinyllama-q4"); err != nil {
		fmt.Fprintf(stderr, "error: load model: %v\n", err)
		return 1
	}
	defer rt.UnloadModel("tinyllama-q4")

	ptab := process.NewProcessTable()
	queue := scheduler.NewMLFQQueue()
	sched := scheduler.NewMLFQScheduler(ptab, queue)
	coord := scheduler.NewEventCoordinator(sched, bus, ptab)
	defer coord.Stop()

	exec := agents.NewExecutor(rt, bus)
	listener := agents.NewEventListener(bus, exec, ptab, areg)
	defer listener.Stop()

	rtr := router.NewRouter(bus, areg, ptab)
	defer rtr.Stop()

	fmt.Fprintf(stderr, "TinyBrain %s | workflow profile=%s backend=%s\n",
		Version, profile.Name, profile.Probe.Backend)

	done := make(chan string)
	
	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())

	// Trace events
	bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
		payload := ev.Payload.(events.TaskCompletedPayload)
		if payload.TaskID == taskID {
			fmt.Fprintln(stderr, "-> Event: TaskCompleted")
			done <- payload.Result
		}
	})

	// Submit task
	fmt.Fprintln(stderr, "--- Submitting task ---")
	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  taskID,
		Input:   *prompt,
		AgentID: "sample-alpha",
	}, time.Now()))

	var result string
	select {
	case result = <-done:
	case <-time.After(5 * time.Minute):
		fmt.Fprintln(stderr, "error: task timed out")
		return 1
	}

	// Extract generated text from structured JSON if possible
	var parsed struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil || parsed.Text == "" {
		parsed.Text = result
	}

	fmt.Fprintf(stdout, "\n[Final Output]\n%s\n", parsed.Text)

	return 0
}
