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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ModelType defines the type of model
type ModelType string

const (
	// ModelTypeHuggingFace represents a Hugging Face model
	ModelTypeHuggingFace ModelType = "HUGGINGFACE"
	// ModelTypeExternal represents an external model
	ModelTypeExternal ModelType = "EXTERNAL"
)

// ModelState defines the state of the model
type ModelState string

const (
	// ModelStateNotDeployed represents a model that is not deployed
	ModelStateNotDeployed ModelState = "NOT_DEPLOYED"
	// ModelStateDeploying represents a model that is being deployed
	ModelStateDeploying ModelState = "DEPLOYING"
	// ModelStateReady represents a model that is ready
	ModelStateReady ModelState = "READY"
	// ModelStateError represents a model that has failed
	ModelStateError ModelState = "ERROR"
)

type ModelDeploymentSpec struct {
	Enabled bool              `json:"enabled"`
	CPU     resource.Quantity `json:"cpu"`
	Memory  resource.Quantity `json:"memory"`
}

// ModelSpec defines the desired state of Model
type ModelSpec struct {
	Type ModelType `json:"type"`
	// Hugging Face model spec
	HuggingFace *HuggingFaceModelSpec `json:"huggingface,omitempty"`
	// Deployment spec
	Deployment *ModelDeploymentSpec `json:"deployment,omitempty"`
}

// HuggingFaceModelSpec defines the desired state of HuggingFaceModel
type HuggingFaceModelSpec struct {
	Organization      string `json:"organization"`
	Name              string `json:"name"`
	MaxSequenceLength int    `json:"maxSequenceLength"`
}

// ModelStatus defines the observed state of Model
type ModelStatus struct {
	State      *ModelState        `json:"state,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Model is the Schema for the models API
type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelSpec   `json:"spec,omitempty"`
	Status ModelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ModelList contains a list of Model
type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Model `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Model{}, &ModelList{})
}
