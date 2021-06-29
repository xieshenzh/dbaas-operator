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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Provider ConfigMap controller", func() {
	operatorDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "operator-deployment",
			Namespace: testNamespace,
			Labels: map[string]string{
				deploymentLabelKey: deploymentLabelValue,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"control-plane": "controller-manager",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "operator-deployment-template",
					Labels: map[string]string{
						"control-plane": "controller-manager",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "operator-deployment-pod",
							Image: "golang:1.16",
						},
					},
				},
			},
		},
	}

	BeforeEach(assertResourceCreation(operatorDeployment))
	AfterEach(assertResourceDeletion(operatorDeployment))

	Describe("reconcile", func() {
		Context("after creating, updating or deleting provider ConfigMap", func() {
			It("should update the operator Deployment", func() {
				testConfigMap := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "crunchy-bridge",
						Namespace: testNamespace,
						Labels:    ConfigMapSelector,
					},
					Data: map[string]string{
						"connection_kind": "CrunchyBridgeConnection",
						"inventory_kind":  "CrunchyBridgeInventory",
						"provider":        "CrunchyBridge",
					},
				}

				By("first creating new provider ConfigMap")
				oldUUID := getOperatorDeploymentAnnotation(operatorDeployment)
				assertResourceCreation(testConfigMap)()
				newUUID := ""
				Eventually(func() bool {
					newUUID = getOperatorDeploymentAnnotation(operatorDeployment)
					return newUUID != oldUUID
				}, timeout, interval).Should(BeTrue())

				By("then updating provider ConfigMap")
				oldUUID = newUUID
				newData := "display_name"
				configMap := testConfigMap.DeepCopy()
				configMap.Data[newData] = "Crunchy Bridge"
				Expect(k8sClient.Update(ctx, configMap)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKeyFromObject(testConfigMap), testConfigMap)
					if err != nil {
						return false
					}
					_, exist := testConfigMap.Data[newData]
					return exist
				}, timeout, interval).Should(BeTrue())
				Eventually(func() bool {
					newUUID = getOperatorDeploymentAnnotation(operatorDeployment)
					return newUUID != oldUUID
				}, timeout, interval).Should(BeTrue())

				By("then deleting provider ConfigMap")
				oldUUID = newUUID
				assertResourceDeletion(testConfigMap)()
				Eventually(func() bool {
					newUUID = getOperatorDeploymentAnnotation(operatorDeployment)
					return newUUID != oldUUID
				}, timeout, interval).Should(BeTrue())
			})
		})
	})
})

func getOperatorDeploymentAnnotation(deployment *appsv1.Deployment) string {
	By("checking the DBaaS resource status has no conditions")
	err := k8sClient.Get(ctx, client.ObjectKeyFromObject(deployment), deployment)
	Expect(err).NotTo(HaveOccurred())

	if deployment.Spec.Template.Annotations != nil {
		if value, exist := deployment.Spec.Template.Annotations[deploymentProviderSyncKey]; exist {
			return value
		}
	}
	return ""
}
