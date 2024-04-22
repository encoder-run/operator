package converters

import (
	"fmt"

	"github.com/encoder-run/operator/api/cloud/v1alpha1"
	"github.com/encoder-run/operator/pkg/graph/model"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PipelineInputToCRD(input model.AddPipelineInput) (*v1alpha1.Pipeline, error) {
	// Convert the input to the CRD.
	pipelineCRD := &v1alpha1.Pipeline{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: "pipeline-",
			Namespace:    "default",
		},
	}

	switch input.Type {
	case model.PipelineTypeRepositoryEmbeddings:
		pipelineCRD.Spec.Type = v1alpha1.PipelineTypeRepositoryEmbeddings
		pipelineCRD.Spec.RepositoryEmbeddings = &v1alpha1.RepositoryEmbeddingsSpec{
			Repository: corev1.ObjectReference{
				Name:      input.RepositoryEmbeddings.RepositoryID,
				Namespace: "default",
			},
			Storage: corev1.ObjectReference{
				Name:      input.RepositoryEmbeddings.StorageID,
				Namespace: "default",
			},
			Model: corev1.ObjectReference{
				Name:      input.RepositoryEmbeddings.ModelID,
				Namespace: "default",
			},
		}
	default:
		return nil, fmt.Errorf("unsupported model type: %s", input.Type)
	}

	pipelineCRD.Spec.Name = input.Name
	pipelineCRD.Spec.Enabled = false
	return pipelineCRD, nil
}

func PipelineCRDToModel(pipelineCRD *v1alpha1.Pipeline) (*model.Pipeline, error) {
	p := &model.Pipeline{}
	p.ID = pipelineCRD.Name
	p.Name = pipelineCRD.Spec.Name

	var pipelineType model.PipelineType
	switch pipelineCRD.Spec.Type {
	case v1alpha1.PipelineTypeRepositoryEmbeddings:
		pipelineType = model.PipelineTypeRepositoryEmbeddings
	default:
		return nil, fmt.Errorf("unknown pipeline type: %s", pipelineCRD.Spec.Type)
	}

	p.Type = pipelineType
	switch pipelineType {
	case model.PipelineTypeRepositoryEmbeddings:
		p.RepositoryEmbeddings = &model.RepositoryEmbeddings{
			RepositoryID: pipelineCRD.Spec.RepositoryEmbeddings.Repository.Name,
			StorageID:    pipelineCRD.Spec.RepositoryEmbeddings.Storage.Name,
			ModelID:      pipelineCRD.Spec.RepositoryEmbeddings.Model.Name,
		}
	}

	var status model.PipelineStatus
	if pipelineCRD.Status.State != nil {
		switch *pipelineCRD.Status.State {
		case v1alpha1.PipelineStateNotDeployed:
			status = model.PipelineStatusNotDeployed
		case v1alpha1.PipelineStateDeploying:
			status = model.PipelineStatusDeploying
		case v1alpha1.PipelineStateReady:
			status = model.PipelineStatusReady
		case v1alpha1.PipelineStateError:
			status = model.PipelineStatusError
		default:
			return nil, fmt.Errorf("unknown pipeline state: %s", *pipelineCRD.Status.State)
		}
	} else {
		status = model.PipelineStatusNotDeployed
	}
	p.Status = status

	p.Enabled = pipelineCRD.Spec.Enabled
	return p, nil
}

func PipelineExecutionCRDToModel(pipelineExecutionCRD *v1alpha1.PipelineExecution) (*model.PipelineExecution, error) {
	p := &model.PipelineExecution{}
	p.ID = pipelineExecutionCRD.Name
	var status model.PipelineExecutionStatus
	if pipelineExecutionCRD.Status.State != nil {
		switch *pipelineExecutionCRD.Status.State {
		case v1alpha1.PipelineExecutionStateActive:
			status = model.PipelineExecutionStatusActive
		case v1alpha1.PipelineExecutionStateSucceeded:
			status = model.PipelineExecutionStatusSucceeded
		case v1alpha1.PipelineExecutionStateFailed:
			status = model.PipelineExecutionStatusFailed
		default:
			return nil, fmt.Errorf("unknown pipeline execution state: %s", *pipelineExecutionCRD.Status.State)
		}
	} else {
		status = model.PipelineExecutionStatusPending
	}
	p.Status = status
	return p, nil
}
