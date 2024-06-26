package converters

import (
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/graph/model"
	corev1 "k8s.io/api/core/v1"
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

	if storageCRD.Spec.Deployment != nil {
		storage.Deployment = &model.StorageDeployment{
			Enabled: storageCRD.Spec.Deployment.Enabled,
			CPU:     storageCRD.Spec.Deployment.CPU.String(),
			Memory:  storageCRD.Spec.Deployment.Memory.String(),
		}
	}

	var status model.StorageStatus
	if storageCRD.Status.State != nil {
		switch *storageCRD.Status.State {
		case v1alpha1.StorageStateNotDeployed:
			status = model.StorageStatusNotDeployed
		case v1alpha1.StorageStateDeploying:
			status = model.StorageStatusDeploying
		case v1alpha1.StorageStateReady:
			status = model.StorageStatusReady
		case v1alpha1.StorageStateError:
			status = model.StorageStatusError
		default:
			return nil, fmt.Errorf("unknown storage state: %s", *storageCRD.Status.State)
		}
	} else {
		// If the state is not set, then the storage is not deployed.
		status = model.StorageStatusNotDeployed
	}

	storage.Status = status
	storage.Type = storageType
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
	storageCRD.Spec.Name = input.Name

	if storageType == v1alpha1.StorageTypePostgres {
		if input.Postgres == nil {
			return nil, fmt.Errorf("postgres spec is required for storage type: %s", input.Type)
		}
		storageCRD.Spec.Postgres = &v1alpha1.PostgresSpec{
			External: input.Postgres.External,
		}
	}

	return storageCRD, nil
}

func PostgresSecretInputToCRD(storageCRD *v1alpha1.Storage, input *model.PostgresInput) (*corev1.Secret, error) {
	secretCRD := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      storageCRD.Name,
			Namespace: storageCRD.Namespace,
		},
	}
	secretCRD.StringData = map[string]string{
		"host":     input.Host,
		"port":     fmt.Sprintf("%d", input.Port),
		"database": input.Database,
		"username": input.Username,
		"password": input.Password,
		"ssl_mode": input.SSLMode,
		"timezone": input.Timezone,
	}
	return secretCRD, nil
}
