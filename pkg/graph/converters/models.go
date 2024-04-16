package converters

import (
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/graph/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ModelCRDToModel(modelCRD *v1alpha1.Model) (*model.Model, error) {
	m := &model.Model{}
	m.ID = modelCRD.Name

	var modelType model.ModelType
	switch modelCRD.Spec.Type {
	case v1alpha1.ModelTypeHuggingFace:
		modelType = model.ModelTypeHuggingface
	case v1alpha1.ModelTypeExternal:
		modelType = model.ModelTypeExternal
	default:
		return nil, fmt.Errorf("unknown model type: %s", modelCRD.Spec.Type)
	}

	m.Type = modelType
	if modelCRD.Spec.Type == v1alpha1.ModelTypeHuggingFace {
		if m.HuggingFace == nil {
			m.HuggingFace = &model.HuggingFace{}
		}
		m.HuggingFace.Name = modelCRD.Spec.HuggingFace.Name
		m.HuggingFace.Organization = modelCRD.Spec.HuggingFace.Organization
		m.DisplayName = fmt.Sprintf("%s/%s", modelCRD.Spec.HuggingFace.Organization, modelCRD.Spec.HuggingFace.Name)
	}

	m.Status = model.ModelStatusNotDeployed

	return m, nil
}

func ModelInputToCRD(input model.AddModelInput) (*v1alpha1.Model, error) {
	modelCRD := &v1alpha1.Model{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: "model-",
			Namespace:    "default",
		},
	}

	switch input.Type {
	case model.ModelTypeHuggingface:
		modelCRD.Spec.Type = v1alpha1.ModelTypeHuggingFace
		modelCRD.Spec.HuggingFace = &v1alpha1.HuggingFaceModelSpec{
			Name:         input.HuggingFace.Name,
			Organization: input.HuggingFace.Organization,
		}
	case model.ModelTypeExternal:
		modelCRD.Spec.Type = v1alpha1.ModelTypeExternal
	default:
		return nil, fmt.Errorf("unsupported model type: %s", input.Type)
	}

	return modelCRD, nil
}
