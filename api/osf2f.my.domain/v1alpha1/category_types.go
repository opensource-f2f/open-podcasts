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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CategorySpec defines the desired state of Category
type CategorySpec struct {
	// Description is the description of Category
	Description string `json:"foo,omitempty"`
}

// CategoryStatus defines the observed state of Category
type CategoryStatus struct {
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:resource:categories="devops",scope="Cluster"
//+kubebuilder:subresource:status

// Category is the Schema for the categories API
type Category struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CategorySpec   `json:"spec,omitempty"`
	Status CategoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CategoryList contains a list of Category
type CategoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Category `json:"items"`
}
