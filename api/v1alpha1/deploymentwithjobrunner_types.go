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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type Job struct {
	PodSpec  v1.PodSpec `json:"podspec"`
	Metadata Metadata   `json:"metadata"`
}

type Deployment struct {
	Selector *metav1.LabelSelector `json:"selector"`
	Replicas *int32                `json:"replicas"`
	PodSpec  v1.PodSpec            `json:"podspec"`
	Metadata Metadata              `json:"metadata"`
}

type Metadata struct {
	Labels    map[string]string `json:"labels,omitempty"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// DeploymentWithJobRunnerSpec defines the desired state of DeploymentWithJobRunner
type DeploymentWithJobRunnerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Deployment Deployment `json:"deployment"`
	Job        Job        `json:"job"`
}

// DeploymentWithJobRunnerStatus defines the observed state of DeploymentWithJobRunner
type DeploymentWithJobRunnerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DeploymentWithJobRunner is the Schema for the deploymentwithjobrunners API
type DeploymentWithJobRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentWithJobRunnerSpec   `json:"spec,omitempty"`
	Status DeploymentWithJobRunnerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeploymentWithJobRunnerList contains a list of DeploymentWithJobRunner
type DeploymentWithJobRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploymentWithJobRunner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploymentWithJobRunner{}, &DeploymentWithJobRunnerList{})
}
