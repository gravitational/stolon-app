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

package cluster

import (
	"path/filepath"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/kubernetes"

	"github.com/gravitational/rigging"
	"github.com/gravitational/stolon/pkg/cluster"
	"github.com/gravitational/stolon/pkg/store"
	"github.com/gravitational/trace"
	"k8s.io/api/core/v1"
)

// getPods returns list of keeper and sentinel pods
func getPods(config Config) ([]v1.Pod, error) {
	client, err := kubernetes.NewClient(config.KubeConfig)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	pods, err := client.Pods(config.KeepersPodSelector, config.Namespace)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	sentinelPods, err := client.Pods(config.SentinelsPodSelector, config.Namespace)
	if err != nil {
		return nil, rigging.ConvertError(err)
	}

	return append(pods, sentinelPods...), nil
}

// GetStatus returns status of stolon cluster
func GetStatus(config Config) (*Status, error) {
	podList, err := getPods(config)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var podsStatus []kubernetes.PodStatus
	for _, pod := range podList {
		podIP := pod.Status.PodIP
		if podIP == "" {
			podIP = "<none>"
		}
		podState, containers, readyContainers := kubernetes.DeterminePodStatus(pod)
		podStatus := kubernetes.PodStatus{
			Name:              pod.ObjectMeta.Name,
			HostIP:            pod.Spec.NodeName,
			PodIP:             podIP,
			Status:            podState,
			TotalContainers:   containers,
			ReadyContainers:   readyContainers,
			CreationTimestamp: pod.Status.StartTime,
		}
		podsStatus = append(podsStatus, podStatus)
	}

	etcdStore, err := store.NewStore(
		store.Backend("etcd"),
		config.EtcdEndpoints,
		config.EtcdCertFile,
		config.EtcdKeyFile,
		config.EtcdCAFile,
	)
	if err != nil {
		return nil, trace.Wrap(err, "error connecting to etcd")
	}

	storePath := filepath.Join(defaults.EtcdClusterBasePath, config.Name)
	storeManager := store.NewStoreManager(etcdStore, storePath)
	clusterData, _, err := storeManager.GetClusterData()
	if err != nil {
		return nil, trace.Wrap(err, "error getting cluster data from etcd")
	}

	return &Status{podsStatus, clusterData}, nil
}

// Status represents status of stolon cluster
type Status struct {
	// PodStatusList is a list of stolon pod statuses
	PodsStatus []kubernetes.PodStatus
	// State of the cluster received from etcd
	ClusterData *cluster.ClusterData
}

func (s *Status) getMasterStatus() (*MasterStatus, error) {
	for _, keeperState := range s.ClusterData.KeepersState {
		if keeperState.PGState.Role.String() == "master" {
			for _, pod := range s.PodsStatus {
				if pod.PodIP == keeperState.ListenAddress {
					return &MasterStatus{
						PodName: pod.Name,
						Healthy: keeperState.Healthy,
						HostIP:  pod.HostIP,
					}, nil
				}
			}
		}
	}
	return nil, trace.NotFound("stolon keeper master not exists")
}

// MasterStatus stores information about stolon master
type MasterStatus struct {
	// PodName defines name of the stolon-keeper master pod
	PodName string
	// Healthy represents information about health of stolon-keeper master PostgreSQL
	Healthy bool
	// HostIP defines IP address of host to which the stolon-keeper master pod is assigned
	HostIP string
}
