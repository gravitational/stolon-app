/*
Copyright (C) 2018 Gravitational, Inc.

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

package crd

import (
	"fmt"
	"time"

	"github.com/gravitational/stolon-app/internal/stolontool/pkg/kubernetes"
	"github.com/gravitational/stolon/pkg/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// StolonUpgradeList is a list of StolonUpgradeResource objects
// with additional information related to kubernetes Custom Resource Definition
type StolonUpgradeList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []StolonUpgradeResource `json:"items"`
}

func (cr *StolonUpgradeList) GetObjectKind() schema.ObjectKind {
	return &cr.TypeMeta
}

// StolonUpgradeResource is the definition of kubernetes custom
// resource for upgrade stolon application
type StolonUpgradeResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StolonUpgradeSpec `json:"spec"`
}

func (cr *StolonUpgradeResource) GetObjectKind() schema.ObjectKind {
	return &cr.TypeMeta
}

func (cr *StolonUpgradeResource) String() string {
	return fmt.Sprintf("StolonUpgradeResource(namespace=%v, name=%v, status=%v)",
		cr.Namespace, cr.Name, cr.Spec.Status)
}

// StolonUpgradeSpec is a specification of Custom Resource Definition
// for upgrade stolon application
type StolonUpgradeSpec struct {
	// Status is a status of stolon upgrade
	Status string `json:"status"`
	// ClusterInfo represents information about stolon cluster
	ClusterInfo ClusterInfo `json:"clusterInfo"`
	// Phases is a list of phases to upgrade stolon
	Phases []StolonUpgradePhase `json:"phases"`
	// CreationTimestamp is a starting time of upgrade
	CreationTimestamp time.Time `json:"startTime"`
	// FinishTimestamp is a time when upgrade finished
	FinishTimestamp time.Time `json:"finishTime"`
}

// StolonUpgradePhase defines phase of upgrade
type StolonUpgradePhase struct {
	// Status is a status of upgrade step(phase)
	Status string `json:"status"`
	// Name is a name of upgrade step
	Name string `json:"name"`
	// NodeName represents node where phase is executed
	NodeName string `json:"nodeName"`
	// Description is a small description of upgrade step
	Description string `json:"description"`
	// CreationTimestamp is a starting time of upgrade step
	CreationTimestamp time.Time `json:"startTime"`
	// UpdatedTimestamp is a time when upgrade step updated
	UpdatedTimestamp time.Time `json:"updatedTime"`
}

// MasterStatus stores information about stolon master
type MasterStatus struct {
	// PodName defines name of the stolon-keeper master pod
	PodName string
	// Healthy indicates whether the postgres master is healthy
	Healthy bool
	// HostIP identifies the node where master pod has been scheduled
	HostIP string
	// PodIP identifies ip address of stolon master pod
	PodIP string
}

// ClusterInfo represents information about stolon cluster
type ClusterInfo struct {
	// ClusterData is an information about the cluster received from etcd
	ClusterData cluster.ClusterData `json:"clusterData"`
	// PodsStatus represents information about stolon pods
	PodsStatus []kubernetes.PodStatus `json:"podStatus"`
	// MasterStatus stores information about stolon master
	MasterStatus MasterStatus `json:"masterStatus"`
}
