package v1alpha1

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestAddToScheme(t *testing.T) {
	scheme := runtime.NewScheme()
	err := AddToScheme(scheme)
	if err != nil {
		t.Fatalf("Failed to add to scheme: %v", err)
	}

	if !scheme.IsGroupRegistered("memory.tinybrain.io") {
		t.Errorf("Expected group memory.tinybrain.io to be registered")
	}
}

func TestDeepCopy(t *testing.T) {
	policy := &SwapPolicy{
		Spec: SwapPolicySpec{
			Threshold: "80%",
		},
	}
	copy := policy.DeepCopy()
	if copy.Spec.Threshold != "80%" {
		t.Errorf("Expected copy to have same threshold, got %s", copy.Spec.Threshold)
	}
}
