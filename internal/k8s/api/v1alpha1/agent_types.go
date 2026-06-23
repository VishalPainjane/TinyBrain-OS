package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgentResources defines resource limits for the Agent
type AgentResources struct {
	MaxVram string `json:"maxVram,omitempty"`
}

// AgentSpec defines the desired state of Agent
type AgentSpec struct {
	Model        string         `json:"model,omitempty"`
	Plugin       string         `json:"plugin,omitempty"`
	SystemPrompt string         `json:"systemPrompt,omitempty"`
	Resources    AgentResources `json:"resources,omitempty"`
}

// AgentStatus defines the observed state of Agent
type AgentStatus struct {
	State       string `json:"state,omitempty"`
	ModelLoaded bool   `json:"modelLoaded,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Agent is the Schema for the agents API
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AgentList contains a list of Agent
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Agent{}, &AgentList{})
}

