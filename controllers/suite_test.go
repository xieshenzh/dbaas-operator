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
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context

var testProviderCM *v1.ConfigMap

const (
	testNamespace = "default"

	testCMName         = "mongodb-atlas"
	testProviderName   = "MongoDBAtlas"
	testInventoryKind  = "MongoDBAtlasInventory"
	testConnectionKind = "MongoDBAtlasConnection"

	timeout  = time.Second * 10
	duration = time.Second * 10
	interval = time.Millisecond * 250
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crd", "bases"),
			filepath.Join("..", "config", "test", "crd"),
		},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = dbaasv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	ctx = context.Background()

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:    scheme.Scheme,
		Namespace: testNamespace,
	})
	Expect(err).ToNot(HaveOccurred())

	By("mocking provider ConfigMap")
	DBaaSReconciler := &DBaaSReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}
	testProviderCM = &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testCMName,
			Namespace: testNamespace,
			Labels:    ConfigMapSelector,
		},
		Data: map[string]string{
			"connection_kind": testConnectionKind,
			"inventory_kind":  testInventoryKind,
			"provider":        testProviderName,
		},
	}
	cmList := v1.ConfigMapList{
		Items: []v1.ConfigMap{*testProviderCM},
	}
	providerList, err := DBaaSReconciler.ParseDBaaSProviderList(cmList)
	Expect(err).ToNot(HaveOccurred())

	err = (&DBaaSInventoryReconciler{
		DBaaSReconciler: DBaaSReconciler,
	}).SetupWithManager(k8sManager, providerList)
	Expect(err).ToNot(HaveOccurred())

	err = (&DBaaSConnectionReconciler{
		DBaaSReconciler: DBaaSReconciler,
	}).SetupWithManager(k8sManager, providerList)
	Expect(err).ToNot(HaveOccurred())

	err = (&ProviderConfigMapReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager, cmList)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func assertProviderConfigMapCreated() {
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
}
