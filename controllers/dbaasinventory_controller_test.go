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

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

var _ = Describe("DBaaSInventory controller", func() {
	BeforeEach(func() {
		By("checking the provider ConfigMap created")
		createdCM := v1.ConfigMap{}
		if err := k8sClient.Get(ctx, client.ObjectKey{Name: testCMName, Namespace: testNamespace}, &createdCM); err != nil {
			if errors.IsNotFound(err) {
				By("creating the provider ConfigMap")
				Expect(k8sClient.Create(ctx, testProviderCM)).Should(Succeed())

				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKey{Name: testCMName, Namespace: testNamespace}, &createdCM)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
			} else {
				Fail(err.Error())
			}
		}
	})

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

		Context("when creating DBaaSInventory", func() {
			It("should create a provider inventory", func() {
				By("checking a provider inventory is created")
				gvk := schema.GroupVersionKind{
					Group:   v1alpha1.GroupVersion.Group,
					Version: v1alpha1.GroupVersion.Version,
					Kind:    testInventoryKind,
				}
				var providerInventory = unstructured.Unstructured{}
				providerInventory.SetGroupVersionKind(gvk)
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
		})

		Context("when updating provider inventory status ", func() {
			It("should update DBaaSInventory status", func() {

			})
		})
	})
})
