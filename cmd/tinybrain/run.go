package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
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

// runRun executes the task via the event pipeline and returns an exit code.
func runRun(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	agentID := fs.String("agent", "", "agent ID from fleet (required)")
	modelIDFallback := fs.String("model", "", "legacy alias for agent")
	prompt := fs.String("prompt", "", "prompt text (required)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	
	targetAgent := *agentID
	if targetAgent == "" {
		targetAgent = *modelIDFallback
	}

	if targetAgent == "" || *prompt == "" {
		fmt.Fprintln(stderr, "usage: tinybrain run --agent ID --prompt TEXT")
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

	// 1. Model Registry
	reg, err := openRegistry()
	if err != nil {
		fmt.Fprintf(stderr, "open registry: %v\n", err)
		return 1
	}
	defer reg.Close()

	// 2. Agent Registry (Fleet)
	areg := registry.NewAgentRegistry()
	if err := registry.LoadAgentsYAML(filepath.Join("testdata", "fleet.yaml"), areg); err != nil {
		fmt.Fprintf(stderr, "load fleet: %v\n", err)
		return 1
	}

	// 3. Inference & Runtime
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

	// 4. Kernel (Process Table)
	ptab := process.NewProcessTable()

	// 5. Scheduler
	queue := scheduler.NewMLFQQueue()
	sched := scheduler.NewMLFQScheduler(ptab, queue)
	coord := scheduler.NewEventCoordinator(sched, bus, ptab)
	defer coord.Stop()

	// 6. Agent Executor
	exec := agents.NewExecutor(rt, bus)
	listener := agents.NewEventListener(bus, exec, ptab, areg)
	defer listener.Stop()

	// 7. Router
	rtr := router.NewRouter(bus, areg, ptab)
	defer rtr.Stop()

	fmt.Fprintf(stderr, "TinyBrain %s | event pipeline agent=%s profile=%s backend=%s\n",
		Version, targetAgent, profile.Name, profile.Probe.Backend)

	done := make(chan struct{})
	
	// Trace events
	bus.Subscribe(events.TypeTaskCreated, func(ev events.Event) { fmt.Fprintln(stderr, "-> Event: TaskCreated") })
	bus.Subscribe(events.TypeProcessSpawned, func(ev events.Event) { fmt.Fprintln(stderr, "-> Event: ProcessSpawned") })
	bus.Subscribe(events.TypeProcessStateChanged, func(ev events.Event) {
		payload := ev.Payload.(events.ProcessStateChangedPayload)
		fmt.Fprintf(stderr, "-> Event: ProcessStateChanged (%s -> %s)\n", payload.OldState, payload.NewState)
	})
	bus.Subscribe(events.TypeAgentStarted, func(ev events.Event) { fmt.Fprintln(stderr, "-> Event: AgentStarted") })
	bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
		fmt.Fprintln(stderr, "-> Event: TaskCompleted")
		close(done)
	})

	// Submit task
	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  taskID,
		Input:   *prompt,
		AgentID: targetAgent,
	}, time.Now()))

	<-done
	return 0
}
