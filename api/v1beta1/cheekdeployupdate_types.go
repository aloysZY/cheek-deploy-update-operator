/*
Copyright 2022 aloys.

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

package v1beta1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// type CheekDeployUpdateSpecDeploy struct {
// 	DeploymentName      string `json:"deploymentName"`
// 	DeploymentNamespace string `json:"deploymentNamespace"`
// 	DeploymentImage     string `json:"deploymentImage"`
// }

// CheekDeployUpdateSpec defines the desired state of CheekDeployUpdate
type CheekDeployUpdateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of CheekDeployUpdate. Edit cheekdeployupdate_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	// Deployment []CheekDeployUpdateSpecDeploy `json:"deployment,omitempty"`

	DeploymentName      string `json:"deploymentName"`
	DeploymentNamespace string `json:"deploymentNamespace"`
	DeploymentImage     string `json:"deploymentImage"`
}

// CheekDeployUpdateStatus defines the observed state of CheekDeployUpdate
type CheekDeployUpdateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 这个更新后测速的时候要重新部署CRD
	// DeploymentReplicas *int32                  `json:"deploymentReplicas"`
	// 这里就为了省事了，直接用变更的deployment信息
	CDUStatus appsv1.DeploymentStatus `json:"cduStatus"`
}

/*// +kubebuilder:printcolumn:JSONPath=".status.deploymentReplicas",name=Replicas,type=integer*/

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories="all",shortName="cdu",scope="Cluster"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.deploymentImage",description="image version"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:JSONPath=".status.cduStatus",name=status,type=string

// CheekDeployUpdate is the Schema for the cheekdeployupdates API
type CheekDeployUpdate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CheekDeployUpdateSpec   `json:"spec,omitempty"`
	Status CheekDeployUpdateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CheekDeployUpdateList contains a list of CheekDeployUpdate
type CheekDeployUpdateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CheekDeployUpdate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CheekDeployUpdate{}, &CheekDeployUpdateList{})
}
