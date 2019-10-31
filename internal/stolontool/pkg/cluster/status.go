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

	"github.com/gravitational/stolon-app/internal/stolontool/pkg/crd"
	"github.com/gravitational/stolon-app/internal/stolontool/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolontool/pkg/kubernetes"

	"github.com/gravitational/rigging"
	"github.com/gravitational/stolon/pkg/cluster"
	"github.com/gravitational/stolon/pkg/store"
	"github.com/gravitational/trace"
	"k8s.io/api/core/v1"
)

const stolonKeeperMaster = "master"

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

func (s *Status) getMasterStatus() (*crd.MasterStatus, error) {
	for _, pod := range s.PodsStatus {
		if master := s.ClusterData.KeepersState[s.ClusterData.ClusterView.Master]; master != nil {
			if pod.PodIP == master.ListenAddress {
				return &crd.MasterStatus{
					PodName: pod.Name,
					Healthy: master.Healthy,
					HostIP:  pod.HostIP,
					PodIP:   pod.PodIP,
				}, nil
			}
		}
	}
	return nil, trace.NotFound("stolon keeper master not found")
}
