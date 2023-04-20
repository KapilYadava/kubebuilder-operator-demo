/*
Copyright 2023.

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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	demov1 "crdemo/api/v1"
)

// MyCustomResourceReconciler reconciles a MyCustomResource object
type MyCustomResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demo.mycustomresource,resources=mycustomresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.mycustomresource,resources=mycustomresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.mycustomresource,resources=mycustomresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyCustomResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MyCustomResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(1).Info("Reconcile")

	log.V(1).Info("Reconcile Resource", "Request", req)

	myCustomResource := demov1.MyCustomResource{}
	if err := r.Get(ctx, req.NamespacedName, &myCustomResource); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	name, namespace := myCustomResource.Name, myCustomResource.Namespace
	labels := map[string]string{
		"foo": "bar",
	}
	log.V(1).Info("Get name and namesapce", "name:", name, "namespace:", namespace)

	// pod to be created
	desiredPod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image: "nginx:latest",
					Name:  "nginx-container",
				},
			},
		},
	}

	// log.V(1).Info("Set Controller Reference")
	// if err := controllerutil.SetControllerReference(&myCustomResource, &desiredPod, r.Scheme); err != nil {
	// 	return ctrl.Result{}, err
	// }

	foundPod := &corev1.Pod{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(&desiredPod), foundPod); err != nil {
		if errors.IsNotFound(err) {
			log.V(1).Info("Desired Pod Created")
			return ctrl.Result{}, r.Create(ctx, &desiredPod)
		}

		return ctrl.Result{}, err
	}

	// pod was found
	// validate the labels of the found pod
	expectedLabels := map[string]string{
		"foo": "bar",
	}

	foundLabels := foundPod.Labels
	// if found labels and expected lables are equal don't do anything else update the pod
	if reflect.DeepEqual(expectedLabels, foundLabels) {
		return ctrl.Result{}, nil
	}

	foundPod.Labels = expectedLabels
	log.V(1).Info("Found Pod Updated")
	return ctrl.Result{}, r.Update(ctx, foundPod)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyCustomResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		//For(&demov1.MyCustomResource{}).
		//Owns(&corev1.Pod{}).
		Named("mycustom-resource").
		Watches(
			&source.Kind{Type: &demov1.MyCustomResource{}},
			&handler.EnqueueRequestForObject{},
		).
		Watches(
			&source.Kind{Type: &corev1.Pod{}},
			&handler.EnqueueRequestForOwner{IsController: true, OwnerType: &demov1.MyCustomResource{}},
		).
		Watches(
			&source.Kind{Type: &demov1.Machine{}},
			handler.Funcs{
				UpdateFunc: func(event event.UpdateEvent, queue workqueue.RateLimitingInterface) {
					machine := event.ObjectNew.(*demov1.Machine)
					if machine.Status.State == demov1.MachineStateRunning {
						queue.Add(ctrl.Request{NamespacedName: types.NamespacedName{Name: machine.Name, Namespace: corev1.NamespaceDefault}})
					}
				},
			},
		).
		Complete(r)
}
