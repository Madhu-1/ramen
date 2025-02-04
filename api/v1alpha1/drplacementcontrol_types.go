/*
Copyright 2021 The RamenDR authors.

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
	plrv1 "github.com/open-cluster-management/multicloud-operators-placementrule/pkg/apis/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DRAction which will be either a failover or failback action
// +kubebuilder:validation:Enum=Failover;Failback;Relocate
type DRAction string

// These are the valid values for DRAction
const (
	// Failover, restore PVs to the TargetCluster
	ActionFailover = DRAction("Failover")

	// Failback, restore PVs to the PreferredCluster
	ActionFailback = DRAction("Failback")

	// Relocate, restore PVs to the designated TargetCluster.  PreferredCluster will change
	// to be the TargetCluster.
	ActionRelocate = DRAction("Relocate")
)

// DRPlacementControlSpec defines the desired state of DRPlacementControl
type DRPlacementControlSpec struct {
	// PlacementRef is the reference to the PlacementRule used by DRPC
	PlacementRef v1.ObjectReference `json:"placementRef"`

	// DRPolicyRef is the reference to the DRPolicy participating in the DR replication for this DRPC
	DRPolicyRef v1.ObjectReference `json:"drPolicyRef"`

	// PreferredCluster is the cluster name that the user preferred to run the application on
	PreferredCluster string `json:"preferredCluster,omitempty"`

	// FailoverCluster is the cluster name that the user wants to failover the application to.
	// If not sepcified, then the DRPC will select the surviving cluster from the DRPolicy
	FailoverCluster string `json:"failoverCluster,omitempty"`

	// Label selector to identify all the PVCs that need DR protection.
	// This selector is assumed to be the same for all subscriptions that
	// need DR protection. It will be passed in to the VRG when it is created
	PVCSelector metav1.LabelSelector `json:"pvcSelector"`

	// Action is either failover or failback operation
	Action DRAction `json:"action,omitempty"`
}

// DRState for keeping track of the DR placement
type DRState string

// These are the valid values for DRState
const (
	// Deploying, state recorded in the DRPC status to indicate that the
	// initial deployment is in progress. Deploying means selecting the
	// preffered cluster and creating a VRG MW for it and waiting for MW
	// to be applied in the managed cluster
	Deploying = DRState("Deploying")

	// Deployed, this is the state that will be recorded in the DRPC status
	// when initial deplyment has been performed successfully
	Deployed = DRState("Deployed")

	// FailingOver, state recorded in the DRPC status when the failover
	// is initiated but has not been completed yet
	FailingOver = DRState("FailingOver")

	// FailedOver, state recorded in the DRPC status when the failover
	// process has completed
	FailedOver = DRState("FailedOver")

	// FailingBack, state recorded in the DRPC status when the failback
	// is initiated but has not been completed yet
	FailingBack = DRState("FailingBack")

	// FailedBack, state recorded in the DRPC status when the failback
	// process has completed
	FailedBack = DRState("FailedBack")

	// Relocating, state recorded in the DRPC status to indicate that the
	// relocation is in progress
	Relocating = DRState("Relocating")

	// Relocated, state recorded in
	Relocated = DRState("Relocated")
)

const (
	ConditionAvailable   = "Available"
	ConditionReconciling = "Reconciling"
)

const (
	ReasonProgressing = "Progressing"
	ReasonCleaning    = "Cleaning"
	ReasonSuccess     = "Success"
	ReasonUnknown     = "Unknown"
)

// VRGResourceMeta represents the VRG resource.
type VRGResourceMeta struct {
	// Kind is the kind of the Kubernetes resource.
	// +optional
	Kind string `json:"kind"`

	// Name is the name of the Kubernetes resource.
	Name string `json:"name"`

	// Name is the namespace of the Kubernetes resource.
	Namespace string `json:"namespace"`

	// LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// VRGConditions represents the conditions of the resources deployed on a
// managed cluster.
type VRGConditions struct {
	// ResourceMeta represents the VRG resoure.
	// +required
	ResourceMeta VRGResourceMeta `json:"resourceMeta,omitempty"`

	// Conditions represents the conditions of this resource on a managed cluster.
	// +required
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// DRPlacementControlStatus defines the observed state of DRPlacementControl
type DRPlacementControlStatus struct {
	Phase              DRState                 `json:"phase,omitempty"`
	PreferredDecision  plrv1.PlacementDecision `json:"preferredDecision,omitempty"`
	Conditions         []metav1.Condition      `json:"conditions,omitempty"`
	ResourceConditions VRGConditions           `json:"resourceConditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=drpc

// DRPlacementControl is the Schema for the drplacementcontrols API
type DRPlacementControl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DRPlacementControlSpec   `json:"spec,omitempty"`
	Status DRPlacementControlStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DRPlacementControlList contains a list of DRPlacementControl
type DRPlacementControlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DRPlacementControl `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DRPlacementControl{}, &DRPlacementControlList{})
}
