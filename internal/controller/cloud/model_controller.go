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

package cloud

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
)

// ModelReconciler reconciles a Model object
type ModelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	ModelDeployerImage string
}

//+kubebuilder:rbac:groups=cloud.encoder.run,resources=models,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=models/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=models/finalizers,verbs=update
//+kubebuilder:rbac:groups=serving.kserve.io,resources=inferenceservices,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *ModelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var model v1alpha1.Model
	if err := r.Get(ctx, req.NamespacedName, &model); err != nil {
		log.Error(err, "unable to fetch Model")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// First time around ensure that status is NOT_DEPLOYED
	if model.Status.State == nil {
		state := v1alpha1.ModelStateNotDeployed
		model.Status.State = &state
		if err := r.Status().Update(ctx, &model); err != nil {
			return ctrl.Result{}, err
		}
		// Request a requeue to ensure that the status is update before proceeding.
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.ensureInferenceService(ctx, model); err != nil {
		log.Error(err, "unable to ensure deployment")
		return ctrl.Result{}, err
	}

	if err := r.ensureInferenceCleanup(ctx, model); err != nil {
		log.Error(err, "unable to ensure deployment cleanup")
		return ctrl.Result{}, err
	}

	if err := r.ensureStatus(ctx, model); err != nil {
		log.Error(err, "unable to ensure status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// ensureInferenceService ensures that the inference service for the model is created
// if deployment spec exists and is enabled.
func (r *ModelReconciler) ensureInferenceService(ctx context.Context, model v1alpha1.Model) error {
	log := log.FromContext(ctx)
	// Check if the deployment spec exists and is enabled.
	if model.Spec.Deployment == nil || !model.Spec.Deployment.Enabled {
		return nil
	}

	// Get the inference service if it exists.
	inferenceService := &v1beta1.InferenceService{}
	if err := r.Get(ctx, client.ObjectKey{Name: model.Name, Namespace: model.Namespace}, inferenceService); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// Create the inference service if it does not exist.
			log.Info("Creating InferenceService", "Model.Namespace", model.Namespace, "Model.Name", model.Name)
			return r.createInferenceService(ctx, model)
		}
		return err
	}

	// Determine if any of the deployment spec fields have changed.
	if inferenceService.Spec.Predictor.PodSpec.Containers[0].Resources.Limits[corev1.ResourceCPU] != model.Spec.Deployment.CPU ||
		inferenceService.Spec.Predictor.PodSpec.Containers[0].Resources.Limits[corev1.ResourceMemory] != model.Spec.Deployment.Memory { // Update the inference service if any of the deployment spec fields have changed.

		log.Info("Updating InferenceService", "Model.Namespace", model.Namespace, "Model.Name", model.Name)

		inferenceService.Spec.Predictor.PodSpec.Containers[0].Resources.Limits[corev1.ResourceCPU] = model.Spec.Deployment.CPU
		inferenceService.Spec.Predictor.PodSpec.Containers[0].Resources.Limits[corev1.ResourceMemory] = model.Spec.Deployment.Memory
		if err := r.Update(ctx, inferenceService); err != nil {
			return err
		}
		// Update the status of the model.
		state := v1alpha1.ModelStateDeploying
		model.Status.State = &state
		// Add condition to the model.
		model.Status.Conditions = append(model.Status.Conditions, metav1.Condition{
			Type:               string(v1alpha1.ModelStateDeploying),
			Status:             metav1.ConditionTrue,
			Reason:             "InferenceServiceUpdated",
			Message:            "Inference service updated successfully",
			LastTransitionTime: metav1.Now(),
		})
		// Update the status of the model.
		if err := r.Status().Update(ctx, &model); err != nil {
			return err
		}
	}

	return nil
}

// ensureInferenceCleanup ensures that the inference service for the model is deleted
// if deployment spec is disabled.
func (r *ModelReconciler) ensureInferenceCleanup(ctx context.Context, model v1alpha1.Model) error {
	log := log.FromContext(ctx)
	// Check if the deployment spec exists and is disabled.
	if model.Spec.Deployment == nil || model.Spec.Deployment.Enabled {
		return nil
	}

	// Get the inference service if it exists.
	inferenceService := &v1beta1.InferenceService{}
	if err := r.Get(ctx, client.ObjectKey{Name: model.Name, Namespace: model.Namespace}, inferenceService); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	// Delete the inference service if it exists.
	log.Info("Deleting InferenceService", "Model.Namespace", model.Namespace, "Model.Name", model.Name)
	if err := r.Delete(ctx, inferenceService); err != nil {
		return err
	}
	// Update the status of the model.
	state := v1alpha1.ModelStateNotDeployed
	model.Status.State = &state
	// Add condition to the model.
	model.Status.Conditions = append(model.Status.Conditions, metav1.Condition{
		Type:               string(v1alpha1.ModelStateNotDeployed),
		Status:             metav1.ConditionTrue,
		Reason:             "InferenceServiceDeleted",
		Message:            "Inference service deleted successfully",
		LastTransitionTime: metav1.Now(),
	})
	// Update the status of the model.
	if err := r.Status().Update(ctx, &model); err != nil {
		return err
	}

	return nil
}

// ensureStatus ensures that the status of the model is updated based on the state of the inference service.
func (r *ModelReconciler) ensureStatus(ctx context.Context, model v1alpha1.Model) error {
	log := log.FromContext(ctx)
	// Get the inference service if it exists. Ignore if not.
	inferenceService := &v1beta1.InferenceService{}
	if err := r.Get(ctx, client.ObjectKey{Name: model.Name, Namespace: model.Namespace}, inferenceService); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	if inferenceService.Status.IsReady() && model.Status.State != nil && *model.Status.State == v1alpha1.ModelStateReady {
		return nil
	}

	// Check if the inference service is ready.
	if inferenceService.Status.IsReady() {
		// Update the status of the model.
		state := v1alpha1.ModelStateReady
		model.Status.State = &state
		// Add condition to the model.
		model.Status.Conditions = append(model.Status.Conditions, metav1.Condition{
			Type:               string(v1alpha1.ModelStateReady),
			Status:             metav1.ConditionTrue,
			Reason:             "InferenceServiceReady",
			Message:            "Inference service is ready",
			LastTransitionTime: metav1.Now(),
		})
		// Update the status of the model.
		log.Info("Model is ready", "Model.Namespace", model.Namespace, "Model.Name", model.Name)
		if err := r.Status().Update(ctx, &model); err != nil {
			return err
		}
	}

	return nil
}

// createInferenceService creates the inference service for the model.
func (r *ModelReconciler) createInferenceService(ctx context.Context, model v1alpha1.Model) error {
	if model.Spec.Type == v1alpha1.ModelTypeHuggingFace {
		// Create the inference service for the Hugging Face model.
		// Create or update the inference service using the controller runtime.
		inferenceService := &v1beta1.InferenceService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      model.Name,
				Namespace: model.Namespace,
			},
			Spec: v1beta1.InferenceServiceSpec{
				Predictor: v1beta1.PredictorSpec{
					PodSpec: v1beta1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "kserve-container",
								Image: r.ModelDeployerImage,
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    model.Spec.Deployment.CPU,
										corev1.ResourceMemory: model.Spec.Deployment.Memory,
									},
								},
								Args: []string{
									fmt.Sprintf("--repo_name=%s", model.Spec.HuggingFace.Name),
									fmt.Sprintf("--org_name=%s", model.Spec.HuggingFace.Organization),
									fmt.Sprintf("--max_sequence_length=%d", model.Spec.HuggingFace.MaxSequenceLength),
								},
							},
						},
					},
				},
			},
		}
		// Set the inference service as the owner of the model.
		if err := controllerutil.SetControllerReference(&model, inferenceService, r.Scheme); err != nil {
			return err
		}
		// Create the inference service.
		if err := r.Create(ctx, inferenceService); err != nil {
			return err
		}
		// Update the status of the model.
		state := v1alpha1.ModelStateDeploying
		model.Status.State = &state
		// Add condition to the model.
		model.Status.Conditions = append(model.Status.Conditions, metav1.Condition{
			Type:               string(v1alpha1.ModelStateDeploying),
			Status:             metav1.ConditionTrue,
			Reason:             "InferenceServiceCreated",
			Message:            "Inference service created successfully",
			LastTransitionTime: metav1.Now(),
		})
		// Update the status of the model.
		if err := r.Status().Update(ctx, &model); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unimplemented model type: %s", model.Spec.Type)
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Model{}).
		Owns(&v1beta1.InferenceService{}).
		Complete(r)
}
