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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type StolonUpgradeList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []StolonUpgradeResource `json:"items"`
}

func (cr *StolonUpgradeList) GetObjectKind() schema.ObjectKind {
	return &cr.TypeMeta
}

type StolonUpgradeResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StolonUpgradeSpec `json:"spec"`
}

func (cr *StolonUpgradeResource) GetObjectKind() schema.ObjectKind {
	return &cr.TypeMeta
}

func (cr *StolonUpgradeResource) String() string {
	return fmt.Sprintf("namespace=%v, name=%v. status=%v", cr.Namespace, cr.Name, cr.Spec.Status)
}

type StolonUpgradeSpec struct {
	UID               string    `json:"uid"`
	Status            string    `json:"status"`
	Step              string    `json:"step"`
	CreationTimestamp time.Time `json:"start_time"`
	FinishTimestamp   time.Time `json:"finish_time"`
}
