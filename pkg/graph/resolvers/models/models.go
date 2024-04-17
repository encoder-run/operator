package models

import (
	"context"
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/common"
	"github.com/encoder-run/operator/pkg/graph/converters"
	"github.com/encoder-run/operator/pkg/graph/model"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func List(ctx context.Context) ([]*model.Model, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// List all the models.
	modelList := &v1alpha1.ModelList{}
	if err := ctrlClient.List(ctx, modelList, &client.ListOptions{Namespace: "default"}); err != nil {
		return nil, err
	}

	// Convert the models to the model.
	models := make([]*model.Model, 0, len(modelList.Items))
	for _, model := range modelList.Items {
		m, err := converters.ModelCRDToModel(&model)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}

	return models, nil
}

func Get(ctx context.Context, id string) (*model.Model, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the model.
	modelCRD := &v1alpha1.Model{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: id}, modelCRD); err != nil {
		return nil, err
	}

	// Convert the model to the model.
	m, err := converters.ModelCRDToModel(modelCRD)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func Add(ctx context.Context, input model.AddModelInput) (*model.Model, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Convert the input to the CRD.
	modelCRD, err := converters.ModelInputToCRD(input)
	if err != nil {
		return nil, err
	}

	// Create the model.
	if err := ctrlClient.Create(ctx, modelCRD); err != nil {
		return nil, err
	}

	// Convert the model to the model.
	m, err := converters.ModelCRDToModel(modelCRD)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func AddDeployment(ctx context.Context, input model.AddModelDeploymentInput) (*model.Model, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	cpu, err := resource.ParseQuantity(input.CPU)
	if err != nil {
		return nil, err
	}

	memory, err := resource.ParseQuantity(input.Memory)
	if err != nil {
		return nil, err
	}

	// Get the model.
	modelCRD := &v1alpha1.Model{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: input.ID}, modelCRD); err != nil {
		return nil, err
	}

	// Update the model with the deployment.
	modelCRD.Spec.Deployment = &v1alpha1.ModelDeploymentSpec{
		Enabled: true,
		CPU:     cpu,
		Memory:  memory,
	}

	// Update the model.
	if err := ctrlClient.Update(ctx, modelCRD); err != nil {
		return nil, err
	}

	// Convert the model to the model.
	m, err := converters.ModelCRDToModel(modelCRD)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func Delete(ctx context.Context, id string) (*model.Model, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the model.
	modelCRD := &v1alpha1.Model{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: id}, modelCRD); err != nil {
		return nil, err
	}

	// Delete the model.
	if err := ctrlClient.Delete(ctx, modelCRD); err != nil {
		return nil, err
	}

	// Convert the model to the model.
	m, err := converters.ModelCRDToModel(modelCRD)
	if err != nil {
		return nil, err
	}

	return m, nil
}
