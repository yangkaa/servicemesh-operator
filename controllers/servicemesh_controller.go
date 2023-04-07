/*
Copyright 2022.

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
	rainbondiov1alpha1 "github.com/goodrain/servicemesh-operator/pkg/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	//"k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/clientset/versioned"

	//"github.com/goodrain/servicemesh-operator/api/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	selectorField = ".spec.selector"
)

// ServiceMeshReconciler reconciles a ServiceMesh object
type ServiceMeshReconciler struct {
	//rainbondClient versioned.Interface
	client.Client
	Scheme *runtime.Scheme
	//rainbondClient versioned.Interfaced
}

//+kubebuilder:rbac:groups=rainbond.io,resources=servicemeshes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rainbond.io,resources=servicemeshes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rainbond.io,resources=servicemeshes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceMesh object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ServiceMeshReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logrus.Infof("reconcile servicemesh: %v", req.NamespacedName)
	var servicemesh rainbondiov1alpha1.ServiceMesh
	err := r.Client.Get(ctx, req.NamespacedName, &servicemesh)
	if err != nil {
		logrus.Infof("get servicemesh error: %v", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var servicemeshClass rainbondiov1alpha1.ServiceMeshClass
	err = r.Client.Get(ctx, client.ObjectKey{Name: servicemesh.Provisioner}, &servicemeshClass)
	if err != nil {
		logrus.Infof("get servicemeshclass error: %v", err)
		return ctrl.Result{RequeueAfter: 3 * time.Second}, client.IgnoreNotFound(err)
	}
	// handle deployment
	var deployments v1.DeploymentList
	err = r.List(ctx, &deployments, client.InNamespace(req.Namespace), client.MatchingLabels(servicemesh.Selector))
	if err != nil {
		logrus.Infof("get deployment error: %v", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	for _, deploy := range deployments.Items {
		logrus.Infof("start patch deployment: %v", deploy.Name)
		for _, item := range servicemeshClass.InjectMethods {
			err = r.Get(ctx, types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, &deploy)
			if err != nil {
				logrus.Infof("get deployment error: %v", err)
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
			if deploy.Labels == nil {
				deploy.Labels = make(map[string]string)
			}
			if deploy.Annotations == nil {
				deploy.Annotations = make(map[string]string)
			}
			var value string
			if item.Value != "" {
				value = item.Value
			} else if item.ValueFrom != nil {
				if item.ValueFrom.Type == "label" {
					value = deploy.Labels[item.ValueFrom.Name]
				}
				if item.ValueFrom.Type == "annotation" {
					value = deploy.Annotations[item.ValueFrom.Name]
				}
				if item.ValueFrom.Type == "name" {
					value = deploy.Name
				}
			}
			if item.Method == "label" {
				if _, ok := deploy.Labels[item.Name]; !ok {
					deploy.Labels[item.Name] = value
				}
				if _, ok := deploy.Spec.Template.Labels[item.Name]; !ok {
					deploy.Spec.Template.Labels[item.Name] = value
				}
				if item.Cover {
					deploy.Labels[item.Name] = value
					deploy.Spec.Template.Labels[item.Name] = value
				}
			}
			if item.Method == "annotation" {
				if _, ok := deploy.Annotations[item.Name]; !ok {
					deploy.Annotations[item.Name] = value
				}
				if _, ok := deploy.Spec.Template.Annotations[item.Name]; !ok {
					deploy.Spec.Template.Annotations[item.Name] = value
				}
				if item.Cover {
					deploy.Annotations[item.Name] = value
					deploy.Spec.Template.Annotations[item.Name] = value
				}
			}
			err = r.Patch(ctx, &deploy, client.Merge)
			if err != nil {
				logrus.Infof("patch deployment error: %v", err)
				return ctrl.Result{RequeueAfter: 3 * time.Second}, client.IgnoreNotFound(err)
			}
		}
	}

	// handle statefulset
	var statefulsets v1.StatefulSetList
	err = r.List(ctx, &statefulsets, client.InNamespace(req.Namespace), client.MatchingLabels(servicemesh.Selector))
	if err != nil {
		logrus.Infof("get statefulset error: %v", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	for _, statefulset := range statefulsets.Items {
		logrus.Infof("start patch statefulset: %v", statefulset.Name)
		for _, item := range servicemeshClass.InjectMethods {
			err = r.Get(ctx, types.NamespacedName{Name: statefulset.Name, Namespace: statefulset.Namespace}, &statefulset)
			if err != nil {
				logrus.Infof("get statefulset error: %v", err)
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
			if item.Method == "label" {
				if statefulset.Labels == nil {
					statefulset.Labels = make(map[string]string)
				}
				statefulset.Labels[item.Name] = item.Value
				statefulset.Spec.Template.Labels[item.Name] = item.Value
			}
			if item.Method == "annotation" {
				if statefulset.Annotations == nil {
					statefulset.Annotations = make(map[string]string)
				}
				statefulset.Annotations[item.Name] = item.Value
				statefulset.Spec.Template.Annotations[item.Name] = item.Value
			}
			err = r.Patch(ctx, &statefulset, client.Merge)
			if err != nil {
				logrus.Infof("patch statefulset error: %v", err)
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
		}
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceMeshReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rainbondiov1alpha1.ServiceMesh{}).
		Watches(&source.Kind{Type: &v1.Deployment{}}, handler.EnqueueRequestsFromMapFunc(r.generateReconcileReq), builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Watches(&source.Kind{Type: &v1.StatefulSet{}}, handler.EnqueueRequestsFromMapFunc(r.generateReconcileReq), builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Complete(r)
}

func (r *ServiceMeshReconciler) generateReconcileReq(object client.Object) []reconcile.Request {
	var smList rainbondiov1alpha1.ServiceMeshList
	err := r.List(context.Background(), &smList, client.InNamespace(object.GetNamespace()))
	if err != nil {
		logrus.Infof("get servicemesh error: %v", err)
		return nil
	}
	if len(smList.Items) == 0 {
		return nil
	}
	// If the labels of object contains selector, reconcile is triggered
	var result []reconcile.Request
	for _, sm := range smList.Items {
		for k, v := range sm.Selector {
			if object.GetLabels()[k] == v {
				logrus.Infof("find servicemesh: %v，reconcile object %v, key [%s], val [%s]", sm.Name, object.GetName(), k, v)
				result = append(result, reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      sm.Name,
					Namespace: sm.Namespace,
				}})
			}
		}
	}
	return result
}
