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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

var _ = Describe("mapstructure utils", func() {
	Describe("decode", func() {
		Describe("decode inventory status", func() {
			var (
				input    map[string]interface{}
				output   v1alpha1.DBaaSInventoryStatus
				expected v1alpha1.DBaaSInventoryStatus
			)

			Context("when decoding succeeds", func() {
				BeforeEach(func() {
					lastTransitionTimeString := "2021-06-18T20:03:20Z"
					lastTransitionTime, err := time.Parse(time.RFC3339, lastTransitionTimeString)
					Expect(err).NotTo(HaveOccurred())

					input = map[string]interface{}{
						"type": "MongoDB",
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
					output = v1alpha1.DBaaSInventoryStatus{}
					expected = v1alpha1.DBaaSInventoryStatus{
						Type: "MongoDB",
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
	})
})
