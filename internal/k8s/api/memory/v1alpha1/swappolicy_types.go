package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwapPolicySpec defines the desired state of SwapPolicy
type SwapPolicySpec struct {
	Threshold string `json:"threshold,omitempty"`
	Eviction  string `json:"eviction,omitempty"`
}

// SwapPolicyStatus defines the observed state of SwapPolicy
type SwapPolicyStatus struct {
	Active bool `json:"active,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SwapPolicy is the Schema for the swappolicy API
type SwapPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwapPolicySpec   `json:"spec,omitempty"`
	Status SwapPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwapPolicyList contains a list of SwapPolicy
type SwapPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwapPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwapPolicy{}, &SwapPolicyList{})
}

