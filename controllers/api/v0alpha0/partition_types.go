/*


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

package v0alpha0

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// corev1 "k8s.io/api/core/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PartitionSpec defines the desired state of Partition
type PartitionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Minimum=0
	RunNumber  *int64       `json:"runNumber"`
	ConfigName string       `json:"configName"`
	Resources  []ModuleSpec `json:"resources"`
}

type ModuleSpec struct {
	Module string              `json:"module"`
	TPC    ModuleResourcesSpec `json:"TPC"`
	PDS    ModuleResourcesSpec `json:"PDS"`
}

type ModuleResourcesSpec struct {
	APAs []string `json:"APAs"`
}

// PartitionStatus defines the observed state of Partition
type PartitionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Run-number",type=integer,JSONPath=`.spec.runNumber`,description="DAQ Run Number"

// Partition is the Schema for the partitions API
type Partition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PartitionSpec   `json:"spec,omitempty"`
	Status PartitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PartitionList contains a list of Partition
type PartitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Partition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Partition{}, &PartitionList{})
}
