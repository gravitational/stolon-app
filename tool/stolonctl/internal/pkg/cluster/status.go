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
	"fmt"
	"path/filepath"

	"github.com/gravitational/stolon-app/tool/stolonctl/internal/pkg/defaults"
	"github.com/gravitational/stolon-app/tool/stolonctl/internal/pkg/kubernetes"

	"github.com/gravitational/stolon/pkg/cluster"
	"github.com/gravitational/stolon/pkg/store"
	"github.com/gravitational/trace"
	"k8s.io/client-go/pkg/api/v1"
)

// Pods returns list of keeper and sentinel pods
func Pods(client *kubernetes.Client, config *Config) ([]v1.Pod, error) {
	pods, err := client.Pods(config.KeepersPodFilterKey, config.KeepersPodFilterValue, config.Namespace)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	sentinelPods, err := client.Pods(config.SentinelsPodFilterKey, config.SentinelsPodFilterValue, config.Namespace)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	pods = append(pods, sentinelPods...)
	return pods, nil
}

// GetStatus returns status of Stolon cluster
func GetStatus(pods []v1.Pod, config *Config) (*Status, error) {
	var podsStatus []kubernetes.PodStatus
	for _, pod := range pods {
		podIP := pod.Status.PodIP
		if podIP == "" {
			podIP = "<none>"
		}
		podState, totalC, readyC := kubernetes.DeterminePodStatus(pod)
		podStatus := kubernetes.PodStatus{
			Name:              pod.ObjectMeta.Name,
			HostIP:            pod.Spec.NodeName,
			PodIP:             podIP,
			Status:            podState,
			TotalContainers:   totalC,
			ReadyContainers:   readyC,
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
		return nil, trace.Wrap(err, "error connecting to Etcd")
	}

	fmt.Printf("%+v\n", etcdStore)

	storePath := filepath.Join(defaults.EtcdClusterBasePath, config.Name)
	storeManager := store.NewStoreManager(etcdStore, storePath)
	clusterData, _, err := storeManager.GetClusterData()
	if err != nil {
		return nil, trace.Wrap(err, "error getting cluster data from Etcd")
	}

	return &Status{podsStatus, clusterData}, nil
}

// Status represents status of Stolon cluster
type Status struct {
	PodsStatus  []kubernetes.PodStatus
	ClusterData *cluster.ClusterData
}
