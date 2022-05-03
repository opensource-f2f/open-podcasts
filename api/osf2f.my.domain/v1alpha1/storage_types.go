/*
Copyright 2022 The open-podcasts Authors. All rights reserved.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StorageSpec defines the desired state of Storage
type StorageSpec struct {
	Images              []Image              `json:"images,omitempty"`
	GitProviderReleases []GitProviderRelease `json:"gitProviderReleases,omitempty"`
}

type Image struct {
	Registry *v1.SecretReference `json:"registry"`
	Name     string              `json:"name,omitempty"`
	// Repo represents a full path of an image, for example: library/name
	Repo     string `json:"repo"`
	Filepath string `json:"filepath,omitempty"`
}

type GitProviderRelease struct {
	Name     string              `json:"name,omitempty"`
	Provider string              `json:"provider,omitempty"`
	Secret   *v1.SecretReference `json:"secret"`
	// Server is necessary for those self-hosted git providers
	Server string `json:"server,omitempty"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
}

// StorageStatus defines the observed state of Storage
type StorageStatus struct {
	State      string            `json:"state,omitempty"`
	Conditions map[string]string `json:"conditions,omitempty"`
}

// +genclient
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

const AnnotationKeyRSS = "rss"
