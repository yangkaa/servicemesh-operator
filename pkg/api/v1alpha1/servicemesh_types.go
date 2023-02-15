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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceMeshStatus defines the observed state of ServiceMesh
type ServiceMeshStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ServiceMesh is the Schema for the servicemeshes API
type ServiceMesh struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Provisioner indicates the type of the provisioner.
	Provisioner string `json:"provisioner"`

	// Selector query the relevant workload and inject it. Deployment and StatefulSet are currently supported
	Selector map[string]string `json:"selector"`

	// Status defines the observed state of ServiceMesh
	Status ServiceMeshStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceMeshList contains a list of ServiceMesh
type ServiceMeshList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceMesh `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceMesh{}, &ServiceMeshList{})
}
