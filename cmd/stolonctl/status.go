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
	"os"
	"time"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/cluster"

	"github.com/gravitational/trace"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of PostgreSQL cluster",
	RunE:  status,
}

func init() {
	stolonctlCmd.AddCommand(statusCmd)
}

func status(ccmd *cobra.Command, args []string) error {
	clusterStatus, err := Status()
	if err != nil {
		return trace.Wrap(err)
	}

	PrintStatus(clusterStatus)
	return nil
}

// Status returns status of Stolon cluster(Pods state and ClusterView)
func Status() (*cluster.Status, error) {
	clusterConfig := &cluster.Config{
		KubeConfig:         kubeConfig,
		Name:               clusterName,
		Namespace:          namespace,
		KeepersPodFilter:   keepersFilter,
		SentinelsPodFilter: sentinelsFilter,
		EtcdEndpoints:      etcdEndpoints,
		EtcdCertFile:       etcdCertFile,
		EtcdKeyFile:        etcdKeyFile,
		EtcdCAFile:         etcdCAFile,
	}

	status, err := cluster.GetStatus(clusterConfig)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return status, nil
}

// PrintStatus prints status of Stolon cluster to stdout
func PrintStatus(status *cluster.Status) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "POD_STATUS", "READY", "IP", "NODE", "AGE", "KEEPER_ID", "KEEPER_HEALTHY", "MASTER"})

	for _, pod := range status.PodsStatus {
		var keeperID string
		for _, keeperState := range status.ClusterData.KeepersState {
			if pod.PodIP == keeperState.ListenAddress {
				keeperID = keeperState.ID
			}
		}
		if keeperID != "" {
			table.Append([]string{pod.Name, pod.Status, string(pod.ReadyContainers), string(pod.TotalContainers), pod.PodIP, pod.HostIP,
				translateTimestamp(*pod.CreationTimestamp), keeperID, convertBoolean(status.ClusterData.KeepersState[keeperID].Healthy),
				status.ClusterData.KeepersState[keeperID].PGState.Role.String()})
		} else {
			table.Append([]string{pod.Name, pod.Status, string(pod.ReadyContainers), string(pod.TotalContainers), pod.PodIP, pod.HostIP,
				translateTimestamp(*pod.CreationTimestamp), "N/A", "N/A", "N/A"})
		}
	}

	table.Render()
}

// shortHumanDuration represents pod creation timestamp in
// human readable format
func shortHumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < -1 {
		return fmt.Sprintf("<invalid>")
	} else if seconds < 0 {
		return fmt.Sprintf("0s")
	} else if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if minutes := int(d.Minutes()); minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	} else if hours := int(d.Hours()); hours < 24 {
		return fmt.Sprintf("%dh", hours)
	} else if hours < 24*364 {
		return fmt.Sprintf("%dd", hours/24)
	}
	return fmt.Sprintf("%dy", int(d.Hours()/24/365))
}

// translateTimestamp returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestamp(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}
	return shortHumanDuration(time.Since(timestamp.Time))
}

func convertBoolean(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}
