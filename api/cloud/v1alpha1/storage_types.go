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

// StorageType defines the type of storage
type StorageType string

const (
	// StorageTypeRedis represents a Redis storage
	StorageTypeRedis StorageType = "REDIS"
	// StorageTypePostgres represents a Postgres storage
	StorageTypePostgres StorageType = "POSTGRES"
	// StorageTypeElasticsearch represents an Elasticsearch storage
	StorageTypeElasticsearch StorageType = "ELASTICSEARCH"
)

// StorageSpec defines the desired state of Storage
type StorageSpec struct {
	// Type of storage
	Type StorageType `json:"type"`
	// Name of the storage
	Name string `json:"name"`
}

// StorageStatus defines the observed state of Storage
type StorageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Storage is the Schema for the storages API
type Storage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageSpec   `json:"spec,omitempty"`
	Status StorageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StorageList contains a list of Storage
type StorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Storage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Storage{}, &StorageList{})
}
