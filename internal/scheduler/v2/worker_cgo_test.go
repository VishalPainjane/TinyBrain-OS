package v2

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCGOWorker(t *testing.T) {
	worker := NewCGOWorker()

	payloads := []StepPayload{
		{
			SeqID:           "test-seq-1",
			IsPrefill:       true,
			Tokens:          []int32{1, 2, 3, 4},
			LogicalKVBlocks: []int{0},
		},
	}

	out, err := worker.ExecuteStep(payloads)
	if err != nil {
		t.Fatalf("ExecuteStep failed: %v", err)
	}

	if out == nil {
		t.Fatalf("Expected output mapping")
	}
}

func TestCGOWorkerStress(t *testing.T) {
	worker := NewCGOWorker()

	numGoroutines := 10
	iterations := 100

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(gId int) {
			defer wg.Done()
			seqID := fmt.Sprintf("stress-seq-%d", gId)
			for j := 0; j < iterations; j++ {
				payloads := []StepPayload{
					{
						SeqID:           seqID,
						IsPrefill:       false,
						Tokens:          []int32{int32(j)},
						LogicalKVBlocks: []int{j % 512},
					},
				}
				_, err := worker.ExecuteStep(payloads)
				if err != nil {
					t.Errorf("ExecuteStep failed: %v", err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	totalTokens := numGoroutines * iterations
	itl := duration.Seconds() * 1000 / float64(totalTokens)

	fmt.Printf("[Telemetry] Stress Test Completed.\n")
	fmt.Printf("Total Tokens Processed: %d\n", totalTokens)
	fmt.Printf("Total Time: %v\n", duration)
	fmt.Printf("Average ITL: %.2f ms/token\n", itl)
}
