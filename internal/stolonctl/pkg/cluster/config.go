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

import "github.com/gravitational/trace"

// Config represents configuration of stolon cluster
type Config struct {
	// KubeConfig defines path to Kubernetes config file
	KubeConfig string
	// Namespace defines Kubernetes namespace for stolon cluster
	Namespace string

	// KeepersPodSelector defines labels to select keeper pods
	KeepersPodSelector string
	// SentinelsPodSelector defines labels to select sentinel pods
	SentinelsPodSelector string

	// EtcdEndpoints defines addresses for connecting to Etcd
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

// CheckConfig checks provided configuration
func (c *Config) CheckConfig() error {
	if c.EtcdCertFile == "" || c.EtcdKeyFile == "" ||
		c.EtcdCAFile == "" || c.EtcdEndpoints == "" {
		return trace.BadParameter("etcd-endpoints, etcd-cert-file, etcd-key-file and etcd-ca-file are required")
	}
	return nil
}
