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

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	consolev1alpha1 "github.com/openshift/api/console/v1alpha1"
	operatorv1 "github.com/openshift/api/operator/v1"
	coreosv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorframwork "github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"github.com/RHEcosystemAppEng/dbaas-operator/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(operatorframwork.AddToScheme(scheme))
	utilruntime.Must(coreosv1.AddToScheme(scheme))
	utilruntime.Must(consolev1alpha1.Install(scheme))
	utilruntime.Must(operatorv1.Install(scheme))

	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "e4addb06.redhat.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	DBaaSReconciler := &controllers.DBaaSReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	if DBaaSReconciler.InstallNamespace, err = controllers.GetInstallNamespace(); err != nil {
		setupLog.Error(err, "unable to retrieve install namespace. default Tenant object cannot be installed")
	}

	connectionCtrl, err := (&controllers.DBaaSConnectionReconciler{
		DBaaSReconciler: DBaaSReconciler,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DBaaSConnection")
		os.Exit(1)
	}
	inventoryCtrl, err := (&controllers.DBaaSInventoryReconciler{
		DBaaSReconciler: DBaaSReconciler,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DBaaSInventory")
		os.Exit(1)
	}
	err = (&controllers.DBaaSProviderReconciler{
		DBaaSReconciler: DBaaSReconciler,
		ConnectionCtrl:  connectionCtrl,
		InventoryCtrl:   inventoryCtrl,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DBaaSProvider")
		os.Exit(1)
	}
	//We'll just make sure to set `ENABLE_WEBHOOKS=false` when we run locally.

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&v1alpha1.DBaaSConnection{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DBaaSConnection")
			os.Exit(1)
		}
	}
	err = (&controllers.DBaaSTenantReconciler{
		DBaaSReconciler: DBaaSReconciler,
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DBaaSTenant")
		os.Exit(1)
	}
	if err = (&controllers.DBaaSPlatformReconciler{
		DBaaSReconciler: DBaaSReconciler,
		Log:             ctrl.Log.WithName("controllers").WithName("DBaaSPlatform"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DBaaSPlatform")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
