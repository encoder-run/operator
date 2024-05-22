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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudv1alpha1 "github.com/encoder-run/operator/api/cloud/v1alpha1"
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelines/finalizers,verbs=update
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelineexecutions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelineexecutions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelineexecutions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pipeline object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile

func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var pipeline cloudv1alpha1.Pipeline
	if err := r.Get(ctx, req.NamespacedName, &pipeline); err != nil {
		logger.Error(err, "Unable to fetch Pipeline")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var pipelineExecutions cloudv1alpha1.PipelineExecutionList
	if err := r.List(ctx, &pipelineExecutions, client.InNamespace(req.Namespace), client.MatchingFields{"spec.pipelineRef.name": pipeline.Name}); err != nil {
		logger.Error(err, "Unable to list child PipelineExecutions")
		return ctrl.Result{}, err
	}

	// Determine the overall state based on child PipelineExecutions
	pipelineStatus := cloudv1alpha1.PipelineStateReady // default to READY
	for _, exec := range pipelineExecutions.Items {
		if exec.Status.State == nil {
			pipelineStatus = cloudv1alpha1.PipelineStateRunning
		} else {
			switch *exec.Status.State {
			case cloudv1alpha1.PipelineExecutionStatePending:
			case cloudv1alpha1.PipelineExecutionStateActive:
				pipelineStatus = cloudv1alpha1.PipelineStateRunning
			case cloudv1alpha1.PipelineExecutionStateFailed:
				pipelineStatus = cloudv1alpha1.PipelineStateError
			}
		}
	}

	// Update status of the Pipeline
	pipeline.Status.State = &pipelineStatus
	if err := r.Status().Update(ctx, &pipeline); err != nil {
		logger.Error(err, "Failed to update Pipeline status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func indexPipelineRefName(obj client.Object) []string {
	pe, ok := obj.(*cloudv1alpha1.PipelineExecution)
	if !ok {
		// Log error, handle case where object is not a PipelineExecution
		return nil
	}
	return []string{pe.Spec.PipelineRef.Name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create an EventHandler for watching PipelineExecution objects
	ownerHandler := handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &cloudv1alpha1.Pipeline{}, handler.OnlyControllerOwner())

	// Add an index on the spec.pipelineRef.name field for PipelineExecution objects
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &cloudv1alpha1.PipelineExecution{}, "spec.pipelineRef.name", indexPipelineRefName); err != nil {
		return fmt.Errorf("failed to set up pipelineRef name indexer: %w", err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha1.Pipeline{}).
		Watches(&cloudv1alpha1.PipelineExecution{}, ownerHandler).
		Complete(r)
}
