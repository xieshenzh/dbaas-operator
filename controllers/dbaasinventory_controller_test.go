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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

var _ = Describe("DBaaSInventory controller", func() {
	BeforeEach(assertProviderConfigMapCreated)

	Describe("reconcile", func() {
		var (
			inventoryName      = "test-inventory"
			credentialsRefName = "test-credentialsRef"
			DBaaSInventorySpec = v1alpha1.DBaaSInventorySpec{
				CredentialsRef: &v1alpha1.NamespacedName{
					Name:      credentialsRefName,
					Namespace: testNamespace,
				},
			}
		)

		Context("when creating DBaaSInventory succeeds", func() {
			BeforeEach(func() {
				By("creating a new DBaaSInventory")
				DBaaSInventory := &v1alpha1.DBaaSInventory{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inventoryName,
						Namespace: testNamespace,
					},
					Spec: v1alpha1.DBaaSOperatorInventorySpec{
						Provider: v1alpha1.DatabaseProvider{
							Name: testProviderName,
						},
						DBaaSInventorySpec: DBaaSInventorySpec,
					},
				}
				Expect(k8sClient.Create(ctx, DBaaSInventory)).Should(Succeed())
				var inventory v1alpha1.DBaaSInventory
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKeyFromObject(DBaaSInventory), &inventory)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
			})

			AfterEach(func() {
				By("deleting the new DBaaSInventory")
				DBaaSInventory := &v1alpha1.DBaaSInventory{
					ObjectMeta: metav1.ObjectMeta{
						Name:      inventoryName,
						Namespace: testNamespace,
					},
				}
				Expect(k8sClient.Delete(ctx, DBaaSInventory)).Should(Succeed())
			})

			It("should create a provider inventory", func() {
				By("checking a provider inventory is created")
				var providerInventory = unstructured.Unstructured{}
				providerInventory.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   v1alpha1.GroupVersion.Group,
					Version: v1alpha1.GroupVersion.Version,
					Kind:    testInventoryKind,
				})
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: inventoryName, Namespace: testNamespace}, &providerInventory)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())

				By("checking a provider inventory spec is correct")
				providerInventorySpec, existSpec := providerInventory.UnstructuredContent()["spec"]
				Expect(existSpec).Should(Equal(true))
				var spec v1alpha1.DBaaSInventorySpec
				err := decode(providerInventorySpec, &spec)
				Expect(err).NotTo(HaveOccurred())

				Expect(spec).Should(Equal(DBaaSInventorySpec))
			})

			Context("when updating provider inventory status ", func() {
				It("should update DBaaSInventory status", func() {
					By("checking the DBaaSInventory status has no instance")
					objectKey := client.ObjectKey{Name: inventoryName, Namespace: testNamespace}
					DBaaSInventory := &v1alpha1.DBaaSInventory{}
					Consistently(func() (int, error) {
						err := k8sClient.Get(ctx, objectKey, DBaaSInventory)
						if err != nil {
							return -1, err
						}
						return len(DBaaSInventory.Status.Instances), nil
					}, duration, interval).Should(Equal(0))

					By("getting the provider inventory")
					var providerInventory = unstructured.Unstructured{}
					providerInventory.SetGroupVersionKind(schema.GroupVersionKind{
						Group:   v1alpha1.GroupVersion.Group,
						Version: v1alpha1.GroupVersion.Version,
						Kind:    testInventoryKind,
					})
					Eventually(func() bool {
						err := k8sClient.Get(ctx, client.ObjectKey{Name: inventoryName, Namespace: testNamespace}, &providerInventory)
						if err != nil {
							return false
						}
						return true
					}, timeout, interval).Should(BeTrue())

					By("updating the provider inventory status")
					lastTransitionTime, err := time.Parse(time.RFC3339, "2021-06-30T22:17:55-04:00")
					Expect(err).NotTo(HaveOccurred())
					status := v1alpha1.DBaaSInventoryStatus{
						Type: testProviderName,
						Instances: []v1alpha1.Instance{
							{
								InstanceID: "testInstanceID",
								Name:       "testInstance",
								InstanceInfo: map[string]string{
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
					providerInventory.UnstructuredContent()["status"] = status
					Expect(k8sClient.Status().Update(ctx, &providerInventory)).Should(Succeed())

					By("checking the DBaaSInventory status is updated")
					Eventually(func() (int, error) {
						err := k8sClient.Get(ctx, objectKey, DBaaSInventory)
						if err != nil {
							return -1, err
						}
						return len(DBaaSInventory.Status.Instances), nil
					}, timeout, interval).Should(Equal(1))
					Expect(DBaaSInventory.Status).Should(Equal(status))
				})
			})
		})
	})
})
