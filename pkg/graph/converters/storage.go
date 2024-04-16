package converters

import (
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/graph/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StorageCRDToModel(storageCRD *v1alpha1.Storage) (*model.Storage, error) {
	storage := &model.Storage{}
	storage.ID = storageCRD.Name
	var storageType model.StorageType
	switch storageCRD.Spec.Type {
	case v1alpha1.StorageTypeRedis:
		storageType = model.StorageTypeRedis
	case v1alpha1.StorageTypePostgres:
		storageType = model.StorageTypePostgres
	case v1alpha1.StorageTypeElasticsearch:
		storageType = model.StorageTypeElasticsearch
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageCRD.Spec.Type)
	}

	storage.Type = storageType
	storage.Status = model.StorageStatusNotDeployed
	storage.Name = storageCRD.Name
	return storage, nil
}

func StorageInputToCRD(input model.AddStorageInput) (*v1alpha1.Storage, error) {
	storageCRD := &v1alpha1.Storage{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: "storage-",
			Namespace:    "default",
		},
	}
	var storageType v1alpha1.StorageType
	switch input.Type {
	case model.StorageTypeRedis:
		storageType = v1alpha1.StorageTypeRedis
	case model.StorageTypePostgres:
		storageType = v1alpha1.StorageTypePostgres
	case model.StorageTypeElasticsearch:
		storageType = v1alpha1.StorageTypeElasticsearch
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", input.Type)
	}

	storageCRD.Spec.Type = storageType
	storageCRD.Name = input.Name

	return storageCRD, nil
}
