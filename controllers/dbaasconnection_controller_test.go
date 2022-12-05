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

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha2"
)

var _ = Describe("DBaaSConnection controller with errors", func() {
	BeforeEach(assertResourceCreationIfNotExists(&testSecret))
	Context("after creating DBaaSConnection without inventory", func() {
		connectionName := "test-connection-no-inventory"
		instanceID := "test-instanceID"
		inventoryRefName := "test-inventory-no-exist-ref"
		DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
			InventoryRef: v1alpha2.NamespacedName{
				Name:      inventoryRefName,
				Namespace: testNamespace,
			},
			DatabaseServiceID:   instanceID,
			DatabaseServiceType: v1alpha2.InstanceDatabaseService,
		}
		createdDBaaSConnection := &v1alpha2.DBaaSConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      connectionName,
				Namespace: testNamespace,
			},
			Spec: *DBaaSConnectionSpec,
		}

		BeforeEach(assertResourceCreation(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSConnection))
		It("reconcile with error", assertDBaaSResourceStatusUpdated(createdDBaaSConnection, metav1.ConditionFalse, v1alpha1.DBaaSInventoryNotFound))
	})
	Context("after creating DBaaSConnection with inventory that is not ready", func() {
		connectionName := "test-connection-not-ready"
		instanceID := "test-instanceID"
		inventoryName := "test-connection-inventory-not-ready"
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
		DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
			InventoryRef: v1alpha2.NamespacedName{
				Name:      inventoryName,
				Namespace: testNamespace,
			},
			DatabaseServiceID:   instanceID,
			DatabaseServiceType: v1alpha2.InstanceDatabaseService,
		}
		createdDBaaSConnection := &v1alpha2.DBaaSConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      connectionName,
				Namespace: testNamespace,
			},
			Spec: *DBaaSConnectionSpec,
		}
		lastTransitionTime := getLastTransitionTimeForTest()
		providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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
					Status:             metav1.ConditionFalse,
					Reason:             "BackendError",
					LastTransitionTime: metav1.Time{Time: lastTransitionTime},
				},
			},
		}

		BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
		BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
		BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))
		BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionFalse, testInventoryKind, providerInventoryStatus))
		BeforeEach(assertResourceCreationIfNotExists(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSInventory))
		It("reconcile with error", assertDBaaSResourceStatusUpdated(createdDBaaSConnection, metav1.ConditionFalse, v1alpha1.DBaaSInventoryNotReady))
	})
	Context("after creating DBaaSConnection in an invalid namespace", func() {
		connectionName := "test-connection"
		instanceID := "test-instanceID"
		inventoryName := "test-connection-inventory"
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
				DBaaSInventoryPolicy: v1alpha2.DBaaSInventoryPolicy{
					ConnectionNamespaces: &[]string{"valid-ns", "random"},
				},
				DBaaSInventorySpec: *DBaaSInventorySpec,
			},
		}
		DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
			InventoryRef: v1alpha2.NamespacedName{
				Name:      inventoryName,
				Namespace: testNamespace,
			},
			DatabaseServiceID:   instanceID,
			DatabaseServiceType: v1alpha2.InstanceDatabaseService,
		}
		otherNS := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "other",
			},
		}
		createdDBaaSConnection := &v1alpha2.DBaaSConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      connectionName,
				Namespace: otherNS.Name,
			},
			Spec: *DBaaSConnectionSpec,
		}
		lastTransitionTime := getLastTransitionTimeForTest()
		providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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

		BeforeEach(assertResourceCreationIfNotExists(&otherNS))
		BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
		BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
		BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))
		BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
		BeforeEach(assertResourceCreationIfNotExists(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSInventory))
		It("reconcile with error", assertDBaaSResourceStatusUpdated(createdDBaaSConnection, metav1.ConditionFalse, v1alpha1.DBaaSInvalidNamespace))
	})
	Context("after creating DBaaSConnection with an invalid instanceRef", func() {
		connectionName := "test-connection"
		instanceName := "test-instance-invalid"
		inventoryName := "test-connection-inventory"
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
				DBaaSInventoryPolicy: v1alpha2.DBaaSInventoryPolicy{
					ConnectionNamespaces: &[]string{"valid-ns", "random"},
				},
				DBaaSInventorySpec: *DBaaSInventorySpec,
			},
		}
		DBaaSInstanceSpec := &v1alpha1.DBaaSInstanceSpec{
			InventoryRef: v1alpha1.NamespacedName{
				Name:      inventoryName,
				Namespace: testNamespace,
			},
			Name:          "test-instance-to-create",
			CloudProvider: "aws",
			CloudRegion:   "test-region",
			OtherInstanceParams: map[string]string{
				"testParam": "test-param",
			},
		}
		createdDBaaSInstance := &v1alpha1.DBaaSInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instanceName,
				Namespace: testNamespace,
			},
			Spec: *DBaaSInstanceSpec,
		}
		lastTransitionTime := getLastTransitionTimeForTest()
		providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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
		DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
			InventoryRef: v1alpha2.NamespacedName{
				Name:      inventoryName,
				Namespace: testNamespace,
			},
			DatabaseServiceRef: &v1alpha2.NamespacedName{
				Name:      createdDBaaSInstance.Name,
				Namespace: createdDBaaSInstance.Namespace,
			},
		}
		createdDBaaSConnection := &v1alpha2.DBaaSConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      connectionName,
				Namespace: testNamespace,
			},
			Spec: *DBaaSConnectionSpec,
		}
		BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
		BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
		BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))
		BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
		BeforeEach(assertResourceCreationIfNotExists(createdDBaaSInstance))
		BeforeEach(assertResourceCreationIfNotExists(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSConnection))
		AfterEach(assertResourceDeletion(createdDBaaSInstance))
		AfterEach(assertResourceDeletion(createdDBaaSInventory))
		It("reconcile with error", assertDBaaSResourceStatusUpdated(createdDBaaSConnection, metav1.ConditionFalse, v1alpha2.DBaaSServiceNotAvailable))
	})
})

