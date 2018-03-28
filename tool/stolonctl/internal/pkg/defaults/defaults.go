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

package defaults

const (
	// ConfigPath defines path to configuration file
	ConfigPath = "/etc/pgupgrade/config.yaml"
	// Namespace defines Kubernetes namespace for Stolon application
	Namespace = "default"
	// KeepersPodFilterValue defines label value to filter keeper pods
	KeepersPodFilterValue = "stolon-keeper"
	// KeepersPodFilterKey defines label key to filter keeper pods
	KeepersPodFilterKey = "name"
	// SentinelsPodFilterValue defines label value to filter sentinel pods
	SentinelsPodFilterValue = "stolon-sentinel"
	// SentinelsPodFilterKey defines label key to filter sentinel pods
	SentinelsPodFilterKey = "name"

	// EtcdEndpoints defines endpoints for connection to Etcd server
	EtcdEndpoints = "127.0.0.1:2379"
	// EtcdClusterBasePath defines base path for cluster data in Etcd
	EtcdClusterBasePath = "stolon/cluster"

	// ClusterName defines name of Stolon cluster
	ClusterName = "kube-stolon"
)
