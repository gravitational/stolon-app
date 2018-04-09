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
	"github.com/gravitational/trace"

	"github.com/spf13/cobra"
)

var (
	upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade stolon application",
		RunE:  upgrade,
	}
)

func init() {
	stolonctlCmd.AddCommand(upgradeCmd)
}

func upgrade(ccmd *cobra.Command, args []string) error {
	if err := clusterConfig.CheckConfig(); err != nil {
		return trace.Wrap(err)
	}

	err := cluster.Upgrade(ctx, &clusterConfig)
	if err != nil {
		return trace.Wrap(err, "error upgrading cluster")
	}
	return nil
}
