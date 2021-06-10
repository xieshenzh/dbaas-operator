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
	"fmt"
	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DBaaSConnectionReconciler reconciles a DBaaSConnection object
type DBaaSConnectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dbaas.redhat.com,resources=dbaasconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dbaas.redhat.com,resources=dbaasconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dbaas.redhat.com,resources=dbaasconnections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *DBaaSConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx, "dbaasconnection", req.NamespacedName)

	var connection v1alpha1.DBaaSConnection
	if err := r.Get(ctx, req.NamespacedName, &connection); err != nil {
		if errors.IsNotFound(err) {
			// CR deleted since request queued, child objects getting GC'd, no requeue
			logger.Info("DBaaSConnection resource not found, has been deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Error fetching DBaaSConnection for reconcile")
		return ctrl.Result{}, err
	}

	if connection.Spec.InventoryRef == nil {
		err := fmt.Errorf("inventory reference is missing for DBaaS connection %s", connection.Name)
		logger.Error(err, "Invalid DBaaSConnection for reconcile")
		return ctrl.Result{}, err
	}

	var inventory v1alpha1.DBaaSInventory
	if err := r.Get(ctx, types.NamespacedName{Namespace: connection.Namespace, Name: connection.Spec.InventoryRef.Name}, &inventory); err != nil {
		logger.Error(err, "Error fetching DBaaSInventory resource reference for DBaaSConnection reconcile", "Inventory", connection.Spec.InventoryRef.Name)
		return ctrl.Result{}, err
	}

	provider, err := getDBaaSProvider(inventory.Spec.Provider, r.Client, ctx, logger)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "Requested DBaaS Provider is not configured in this environment", "Provider", provider.Provider)
			return ctrl.Result{}, err
		}
		logger.Error(err, "Error reading configured DBaaS providers")
		return ctrl.Result{}, err
	}

	logger.Info("Found DBaaS provider", "provider", provider)
	providerConnection := &unstructured.Unstructured{}
	providerConnection.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   connection.GroupVersionKind().Group,
		Version: connection.GroupVersionKind().Version,
		Kind:    provider.ConnectionKind,
	})
	providerConnection.SetNamespace(connection.GetNamespace())
	providerConnection.SetName(connection.GetName())
	providerConnection.UnstructuredContent()["spec"] = connection.Spec.DeepCopy()
	logger.Info("Connection resource created as ", "providerConnection", providerConnection)
	if err = r.Create(ctx, providerConnection); err != nil {
		logger.Error(err, "Error creating a provider connection", "providerConnection", providerConnection)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBaaSConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DBaaSConnection{}).
		Complete(r)
}
