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

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
)

// PipelineExecutionReconciler reconciles a PipelineExecution object
type PipelineExecutionReconciler struct {
	client.Client
	Scheme                  *runtime.Scheme
	RepositoryEmbedderImage string
}

//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelineexecutions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelineexecutions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelineexecutions/finalizers,verbs=update
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=pipelines,verbs=get;list;watch
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *PipelineExecutionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var pe v1alpha1.PipelineExecution
	if err := r.Get(ctx, req.NamespacedName, &pe); err != nil {
		log.Error(err, "unable to fetch PipelineExecution")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var pipeline v1alpha1.Pipeline
	if err := r.Get(ctx, types.NamespacedName{Name: pe.Spec.PipelineRef.Name, Namespace: pe.Namespace}, &pipeline); err != nil {
		log.Error(err, "unable to fetch Pipeline")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the pipeline is enabled.
	if !pipeline.Spec.Enabled {
		log.Info("pipeline is not enabled")
		// Update the state to pending
		state := v1alpha1.PipelineExecutionStatePending
		pe.Status.State = &state
		if err := r.Status().Update(ctx, &pe); err != nil {
			log.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := r.ensureJob(ctx, &pe, &pipeline); err != nil {
		log.Error(err, "unable to ensure job")
		return ctrl.Result{}, err
	}

	if err := r.ensureStatus(ctx, &pe); err != nil {
		log.Error(err, "unable to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// ensureJob checks if a Job for the PipelineExecution has been created;
// if not, it creates one.
func (r *PipelineExecutionReconciler) ensureJob(ctx context.Context, pe *v1alpha1.PipelineExecution, pipeline *v1alpha1.Pipeline) error {
	job := &batchv1.Job{}
	jobName := pe.Name
	err := r.Get(ctx, types.NamespacedName{Name: jobName, Namespace: pe.Namespace}, job)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		// Define the job
		job = &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: pe.Namespace,
				Labels: map[string]string{
					"pipelineId": pe.Spec.PipelineRef.Name,
				},
			},
			Spec: batchv1.JobSpec{
				BackoffLimit: ptr.Int32(0),
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						ServiceAccountName: "pipeline-worker",
						Containers: []v1.Container{
							{
								Name:    "repoembedder-container",
								Image:   r.RepositoryEmbedderImage,
								Command: []string{"./main"},
								Args: []string{
									fmt.Sprintf("--storageId=%s", pipeline.Spec.RepositoryEmbeddings.Storage.Name),
									fmt.Sprintf("--repositoryId=%s", pipeline.Spec.RepositoryEmbeddings.Repository.Name),
									fmt.Sprintf("--modelId=%s", pipeline.Spec.RepositoryEmbeddings.Model.Name),
								},
							},
						},
						RestartPolicy: v1.RestartPolicyNever,
					},
				},
			},
		}
		// Set PipelineExecution instance as the owner and controller
		controllerutil.SetControllerReference(pe, job, r.Scheme)
		err = r.Create(ctx, job)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PipelineExecutionReconciler) ensureStatus(ctx context.Context, pe *v1alpha1.PipelineExecution) error {
	// Get the job if it exists
	job := &batchv1.Job{}
	jobName := pe.Name
	err := r.Get(ctx, types.NamespacedName{Name: jobName, Namespace: pe.Namespace}, job)
	if err != nil {
		return client.IgnoreNotFound(err)
	}

	// Update the status based on the job
	var state v1alpha1.PipelineExecutionState
	if job.Status.Succeeded > 0 {
		state = v1alpha1.PipelineExecutionStateSucceeded
	} else if job.Status.Failed > 0 {
		state = v1alpha1.PipelineExecutionStateFailed
	} else if job.Status.Active > 0 {
		state = v1alpha1.PipelineExecutionStateActive
	}
	pe.Status.State = &state

	// Update the PipelineExecution status
	if err := r.Status().Update(ctx, pe); err != nil {
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineExecutionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PipelineExecution{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
