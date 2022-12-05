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
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var dbaasconnectionlog = logf.Log.WithName("dbaasconnection-resource")

func (r *DBaaSConnection) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-dbaas-redhat-com-v1alpha2-dbaasconnection,mutating=false,failurePolicy=fail,sideEffects=None,groups=dbaas.redhat.com,resources=dbaasconnections,verbs=create;update,versions=v1alpha2,name=vdbaasconnection.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &DBaaSConnection{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *DBaaSConnection) ValidateCreate() error {
	dbaasconnectionlog.Info("validate create", "name", r.Name)
	return r.validateCreateDBaaSConnectionSpec()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *DBaaSConnection) ValidateUpdate(old runtime.Object) error {
	dbaasconnectionlog.Info("validate update", "name", r.Name)
	return r.validateUpdateDBaaSConnectionSpec(old.(*DBaaSConnection))
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *DBaaSConnection) ValidateDelete() error {
	dbaasconnectionlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *DBaaSConnection) validateCreateDBaaSConnectionSpec() error {
	if len(r.Spec.DatabaseServiceID) > 0 && r.Spec.DatabaseServiceRef != nil && len(r.Spec.DatabaseServiceRef.Name) > 0 {
		return field.Invalid(field.NewPath("spec").Child("databaseServiceID"), r.Spec.DatabaseServiceID, "both databaseServiceID and databaseServiceRef are specified")
	}
	if len(r.Spec.DatabaseServiceID) == 0 && (r.Spec.DatabaseServiceRef == nil || len(r.Spec.DatabaseServiceRef.Name) == 0) {
		return field.Invalid(field.NewPath("spec").Child("databaseServiceID"), r.Spec.DatabaseServiceID, "either databaseServiceID or databaseServiceRef must be specified")
	}
	return nil
}

func (r *DBaaSConnection) validateUpdateDBaaSConnectionSpec(old *DBaaSConnection) error {
	if r.Spec.DatabaseServiceID != old.Spec.DatabaseServiceID {
		return field.Invalid(field.NewPath("spec").Child("databaseServiceID"), r.Spec.DatabaseServiceID, "databaseServiceID is immutable")
	}

	if !reflect.DeepEqual(r.Spec.InventoryRef, old.Spec.InventoryRef) {
		return field.Invalid(field.NewPath("spec").Child("inventoryRef"), r.Spec.InventoryRef, "inventoryRef is immutable")
	}

	if !reflect.DeepEqual(r.Spec.DatabaseServiceRef, old.Spec.DatabaseServiceRef) {
		return field.Invalid(field.NewPath("spec").Child("databaseServiceRef"), r.Spec.DatabaseServiceRef, "databaseServiceRef is immutable")
	}

	if !reflect.DeepEqual(r.Spec.DatabaseServiceType, old.Spec.DatabaseServiceType) {
		return field.Invalid(field.NewPath("spec").Child("databaseServiceType"), r.Spec.DatabaseServiceType, "databaseServiceType is immutable")
	}

	return nil
}
