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

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+genclient

// ServiceMeshClass is the Schema for the servicemeshclasses API
type ServiceMeshClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Description For front-end display
	Description string `json:"description,omitempty"`

	// InjectMethods is the method of injecting sidecar
	// +kubebuilder:validation:Enum=label;annotation
	InjectMethods []InjectMethod `json:"injectMethods"`
}

// InjectMethod is the method of injecting sidecar
type InjectMethod struct {
	// +kubebuilder:validation:Enum=label;annotation
	Method string `json:"method"`

	// Name is the name of the label or annotation
	Name string `json:"name"`

	// Value is the value of the label or annotation
	Value string `json:"value"`
}

//+kubebuilder:object:root=true

// ServiceMeshClassList contains a list of ServiceMeshClass
type ServiceMeshClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceMeshClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceMeshClass{}, &ServiceMeshClassList{})
}
