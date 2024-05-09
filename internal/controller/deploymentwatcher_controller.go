/*
Copyright 2024.

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

package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	watcherv1alpha1 "github.com/kcddhaka/kcd2024-deployment-watcher/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
)

// DeploymentWatcherReconciler reconciles a DeploymentWatcher object
type DeploymentWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=watcher.operators.kcddhaka.org,resources=deploymentwatchers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=watcher.operators.kcddhaka.org,resources=deploymentwatchers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=watcher.operators.kcddhaka.org,resources=deploymentwatchers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeploymentWatcher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
var oldSpecs = make(map[string]v1.DeploymentSpec)

func (r *DeploymentWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	Log := log.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	//Log.Info("Reconcile Called")

	dw := &watcherv1alpha1.DeploymentWatcher{}

	err := r.Get(ctx, req.NamespacedName, dw)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, deploy := range dw.Spec.Deployments {
		deployment := &v1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: deploy.Namespace,
			Name:      deploy.Name,
		}, deployment)
		if err != nil {
			return ctrl.Result{}, err
		}

		key := deploy.Namespace + "/" + deploy.Name
		oldSpec, ok := oldSpecs[key]
		if ok && !reflect.DeepEqual(oldSpec, deployment.Spec) {
			Log.Info("Deployment spec has changed", "namespace", deploy.Namespace, "name", deploy.Name)

			if !reflect.DeepEqual(oldSpec.Template.Spec.Containers[0].Image, deployment.Spec.Template.Spec.Containers[0].Image) {
				Log.Info("Deployment image has changed", "namespace", deploy.Namespace, "name", deploy.Name)

				err = notify(deploy.Namespace, deploy.Name, deployment.Spec.Template.Spec.Containers[0].Image, dw.Spec.Slack.Token, dw.Spec.Slack.Channel)
				if err != nil {
					return ctrl.Result{}, err
				}
			} else if !reflect.DeepEqual(*oldSpec.Replicas, *deployment.Spec.Replicas) {
				Log.Info("Deployment replicas have changed", "namespace", deploy.Namespace, "name", deploy.Name)
			}
		}

		oldSpecs[key] = deployment.Spec
	}

	return ctrl.Result{RequeueAfter: time.Duration(10 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&watcherv1alpha1.DeploymentWatcher{}).
		Complete(r)
}
