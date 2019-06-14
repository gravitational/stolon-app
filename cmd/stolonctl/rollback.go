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

import "github.com/spf13/cobra"

var (
	rollbackCmd = &cobra.Command{
		Use:          "rollback",
		Short:        "Rollback stolon database schema upgrade",
		SilenceUsage: true,
		RunE:         rollback,
	}
)

func init() {
	stolonctlCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().StringVar(&clusterConfig.Upgrade.NewAppVersion, "app-version",
		"", "Version of application to upgrade to")
	rollbackCmd.Flags().StringVar(&clusterConfig.Upgrade.Changeset, "changeset",
		"", "Changeset of upgrade")
	rollbackCmd.Flags().BoolVar(&clusterConfig.Upgrade.Force, "force",
		false, "Force allows to force phase execution")

	bindFlagEnv(rollbackCmd.Flags())
}

func rollback(ccmd *cobra.Command, args []string) error {
	return nil
}
