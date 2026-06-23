package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KVCacheSpec defines the desired state of KVCache
type KVCacheSpec struct {
	Size     string `json:"size,omitempty"`
	Strategy string `json:"strategy,omitempty"`
}

// KVCacheStatus defines the observed state of KVCache
type KVCacheStatus struct {
	UsedSize string `json:"usedSize,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KVCache is the Schema for the kvcache API
type KVCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KVCacheSpec   `json:"spec,omitempty"`
	Status KVCacheStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KVCacheList contains a list of KVCache
type KVCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KVCache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KVCache{}, &KVCacheList{})
}

