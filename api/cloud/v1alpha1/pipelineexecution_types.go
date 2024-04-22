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

type PipelineExecutionState string

const (
	// Active represents an active pipeline execution
	PipelineExecutionStateActive PipelineExecutionState = "ACTIVE"
	// Succeeded represents a succeeded pipeline execution
	PipelineExecutionStateSucceeded PipelineExecutionState = "SUCCEEDED"
	// Failed represents a failed pipeline execution
	PipelineExecutionStateFailed PipelineExecutionState = "FAILED"
	// Pending represents a pending pipeline execution
	PipelineExecutionStatePending PipelineExecutionState = "PENDING"
)

// PipelineExecutionSpec defines the desired state of PipelineExecution
type PipelineExecutionSpec struct {
	// PipelineRef is a reference to the pipeline
	PipelineRef v1.ObjectReference `json:"pipelineRef"`
	// Metadata is a map of metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}

// PipelineExecutionStatus defines the observed state of PipelineExecution
type PipelineExecutionStatus struct {
	State      *PipelineExecutionState `json:"state,omitempty"`
	Conditions []metav1.Condition      `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PipelineExecution is the Schema for the pipelineexecutions API
type PipelineExecution struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineExecutionSpec   `json:"spec,omitempty"`
	Status PipelineExecutionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineExecutionList contains a list of PipelineExecution
type PipelineExecutionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineExecution `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PipelineExecution{}, &PipelineExecutionList{})
}
