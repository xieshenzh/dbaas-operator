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
	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

var _ = Describe("DBaaSInventory controller with errors", func() {
	Context("after creating DBaaSInventory without policy in the target namespace", func() {
		inventoryName := "test-inventory-no-policy"
		ns := "testns-no-policy"
		nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}

		DBaaSInventorySpec := &v1alpha2.DBaaSInventorySpec{
			CredentialsRef: &v1alpha2.LocalObjectReference{
				Name: testSecret.Name,
			},
		}
		createdDBaaSInventory := &v1alpha2.DBaaSInventory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      inventoryName,
				Namespace: ns,
			},
			Spec: v1alpha2.DBaaSOperatorInventorySpec{
				ProviderRef: v1alpha2.NamespacedName{
					Name: testProviderName,
				},
				DBaaSInventorySpec: *DBaaSInventorySpec,
			},
		}

		BeforeEach(assertResourceCreationIfNotExists(&testSecret))
		BeforeEach(assertResourceCreationIfNotExists(nsSpec))
		BeforeEach(assertResourceCreationIfNotExists(createdDBaaSInventory))
		It("reconcile with error", assertDBaaSResourceStatusUpdated(createdDBaaSInventory, metav1.ConditionFalse, v1alpha1.DBaaSPolicyNotFound))
	})

	Context("after creating DBaaSInventory without valid provider", func() {
		inventoryName := "test-inventory-no-provider"
		providerName := "provider-no-exist"
		DBaaSInventorySpec := &v1alpha2.DBaaSInventorySpec{
			CredentialsRef: &v1alpha2.LocalObjectReference{
				Name: testSecret.Name,
			},
		}
		createdDBaaSInventory := &v1alpha2.DBaaSInventory{
			ObjectMeta: metav1.ObjectMeta{
				Name:      inventoryName,
				Namespace: testNamespace,
			},
			Spec: v1alpha2.DBaaSOperatorInventorySpec{
				ProviderRef: v1alpha2.NamespacedName{
					Name: providerName,
				},
				DBaaSInventorySpec: *DBaaSInventorySpec,
			},
		}
		BeforeEach(assertResourceCreationIfNotExists(&testSecret))
		BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
		BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))
		BeforeEach(assertResourceCreationIfNotExists(createdDBaaSInventory))
		It("reconcile with error", assertDBaaSResourceStatusUpdated(createdDBaaSInventory, metav1.ConditionFalse, v1alpha1.DBaaSProviderNotFound))
	})
})

var _ = Describe("DBaaSInventory controller - nominal", func() {
	BeforeEach(assertResourceCreationIfNotExists(&testSecret))
	BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
	BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
	BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))

	Describe("reconcile", func() {
		Context("after creating DBaaSInventory", func() {
			inventoryName := "test-inventory"
			DBaaSInventorySpec := &v1alpha2.DBaaSInventorySpec{
				CredentialsRef: &v1alpha2.LocalObjectReference{
					Name: testSecret.Name,
				},
			}
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventorySpec: *DBaaSInventorySpec,
				},
			}

			BeforeEach(assertResourceCreationIfNotExists(&testSecret))
			BeforeEach(assertResourceCreation(createdDBaaSInventory))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))

			It("should create a provider inventory", assertProviderResourceCreated(createdDBaaSInventory, testInventoryKind, DBaaSInventorySpec))

			Context("when updating provider inventory status", func() {
				lastTransitionTime := getLastTransitionTimeForTest()
				status := &v1alpha2.DBaaSInventoryStatus{
					DatabaseServices: []v1alpha2.DatabaseService{
						{
							ServiceID:   "testInstanceID",
							ServiceName: "testInstance",
							ServiceType: v1alpha2.InstanceDatabaseService,
							ServiceInfo: map[string]string{
								"testInstanceInfo": "testInstanceInfo",
							},
						},
					},
					Conditions: []metav1.Condition{
						{
							Type:               "SpecSynced",
							Status:             metav1.ConditionTrue,
							Reason:             "SyncOK",
							LastTransitionTime: metav1.Time{Time: lastTransitionTime},
						},
					},
				}
				BeforeEach(assertResourceCreationIfNotExists(&testSecret))
				It("should update DBaaSInventory status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, status))
			})

			Context("when updating DBaaSInventory spec", func() {
				updatedTestSecret := v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "updated-test-credentials",
						Namespace: testNamespace,
					},
				}
				DBaaSInventorySpec := &v1alpha2.DBaaSInventorySpec{
					CredentialsRef: &v1alpha2.LocalObjectReference{
						Name: updatedTestSecret.Name,
					},
				}
				BeforeEach(assertResourceCreationIfNotExists(&updatedTestSecret))
				It("should update provider inventory spec", assertProviderResourceSpecUpdated(createdDBaaSInventory, testInventoryKind, DBaaSInventorySpec))
				It("should return the secret without error and with proper label", func() {
					getSecret := v1.Secret{}
					err := dRec.Get(ctx, client.ObjectKeyFromObject(&updatedTestSecret), &getSecret)
					Expect(err).NotTo(HaveOccurred())
					labels := getSecret.GetLabels()
					Expect(labels).Should(Not(BeNil()))
					Expect(labels[v1alpha1.TypeLabelKeyMongo]).Should(Equal(v1alpha1.TypeLabelValue))
				})
			})
		})
	})
})