var _ = Describe("DBaaSConnection controller - nominal", func() {
	BeforeEach(assertResourceCreationIfNotExists(&testSecret))
	BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
	BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
	BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))

	Describe("reconcile", func() {
		Context("after creating DBaaSInventory", func() {
			inventoryRefName := "test-inventory-ref"
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventorySpec: v1alpha2.DBaaSInventorySpec{
						CredentialsRef: &v1alpha2.LocalObjectReference{
							Name: testSecret.Name,
						},
					},
				},
			}
			lastTransitionTime := getLastTransitionTimeForTest()
			providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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

			Context("after creating DBaaSConnection", func() {
				connectionName := "test-connection-1"
				instanceID := "test-instanceID"
				DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
					InventoryRef: v1alpha2.NamespacedName{
						Name:      inventoryRefName,
						Namespace: testNamespace,
					},
					DatabaseServiceID:   instanceID,
					DatabaseServiceType: v1alpha2.InstanceDatabaseService,
				}
				createdDBaaSConnection := &v1alpha2.DBaaSConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      connectionName,
						Namespace: testNamespace,
					},
					Spec: *DBaaSConnectionSpec,
				}
				BeforeEach(assertResourceCreation(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdDBaaSConnection))

				It("should create a provider connection", func() {
					assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec)()

					By("checking if the Deployment is created")
					deployment := &appv1.Deployment{
						ObjectMeta: metav1.ObjectMeta{
							Name:      connectionName,
							Namespace: testNamespace,
						},
					}
					Eventually(func() bool {
						err := dRec.Get(ctx, client.ObjectKeyFromObject(deployment), deployment)
						if err != nil {
							return false
						}
						Expect(deployment.Spec.Replicas).ShouldNot(BeNil())
						Expect(*deployment.Spec.Replicas).ShouldNot(Equal(0))

						Expect(deployment.Labels).Should(BeNil())
						Expect(deployment.Annotations).ShouldNot(BeNil())
						mb, mbOk := deployment.Annotations["managed-by"]
						Expect(mbOk).Should(BeTrue())
						Expect(mb).Should(Equal("dbaas-operator"))
						owner, ownerOk := deployment.Annotations["owner"]
						Expect(ownerOk).Should(BeTrue())
						Expect(owner).Should(Equal(connectionName))
						ownerKind, ownerKindOk := deployment.Annotations["owner.kind"]
						Expect(ownerKindOk).Should(BeTrue())
						Expect(ownerKind).Should(Equal("DBaaSConnection"))
						ownerNS, ownerNSOk := deployment.Annotations["owner.namespace"]
						Expect(ownerNSOk).Should(BeTrue())
						Expect(ownerNS).Should(Equal(testNamespace))

						deploymentOwner := metav1.GetControllerOf(deployment)
						Expect(deploymentOwner).ShouldNot(BeNil())
						Expect(deploymentOwner.Kind).Should(Equal("DBaaSConnection"))
						Expect(deploymentOwner.Name).Should(Equal(connectionName))
						Expect(deploymentOwner.Controller).ShouldNot(BeNil())
						Expect(*deploymentOwner.Controller).Should(BeTrue())
						Expect(deploymentOwner.BlockOwnerDeletion).ShouldNot(BeNil())
						Expect(*deploymentOwner.BlockOwnerDeletion).Should(BeTrue())
						return true
					}, timeout).Should(BeTrue())
				})
				Context("when updating provider connection status", func() {
					lastTransitionTime := getLastTransitionTimeForTest()
					status := &v1alpha2.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								Type:               "ReadyForBinding",
								Status:             metav1.ConditionTrue,
								Reason:             "SyncOK",
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: testSecret.Name,
						},
						ConnectionInfoRef: &v1.LocalObjectReference{
							Name: "testConnectionInfoRef",
						},
					}
					It("should update DBaaSConnection status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSConnection, metav1.ConditionTrue, testConnectionKind, status))
				})

				Context("when updating DBaaSConnection spec", func() {
					DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
						InventoryRef: v1alpha2.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						DatabaseServiceID:   "updated-test-instanceID",
						DatabaseServiceType: v1alpha2.InstanceDatabaseService,
					}
					It("should update provider connection spec", assertProviderResourceSpecUpdated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				})

			})

			BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))
		})
	})

})

