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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DAQApplicationSpec defines the desired state of DAQApplication
type DAQApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	// +kubebuilder:validation:Default=STARTED
	DesiredState DAQFSMState `json:"desiredState"`
	PodName      string      `json:"podName"`
}

// +kubebuilder:validation:Enum=INIT;CONFIGURED;STARTED
type DAQFSMState string

const (
	DAQInitState       DAQFSMState = "INIT"
	DAQConfiguredState DAQFSMState = "CONFIGURED"
	DAQStartedState    DAQFSMState = "STARTED"
)

// DAQApplicationStatus defines the observed state of DAQApplication
type DAQApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastSeenState   string `json:"lastSeenState"`
	LastCommandSent string `json:"lastCommandSent"`
	Status          string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DAQ State",type=string,JSONPath=`.status.lastSeenState`,description="Last observed DAQ FSM State"
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`,description="Controller Status"

// DAQApplication is the Schema for the daqapplications API
type DAQApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DAQApplicationSpec   `json:"spec,omitempty"`
	Status DAQApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DAQApplicationList contains a list of DAQApplication
type DAQApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DAQApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DAQApplication{}, &DAQApplicationList{})
}
