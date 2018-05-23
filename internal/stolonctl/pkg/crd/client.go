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
	"encoding/json"
	"time"

	"github.com/gravitational/rigging"
	"github.com/gravitational/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// NewClient creates CRD client interface
func NewClient(cfg *rest.Config, namespace string) (*Client, error) {
	cfg.APIPath = "/apis"
	if cfg.UserAgent == "" {
		cfg.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	cfg.ContentType = runtime.ContentTypeJSON
	cfg.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	cfg.GroupVersion = &schema.GroupVersion{Group: StolonUpgradeGroup, Version: StolonUpgradeVersion}

	clt, err := rest.RESTClientFor(cfg)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	return &Client{
		Client:    clt,
		namespace: namespace,
		plural:    StolonUpgradePlural,
	}, nil
}

// Client defines client to kubernetes custom resource definition API
type Client struct {
	Client    *rest.RESTClient
	namespace string
	plural    string
}

func (c *Client) CreateOrRead(objMeta metav1.ObjectMeta) (*StolonUpgradeResource, error) {
	res := &StolonUpgradeResource{
		TypeMeta: metav1.TypeMeta{
			Kind:       StolonUpgradeKind,
			APIVersion: StolonUpgradeAPIVersion,
		},
		ObjectMeta: objMeta,
		Spec: StolonUpgradeSpec{
			Status:            StolonUpgradeStatusInProgress,
			Phases:            stolonUpgradePhases(time.Now().UTC()),
			CreationTimestamp: time.Now().UTC(),
		},
	}
	out, err := c.create(res)
	if err == nil {
		return out, nil
	}
	if !trace.IsAlreadyExists(err) {
		return nil, trace.Wrap(err)
	}
	return c.get(objMeta.Name)
}

func (c *Client) create(obj *StolonUpgradeResource) (*StolonUpgradeResource, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var raw runtime.Unknown
	err = c.Client.Post().
		Namespace(c.namespace).
		Resource(c.plural).
		Body(data).
		Do().
		Into(&raw)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	var result StolonUpgradeResource
	if err := json.Unmarshal(raw.Raw, &result); err != nil {
		return nil, trace.Wrap(err)
	}
	return &result, nil
}

func (c *Client) List() (*StolonUpgradeList, error) {
	var raw runtime.Unknown
	err := c.Client.Get().
		Namespace(c.namespace).
		Resource(c.plural).
		Do().
		Into(&raw)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	var result StolonUpgradeList
	if err := json.Unmarshal(raw.Raw, &result); err != nil {
		return nil, trace.Wrap(err)
	}
	return &result, nil
}

func (c *Client) get(name string) (*StolonUpgradeResource, error) {
	var raw runtime.Unknown
	err := c.Client.Get().
		Namespace(c.namespace).
		Resource(c.plural).
		Name(name).
		Do().
		Into(&raw)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	var result StolonUpgradeResource
	if err := json.Unmarshal(raw.Raw, &result); err != nil {
		return nil, trace.Wrap(err)
	}
	return &result, nil
}

func (c *Client) update(obj *StolonUpgradeResource) (*StolonUpgradeResource, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var raw runtime.Unknown
	err = c.Client.Put().
		Namespace(c.namespace).
		Resource(c.plural).
		Name(obj.Name).
		Body(data).
		Do().
		Into(&raw)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	var result StolonUpgradeResource
	if err := json.Unmarshal(raw.Raw, &result); err != nil {
		return nil, trace.Wrap(err)
	}
	return &result, nil
}

// MarkPhase marks phase
func (c *Client) MarkPhase(obj *StolonUpgradeResource, phaseName string, phaseStatus string) (*StolonUpgradeResource, error) {
	var phases []StolonUpgradePhase
	for _, phase := range obj.Spec.Phases {
		if phaseName == phase.Name {
			if phaseStatus != phase.Status {
				phase.Status = phaseStatus
				if phaseStatus == StolonUpgradeStatusInProgress {
					phase.CreationTimestamp = time.Now().UTC()
				}
				if phaseStatus == StolonUpgradeStatusCompleted {
					phase.FinishTimestamp = time.Now().UTC()
				}
			}
		}
		phases = append(phases, phase)
	}
	obj.Spec.Phases = phases
	return c.update(obj)
}

// UpdateClusterInfo updates information about stolon cluster
func (c *Client) UpdateClusterInfo(obj *StolonUpgradeResource, clusterInfo ClusterInfo) (*StolonUpgradeResource, error) {
	obj.Spec.ClusterInfo = clusterInfo
	return c.update(obj)
}

// IsPhaseCompleted checks phase for a completed status
func (c *Client) IsPhaseCompleted(obj *StolonUpgradeResource, phaseName string) bool {
	for _, phase := range obj.Spec.Phases {
		if phaseName == phase.Name && phase.Status == StolonUpgradeStatusCompleted {
			return true
		}
	}
	return false
}

// CompleteUpgrade marks upgrade as completed
func (c *Client) CompleteUpgrade(obj *StolonUpgradeResource) (*StolonUpgradeResource, error) {
	obj.Spec.Status = StolonUpgradeStatusCompleted
	obj.Spec.FinishTimestamp = time.Now().UTC()
	return c.update(obj)
}

func stolonUpgradePhases(creationTime time.Time) []StolonUpgradePhase {
	return []StolonUpgradePhase{
		StolonUpgradePhase{
			Status:            StolonUpgradeStatusInProgress,
			Name:              StolonUpgradePhaseInit,
			Description:       "Initialize upgrade operation",
			CreationTimestamp: creationTime,
		},
		StolonUpgradePhase{
			Status:      StolonUpgradeStatusUnstarted,
			Name:        StolonUpgradePhaseChecks,
			Description: "Checks status of cluster before starting upgrade",
		},
		StolonUpgradePhase{
			Status:      StolonUpgradeStatusUnstarted,
			Name:        StolonUpgradePhaseBackupPostgres,
			Description: "Backup stolon PostgreSQL",
		},
		StolonUpgradePhase{
			Status:      StolonUpgradeStatusUnstarted,
			Name:        StolonUpgradePhaseDeleteDeployment,
			Description: "Delete stolon-sentinel deployment",
		},
		StolonUpgradePhase{
			Status:      StolonUpgradeStatusUnstarted,
			Name:        StolonUpgradePhaseDeleteDaemonset,
			Description: "Delete stolon-keeper daemonset",
		},
		StolonUpgradePhase{
			Status:      StolonUpgradeStatusUnstarted,
			Name:        StolonUpgradePhaseUpgradePostgresSchema,
			Description: "Upgrade PostgreSQL schema on master pod",
		},
	}
}
