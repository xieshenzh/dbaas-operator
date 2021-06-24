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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/google/uuid"

	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
)

const (
	deploymentProviderSyncKey = "provider-sync-uuid"

	deploymentLabelKey   = "dbaas-operator-resource"
	deploymentLabelValue = "dbaas-operator-deployment"
)

// DBaaSProviderReconciler reconciles a DBaaSProvider object
type DBaaSProviderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dbaas.redhat.com,resources=*,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dbaas.redhat.com,resources=*/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dbaas.redhat.com,resources=*/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *DBaaSProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx, "DBaaS Provider", req.NamespacedName)

	ns, err := getInstallNamespace()
	if err != nil {
		logger.Error(err, "Error reading operator install namespace from env")
		return ctrl.Result{}, err
	}

	var deploymentList appsv1.DeploymentList
	opts := []client.ListOption{
		client.InNamespace(ns),
		client.MatchingLabels{
			deploymentLabelKey: deploymentLabelValue,
		},
	}

	if err := r.List(ctx, &deploymentList, opts...); err != nil {
		logger.Error(err, "Error fetching the operator's Deployment")
		return ctrl.Result{}, err
	} else if len(deploymentList.Items) != 1 {
		err = fmt.Errorf("found %d Deployments that match operator's Deployment", len(deploymentList.Items))
		return ctrl.Result{}, err
	}

	deployment := deploymentList.Items[0]
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = map[string]string{deploymentProviderSyncKey: uuid.New().String()}
	} else {
		deployment.Spec.Template.Annotations[deploymentProviderSyncKey] = uuid.New().String()
	}
	if err := r.Update(ctx, &deployment); err != nil {
		logger.Error(err, "Error updating the operator's Deployment to reload providers")
	}
	logger.Info("Operator's Deployment updated, operator will restart")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBaaSProviderReconciler) SetupWithManager(mgr ctrl.Manager, providerList v1alpha1.DBaaSProviderList) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DBaaSProvider{}).
		WithEventFilter(filterEventPredicate(providerList)).
		Complete(r)
}

func filterEventPredicate(providerList v1alpha1.DBaaSProviderList) predicate.Predicate {
	providerNames := make([]string, len(providerList.Items))
	for i, cm := range providerList.Items {
		providerNames[i] = cm.Name
	}

	return predicate.Funcs{
		// Ignore CreateEvent for existing provider
		CreateFunc: func(createEvent event.CreateEvent) bool {
			for _, n := range providerNames {
				if n == createEvent.Object.GetName() {
					return false
				}
			}
			return true
		},
		// Ignore UpdateEvent if the provider is not updated
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return updateEvent.ObjectNew.GetGeneration() != updateEvent.ObjectOld.GetGeneration()
		},
		// Reload the operator for DeleteEvent
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return true
		},
		// Ignore GenericEvent
		GenericFunc: func(genericEvent event.GenericEvent) bool {
			return false
		},
	}
}
