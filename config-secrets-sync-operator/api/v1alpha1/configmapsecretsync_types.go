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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigMapSecretSyncSpec defines the desired state of ConfigMapSecretSync
type ConfigMapSecretSyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ConfigMapSecretSync. Edit configmapsecretsync_types.go to remove/update
	SourceNamespace  string   `json:"sourceNamespace"`
	TargetNamespaces []string `json:"targetNamespaces"`
	ConfigMapNames   []string `json:"configMapNames"`
	SecretNames      []string `json:"secretNames"`
}

// ConfigMapSecretSyncStatus defines the observed state of ConfigMapSecretSync
type ConfigMapSecretSyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConfigMapSecretSync is the Schema for the configmapsecretsyncs API
type ConfigMapSecretSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigMapSecretSyncSpec   `json:"spec,omitempty"`
	Status ConfigMapSecretSyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigMapSecretSyncList contains a list of ConfigMapSecretSync
type ConfigMapSecretSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigMapSecretSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigMapSecretSync{}, &ConfigMapSecretSyncList{})
}
