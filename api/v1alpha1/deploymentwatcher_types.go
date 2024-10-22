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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeploymentWatcherSpec defines the desired state of DeploymentWatcher
type DeploymentWatcherSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Deployments []NamespacedName `json:"deployments"`
	Slack       SlackConfig      `json:"slack"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type SlackConfig struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

// DeploymentWatcherStatus defines the observed state of DeploymentWatcher
type DeploymentWatcherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DeploymentWatcher is the Schema for the deploymentwatchers API
type DeploymentWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentWatcherSpec   `json:"spec,omitempty"`
	Status DeploymentWatcherStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeploymentWatcherList contains a list of DeploymentWatcher
type DeploymentWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploymentWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploymentWatcher{}, &DeploymentWatcherList{})
}
