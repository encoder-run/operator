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
		m.HuggingFace = &model.HuggingFace{
			Name:              modelCRD.Spec.HuggingFace.Name,
			Organization:      modelCRD.Spec.HuggingFace.Organization,
			MaxSequenceLength: modelCRD.Spec.HuggingFace.MaxSequenceLength,
		}
		m.DisplayName = fmt.Sprintf("%s/%s", modelCRD.Spec.HuggingFace.Organization, modelCRD.Spec.HuggingFace.Name)
	}

	if modelCRD.Spec.Deployment != nil {
		m.Deployment = &model.ModelDeployment{
			Enabled: modelCRD.Spec.Deployment.Enabled,
			CPU:     modelCRD.Spec.Deployment.CPU.String(),
			Memory:  modelCRD.Spec.Deployment.Memory.String(),
		}
	}

	var status model.ModelStatus
	if modelCRD.Status.State != nil {
		switch *modelCRD.Status.State {
		case v1alpha1.ModelStateNotDeployed:
			status = model.ModelStatusNotDeployed
		case v1alpha1.ModelStateDeploying:
			status = model.ModelStatusDeploying
		case v1alpha1.ModelStateReady:
			status = model.ModelStatusReady
		case v1alpha1.ModelStateError:
			status = model.ModelStatusError
		default:
			return nil, fmt.Errorf("unknown model state: %s", *modelCRD.Status.State)
		}
	} else {
		// If the state is not set, then the model is not deployed.
		status = model.ModelStatusNotDeployed
	}
	m.Status = status

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
			Name:              input.HuggingFace.Name,
			Organization:      input.HuggingFace.Organization,
			MaxSequenceLength: input.HuggingFace.MaxSequenceLength,
		}
	case model.ModelTypeExternal:
		modelCRD.Spec.Type = v1alpha1.ModelTypeExternal
	default:
		return nil, fmt.Errorf("unsupported model type: %s", input.Type)
	}

	return modelCRD, nil
}
