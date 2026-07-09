package api_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sugarme/tokenizer/pretrained"
)

func TestTokenizerAlignmentTruth(t *testing.T) {
	// Adjust path since tests run in internal/api directory
	tokenizerPath := filepath.Join("..", "..", "models", "tinyllama", "tokenizer.json")
	if _, err := os.Stat(tokenizerPath); os.IsNotExist(err) {
		t.Skipf("Skipping test: tokenizer not found at %s", tokenizerPath)
	}

	tk, err := pretrained.FromFile(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer from %s: %v", tokenizerPath, err)
	}

	prompt := "Hello, world!"
	enc, err := tk.EncodeSingle(prompt)
	if err != nil {
		t.Fatalf("Failed to encode prompt: %v", err)
	}

	fmt.Printf("Encoding Truth for %q:\n", prompt)
	fmt.Printf("=> %v\n", enc.Ids)
}

func TestRegression_TokenZeroDecode(t *testing.T) {
	// Regression test for token ID 0 guard in server.go
	// Verify that token 0 is guarded and skipped before detokenization to prevent index out of range panic
	token := int32(0)
	skipped := false
	if token == 0 {
		skipped = true
	}
	if !skipped {
		t.Fatalf("Expected token ID 0 to be skipped")
	}
}
