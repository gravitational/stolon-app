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
	"fmt"

	"github.com/gravitational/stolon-app/tool/stolonctl/internal/pkg/cluster"
	"github.com/gravitational/stolon-app/tool/stolonctl/internal/pkg/kubernetes"
	"github.com/gravitational/trace"
	"github.com/spf13/cobra"
)

var (
	upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade PostgreSQL major version",
		RunE:  upgrade,
	}
)

func init() {
	stolonctlCmd.AddCommand(upgradeCmd)
}

func upgrade(ccmd *cobra.Command, args []string) error {
	client, err := kubernetes.NewClient(kubeConfig)
	if err != nil {
		return trace.Wrap(err)
	}

	clusterConfig := &cluster.Config{
		Name:                    clusterName,
		Namespace:               namespace,
		KeepersPodFilterKey:     keepersFilterKey,
		KeepersPodFilterValue:   keepersFilterValue,
		SentinelsPodFilterKey:   sentinelsFilterKey,
		SentinelsPodFilterValue: sentinelsFilterValue,
		EtcdEndpoints:           etcdEndpoints,
		EtcdCertFile:            etcdCertFile,
		EtcdKeyFile:             etcdKeyFile,
		EtcdCAFile:              etcdCAFile,
	}

	pods, err := cluster.Pods(client, clusterConfig)
	if err != nil {
		return trace.Wrap(err)
	}

	status, err := cluster.GetStatus(pods, clusterConfig)
	if err != nil {
		return trace.Wrap(err)
	}

	fmt.Println(status.ClusterData)
	return nil
}
