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

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

var _ = Describe("mapstructure utils", func() {
	Describe("decode", func() {
		Describe("decode inventory status", func() {
			var (
				input    map[string]interface{}
				output   *v1alpha1.DBaaSInventoryStatus
				expected *v1alpha1.DBaaSInventoryStatus
			)

			Context("when decoding succeeds", func() {
				BeforeEach(func() {
					lastTransitionTimeString := "2021-06-18T20:03:20Z"
					lastTransitionTime, err := time.Parse(time.RFC3339, lastTransitionTimeString)
					Expect(err).NotTo(HaveOccurred())

					input = map[string]interface{}{
						"conditions": []map[string]interface{}{
							{
								"lastTransitionTime": lastTransitionTimeString,
								"message":            "Secret not found",
								"reason":             "InputError",
								"status":             "False",
								"type":               "SpecSynced",
							},
						},
					}
					output = &v1alpha1.DBaaSInventoryStatus{}
					expected = &v1alpha1.DBaaSInventoryStatus{
						Conditions: []metav1.Condition{
							{
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
								Message:            "Secret not found",
								Reason:             "InputError",
								Status:             "False",
								Type:               "SpecSynced",
							},
						},
					}
				})

				It("should populate the fields of DBaaSInventoryStatus", func() {
					err := decode(input, &output)
					Expect(err).NotTo(HaveOccurred())
					Expect(output).Should(Equal(expected))
				})
			})
		})

		Describe("decode inventory spec", func() {
			var (
				input    map[string]interface{}
				output   *v1alpha1.DBaaSInventorySpec
				expected *v1alpha1.DBaaSInventorySpec
			)

			Context("when decoding succeeds", func() {
				BeforeEach(func() {
					input = map[string]interface{}{
						"credentialsRef": map[string]interface{}{
							"name":      "testName",
							"namespace": "testNamespace",
						},
					}
					output = &v1alpha1.DBaaSInventorySpec{}
					expected = &v1alpha1.DBaaSInventorySpec{
						CredentialsRef: &v1alpha1.NamespacedName{
							Name:      "testName",
							Namespace: "testNamespace",
						},
					}
				})

				It("should populate the fields of DBaaSInventorySpec", func() {
					err := decode(input, output)
					Expect(err).NotTo(HaveOccurred())
					Expect(output).Should(Equal(expected))
				})
			})
		})

		Describe("decode connection status", func() {
			var (
				input    map[string]interface{}
				output   *v1alpha1.DBaaSConnectionStatus
				expected *v1alpha1.DBaaSConnectionStatus
			)

			Context("when decoding succeeds", func() {
				BeforeEach(func() {
					lastTransitionTimeString := "2021-06-18T20:03:20Z"
					lastTransitionTime, err := time.Parse(time.RFC3339, lastTransitionTimeString)
					Expect(err).NotTo(HaveOccurred())

					input = map[string]interface{}{
						"conditions": []map[string]interface{}{
							{
								"lastTransitionTime": lastTransitionTimeString,
								"message":            "Secret not found",
								"reason":             "InputError",
								"status":             "False",
								"type":               "SpecSynced",
							},
						},
						"credentialsRef": map[string]interface{}{
							"name": "testCredentialsRef",
						},
						"connectionInfo": map[string]interface{}{
							"name": "testConnectionInfo",
						},
					}
					output = &v1alpha1.DBaaSConnectionStatus{}
					expected = &v1alpha1.DBaaSConnectionStatus{
						Conditions: []metav1.Condition{
							{
								LastTransitionTime: metav1.Time{Time: lastTransitionTime},
								Message:            "Secret not found",
								Reason:             "InputError",
								Status:             "False",
								Type:               "SpecSynced",
							},
						},
						CredentialsRef: &v1.LocalObjectReference{
							Name: "testCredentialsRef",
						},
						ConnectionInfo: &v1.LocalObjectReference{
							Name: "testConnectionInfo",
						},
					}
				})

				It("should populate the fields of DBaaSConnectionStatus", func() {
					err := decode(input, output)
					Expect(err).NotTo(HaveOccurred())
					Expect(output).Should(Equal(expected))
				})
			})
		})

		Describe("decode connection spec", func() {
			var (
				input    map[string]interface{}
				output   *v1alpha1.DBaaSConnectionSpec
				expected *v1alpha1.DBaaSConnectionSpec
			)

			Context("when decoding succeeds", func() {
				BeforeEach(func() {
					input = map[string]interface{}{
						"inventoryRef": map[string]interface{}{
							"name": "testName",
						},
						"instanceID": "testInstanceID",
					}
					output = &v1alpha1.DBaaSConnectionSpec{}
					expected = &v1alpha1.DBaaSConnectionSpec{
						InventoryRef: &v1.LocalObjectReference{
							Name: "testName",
						},
						InstanceID: "testInstanceID",
					}
				})

				It("should populate the fields of DBaaSConnectionSpec", func() {
					err := decode(input, output)
					Expect(err).NotTo(HaveOccurred())
					Expect(output).Should(Equal(expected))
				})
			})
		})
	})
})