var _ = Describe("DBaaSConnection controller - nominal with instance reference", func() {
	BeforeEach(assertResourceCreationIfNotExists(&testSecret))
	BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
	BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
	BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))

	Describe("reconcile", func() {
		Context("after creating DBaaSInventory", func() {
			instanceID := "test-instance-ref-ID"
			inventoryRefName := "test-instance-ref-inventory-ref"
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventorySpec: v1alpha2.DBaaSInventorySpec{
						CredentialsRef: &v1alpha2.LocalObjectReference{
							Name: testSecret.Name,
						},
					},
				},
			}
			lastTransitionTime := getLastTransitionTimeForTest()
			providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
				DatabaseServices: []v1alpha2.DatabaseService{
					{
						ServiceID:   instanceID,
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
			BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))

			Context("after creating DBaaSInstance", func() {
				instanceRefName := "test-instance-ref-instance"
				createdDBaaSInstance := &v1alpha1.DBaaSInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      instanceRefName,
						Namespace: testNamespace,
					},
					Spec: v1alpha1.DBaaSInstanceSpec{
						InventoryRef: v1alpha1.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						Name: instanceRefName,
					},
				}
				lastTransitionTime := getLastTransitionTimeForTest()
				instanceStatus := &v1alpha1.DBaaSInstanceStatus{
					InstanceID: instanceID,
					Phase:      v1alpha1.InstancePhaseReady,
					Conditions: []metav1.Condition{
						{
							Type:               "ProvisionReady",
							Status:             metav1.ConditionTrue,
							Reason:             "SyncOK",
							LastTransitionTime: metav1.Time{Time: lastTransitionTime},
						},
					},
				}

				BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInstance, metav1.ConditionTrue, testInstanceKind, instanceStatus))
				AfterEach(assertResourceDeletion(createdDBaaSInstance))

				Context("after creating DBaaSConnection", func() {
					connectionName := "test-instance-ref-connection-1"
					DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
						InventoryRef: v1alpha2.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						DatabaseServiceRef: &v1alpha2.NamespacedName{
							Name:      instanceRefName,
							Namespace: testNamespace,
						},
					}
					createdDBaaSConnection := &v1alpha2.DBaaSConnection{
						ObjectMeta: metav1.ObjectMeta{
							Name:      connectionName,
							Namespace: testNamespace,
						},
						Spec: *DBaaSConnectionSpec,
					}
					BeforeEach(assertResourceCreation(createdDBaaSConnection))
					AfterEach(assertResourceDeletion(createdDBaaSConnection))

					It("should create a provider connection", func() {
						expectedDBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
							InventoryRef: v1alpha2.NamespacedName{
								Name:      inventoryRefName,
								Namespace: testNamespace,
							},
							DatabaseServiceID:   instanceID,
							DatabaseServiceType: v1alpha2.InstanceDatabaseService,
						}
						assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, expectedDBaaSConnectionSpec)()

						By("checking if the Deployment is created")
						deployment := &appv1.Deployment{
							ObjectMeta: metav1.ObjectMeta{
								Name:      connectionName,
								Namespace: testNamespace,
							},
						}
						Eventually(func() bool {
							err := dRec.Get(ctx, client.ObjectKeyFromObject(deployment), deployment)
							if err != nil {
								return false
							}
							Expect(deployment.Spec.Replicas).ShouldNot(BeNil())
							Expect(*deployment.Spec.Replicas).ShouldNot(Equal(0))

							Expect(deployment.Labels).Should(BeNil())
							Expect(deployment.Annotations).ShouldNot(BeNil())
							mb, mbOk := deployment.Annotations["managed-by"]
							Expect(mbOk).Should(BeTrue())
							Expect(mb).Should(Equal("dbaas-operator"))
							owner, ownerOk := deployment.Annotations["owner"]
							Expect(ownerOk).Should(BeTrue())
							Expect(owner).Should(Equal(connectionName))
							ownerKind, ownerKindOk := deployment.Annotations["owner.kind"]
							Expect(ownerKindOk).Should(BeTrue())
							Expect(ownerKind).Should(Equal("DBaaSConnection"))
							ownerNS, ownerNSOk := deployment.Annotations["owner.namespace"]
							Expect(ownerNSOk).Should(BeTrue())
							Expect(ownerNS).Should(Equal(testNamespace))

							deploymentOwner := metav1.GetControllerOf(deployment)
							Expect(deploymentOwner).ShouldNot(BeNil())
							Expect(deploymentOwner.Kind).Should(Equal("DBaaSConnection"))
							Expect(deploymentOwner.Name).Should(Equal(connectionName))
							Expect(deploymentOwner.Controller).ShouldNot(BeNil())
							Expect(*deploymentOwner.Controller).Should(BeTrue())
							Expect(deploymentOwner.BlockOwnerDeletion).ShouldNot(BeNil())
							Expect(*deploymentOwner.BlockOwnerDeletion).Should(BeTrue())
							return true
						}, timeout).Should(BeTrue())
					})
					Context("when updating provider connection status", func() {
						lastTransitionTime := getLastTransitionTimeForTest()
						status := &v1alpha2.DBaaSConnectionStatus{
							Conditions: []metav1.Condition{
								{
									Type:               "ReadyForBinding",
									Status:             metav1.ConditionTrue,
									Reason:             "SyncOK",
									LastTransitionTime: metav1.Time{Time: lastTransitionTime},
								},
							},
							CredentialsRef: &v1.LocalObjectReference{
								Name: testSecret.Name,
							},
							ConnectionInfoRef: &v1.LocalObjectReference{
								Name: "testConnectionInfoRef",
							},
						}
						It("should update DBaaSConnection status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSConnection, metav1.ConditionTrue, testConnectionKind, status))
					})
				})
			})
		})
	})
})

