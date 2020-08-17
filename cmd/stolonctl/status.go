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
	"io"
	"math"
	"os"
	"text/tabwriter"
	"time"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/cluster"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"

	"github.com/gravitational/stolon/common"
	"github.com/gravitational/trace"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	statusCmd = &cobra.Command{
		Use:          "status",
		Short:        "Show status of PostgreSQL cluster",
		SilenceUsage: true,
		RunE:         status,
	}
	shortOutput bool
	// lagThreshold is the threshold for the lag between master and standby
	// reaching which cluster status will be considered as failed
	lagThreshold uint64
)

func init() {
	const defaultShortOutput = false

	stolonctlCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().BoolVarP(&shortOutput, "short", "s", defaultShortOutput, "Output only overall cluster status and reason if unhealthy.")
	statusCmd.PersistentFlags().Uint64Var(&lagThreshold, "lag-threshold", defaults.LagThreshold, "Threshold for lag between master node and standby.")
}

func status(ccmd *cobra.Command, args []string) error {
	if err := clusterConfig.CheckAndSetDefaults(); err != nil {
		return trace.Wrap(err)
	}

	clusterStatus, err := Status()
	if err != nil {
		return trace.Wrap(err)
	}

	if err := PrintStatus(clusterStatus); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

// Status returns status of stolon cluster(Pods state and ClusterView)
func Status() (*cluster.Status, error) {
	status, err := cluster.GetStatus(clusterConfig)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return status, nil
}

// PrintStatus prints status of stolon cluster to stdout
func PrintStatus(status *cluster.Status) error {
	output := os.Stdout
	if !shortOutput {
		printClusterStatus(status, output)
		fmt.Fprintln(output)
	}
	if err := printOverallStatus(status, output); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func printClusterStatus(status *cluster.Status, output io.Writer) {
	w := new(tabwriter.Writer)

	var (
		minwidth int
		tabwidth = 8
		padding  = 2
		flags    uint
		padchar  byte = '\t'
	)
	w.Init(output, minwidth, tabwidth, padding, padchar, flags)
	fmt.Fprintln(w, "NAME\tREADY\tSTATUS\tIP\tNODE\tAGE\tKEEPER_ID\tHEALTHY\tROLE")

	for _, pod := range status.PodsStatus {
		var keeperID string
		for _, keeperState := range status.ClusterData.KeepersState {
			if pod.PodIP == keeperState.ListenAddress {
				keeperID = keeperState.ID
			}
		}
		if keeperID != "" && status.ClusterData.KeepersState[keeperID].PGState != nil {
			fmt.Fprintf(w, "%s\t%v/%v\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", pod.Name,
				pod.ReadyContainers, pod.TotalContainers, pod.Status, pod.PodIP, pod.HostIP,
				translateTimestamp(*pod.CreationTimestamp), keeperID,
				translateBoolean(status.ClusterData.KeepersState[keeperID].Healthy),
				status.ClusterData.KeepersState[keeperID].PGState.Role.String())
		} else {
			fmt.Fprintf(w, "%s\t%v/%v\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", pod.Name,
				pod.ReadyContainers, pod.TotalContainers, pod.Status, pod.PodIP, pod.HostIP,
				translateTimestamp(*pod.CreationTimestamp), "N/A", "N/A", "N/A")
		}
	}

	w.Flush()
}

func printOverallStatus(status *cluster.Status, output io.Writer) error {
	reason, isHealthy := isClusterHealthy(status)
	fmt.Fprintf(output, "Cluster status: %s\n", getStatusString(isHealthy))
	if !isHealthy {
		return trace.Errorf("Cluster is unhealthy. Reason: %s", reason)
	}

	return nil
}

// shortHumanDuration represents pod creation timestamp in
// human readable format
func shortHumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < -1 {
		return "<invalid>"
	} else if seconds < 0 {
		return "0s"
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

func translateBoolean(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func isClusterHealthy(status *cluster.Status) (unhealthyReason string, healthy bool) {
	metrics := collectMetricsFromStatus(status)
	unhealthyReason = computeReasonUnhealthy(metrics)
	if unhealthyReason != "" {
		return unhealthyReason, false
	}
	return "", true
}

func computeReasonUnhealthy(metrics statusMetrics) (reason string) {
	if metrics.runningPods <= 1 {
		return "cluster is running with less than 2 nodes"
	}
	if metrics.runningMasters == 0 {
		return "cluster has no master node"
	}
	if metrics.runningMasters > 1 {
		return "cluster has more than one master node"
	}
	if !metrics.masterHealthy {
		return "master is unhealthy"
	}
	if metrics.runningStandbys == 0 {
		return "cluster has no standby servers"
	}
	if metrics.healthyNodes < 2 {
		return fmt.Sprintf("there are only %v healthy keeper nodes in cluster", metrics.healthyNodes)
	}
	if metrics.replicationLag {
		return "high replication lag between master and stanby(s)"
	}

	return ""
}

func collectMetricsFromStatus(status *cluster.Status) (result statusMetrics) {
	var masterID string
	for _, pod := range status.PodsStatus {
		var keeperID string
		for _, keeperState := range status.ClusterData.KeepersState {
			if pod.PodIP == keeperState.ListenAddress {
				keeperID = keeperState.ID
			}
		}
		if keeperID != "" && status.ClusterData.KeepersState[keeperID].PGState != nil {
			result.runningPods++
			if status.ClusterData.KeepersState[keeperID].PGState.Role == common.MasterRole {
				if status.ClusterData.KeepersState[keeperID].Healthy {
					result.masterHealthy = true
					masterID = keeperID
				}
				// count amount of running master nodes
				result.runningMasters++
			} else {
				// count amount of running standby nodes
				result.runningStandbys++
				if masterID != "" {
					if math.Abs(float64(status.ClusterData.KeepersState[masterID].PGState.XLogPos-status.ClusterData.KeepersState[keeperID].PGState.XLogPos)) > float64(lagThreshold) {
						result.replicationLag = true
					}
				}
			}
			if status.ClusterData.KeepersState[keeperID].Healthy {
				// count amount of running healthy nodes
				result.healthyNodes++
			}
		}
	}
	return result
}

func getStatusString(status bool) string {
	if status {
		return "Healthy"
	}
	return "Unhealthy"
}

type statusMetrics struct {
	runningPods     int
	masterHealthy   bool
	runningMasters  int
	runningStandbys int
	healthyNodes    int
	replicationLag  bool
}
