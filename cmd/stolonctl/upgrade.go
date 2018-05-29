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

package main

import (
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/cluster"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"

	"github.com/gravitational/trace"
	"github.com/spf13/cobra"
)

var (
	upgradeCmd = &cobra.Command{
		Use:          "upgrade",
		Short:        "Upgrade stolon application",
		SilenceUsage: true,
		RunE:         upgrade,
	}
)

func init() {
	stolonctlCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().StringVar(&clusterConfig.Postgres.Host, "postgres-host",
		defaults.PostgresHost, "Hostname for connection to stolon PostreSQL host")
	upgradeCmd.Flags().StringVar(&clusterConfig.Postgres.Port, "postgres-port",
		defaults.PostgresPort, "Port for connection to stolon PostgreSQL host")
	upgradeCmd.Flags().StringVar(&clusterConfig.Postgres.User, "postgres-user",
		defaults.PostgresUser, "Username for connection to stolon PostgreSQL host")
	upgradeCmd.Flags().StringVar(&clusterConfig.Postgres.Password, "postgres-password",
		"", "Password for connection to stolon PostgreSQL host")
	upgradeCmd.Flags().StringVar(&clusterConfig.Postgres.BackupPath,
		"postgres-backup-path", defaults.PostgresBackupPath,
		"Path to store backup of stolon PostgreSQL data")
	upgradeCmd.Flags().StringVar(&clusterConfig.Postgres.PgPassPath, "postgres-pgpass-path",
		defaults.PostgresPgPassPath, "Path to store the password file for PostgresQL")
	upgradeCmd.Flags().StringVar(&clusterConfig.Upgrade.NewAppVersion, "app-version",
		"", "Version of application to upgrade to")
	upgradeCmd.Flags().StringVar(&clusterConfig.Upgrade.Changeset, "changeset",
		"", "Changeset of upgrade")
	upgradeCmd.Flags().StringVar(&clusterConfig.Upgrade.NodeName, "nodename", "",
		"Node name where pod is started")
	upgradeCmd.Flags().BoolVar(&clusterConfig.Upgrade.Force, "force",
		false, "Force allows to force phase execution")

	bindFlagEnv(upgradeCmd.Flags())
}

func upgrade(ccmd *cobra.Command, args []string) error {
	if err := clusterConfig.CheckAndSetDefaults(); err != nil {
		return trace.Wrap(err)
	}

	err := cluster.Upgrade(ctx, clusterConfig)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}