var _ = Describe("DBaaSConnection controller - valid dev namespaces", func() {
	BeforeEach(assertResourceCreationIfNotExists(&testSecret))
	BeforeEach(assertResourceCreationIfNotExists(mongoProvider))
	BeforeEach(assertResourceCreationIfNotExists(&defaultPolicy))
	BeforeEach(assertDBaaSResourceStatusUpdated(&defaultPolicy, metav1.ConditionTrue, v1alpha1.Ready))

	Describe("reconcile", func() {
		Context("after creating DBaaSInventory w/ addtl dev namespace set", func() {
			otherNS := v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "other2",
				},
			}
			inventoryRefName := "test-inventory-ref-2"
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventoryPolicy: v1alpha2.DBaaSInventoryPolicy{
						ConnectionNamespaces: &[]string{otherNS.Name},
					},
					DBaaSInventorySpec: v1alpha2.DBaaSInventorySpec{
						CredentialsRef: &v1alpha2.LocalObjectReference{
							Name: testSecret.Name,
						},
					},
				},
			}
			lastTransitionTime := getLastTransitionTimeForTest()
			providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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

			Context("after creating DBaaSConnections in separate, valid dev namespaces", func() {
				connectionName := "test-connection-2"
				instanceID := "test-instanceID"
				DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
					InventoryRef: v1alpha2.NamespacedName{
						Name:      inventoryRefName,
						Namespace: testNamespace,
					},
					DatabaseServiceID:   instanceID,
					DatabaseServiceType: v1alpha2.InstanceDatabaseService,
				}
				createdDBaaSConnection := &v1alpha2.DBaaSConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      connectionName,
						Namespace: otherNS.Name,
					},
					Spec: *DBaaSConnectionSpec,
				}
				BeforeEach(assertResourceCreation(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdDBaaSConnection))

				It("should create a provider connection", assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				Context("when updating provider connection status", func() {
					lastTransitionTime := getLastTransitionTimeForTest()
					status := &v1alpha2.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								Type:               "ReadyForBinding",
								Status:             metav1.ConditionTrue,
								Reason:             "SyncOK",
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: testSecret.Name,
						},
						ConnectionInfoRef: &v1.LocalObjectReference{
							Name: "testConnectionInfoRef",
						},
					}
					It("should update DBaaSConnection status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSConnection, metav1.ConditionTrue, testConnectionKind, status))
				})

				Context("when updating DBaaSConnection spec", func() {
					DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
						InventoryRef: v1alpha2.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						DatabaseServiceID:   "updated-test-instanceID",
						DatabaseServiceType: v1alpha2.InstanceDatabaseService,
					}
					It("should update provider connection spec", assertProviderResourceSpecUpdated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				})

			})

			BeforeEach(assertResourceCreationIfNotExists(&otherNS))
			BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))
		})

		Context("after creating DBaaSInventory w/ wildcard dev namespace set", func() {
			otherNS := v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "other3",
				},
			}
			inventoryRefName := "test-inventory-ref-3"
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventoryPolicy: v1alpha2.DBaaSInventoryPolicy{
						ConnectionNamespaces: &[]string{"*"},
					},
					DBaaSInventorySpec: v1alpha2.DBaaSInventorySpec{
						CredentialsRef: &v1alpha2.LocalObjectReference{
							Name: testSecret.Name,
						},
					},
				},
			}
			lastTransitionTime := getLastTransitionTimeForTest()
			providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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

			Context("after creating DBaaSConnections in separate, valid dev namespaces", func() {
				connectionName := "test-connection-3"
				instanceID := "test-instanceID"
				DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
					InventoryRef: v1alpha2.NamespacedName{
						Name:      inventoryRefName,
						Namespace: testNamespace,
					},
					DatabaseServiceID:   instanceID,
					DatabaseServiceType: v1alpha2.InstanceDatabaseService,
				}
				createdDBaaSConnection := &v1alpha2.DBaaSConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      connectionName,
						Namespace: otherNS.Name,
					},
					Spec: *DBaaSConnectionSpec,
				}
				BeforeEach(assertResourceCreation(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdDBaaSConnection))

				It("should create a provider connection", assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				Context("when updating provider connection status", func() {
					lastTransitionTime := getLastTransitionTimeForTest()
					status := &v1alpha2.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								Type:               "ReadyForBinding",
								Status:             metav1.ConditionTrue,
								Reason:             "SyncOK",
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: testSecret.Name,
						},
						ConnectionInfoRef: &v1.LocalObjectReference{
							Name: "testConnectionInfoRef",
						},
					}
					It("should update DBaaSConnection status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSConnection, metav1.ConditionTrue, testConnectionKind, status))
				})

				Context("when updating DBaaSConnection spec", func() {
					DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
						InventoryRef: v1alpha2.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						DatabaseServiceID:   "updated-test-instanceID",
						DatabaseServiceType: v1alpha2.InstanceDatabaseService,
					}
					It("should update provider connection spec", assertProviderResourceSpecUpdated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				})

			})

			BeforeEach(assertResourceCreationIfNotExists(&otherNS))
			BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))
		})

		Context("after creating DBaaSInventory w/ connection namespace matchExpressions selector", func() {
			labels := map[string]string{"testlabel": "foo", "testagain": "bar"}
			otherNS := v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "expr-selectorns",
					Labels: labels,
				},
			}
			inventoryRefName := "test-inventory-ref-expr-selector"
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventoryPolicy: v1alpha2.DBaaSInventoryPolicy{
						ConnectionNsSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "testlabel",
									Operator: metav1.LabelSelectorOpExists,
								},
								{
									Key:      "testagain",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{"bar", "none"},
								},
								{
									Key:      "blah",
									Operator: metav1.LabelSelectorOpDoesNotExist,
								},
							},
						},
					},
					DBaaSInventorySpec: v1alpha2.DBaaSInventorySpec{
						CredentialsRef: &v1alpha2.LocalObjectReference{
							Name: testSecret.Name,
						},
					},
				},
			}
			lastTransitionTime := getLastTransitionTimeForTest()
			providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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

			Context("after creating DBaaSConnections in separate, valid dev namespaces", func() {
				connectionName := "test-connection-expr-selector"
				instanceID := "test-instanceID"
				DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
					InventoryRef: v1alpha2.NamespacedName{
						Name:      inventoryRefName,
						Namespace: testNamespace,
					},
					DatabaseServiceID:   instanceID,
					DatabaseServiceType: v1alpha2.InstanceDatabaseService,
				}
				createdDBaaSConnection := &v1alpha2.DBaaSConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      connectionName,
						Namespace: otherNS.Name,
					},
					Spec: *DBaaSConnectionSpec,
				}
				BeforeEach(assertResourceCreation(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdDBaaSConnection))

				It("should create a provider connection", assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				Context("when updating provider connection status", func() {
					lastTransitionTime := getLastTransitionTimeForTest()
					status := &v1alpha2.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								Type:               "ReadyForBinding",
								Status:             metav1.ConditionTrue,
								Reason:             "SyncOK",
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: testSecret.Name,
						},
						ConnectionInfoRef: &v1.LocalObjectReference{
							Name: "testConnectionInfoRef",
						},
					}
					It("should update DBaaSConnection status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSConnection, metav1.ConditionTrue, testConnectionKind, status))
				})

				Context("when updating DBaaSConnection spec", func() {
					DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
						InventoryRef: v1alpha2.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						DatabaseServiceID:   "updated-test-instanceID",
						DatabaseServiceType: v1alpha2.InstanceDatabaseService,
					}
					It("should update provider connection spec", assertProviderResourceSpecUpdated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				})

			})

			BeforeEach(assertResourceCreationIfNotExists(&otherNS))
			BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))
		})

		Context("after creating DBaaSInventory w/ connection namespace matchLabels selector", func() {
			labels := map[string]string{"testlabel": "foo", "testagain": "bar", "blah": "blah"}
			otherNS := v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "selectorns",
					Labels: labels,
				},
			}
			inventoryRefName := "test-inventory-ref-selector"
			createdDBaaSInventory := &v1alpha2.DBaaSInventory{
				ObjectMeta: metav1.ObjectMeta{
					Name:      inventoryRefName,
					Namespace: testNamespace,
				},
				Spec: v1alpha2.DBaaSOperatorInventorySpec{
					ProviderRef: v1alpha2.NamespacedName{
						Name: testProviderName,
					},
					DBaaSInventoryPolicy: v1alpha2.DBaaSInventoryPolicy{
						ConnectionNsSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
					},
					DBaaSInventorySpec: v1alpha2.DBaaSInventorySpec{
						CredentialsRef: &v1alpha2.LocalObjectReference{
							Name: testSecret.Name,
						},
					},
				},
			}
			lastTransitionTime := getLastTransitionTimeForTest()
			providerInventoryStatus := &v1alpha2.DBaaSInventoryStatus{
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

			Context("after creating DBaaSConnections in separate, valid dev namespaces", func() {
				connectionName := "test-connection-selector"
				instanceID := "test-instanceID"
				DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
					InventoryRef: v1alpha2.NamespacedName{
						Name:      inventoryRefName,
						Namespace: testNamespace,
					},
					DatabaseServiceID:   instanceID,
					DatabaseServiceType: v1alpha2.InstanceDatabaseService,
				}
				createdDBaaSConnection := &v1alpha2.DBaaSConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      connectionName,
						Namespace: otherNS.Name,
					},
					Spec: *DBaaSConnectionSpec,
				}
				BeforeEach(assertResourceCreation(createdDBaaSConnection))
				AfterEach(assertResourceDeletion(createdDBaaSConnection))

				It("should create a provider connection", assertProviderResourceCreated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				Context("when updating provider connection status", func() {
					lastTransitionTime := getLastTransitionTimeForTest()
					status := &v1alpha2.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								Type:               "ReadyForBinding",
								Status:             metav1.ConditionTrue,
								Reason:             "SyncOK",
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: testSecret.Name,
						},
						ConnectionInfoRef: &v1.LocalObjectReference{
							Name: "testConnectionInfoRef",
						},
					}
					It("should update DBaaSConnection status", assertDBaaSResourceProviderStatusUpdated(createdDBaaSConnection, metav1.ConditionTrue, testConnectionKind, status))
				})

				Context("when updating DBaaSConnection spec", func() {
					DBaaSConnectionSpec := &v1alpha2.DBaaSConnectionSpec{
						InventoryRef: v1alpha2.NamespacedName{
							Name:      inventoryRefName,
							Namespace: testNamespace,
						},
						DatabaseServiceID:   "updated-test-instanceID",
						DatabaseServiceType: v1alpha2.InstanceDatabaseService,
					}
					It("should update provider connection spec", assertProviderResourceSpecUpdated(createdDBaaSConnection, testConnectionKind, DBaaSConnectionSpec))
				})

			})

			BeforeEach(assertResourceCreationIfNotExists(&otherNS))
			BeforeEach(assertResourceCreationWithProviderStatus(createdDBaaSInventory, metav1.ConditionTrue, testInventoryKind, providerInventoryStatus))
			AfterEach(assertResourceDeletion(createdDBaaSInventory))
		})
	})
})
