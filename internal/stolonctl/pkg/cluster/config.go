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

// Config represents configuration of stolon cluster
type Config struct {
	// KubeConfig defines path to Kubernetes config file
	KubeConfig string
	// Namespace defines Kubernetes namespace for Stolon cluster
	Namespace string

	// KeepersPodFilterValue defines label value to filter keeper pods
	KeepersPodFilterValue string
	// KeepersPodFilterKey defines label key to filter keeper pods
	KeepersPodFilterKey string
	// SentinelsPodFilterValue defines label value to filter sentinel pods
	SentinelsPodFilterValue string
	// SentinelsPodFilterKey defines label key to filter sentinel pods
	SentinelsPodFilterKey string

	// EtcdEnpoints defines addresses for connecting to Etcd
	EtcdEndpoints string
	// EtcdCertFile defines path to TLS cert for connecting to etcd
	EtcdCertFile string
	// EtcdKeyFile defines path to TLS key for connecting to etcd
	EtcdKeyFile string
	// EtcdCAFile defines path to TLS CA for connecting to etcd
	EtcdCAFile string

	// Name defines name of stolon cluster
	Name string
}
