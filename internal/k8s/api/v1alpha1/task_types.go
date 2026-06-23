package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TaskSpec defines the desired state of Task
type TaskSpec struct {
	AgentRef string `json:"agentRef,omitempty"`
	Prompt   string `json:"prompt,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

// TaskStatus defines the observed state of Task
type TaskStatus struct {
	State  string `json:"state,omitempty"`
	Result string `json:"result,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Task is the Schema for the tasks API
type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskSpec   `json:"spec,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TaskList contains a list of Task
type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Task `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Task{}, &TaskList{})
}

