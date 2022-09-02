/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MiddlewarePolicySpec defines the desired state of MiddlewarePolicy
type MiddlewarePolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Selector labels.Selector `json:"selector,omitempty"`

	MiddlewareName string `json:"middlewareName,omitempty"`
}

// MiddlewarePolicyStatus defines the observed state of MiddlewarePolicy
type MiddlewarePolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MiddlewarePolicy is the Schema for the middlewarepolicies API
type MiddlewarePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiddlewarePolicySpec   `json:"spec,omitempty"`
	Status MiddlewarePolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MiddlewarePolicyList contains a list of MiddlewarePolicy
type MiddlewarePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MiddlewarePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MiddlewarePolicy{}, &MiddlewarePolicyList{})
}
