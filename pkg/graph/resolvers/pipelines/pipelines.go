package pipelines

import (
	"context"
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/common"
	"github.com/encoder-run/operator/pkg/graph/converters"
	"github.com/encoder-run/operator/pkg/graph/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Add(ctx context.Context, input model.AddPipelineInput) (*model.Pipeline, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Convert the input to the CRD.
	pipelineCRD, err := converters.PipelineInputToCRD(input)
	if err != nil {
		return nil, err
	}

	// Create the pipeline.
	if err := ctrlClient.Create(ctx, pipelineCRD); err != nil {
		return nil, err
	}

	// Convert the pipeline to the model.
	p, err := converters.PipelineCRDToModel(pipelineCRD)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func Get(ctx context.Context, id string) (*model.Pipeline, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the pipeline.
	pipelineCRD := &v1alpha1.Pipeline{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: id}, pipelineCRD); err != nil {
		return nil, err
	}

	// Convert the pipeline to the model.
	p, err := converters.PipelineCRDToModel(pipelineCRD)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func AddDeployment(ctx context.Context, input model.AddPipelineDeploymentInput) (*model.Pipeline, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the pipeline.
	pipelineCRD := &v1alpha1.Pipeline{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: input.ID}, pipelineCRD); err != nil {
		return nil, err
	}

	// Set the enabled flag.
	pipelineCRD.Spec.Enabled = input.Enabled

	// Update the pipeline.
	if err := ctrlClient.Update(ctx, pipelineCRD); err != nil {
		return nil, err
	}

	// Create the first pipeline execution.
	if input.Enabled {
		executionCRD := &v1alpha1.PipelineExecution{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "pipeline-execution-",
				Namespace:    "default",
				Labels: map[string]string{
					"pipelineId": input.ID,
				},
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(pipelineCRD, v1alpha1.GroupVersion.WithKind("Pipeline")),
				},
			},
			Spec: v1alpha1.PipelineExecutionSpec{
				PipelineRef: v1.ObjectReference{
					Name:      input.ID,
					Namespace: "default",
				},
			},
		}

		if err := ctrlClient.Create(ctx, executionCRD); err != nil {
			return nil, err
		}
	}

	// Convert the pipeline to the model.
	p, err := converters.PipelineCRDToModel(pipelineCRD)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func List(ctx context.Context) ([]*model.Pipeline, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// List all the pipelines.
	pipelineList := &v1alpha1.PipelineList{}
	if err := ctrlClient.List(ctx, pipelineList, &client.ListOptions{Namespace: "default"}); err != nil {
		return nil, err
	}

	// Convert the pipelines to the model.
	pipelines := make([]*model.Pipeline, 0, len(pipelineList.Items))
	for _, pipeline := range pipelineList.Items {
		p, err := converters.PipelineCRDToModel(&pipeline)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, p)
	}

	return pipelines, nil
}

func Delete(ctx context.Context, id string) (*model.Pipeline, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the pipeline.
	pipelineCRD := &v1alpha1.Pipeline{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: id}, pipelineCRD); err != nil {
		return nil, err
	}

	// Delete the pipeline.
	if err := ctrlClient.Delete(ctx, pipelineCRD); err != nil {
		return nil, err
	}

	// Convert the pipeline to the model.
	p, err := converters.PipelineCRDToModel(pipelineCRD)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func Executions(ctx context.Context, id string) ([]*model.PipelineExecution, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"pipelineId": id, // This assumes your label for the pipeline ID is 'pipelineId'
		},
	}

	// Convert LabelSelector to a labels.Selector
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, fmt.Errorf("error creating label selector: %w", err)
	}

	// Prepare ListOptions with the label selector
	listOpts := &client.ListOptions{
		Namespace:     "default",
		LabelSelector: selector,
	}

	// List all the pipeline executions with the label filter
	executionList := &v1alpha1.PipelineExecutionList{}
	if err := ctrlClient.List(ctx, executionList, listOpts); err != nil {
		return nil, fmt.Errorf("error listing pipeline executions: %w", err)
	}

	// Filter the pipeline executions by the pipeline.
	executions := make([]*model.PipelineExecution, 0)
	for _, execution := range executionList.Items {
		if execution.Spec.PipelineRef.Name == id {
			e, err := converters.PipelineExecutionCRDToModel(&execution)
			if err != nil {
				return nil, err
			}
			executions = append(executions, e)
		}
	}

	return executions, nil
}
