package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/VishalPainjane/TinyBrain-OS/internal/api"
	v2 "github.com/VishalPainjane/TinyBrain-OS/internal/scheduler/v2"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func main() {
	var modelPath string
	var tokenizerPath string
	var maxBlocks int
	flag.StringVar(&modelPath, "model", "", "Path to safetensors model")
	flag.StringVar(&tokenizerPath, "tokenizer", "", "Path to tokenizer.json")
	flag.IntVar(&maxBlocks, "max-blocks", 512, "Maximum KV cache blocks")
	flag.Parse()

	prefillChunkSize := 512

	// Initialize the CGO Data Plane Worker
	worker := v2.NewCGOWorker()
	if modelPath != "" {
		if err := worker.Init(modelPath); err != nil {
			log.Fatalf("Failed to initialize worker with model %s: %v", modelPath, err)
		}
	}

	// Initialize the Scheduler Control Plane
	reqChan := make(chan *v2.SequenceGroup, 1000)
	scheduler := v2.NewScheduler(worker, maxBlocks, reqChan, prefillChunkSize)

	// Initialize the Master Orchestrator Engine
	engine := v2.NewEngine(scheduler, reqChan)

	// Start Engine Loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go engine.Start(ctx)

	// Initialize Tokenizer
	var tk *tokenizer.Tokenizer
	if tokenizerPath != "" {
		var err error
		tk, err = pretrained.FromFile(tokenizerPath)
		if err != nil {
			log.Fatalf("Failed to load tokenizer: %v", err)
		}
	}

	// Initialize API Layer
	server := api.NewHTTPServer(engine, tk)
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/generate", server.HandleGenerate)

	log.Println("TinyBrain OS v2 Ingress Server started on :8080")
	
	// Start Server
	errChan := make(chan error, 1)
	go func() {
		errChan <- http.ListenAndServe(":8080", mux)
	}()

	// Graceful shutdown on interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down...", sig)
		cancel()
	}
}
