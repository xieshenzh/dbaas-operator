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
	"reflect"

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

const (
	testProviderName   = "MongoDBAtlas"
	testInventoryKind  = "MongoDBAtlasInventory"
	testConnectionKind = "MongoDBAtlasConnection"
)

var providerConfigMap = &v1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "mongodb-atlas",
		Namespace: testNamespace,
		Labels:    ConfigMapSelector,
	},
	Data: map[string]string{
		"connection_kind": testConnectionKind,
		"inventory_kind":  testInventoryKind,
		"provider":        testProviderName,
	},
}

func assertProviderConfigMapCreation() func() {
	return func() {
		By("checking the provider ConfigMap created")
		createdCM := v1.ConfigMap{}
		if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(providerConfigMap), &createdCM); err != nil {
			if errors.IsNotFound(err) {
				By("creating the provider ConfigMap")
				Expect(k8sClient.Create(ctx, providerConfigMap)).Should(Succeed())

				By("checking the provider ConfigMap created")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, client.ObjectKeyFromObject(providerConfigMap), &createdCM)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
			} else {
				Fail(err.Error())
			}
		}
	}
}

func assertResourceCreation(object client.Object) func() {
	return func() {
		By("creating resource")
		object.SetResourceVersion("")
		Expect(k8sClient.Create(ctx, object)).Should(Succeed())

		By("checking the resource created")
		Eventually(func() bool {
			err := k8sClient.Get(ctx, client.ObjectKeyFromObject(object), object)
			if err != nil {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())
	}
}

func assertResourceDeletion(object client.Object) func() {
	return func() {
		By("deleting resource")
		Expect(k8sClient.Delete(ctx, object)).Should(Succeed())

		By("checking the resource deleted")
		Eventually(func() bool {
			err := k8sClient.Get(ctx, client.ObjectKeyFromObject(object), object)
			if err != nil && errors.IsNotFound(err) {
				return true
			}
			return false
		}, timeout, interval).Should(BeTrue())
	}
}

func assertProviderResourceCreated(object client.Object, providerResourceKind string, DBaaSResourceSpec interface{}) func() {
	return func() {
		By("checking a provider resource created")
		objectKey := client.ObjectKeyFromObject(object)
		providerResource := &unstructured.Unstructured{}
		providerResource.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   v1alpha1.GroupVersion.Group,
			Version: v1alpha1.GroupVersion.Version,
			Kind:    providerResourceKind,
		})
		Eventually(func() bool {
			err := k8sClient.Get(ctx, objectKey, providerResource)
			if err != nil {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())

		By("checking the provider resource spec is correct")
		providerResourceSpec, existSpec := providerResource.UnstructuredContent()["spec"]
		Expect(existSpec).Should(Equal(true))
		switch v := object.(type) {
		case *v1alpha1.DBaaSInventory:
			spec := &v1alpha1.DBaaSInventorySpec{}
			err := decode(providerResourceSpec, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(spec).Should(Equal(DBaaSResourceSpec))
		case *v1alpha1.DBaaSConnection:
			spec := &v1alpha1.DBaaSConnectionSpec{}
			err := decode(providerResourceSpec, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(spec).Should(Equal(DBaaSResourceSpec))
		default:
			_ = v.GetName() // to avoid syntax error
			Fail("invalid test object")
		}
	}
}

func assertDBaaSResourceStatusUpdated(object client.Object, providerResourceKind string, providerResourceStatus interface{}) func() {
	return func() {
		By("checking the DBaaS resource status has no conditions")
		objectKey := client.ObjectKeyFromObject(object)
		Consistently(func() (int, error) {
			err := k8sClient.Get(ctx, objectKey, object)
			if err != nil {
				return -1, err
			}
			switch v := object.(type) {
			case *v1alpha1.DBaaSInventory:
				return len(v.Status.Conditions), nil
			case *v1alpha1.DBaaSConnection:
				return len(v.Status.Conditions), nil
			default:
				Fail("invalid test object")
				return -1, err
			}
		}, duration, interval).Should(Equal(0))

		By("getting the provider resource")
		providerResource := &unstructured.Unstructured{}
		providerResource.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   v1alpha1.GroupVersion.Group,
			Version: v1alpha1.GroupVersion.Version,
			Kind:    providerResourceKind,
		})
		Eventually(func() bool {
			err := k8sClient.Get(ctx, objectKey, providerResource)
			if err != nil {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())

		By("updating the provider resource status")
		providerResource.UnstructuredContent()["status"] = providerResourceStatus
		Expect(k8sClient.Status().Update(ctx, providerResource)).Should(Succeed())

		By("checking the DBaaS resource status updated")
		Eventually(func() (int, error) {
			err := k8sClient.Get(ctx, objectKey, object)
			if err != nil {
				return -1, err
			}
			switch v := object.(type) {
			case *v1alpha1.DBaaSInventory:
				return len(v.Status.Conditions), nil
			case *v1alpha1.DBaaSConnection:
				return len(v.Status.Conditions), nil
			default:
				Fail("invalid test object")
				return -1, err
			}
		}, timeout, interval).Should(Equal(1))
		switch v := object.(type) {
		case *v1alpha1.DBaaSInventory:
			Expect(&v.Status).Should(Equal(providerResourceStatus))
		case *v1alpha1.DBaaSConnection:
			Expect(&v.Status).Should(Equal(providerResourceStatus))
		default:
			Fail("invalid test object")
		}
	}
}

func assertProviderResourceSpecUpdated(object client.Object, providerResourceKind string, DBaaSResourceSpec interface{}) func() {
	return func() {
		By("updating the DBaaS resource spec")
		objectKey := client.ObjectKeyFromObject(object)
		err := k8sClient.Get(ctx, objectKey, object)
		Expect(err).NotTo(HaveOccurred())

		switch v := object.(type) {
		case *v1alpha1.DBaaSInventory:
			v.Spec.DBaaSInventorySpec = *DBaaSResourceSpec.(*v1alpha1.DBaaSInventorySpec)
		case *v1alpha1.DBaaSConnection:
			v.Spec = *DBaaSResourceSpec.(*v1alpha1.DBaaSConnectionSpec)
		default:
			Fail("invalid test object")
		}
		Expect(k8sClient.Update(ctx, object)).Should(Succeed())

		By("checking the provider resource status updated")
		providerResource := &unstructured.Unstructured{}
		providerResource.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   v1alpha1.GroupVersion.Group,
			Version: v1alpha1.GroupVersion.Version,
			Kind:    providerResourceKind,
		})
		Eventually(func() bool {
			err := k8sClient.Get(ctx, objectKey, providerResource)
			if err != nil {
				return false
			}

			providerResourceSpec, existSpec := providerResource.UnstructuredContent()["spec"]
			Expect(existSpec).Should(Equal(true))
			switch v := object.(type) {
			case *v1alpha1.DBaaSInventory:
				spec := &v1alpha1.DBaaSInventorySpec{}
				err = decode(providerResourceSpec, spec)
				Expect(err).NotTo(HaveOccurred())
				return reflect.DeepEqual(spec, DBaaSResourceSpec)
			case *v1alpha1.DBaaSConnection:
				spec := &v1alpha1.DBaaSConnectionSpec{}
				err = decode(providerResourceSpec, spec)
				Expect(err).NotTo(HaveOccurred())
				return reflect.DeepEqual(spec, DBaaSResourceSpec)
			default:
				_ = v.GetName() // to avoid syntax error
				Fail("invalid test object")
				return false
			}
		}, timeout, interval).Should(BeTrue())
	}
}
