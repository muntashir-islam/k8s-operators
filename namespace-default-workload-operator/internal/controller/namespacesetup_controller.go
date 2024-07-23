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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	resourcev1 "github.com/muntashir-islam/k8s-operators/namespace-default-workload-operator/api/v1"
)

// NamespaceSetupReconciler reconciles a NamespaceSetup object
type NamespaceSetupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=resource.muntashirislam.com,resources=namespacesetups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=resource.muntashirislam.com,resources=namespacesetups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=resource.muntashirislam.com,resources=namespacesetups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NamespaceSetup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NamespaceSetupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling NamespaceSetup")

	// Fetch the NamespaceSetup instance
	namespaceSetup := &resourcev1.NamespaceSetup{}
	err := r.Get(ctx, req.NamespacedName, namespaceSetup)
	if err != nil {
		log.Error(err, "Unable to fetch NamespaceSetup")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get all namespaces
	namespaceList := &corev1.NamespaceList{}
	err = r.List(ctx, namespaceList)
	if err != nil {
		log.Error(err, "Unable to list namespaces")
		return ctrl.Result{}, err
	}

	// Iterate over namespaces and create resources if not already present
	for _, ns := range namespaceList.Items {
		// Check if deployment already exists
		deployment := &appsv1.Deployment{}
		err = r.Get(ctx, client.ObjectKey{Name: "example-deployment", Namespace: ns.Name}, deployment)
		if err != nil && client.IgnoreNotFound(err) != nil {
			log.Error(err, "Unable to get Deployment")
			return ctrl.Result{}, err
		}
		if deployment.Name == "" {
			// Create Deployment
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "initial-deployment",
					Namespace: ns.Name,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: int32Ptr(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "initial"},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": "initial"},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  "initial-container",
								Image: namespaceSetup.Spec.ImageRepo,
							}},
						},
					},
				},
			}

			if err := r.Create(ctx, deployment); err != nil {
				log.Error(err, "Unable to create Deployment for Namespace", "namespace", ns.Name)
				return ctrl.Result{}, err
			}

			// Create Service
			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "initial-service",
					Namespace: ns.Name,
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{"app": "initial"},
					Ports: []corev1.ServicePort{{
						Protocol:   corev1.ProtocolTCP,
						Port:       int32(namespaceSetup.Spec.ServicePort),
						TargetPort: intstr.IntOrString{IntVal: int32(namespaceSetup.Spec.TargetPort)},
					}},
				},
			}

			if err := r.Create(ctx, service); err != nil {
				log.Error(err, "Unable to create Service for Namespace", "namespace", ns.Name)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *NamespaceSetupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&resourcev1.NamespaceSetup{}).
		Complete(r)
}

func int32Ptr(i int32) *int32 { return &i }
