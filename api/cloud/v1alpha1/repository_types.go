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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RepositoryType defines the type of repository
type RepositoryType string

const (
	// RepositoryTypeGithub represents a Github repository
	RepositoryTypeGithub RepositoryType = "GITHUB"
	// RepositoryTypeGitlab represents a Gitlab repository
	RepositoryTypeGitlab RepositoryType = "GITLAB"
	// RepositoryTypeBitbucket represents a Bitbucket repository
	RepositoryTypeBitbucket RepositoryType = "BITBUCKET"
)

// RepositorySpec defines the desired state of Repository
type RepositorySpec struct {
	// Type of repository
	Type RepositoryType `json:"type"`
	// Github repository spec
	Github *GithubRepositorySpec `json:"github,omitempty"`
}

// GithubRepositorySpec defines the desired state of a Github repository
type GithubRepositorySpec struct {
	// Owner of the repository
	Owner string `json:"owner"`
	// Name of the repository
	Name string `json:"name"`
	// URL of the repository
	URL string `json:"url"`
	// Branch of the repository
	Branch string `json:"branch,omitempty"`
}

// RepositoryStatus defines the observed state of Repository
type RepositoryStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Repository is the Schema for the repositories API
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec,omitempty"`
	Status RepositoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RepositoryList contains a list of Repository
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}
