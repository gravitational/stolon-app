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
	// Namespace defines Kubernetes namespace for stolon application
	Namespace = "default"
	// KeepersPodSelector defines label selector to select keeper pods
	KeepersPodSelector = "name=stolon-keeper"
	// SentinelsPodSelector defines label selector to select sentinel pods
	SentinelsPodSelector = "name=stolon-sentinel"

	// EtcdEndpoints defines endpoints for connecting to etcd
	EtcdEndpoints = "127.0.0.1:2379"
	// EtcdClusterBasePath defines the path prefix for stolon-specific cluster data in etcd
	EtcdClusterBasePath = "stolon/cluster"

	// ClusterName defines name of stolon cluster
	ClusterName = "kube-stolon"
)
