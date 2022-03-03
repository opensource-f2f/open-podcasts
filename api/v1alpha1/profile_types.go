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

// ProfileSpec defines the desired state of Profile
type ProfileSpec struct {
	// DisplayName is the displayName of Profile. Edit profile_types.go to remove/update
	DisplayName   string                    `json:"displayName,omitempty"`
	SocialLinks   map[string]string         `json:"socialLinks,omitempty"`
	LaterPlayList []PlayTodo                `json:"laterPlayList,omitempty"`
	WatchedList   []v1.LocalObjectReference `json:"watchedList,omitempty"`
	Notifier      v1.LocalObjectReference   `json:"notifier,omitempty"`
	Subscription  v1.LocalObjectReference   `json:"subscription,omitempty"`
}

// PlayTodo represents a later play item
type PlayTodo struct {
	v1.LocalObjectReference `json:",inline"`
	// DisplayName which comes from an Episode
	DisplayName string `json:"displayName,omitempty"`
	Index       int    `json:"index,omitempty"`
	Location    int    `json:"location,omitempty"`
}

// ProfileStatus defines the observed state of Profile
type ProfileStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:printcolumn:name="DisplayName",type=string,JSONPath=`.spec.displayName`,description="The displayName of a profile"
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Profile is the Schema for the profiles API
type Profile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileSpec   `json:"spec,omitempty"`
	Status ProfileStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProfileList contains a list of Profile
type ProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Profile `json:"items"`
}
