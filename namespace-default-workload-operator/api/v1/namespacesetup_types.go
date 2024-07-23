/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// NamespaceSetupSpec defines the desired state of NamespaceSetup
type NamespaceSetupSpec struct {
	ImageRepo   string `json:"imageRepo,omitempty"`
	ServicePort int    `json:"servicePort,omitempty"`
	TargetPort  int    `json:"targetPort,omitempty"`
}

// NamespaceSetupStatus defines the observed state of NamespaceSetup
type NamespaceSetupStatus struct {
	// Define observed state of cluster
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NamespaceSetup is the Schema for the namespacesetups API
type NamespaceSetup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceSetupSpec   `json:"spec,omitempty"`
	Status NamespaceSetupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NamespaceSetupList contains a list of NamespaceSetup
type NamespaceSetupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceSetup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceSetup{}, &NamespaceSetupList{})
}
