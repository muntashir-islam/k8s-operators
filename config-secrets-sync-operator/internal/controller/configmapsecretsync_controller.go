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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/muntashir-islam/k8s-operators/config-secrets-sync-operator/api/v1alpha1"
)

// ConfigMapSecretSyncReconciler reconciles a ConfigMapSecretSync object
type ConfigMapSecretSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.muntashirislam.com,resources=configmapsecretsyncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.muntashirislam.com,resources=configmapsecretsyncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.muntashirislam.com,resources=configmapsecretsyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMapSecretSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ConfigMapSecretSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := log.FromContext(ctx)
	log.Info("Reconciling Configmaps Secret Sync")
	// Fetch the ConfigMapSecretSync instance
	configMapSecretSync := &appsv1alpha1.ConfigMapSecretSync{}
	err := r.Get(ctx, req.NamespacedName, configMapSecretSync)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "Request object not found, could have been deleted after reconcile request.")
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	sourceNamespace := configMapSecretSync.Spec.SourceNamespace
	targetNamespaces := configMapSecretSync.Spec.TargetNamespaces
	configMapNames := configMapSecretSync.Spec.ConfigMapNames
	secretNames := configMapSecretSync.Spec.SecretNames

	// Copy ConfigMaps
	if len(configMapNames) > 0 {
		for _, cmName := range configMapNames {
			sourceCM := &corev1.ConfigMap{}
			if err := r.Get(ctx, types.NamespacedName{Name: cmName, Namespace: sourceNamespace}, sourceCM); err != nil {
				if errors.IsNotFound(err) {
					continue
				}
				return reconcile.Result{}, err
			}

			for _, targetNamespace := range targetNamespaces {
				newCM := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      sourceCM.Name,
						Namespace: targetNamespace,
					},
					Data:       sourceCM.Data,
					BinaryData: sourceCM.BinaryData,
				}
				err := r.Create(ctx, newCM)
				if err != nil && errors.IsAlreadyExists(err) {
					// If the ConfigMap already exists, update it
					existingCM := &corev1.ConfigMap{}
					if err := r.Get(ctx, types.NamespacedName{Name: sourceCM.Name, Namespace: targetNamespace}, existingCM); err != nil {
						return reconcile.Result{}, err
					}
					existingCM.Data = sourceCM.Data
					existingCM.BinaryData = sourceCM.BinaryData
					if err := r.Update(ctx, existingCM); err != nil {
						return reconcile.Result{}, err
					}
				} else if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	// Copy Secrets
	if len(secretNames) > 0 {
		for _, secretName := range secretNames {
			sourceSecret := &corev1.Secret{}
			if err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: sourceNamespace}, sourceSecret); err != nil {
				if errors.IsNotFound(err) {
					continue
				}
				return reconcile.Result{}, err
			}

			for _, targetNamespace := range targetNamespaces {
				newSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      sourceSecret.Name,
						Namespace: targetNamespace,
					},
					Data:       sourceSecret.Data,
					StringData: sourceSecret.StringData,
					Type:       sourceSecret.Type,
				}
				err := r.Create(ctx, newSecret)
				if err != nil && errors.IsAlreadyExists(err) {
					// If the Secret already exists, update it
					existingSecret := &corev1.Secret{}
					if err := r.Get(ctx, types.NamespacedName{Name: sourceSecret.Name, Namespace: targetNamespace}, existingSecret); err != nil {
						return reconcile.Result{}, err
					}
					existingSecret.Data = sourceSecret.Data
					existingSecret.StringData = sourceSecret.StringData
					existingSecret.Type = sourceSecret.Type
					if err := r.Update(ctx, existingSecret); err != nil {
						return reconcile.Result{}, err
					}
				} else if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	// Requeue after 5 minutes to check for updates
	return reconcile.Result{RequeueAfter: time.Minute * 5}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapSecretSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.ConfigMapSecretSync{}).
		Complete(r)
}
