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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Estimate struct {
	QPS     int32  `json:"qps,omitempty"`
	Storage string `json:"storage,omitempty"`
}

type MySQL struct {
	DBName string `json:"dbname,omitempty"`

	Resources v1.ResourceRequirements `json:"resource,omitempty"`
}

type Redis struct {
	DBName string `json:"dbname,omitempty"`

	Resources v1.ResourceRequirements `json:"resource,omitempty"`
}

// MiddlewareClaimSpec defines the desired state of MiddlewareClaim
type MiddlewareClaimSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// MiddlewareName string `json:"middleware_name,omitempty"`
	Estimate Estimate `json:"estimate,omitempty"`

	Mysql MySQL `json:"mysql,omitempty"`

	Redis Redis `json:"redis,omitempty"`
}

// MiddlewareClaimStatus defines the observed state of MiddlewareClaim
type MiddlewareClaimStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MiddlewareClaim is the Schema for the middlewareclaims API
type MiddlewareClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiddlewareClaimSpec   `json:"spec,omitempty"`
	Status MiddlewareClaimStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MiddlewareClaimList contains a list of MiddlewareClaim
type MiddlewareClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MiddlewareClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MiddlewareClaim{}, &MiddlewareClaimList{})
}

// +enum
type MiddlewareType string

const (
	TypeMySQL MiddlewareType = "MySQL"

	VolumeAvailable MiddlewareType = "Available"

	VolumeBound MiddlewareType = "Bound"
	// used for PersistentVolumes where the bound PersistentVolumeClaim was deleted
	// this phase is used by the persistent volume claim binder to signal to another process to reclaim the resource
	VolumeReleased MiddlewareType = "Released"
	// used for PersistentVolumes that failed to be correctly recycled or deleted after being released from a claim
	VolumeFailed MiddlewareType = "Failed"
)
