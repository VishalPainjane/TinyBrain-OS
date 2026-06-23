package controllers

import (
	"context"
	"testing"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	corev1alpha1 "github.com/VishalPainjane/TinyBrain-OS/internal/k8s/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"k8s.io/apimachinery/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAgentReconciler(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1alpha1.AddToScheme(scheme)

	agent := &corev1alpha1.Agent{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core.tinybrain.io/v1alpha1",
			Kind:       "Agent",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-agent",
			Namespace: "default",
		},
		Spec: corev1alpha1.AgentSpec{
			Model: "test-model",
		},
	}

	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(agent).WithStatusSubresource(agent).Build()

	r := &AgentReconciler{
		Client: client,
		Scheme: scheme,
	}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-agent",
			Namespace: "default",
		},
	}

	_, err := r.Reconcile(context.TODO(), req)
	if err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	// Fetch updated agent
	var updatedAgent corev1alpha1.Agent
	if err := client.Get(context.TODO(), req.NamespacedName, &updatedAgent); err != nil {
		t.Fatalf("Failed to get updated agent: %v", err)
	}

	if updatedAgent.Status.State != "Pending" {
		t.Errorf("Expected state Pending, got %s", updatedAgent.Status.State)
	}
}

