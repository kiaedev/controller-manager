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

package v1beta1

import (
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MySQLSpec defines the desired state of MySQL
type MySQLSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Database         string `json:"database"`
	Username         string `json:"username,omitempty"`
	EstimateQPS      int    `json:"estimate_qps,omitempty"`
	EstimateCapacity int    `json:"estimate_capacity,omitempty"`
}

// MySQLStatus defines the observed state of MySQL
type MySQLStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ConnectionInfo `json:",inline"`

	Phase string `json:"phase"`
}

type ConnectionInfo struct {
	Address   string `json:"address"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	PwdSecret string `json:"pwd_secret"`
}

func (info *ConnectionInfo) BuildDSN(password string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", info.Username, password, info.Address, info.Database)
}

func (info *ConnectionInfo) Host() string {
	host, _, _ := net.SplitHostPort(info.Address)
	return host
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`

// MySQL is the Schema for the mysqls API
type MySQL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MySQLSpec   `json:"spec,omitempty"`
	Status MySQLStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MySQLList contains a list of MySQL
type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySQL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MySQL{}, &MySQLList{})
}
