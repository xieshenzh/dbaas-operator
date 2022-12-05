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

package v1alpha2

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	rdsRegistration = "rds-registration"
	providerNameKey = "spec.providerRef.name"
)

// log is for logging in this package.
var dbaasinventorylog = logf.Log.WithName("dbaasinventory-resource")
var inventoryWebhookAPIClient client.Client

func (r *DBaaSInventory) SetupWebhookWithManager(mgr ctrl.Manager) error {
	if inventoryWebhookAPIClient == nil {
		inventoryWebhookAPIClient = mgr.GetClient()
	}
	// index inventory by `spec.providerRef.name`
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &DBaaSInventory{}, providerNameKey, func(rawObj client.Object) []string {
		inventory := rawObj.(*DBaaSInventory)
		return []string{inventory.Spec.ProviderRef.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-dbaas-redhat-com-v1alpha2-dbaasinventory,mutating=false,failurePolicy=fail,sideEffects=None,groups=dbaas.redhat.com,resources=dbaasinventories,verbs=create;update,versions=v1alpha2,name=vdbaasinventory.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DBaaSInventory{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DBaaSInventory) ValidateCreate() error {
	dbaasinventorylog.Info("validate create", "name", r.Name)
	return validateInventory(r, nil)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DBaaSInventory) ValidateUpdate(old runtime.Object) error {
	dbaasinventorylog.Info("validate update", "name", r.Name)
	return validateInventory(r, old.(*DBaaSInventory))
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DBaaSInventory) ValidateDelete() error {
	dbaasinventorylog.Info("validate delete", "name", r.Name)
	return nil
}

func validateInventory(inv *DBaaSInventory, oldInv *DBaaSInventory) error {
	// Provider name is immutable
	if oldInv != nil && oldInv.Spec.ProviderRef.Name != inv.Spec.ProviderRef.Name {
		msg := "provider name is immutable for provider accounts"
		return field.Invalid(field.NewPath("spec").Child("providerRef").Child("name"), inv.Spec.ProviderRef.Name, msg)
	}
	// Retrieve the secret object
	secret := &corev1.Secret{}
	if err := inventoryWebhookAPIClient.Get(context.TODO(), types.NamespacedName{Name: inv.Spec.DBaaSInventorySpec.CredentialsRef.Name, Namespace: inv.Namespace}, secret); err != nil {
		return err
	}
	// Retrieve the provider object
	provider := &unstructured.Unstructured{}
	provider.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   v1alpha1GroupVersion.Group,
		Version: v1alpha1GroupVersion.Version,
		Kind:    "DBaaSProvider",
	})
	provider.SetName(inv.Spec.ProviderRef.Name)
	if err := inventoryWebhookAPIClient.Get(context.TODO(), types.NamespacedName{Name: inv.Spec.ProviderRef.Name, Namespace: ""}, provider); err != nil {
		return err
	}
	// Check RDS
	if oldInv == nil && inv.Spec.ProviderRef.Name == rdsRegistration {
		if err := validateRDS(); err != nil {
			return err
		}
	}
	// Check ns selector
	if inv.Spec.ConnectionNsSelector != nil {
		if _, err := metav1.LabelSelectorAsSelector(inv.Spec.ConnectionNsSelector); err != nil {
			return err
		}
	}
	return validateInventoryMandatoryFields(inv, secret, provider)
}

func validateInventoryMandatoryFields(inv *DBaaSInventory, secret *corev1.Secret, provider *unstructured.Unstructured) error {
	providerObj := &Provider{}

	b, err := provider.MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, providerObj)
	if err != nil {
		return err
	}

	for _, credField := range providerObj.Spec.CredentialFields {
		if credField.Required {
			if value, ok := secret.Data[credField.Key]; !ok || len(value) == 0 {
				//Required key is missing
				msg := fmt.Sprintf("credentialsRef is invalid: %s is required in secret %s", credField.Key, secret.Name)
				return field.Invalid(field.NewPath("spec").Child("credentialsRef"), *(inv.Spec.CredentialsRef), msg)
			}
		}
	}
	return nil
}

func validateRDS() error {
	rdsInventoryList := &DBaaSInventoryList{}
	if err := inventoryWebhookAPIClient.List(context.TODO(), rdsInventoryList, client.MatchingFields{providerNameKey: rdsRegistration}); err != nil {
		return err
	}
	if len(rdsInventoryList.Items) > 0 {
		return fmt.Errorf("only one provider account for RDS can exist in a cluster, but there is already a provider account %s created", rdsInventoryList.Items[0].Name)
	}
	return nil
}

// +kubebuilder:skipversion

type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec   `json:"spec,omitempty"`
	Status ProviderStatus `json:"status,omitempty"`
}

type ProviderSpec struct {
	Provider                     DatabaseProvider        `json:"provider"`
	InventoryKind                string                  `json:"inventoryKind"`
	ConnectionKind               string                  `json:"connectionKind"`
	InstanceKind                 string                  `json:"instanceKind"`
	CredentialFields             []CredentialField       `json:"credentialFields"`
	AllowsFreeTrial              bool                    `json:"allowsFreeTrial"`
	ExternalProvisionURL         string                  `json:"externalProvisionURL"`
	ExternalProvisionDescription string                  `json:"externalProvisionDescription"`
	InstanceParameterSpecs       []InstanceParameterSpec `json:"instanceParameterSpecs"`
}

type ProviderStatus struct {
}

type CredentialField struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	HelpText    string `json:"helpText,omitempty"`
}

type DatabaseProvider struct {
	Name               string       `json:"name"`
	DisplayName        string       `json:"displayName"`
	DisplayDescription string       `json:"displayDescription"`
	Icon               ProviderIcon `json:"icon"`
}

type ProviderIcon struct {
	Data      string `json:"base64data"`
	MediaType string `json:"mediatype"`
}

type InstanceParameterSpec struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
}
