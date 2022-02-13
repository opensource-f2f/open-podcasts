/*
Copyright 2022 The osf2f Authors. All rights reserved.

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

// EpisodeSpec defines the desired state of Episode
type EpisodeSpec struct {
	Title       string `json:"title,omitempty"`
	CoverImage  string `json:"coverImage,omitempty"`
	AudioSource string `json:"audioSource,omitempty"`
	// Link is the link of Episode. Edit episode_types.go to remove/update
	Link string `json:"link,omitempty"`
}

// EpisodeStatus defines the observed state of Episode
type EpisodeStatus struct {
	Hints int64 `json:"hints,omitempty"`
}

// +kubebuilder:printcolumn:name="Title",type=string,JSONPath=`.spec.title`,description="The title of an episode"
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Episode is the Schema for the episodes API
type Episode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EpisodeSpec   `json:"spec,omitempty"`
	Status EpisodeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EpisodeList contains a list of Episode
type EpisodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Episode `json:"items"`
}
