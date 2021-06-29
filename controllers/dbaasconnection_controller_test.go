/*
Copyright 2021.

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

package controllers

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

var _ = Describe("DBaaSConnection controller", func() {
	BeforeEach(assertProviderConfigMapCreation())

	Describe("reconcile", func() {
		Context("after creating DBaaSInventory", func() {
			inventoryRefName := "test-inventory-ref"
			createdDBaaSInventory := &v1alpha1.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha1.DBaaSOperatorInventorySpec{
					Provider: v1alpha1.DatabaseProvider{
						Name: testProviderName,
					},
					DBaaSInventorySpec: v1alpha1.DBaaSInventorySpec{
						CredentialsRef: &v1alpha1.NamespacedName{
							Name:      "test-credentialsRef",
							Namespace: testNamespace,
						},
					},
				},
			}
			createdProviderInventory := &unstructured.Unstructured{}
			createdProviderInventory.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   v1alpha1.GroupVersion.Group,
				Version: v1alpha1.GroupVersion.Version,
				Kind:    testInventoryKind,
			})
			createdProviderInventory.SetNamespace(testNamespace)
			createdProviderInventory.SetName(inventoryRefName)

			BeforeEach(assertResourceCreation(createdDBaaSInventory))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))
			AfterEach(assertResourceDeletion(createdProviderInventory))

			Context("after creating DBaaSConnection", func() {
				connectionName := "test-connection"
				instanceID := "test-instanceID"
				DBaaSConnectionSpec := &v1alpha1.DBaaSConnectionSpec{
					InventoryRef: &v1.LocalObjectReference{
						Name: inventoryRefName,
					},
					InstanceID: instanceID,
				}
				createdDBaaSConnection := &v1alpha1.DBaaSConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      connectionName,
						Namespace: testNamespace,
					},
					Spec: *DBaaSConnectionSpec,
				}
				createdProviderConnection := &unstructured.Unstructured{}
				createdProviderConnection.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   v1alpha1.GroupVersion.Group,
					Version: v1alpha1.GroupVersion.Version,
					Kind:    testConnectionKind,
				})
				createdProviderConnection.SetNamespace(testNamespace)
				createdProviderConnection.SetName(connectionName)

				BeforeEach(assertResourceCreation(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdProviderConnection))

				It("should create a provider connection", assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))

				Context("when updating provider connection status", func() {
					lastTransitionTime, err := time.Parse(time.RFC3339, "2021-06-30T22:17:55-04:00")
					Expect(err).NotTo(HaveOccurred())
					status := &v1alpha1.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								Type:               "SpecSynced",
								Status:             metav1.ConditionTrue,
								Reason:             "SyncOK",
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: "testCredentialsRef",
						},
						ConnectionInfo: &v1.LocalObjectReference{
							Name: "testConnectionInfo",
						},
					}
					It("should update DBaaSConnection status", assertDBaaSResourceStatusUpdated(createdDBaaSConnection, testConnectionKind, status))
				})

				Context("when updating DBaaSConnection spec", func() {
					DBaaSConnectionSpec := &v1alpha1.DBaaSConnectionSpec{
						InventoryRef: &v1.LocalObjectReference{
							Name: inventoryRefName,
						},
						InstanceID: "updated-test-instanceID",
					}
					It("should update provider connection spec", assertProviderResourceSpecUpdated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				})
			})
		})
	})
})
