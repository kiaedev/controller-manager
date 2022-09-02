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

type Database struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Capacity v1.ResourceList `json:"capacity,omitempty"`
	DSN      string          `json:"dsn"`
}

type Cache struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Capacity v1.ResourceList `json:"capacity,omitempty"`
	DSN      string          `json:"dsn"`
}

type MessageQueue struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// MiddlewareSpec defines the desired state of Middleware
type MiddlewareSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Type string `json:"mode,omitempty"` // Cloud: 控制权交给MiddlewareCloudManager,自动创建实例；Manual，自行填写DSN

	Class string `json:"class,omitempty"`

	Database Database `json:"database,omitempty"`

	Cache Cache `json:"cache,omitempty"`

	MQ MessageQueue `json:"mq,omitempty"`
}

// MiddlewareStatus defines the observed state of Middleware
type MiddlewareStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Middleware is the Schema for the middlewares API
type Middleware struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiddlewareSpec   `json:"spec,omitempty"`
	Status MiddlewareStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MiddlewareList contains a list of Middleware
type MiddlewareList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Middleware `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Middleware{}, &MiddlewareList{})
}
