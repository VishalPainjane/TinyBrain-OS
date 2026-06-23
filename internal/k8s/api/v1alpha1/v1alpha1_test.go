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

	if !scheme.IsGroupRegistered("core.tinybrain.io") {
		t.Errorf("Expected group core.tinybrain.io to be registered")
	}
}

func TestDeepCopy(t *testing.T) {
	agent := &Agent{
		Spec: AgentSpec{
			Model: "test-model",
		},
	}
	copy := agent.DeepCopy()
	if copy.Spec.Model != "test-model" {
		t.Errorf("Expected copy to have same model, got %s", copy.Spec.Model)
	}
}

