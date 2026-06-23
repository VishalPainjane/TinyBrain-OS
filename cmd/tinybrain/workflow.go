package main

import (
	"encoding/json"
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

	fmt.Fprintf(stderr, "TinyBrain %s | two-agent workflow profile=%s backend=%s\n",
		Version, profile.Name, profile.Probe.Backend)

	done1 := make(chan string)
	done2 := make(chan string)
	
	task1ID := fmt.Sprintf("task-planner-%d", time.Now().UnixNano())
	task2ID := fmt.Sprintf("task-coder-%d", time.Now().UnixNano())

	// Trace events
	bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
		payload := ev.Payload.(events.TaskCompletedPayload)
		if payload.TaskID == task1ID {
			fmt.Fprintln(stderr, "-> Event: Planner TaskCompleted")
			done1 <- payload.Result
		} else if payload.TaskID == task2ID {
			fmt.Fprintln(stderr, "-> Event: Coder TaskCompleted")
			done2 <- payload.Result
		}
	})

	// Submit task 1 (Planner)
	fmt.Fprintln(stderr, "--- Submitting to sample-alpha (Planner) ---")
	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  task1ID,
		Input:   *prompt,
		AgentID: "sample-alpha",
	}, time.Now()))

	var plannerResult string
	select {
	case plannerResult = <-done1:
	case <-time.After(5 * time.Minute):
		fmt.Fprintln(stderr, "error: planner task timed out")
		return 1
	}

	// Extract generated text from Planner's structured JSON
	var parsed1 struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(plannerResult), &parsed1); err != nil {
		fmt.Fprintf(stderr, "warning: planner output not valid JSON: %v\n", err)
		parsed1.Text = plannerResult
	}

	fmt.Fprintf(stderr, "\n[Planner Output]\n%s\n\n", parsed1.Text)

	// Submit task 2 (Coder) using planner's output
	coderPrompt := fmt.Sprintf("Based on the following plan, write the code:\n%s", parsed1.Text)
	fmt.Fprintln(stderr, "--- Submitting to sample-beta (Coder) ---")
	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  task2ID,
		Input:   coderPrompt,
		AgentID: "sample-beta",
	}, time.Now()))

	var coderResult string
	select {
	case coderResult = <-done2:
	case <-time.After(5 * time.Minute):
		fmt.Fprintln(stderr, "error: coder task timed out")
		return 1
	}

	// Extract generated text from Coder's structured JSON
	var parsed2 struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(coderResult), &parsed2); err != nil {
		fmt.Fprintf(stderr, "warning: coder output not valid JSON: %v\n", err)
		parsed2.Text = coderResult
	}

	fmt.Fprintf(stdout, "\n[Final Output]\n%s\n", parsed2.Text)

	return 0
}
