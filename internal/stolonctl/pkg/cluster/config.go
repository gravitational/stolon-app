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
	"strings"

	"github.com/gravitational/trace"
)

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

	// Postgres stores configuration of PostgreSQL-related parameters
	Postgres PostgresConfig
	// Upgrade stores configuration for upgrading the cluster
	Upgrade UpgradeConfig
}

// PostgresConfig stores configuration of PostgreSQL-related parameters
type PostgresConfig struct {
	// Host defines hostname for connecting to stolon PostgreSQL
	Host string
	// Port defines port for connecting to stolon PostgreSQL
	Port string
	// User defines username for connecting to stolon PostgreSQL
	User string
	// Password defines password for connecting to stolon PostgreSQL
	Password string
	// BackupPath defines path for storing backup of PostgreSQL data
	BackupPath string
	// PgPassPath defines path to the password file
	PgPassPath string
}

// UpgradeConfig defines configuration used for upgrade
type UpgradeConfig struct {
	// NewAppVersion defines new version of application
	NewAppVersion string
	// Changeset defines the name of the changeset for upgrade as used by rig
	Changeset string
	// NodeName defines the name of kubernetes node where stolonctl pod is started
	NodeName string
	// Force allows to force phase execution
	Force bool
}

// CheckAndSetDefaults validates this configuration object and sets defaults
func (c *Config) CheckAndSetDefaults() error {
	var errors []error
	if c.EtcdCertFile == "" {
		errors = append(errors, trace.BadParameter("etcd-cert-file (env 'ETCD_CERT') is required"))
	}
	if c.EtcdKeyFile == "" {
		errors = append(errors, trace.BadParameter("etcd-key-file (env 'ETCD_KEY') is required"))
	}
	if c.EtcdCAFile == "" {
		errors = append(errors, trace.BadParameter("etcd-ca-file (env 'ETCD_CACERT') is required"))
	}
	if c.EtcdEndpoints == "" {
		errors = append(errors, trace.BadParameter("etcd-endpoints (env 'ETCD_ENDPOINTS') is required"))
	}
	if err := c.Postgres.Check(); err != nil {
		errors = append(errors, err)
	}
	if err := c.Upgrade.CheckAndSetDefaults(); err != nil {
		errors = append(errors, err)
	}
	return trace.NewAggregate(errors...)
}

// CheckAndSetDefaults validates this configuration object and sets defaults
func (c *UpgradeConfig) CheckAndSetDefaults() error {
	if c == nil {
		return nil
	}
	if c.NewAppVersion == "" {
		return trace.BadParameter("app-version (env 'APP_VERSION') is required for upgrade")
	}
	if c.NodeName == "" {
		return trace.BadParameter("nodename (env 'NODE_NAME') is required for upgrade")
	}
	if c.Changeset == "" {
		c.Changeset = strings.Replace(c.NewAppVersion, ".", "", -1)
	}
	return nil
}

// Check checks provided configuration for PostgreSQL parameters
func (c *PostgresConfig) Check() error {
	if c == nil {
		return nil
	}
	var errors []error
	if c.Host == "" {
		errors = append(errors, trace.BadParameter("postgres-host is required"))
	}
	if c.Port == "" {
		errors = append(errors, trace.BadParameter("postgres-port is required"))
	}
	if c.User == "" {
		errors = append(errors, trace.BadParameter("postgres-user is required"))
	}
	if c.BackupPath == "" {
		errors = append(errors, trace.BadParameter("postgres-backup-path is required"))
	}
	if c.PgPassPath == "" {
		errors = append(errors, trace.BadParameter("postgres-pgpass-path is required"))
	}
	return trace.NewAggregate(errors...)
}
