//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialField) DeepCopyInto(out *CredentialField) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialField.
func (in *CredentialField) DeepCopy() *CredentialField {
	if in == nil {
		return nil
	}
	out := new(CredentialField)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSConnection) DeepCopyInto(out *DBaaSConnection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSConnection.
func (in *DBaaSConnection) DeepCopy() *DBaaSConnection {
	if in == nil {
		return nil
	}
	out := new(DBaaSConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBaaSConnection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSConnectionList) DeepCopyInto(out *DBaaSConnectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DBaaSConnection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSConnectionList.
func (in *DBaaSConnectionList) DeepCopy() *DBaaSConnectionList {
	if in == nil {
		return nil
	}
	out := new(DBaaSConnectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBaaSConnectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSConnectionSpec) DeepCopyInto(out *DBaaSConnectionSpec) {
	*out = *in
	out.InventoryRef = in.InventoryRef
	if in.DatabaseServiceRef != nil {
		in, out := &in.DatabaseServiceRef, &out.DatabaseServiceRef
		*out = new(NamespacedName)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSConnectionSpec.
func (in *DBaaSConnectionSpec) DeepCopy() *DBaaSConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(DBaaSConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSConnectionStatus) DeepCopyInto(out *DBaaSConnectionStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CredentialsRef != nil {
		in, out := &in.CredentialsRef, &out.CredentialsRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.ConnectionInfoRef != nil {
		in, out := &in.ConnectionInfoRef, &out.ConnectionInfoRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSConnectionStatus.
func (in *DBaaSConnectionStatus) DeepCopy() *DBaaSConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(DBaaSConnectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSInventory) DeepCopyInto(out *DBaaSInventory) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSInventory.
func (in *DBaaSInventory) DeepCopy() *DBaaSInventory {
	if in == nil {
		return nil
	}
	out := new(DBaaSInventory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBaaSInventory) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSInventoryList) DeepCopyInto(out *DBaaSInventoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DBaaSInventory, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSInventoryList.
func (in *DBaaSInventoryList) DeepCopy() *DBaaSInventoryList {
	if in == nil {
		return nil
	}
	out := new(DBaaSInventoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBaaSInventoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSInventoryPolicy) DeepCopyInto(out *DBaaSInventoryPolicy) {
	*out = *in
	if in.DisableProvisions != nil {
		in, out := &in.DisableProvisions, &out.DisableProvisions
		*out = new(bool)
		**out = **in
	}
	if in.ConnectionNamespaces != nil {
		in, out := &in.ConnectionNamespaces, &out.ConnectionNamespaces
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	if in.ConnectionNsSelector != nil {
		in, out := &in.ConnectionNsSelector, &out.ConnectionNsSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSInventoryPolicy.
func (in *DBaaSInventoryPolicy) DeepCopy() *DBaaSInventoryPolicy {
	if in == nil {
		return nil
	}
	out := new(DBaaSInventoryPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSInventorySpec) DeepCopyInto(out *DBaaSInventorySpec) {
	*out = *in
	if in.CredentialsRef != nil {
		in, out := &in.CredentialsRef, &out.CredentialsRef
		*out = new(LocalObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSInventorySpec.
func (in *DBaaSInventorySpec) DeepCopy() *DBaaSInventorySpec {
	if in == nil {
		return nil
	}
	out := new(DBaaSInventorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSInventoryStatus) DeepCopyInto(out *DBaaSInventoryStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DatabaseServices != nil {
		in, out := &in.DatabaseServices, &out.DatabaseServices
		*out = make([]DatabaseService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSInventoryStatus.
func (in *DBaaSInventoryStatus) DeepCopy() *DBaaSInventoryStatus {
	if in == nil {
		return nil
	}
	out := new(DBaaSInventoryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSOperatorInventorySpec) DeepCopyInto(out *DBaaSOperatorInventorySpec) {
	*out = *in
	out.ProviderRef = in.ProviderRef
	in.DBaaSInventorySpec.DeepCopyInto(&out.DBaaSInventorySpec)
	in.DBaaSInventoryPolicy.DeepCopyInto(&out.DBaaSInventoryPolicy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSOperatorInventorySpec.
func (in *DBaaSOperatorInventorySpec) DeepCopy() *DBaaSOperatorInventorySpec {
	if in == nil {
		return nil
	}
	out := new(DBaaSOperatorInventorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSProviderConnection) DeepCopyInto(out *DBaaSProviderConnection) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSProviderConnection.
func (in *DBaaSProviderConnection) DeepCopy() *DBaaSProviderConnection {
	if in == nil {
		return nil
	}
	out := new(DBaaSProviderConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBaaSProviderInventory) DeepCopyInto(out *DBaaSProviderInventory) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBaaSProviderInventory.
func (in *DBaaSProviderInventory) DeepCopy() *DBaaSProviderInventory {
	if in == nil {
		return nil
	}
	out := new(DBaaSProviderInventory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseProvider) DeepCopyInto(out *DatabaseProvider) {
	*out = *in
	out.Icon = in.Icon
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseProvider.
func (in *DatabaseProvider) DeepCopy() *DatabaseProvider {
	if in == nil {
		return nil
	}
	out := new(DatabaseProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseService) DeepCopyInto(out *DatabaseService) {
	*out = *in
	if in.ServiceInfo != nil {
		in, out := &in.ServiceInfo, &out.ServiceInfo
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseService.
func (in *DatabaseService) DeepCopy() *DatabaseService {
	if in == nil {
		return nil
	}
	out := new(DatabaseService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceParameterSpec) DeepCopyInto(out *InstanceParameterSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceParameterSpec.
func (in *InstanceParameterSpec) DeepCopy() *InstanceParameterSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceParameterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalObjectReference) DeepCopyInto(out *LocalObjectReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalObjectReference.
func (in *LocalObjectReference) DeepCopy() *LocalObjectReference {
	if in == nil {
		return nil
	}
	out := new(LocalObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Provider) DeepCopyInto(out *Provider) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Provider.
func (in *Provider) DeepCopy() *Provider {
	if in == nil {
		return nil
	}
	out := new(Provider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProviderIcon) DeepCopyInto(out *ProviderIcon) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProviderIcon.
func (in *ProviderIcon) DeepCopy() *ProviderIcon {
	if in == nil {
		return nil
	}
	out := new(ProviderIcon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProviderSpec) DeepCopyInto(out *ProviderSpec) {
	*out = *in
	out.Provider = in.Provider
	if in.CredentialFields != nil {
		in, out := &in.CredentialFields, &out.CredentialFields
		*out = make([]CredentialField, len(*in))
		copy(*out, *in)
	}
	if in.InstanceParameterSpecs != nil {
		in, out := &in.InstanceParameterSpecs, &out.InstanceParameterSpecs
		*out = make([]InstanceParameterSpec, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProviderSpec.
func (in *ProviderSpec) DeepCopy() *ProviderSpec {
	if in == nil {
		return nil
	}
	out := new(ProviderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProviderStatus) DeepCopyInto(out *ProviderStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProviderStatus.
func (in *ProviderStatus) DeepCopy() *ProviderStatus {
	if in == nil {
		return nil
	}
	out := new(ProviderStatus)
	in.DeepCopyInto(out)
	return out
}
