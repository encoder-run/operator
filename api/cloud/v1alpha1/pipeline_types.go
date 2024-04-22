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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineType string

const (
	PipelineTypeRepositoryEmbeddings PipelineType = "REPOSITORY_EMBEDDINGS"
)

type PipelineState string

const (
	PipelineStateNotDeployed PipelineState = "NOT_DEPLOYED"
	PipelineStateDeploying   PipelineState = "DEPLOYING"
	PipelineStateReady       PipelineState = "READY"
	PipelineStateError       PipelineState = "ERROR"
)

type RepositoryEmbeddingsSpec struct {
	// Repository spec
	Repository v1.ObjectReference `json:"repository"`
	// Model spec
	Model v1.ObjectReference `json:"model"`
	// Storage spec
	Storage v1.ObjectReference `json:"storage"`
}

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// Name of the pipeline
	Name string `json:"name"`
	// Type of the pipeline
	Type PipelineType `json:"type"`
	// Enabled flag
	Enabled bool `json:"enabled"`
	// RepositoryEmbeddings pipeline spec
	RepositoryEmbeddings *RepositoryEmbeddingsSpec `json:"repositoryembeddings,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	State *PipelineState `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pipeline is the Schema for the pipelines API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
