package v2

import (
	"testing"
	"fmt"
)

func TestRealModel(t *testing.T) {
	fmt.Println("Starting real model test")
	worker := NewCGOWorker()
	if err := worker.Init("/app/models/tinyllama/model.safetensors"); err != nil {
		t.Fatalf("Failed to load model: %v", err)
	}
	
	fmt.Println("Loaded model successfully. Running step...")
	
	tokens := []int32{1, 450, 7483, 310, 3444, 338} // "The capital of France is"
	payloads := []StepPayload{
		{
            Tokens: tokens, 
            IsPrefill: true,
            LogicalKVBlocks: []int{0}, // Allocate block 0 for this sequence (assuming seq len < block_size)
        },
	}
	
	out, err := worker.ExecuteStep(payloads)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}
	
	fmt.Printf("Generated token: %d\n", out[""])
}
