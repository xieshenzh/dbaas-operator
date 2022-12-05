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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Defines the desired state of a DBaaSConnection object.
type DBaaSConnectionSpec struct {
	// A reference to the relevant DBaaSInventory custom resource (CR).
	InventoryRef NamespacedName `json:"inventoryRef"`

	// The ID of the database service to connect to, as seen in the status of the referenced DBaaSInventory.
	DatabaseServiceID string `json:"databaseServiceID,omitempty"`

	// A reference to the database service CR used, if the DatabaseServiceID is not specified.
	DatabaseServiceRef *NamespacedName `json:"databaseServiceRef,omitempty"`

	// The type of the database service to connect to, as seen in the status of the referenced DBaaSInventory.
	DatabaseServiceType DatabaseServiceType `json:"databaseServiceType,omitempty"`
}

// Defines the observed state of a DBaaSConnection object.
type DBaaSConnectionStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// The secret holding account credentials for accessing the database instance.
	CredentialsRef *corev1.LocalObjectReference `json:"credentialsRef,omitempty"`

	// A ConfigMap object holding non-sensitive information for connecting to the database instance.
	ConnectionInfoRef *corev1.LocalObjectReference `json:"connectionInfoRef,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// The schema for the DBaaSConnection API.
//+operator-sdk:csv:customresourcedefinitions:displayName="DBaaSConnection"
type DBaaSConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBaaSConnectionSpec   `json:"spec,omitempty"`
	Status DBaaSConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Contains a list of DBaaSConnections.
type DBaaSConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBaaSConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBaaSConnection{}, &DBaaSConnectionList{})
}

// The schema for a provider's connection status.
type DBaaSProviderConnection struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBaaSConnectionSpec   `json:"spec,omitempty"`
	Status DBaaSConnectionStatus `json:"status,omitempty"`
}
