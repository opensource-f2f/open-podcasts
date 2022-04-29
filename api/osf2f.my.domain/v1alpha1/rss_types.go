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

// RSSSpec defines the desired state of RSS
type RSSSpec struct {
	// Title is the title of RSS
	Title       string   `json:"title,omitempty"`
	Language    string   `json:"language,omitempty"`
	Author      string   `json:"author,omitempty"`
	Address     string   `json:"address,omitempty"`
	Image       string   `json:"image,omitempty"`
	Link        string   `json:"link,omitempty"`
	Description string   `json:"description,omitempty"`
	Categories  []string `json:"categories,omitempty"`
}

// RSSStatus defines the observed state of RSS
type RSSStatus struct {
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// +genclient
// +kubebuilder:printcolumn:name="Title",type=string,JSONPath=`.spec.title`,description="The title of an episode"
//+kubebuilder:printcolumn:name="LastUpdate",type=date,JSONPath=`.status.lastUpdateTime`,description="Last update time"
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RSS is the Schema for the rsses API
type RSS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RSSSpec   `json:"spec,omitempty"`
	Status RSSStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RSSList contains a list of RSS
type RSSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RSS `json:"items"`
}
