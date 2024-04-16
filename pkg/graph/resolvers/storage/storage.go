package storage

import (
	"context"
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/common"
	"github.com/encoder-run/operator/pkg/graph/converters"
	"github.com/encoder-run/operator/pkg/graph/model"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func List(ctx context.Context) ([]*model.Storage, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// List all the storages.
	storageList := &v1alpha1.StorageList{}
	if err := ctrlClient.List(ctx, storageList, &client.ListOptions{Namespace: "default"}); err != nil {
		return nil, err
	}

	// Convert the storages to the model.
	storages := make([]*model.Storage, 0, len(storageList.Items))
	for _, storage := range storageList.Items {
		s, err := converters.StorageCRDToModel(&storage)
		if err != nil {
			return nil, err
		}
		storages = append(storages, s)
	}

	return storages, nil
}

func Add(ctx context.Context, input model.AddStorageInput) (*model.Storage, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Convert the input to the CRD.
	storageCRD, err := converters.StorageInputToCRD(input)
	if err != nil {
		return nil, err
	}

	// Create the storage.
	if err := ctrlClient.Create(ctx, storageCRD); err != nil {
		return nil, err
	}

	// Convert the storage to the model.
	s, err := converters.StorageCRDToModel(storageCRD)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func Delete(ctx context.Context, id string) (*model.Storage, error) {
	// Get the controller-runtime client from the context.
	ctrlClient, ok := ctx.Value(common.AdminClientKey).(client.Client)
	if !ok {
		return nil, fmt.Errorf("controller-runtime client not found in context")
	}

	// Get the storage.
	storageCRD := &v1alpha1.Storage{}
	if err := ctrlClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: id}, storageCRD); err != nil {
		return nil, err
	}

	// Delete the storage.
	if err := ctrlClient.Delete(ctx, storageCRD); err != nil {
		return nil, err
	}

	// Convert the storage to the model.
	s, err := converters.StorageCRDToModel(storageCRD)
	if err != nil {
		return nil, err
	}

	return s, nil
}
