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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Sets the inventory policy.
type DBaaSInventoryPolicy struct {
	// Disables provisioning on inventory accounts.
	DisableProvisions *bool `json:"disableProvisions,omitempty"`

	// Namespaces where DBaaSConnection and DBaaSInstance objects are only allowed to reference a policy's inventories.
	// Each inventory can individually override this.
	// Using an asterisk surrounded by single quotes ('*'), allows all namespaces.
	// If not set in the policy or by an inventory object, connections are only allowed in the inventory's namespace.
	ConnectionNamespaces *[]string `json:"connectionNamespaces,omitempty"`

	// Use a label selector to determine the namespaces where DBaaSConnection and DBaaSInstance objects are only allowed to reference a policy's inventories.
	// Each inventory can individually override this.
	// A label selector is a label query over a set of resources.
	// Results use a logical AND from matchExpressions and matchLabels queries.
	// An empty label selector matches all objects.
	// A null label selector matches no objects.
	ConnectionNsSelector *metav1.LabelSelector `json:"connectionNsSelector,omitempty"`
}

// DBaaSInventorySpec defines the Inventory Spec to be used by provider operators
type DBaaSInventorySpec struct {
	// The secret containing the provider-specific connection credentials to use with the provider's API endpoint.
	// The format specifies the secret in the provider’s operator for its DBaaSProvider custom resource (CR), such as the CredentialFields key.
	// The secret must exist within the same namespace as the inventory.
	CredentialsRef *LocalObjectReference `json:"credentialsRef"`
}

// This object defines the desired state of a DBaaSInventory object.
type DBaaSOperatorInventorySpec struct {
	// A reference to a DBaaSProvider custom resource (CR).
	ProviderRef NamespacedName `json:"providerRef"`

	// The properties that will be copied into the provider’s inventory.
	DBaaSInventorySpec `json:",inline"`

	// The policy for this inventory.
	DBaaSInventoryPolicy `json:",inline"`
}

// Defines the inventory status that the provider's operator uses.
type DBaaSInventoryStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// A list of database services returned from querying the database provider.
	DatabaseServices []DatabaseService `json:"databaseServices,omitempty"`
}

// Defines the information of a database service.
type DatabaseService struct {
	// A provider-specific identifier for the database service.
	// It can contain one or more pieces of information used by the provider's operator to identify the database service.
	ServiceID string `json:"serviceID"`

	// The name of the database service.
	ServiceName string `json:"serviceName,omitempty"`

	// The type of the database service.
	ServiceType *DatabaseServiceType `json:"serviceType,omitempty"`

	// Any other provider-specific information related to this service.
	ServiceInfo map[string]string `json:"serviceInfo,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// The schema for the DBaaSInventory API.
// Inventory objects must be created in a valid namespace, determined by the existence of a DBaaSPolicy object.
//+operator-sdk:csv:customresourcedefinitions:displayName="Provider Account"
type DBaaSInventory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBaaSOperatorInventorySpec `json:"spec,omitempty"`
	Status DBaaSInventoryStatus       `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Contains a list of DBaaSInventories.
type DBaaSInventoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBaaSInventory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBaaSInventory{}, &DBaaSInventoryList{})
}

// The schema for a provider's inventory status.
type DBaaSProviderInventory struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBaaSInventorySpec   `json:"spec,omitempty"`
	Status DBaaSInventoryStatus `json:"status,omitempty"`
}
