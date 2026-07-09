package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
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

// TaskRequest is the JSON payload for submitting a task to the daemon.
type TaskRequest struct {
	AgentID string `json:"agent_id"`
	Prompt  string `json:"prompt"`
}

// ChatCompletionMessage represents a single message in the chat array.
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the payload for /v1/chat/completions.
type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	Stream   bool                    `json:"stream"`
}

// ChatCompletionResponse is the standard JSON response.
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      ChatCompletionMessage `json:"message"`
		FinishReason string                `json:"finish_reason"`
		Index        int                   `json:"index"`
	} `json:"choices"`
}

// ChatCompletionStreamResponse is the SSE chunk response.
type ChatCompletionStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
		Index        int     `json:"index"`
	} `json:"choices"`
}

// runDaemon starts the persistent inference engine and an HTTP API for IPC.
func runDaemon(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("daemon", flag.ContinueOnError)
	fs.SetOutput(stderr)
	port := fs.Int("port", 8080, "HTTP port for IPC")
	agentID := fs.String("agent", "sample-alpha", "default agent to preload")
	if err := fs.Parse(args); err != nil {
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
	// Use MMAP to map weights into memory without full copy
	cfg.UseMMAP = true
	cfg.MaxTokens = 4096
	cfg.ContextSize = 8192
	cfg.BatchSize = 1024

	provider := llama.NewLlamaProvider(resolver, cfg)
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(128)

	rt := runtime.NewIntegratedModelRuntime(provider, ld, resolver, bus)

	def, err := areg.GetAgent(*agentID)
	if err != nil {
		fmt.Fprintf(stderr, "error: default agent %s not found in fleet\n", *agentID)
		return 1
	}

	fmt.Fprintf(stderr, "Pre-loading model %s into memory...\n", def.ModelProfile)
	if err := rt.LoadModel(def.ModelProfile); err != nil {
		fmt.Fprintf(stderr, "error loading model %s: %v\n", def.ModelProfile, err)
		return 1
	}

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

	fmt.Fprintf(stderr, "TinyBrain %s | Daemonized Engine | profile=%s backend=%s\n",
		Version, profile.Name, profile.Probe.Backend)

	// API Handler: Submit tasks
	http.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
		
		// Setup response channel for this task
		done := make(chan string)
		unsub := bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
			if payload, ok := ev.Payload.(events.TaskCompletedPayload); ok && payload.TaskID == taskID {
				done <- payload.Result
			}
		})
		defer unsub()

		// Submit the task
		start := time.Now()
		bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
			TaskID:  taskID,
			Input:   req.Prompt,
			AgentID: req.AgentID,
		}, start))

		// Wait for completion
		select {
		case result := <-done:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"task_id": taskID,
				"result":  result,
				"elapsed": time.Since(start).Seconds(),
			})
		case <-time.After(5 * time.Minute):
			http.Error(w, "Task timeout", http.StatusGatewayTimeout)
		}
	})

	// API Handler: OpenAI Chat Completions
	http.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Pass the entire chat history as a JSON string so agents can parse it
		// and maintain conversational context instead of suffering amnesia.
		var prompt string
		if msgBytes, err := json.Marshal(req.Messages); err == nil {
			prompt = string(msgBytes)
		} else if len(req.Messages) > 0 {
			prompt = req.Messages[len(req.Messages)-1].Content
		}

		taskID := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
		
		done := make(chan string)
		unsub := bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
			if payload, ok := ev.Payload.(events.TaskCompletedPayload); ok && payload.TaskID == taskID {
				done <- payload.Result
			}
		})
		defer unsub()

		start := time.Now()
		bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
			TaskID:  taskID,
			Input:   prompt,
			AgentID: req.Model,
		}, start))

		select {
		case resultJSON := <-done:
			// Parse the agent result JSON
			var agentRes map[string]any
			var textOutput string
			if err := json.Unmarshal([]byte(resultJSON), &agentRes); err == nil {
				if t, ok := agentRes["text"].(string); ok {
					textOutput = t
				} else {
					textOutput = resultJSON
				}
			} else {
				textOutput = resultJSON
			}

			if req.Stream {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				// Simulate streaming chunk by chunk
				chunkSize := 4
				for i := 0; i < len(textOutput); i += chunkSize {
					end := i + chunkSize
					if end > len(textOutput) {
						end = len(textOutput)
					}
					
					chunk := ChatCompletionStreamResponse{
						ID:      taskID,
						Object:  "chat.completion.chunk",
						Created: start.Unix(),
						Model:   req.Model,
					}
					chunk.Choices = make([]struct {
						Delta struct {
							Role    string `json:"role,omitempty"`
							Content string `json:"content,omitempty"`
						} `json:"delta"`
						FinishReason *string `json:"finish_reason"`
						Index        int     `json:"index"`
					}, 1)
					
					if i == 0 {
						chunk.Choices[0].Delta.Role = "assistant"
					}
					chunk.Choices[0].Delta.Content = textOutput[i:end]
					
					chunkBytes, _ := json.Marshal(chunk)
					fmt.Fprintf(w, "data: %s\n\n", string(chunkBytes))
					
					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}
					
					// Optional: sleep to simulate typing delay
					time.Sleep(10 * time.Millisecond)
				}
				
				fmt.Fprintf(w, "data: [DONE]\n\n")
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			} else {
				w.Header().Set("Content-Type", "application/json")
				resp := ChatCompletionResponse{
					ID:      taskID,
					Object:  "chat.completion",
					Created: start.Unix(),
					Model:   req.Model,
				}
				resp.Choices = make([]struct {
					Message      ChatCompletionMessage `json:"message"`
					FinishReason string                `json:"finish_reason"`
					Index        int                   `json:"index"`
				}, 1)
				resp.Choices[0].Message = ChatCompletionMessage{
					Role:    "assistant",
					Content: textOutput,
				}
				resp.Choices[0].FinishReason = "stop"
				
				json.NewEncoder(w).Encode(resp)
			}

		case <-time.After(5 * time.Minute):
			http.Error(w, "Task timeout", http.StatusGatewayTimeout)
		}
	})

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	fmt.Fprintf(stderr, "Listening for IPC tasks on http://%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(stderr, "server error: %v\n", err)
		return 1
	}

	return 0
}
