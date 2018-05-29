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
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/crd"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/kubernetes"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/utils"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

// Plan shows plan of upgrade operation
func Plan(config Config) error {
	cfg, err := kubernetes.GetClientConfig(config.KubeConfig)
	if err != nil {
		return trace.Wrap(err)
	}

	crdclient, err := crd.NewClient(cfg, config.Namespace)
	if err != nil {
		return trace.Wrap(err)
	}

	resourceName := fmt.Sprintf("%s-%s", defaults.CRDName, config.Upgrade.Changeset)
	res, err := crdclient.Get(resourceName)
	if trace.IsNotFound(err) {
		log.Infof("Stolon upgrade %s is not started yet. Check with 'kubectl get stolonupgrades'", resourceName)
		return nil
	}
	if err != nil {
		return trace.Wrap(err)
	}
	formatPlanText(os.Stdout, *res)

	return nil
}

func formatPlanText(w io.Writer, upgradeResource crd.StolonUpgradeResource) {
	var t tabwriter.Writer
	t.Init(w, 0, 10, 5, ' ', 0)
	utils.PrintTableHeader(&t, []string{"Phase", "Description", "State", "Node", "Updated"})
	for _, phase := range upgradeResource.Spec.Phases {
		printPhase(&t, phase)
	}
	t.Flush()
}

func printPhase(w io.Writer, phase crd.StolonUpgradePhase) {
	marker := "*"
	if phase.Status == crd.StolonUpgradeStatusInProgress {
		marker = "→"
	} else if phase.Status == crd.StolonUpgradeStatusCompleted {
		marker = "✓"
	} else if phase.Status == crd.StolonUpgradeStatusFailed {
		marker = "⚠"
	}

	fmt.Fprintf(w, "%v %v\t%v\t%v\t%v\t%v\n",
		marker,
		phase.Name,
		phase.Description,
		formatStatus(phase.Status),
		formatNodeName(phase.NodeName),
		formatTimestamp(phase.UpdatedTimestamp))
}

func formatStatus(status string) string {
	switch status {
	case crd.StolonUpgradeStatusUnstarted:
		return "Unstarted"
	case crd.StolonUpgradeStatusInProgress:
		return "In Progress"
	case crd.StolonUpgradeStatusFailed:
		return "Failed"
	case crd.StolonUpgradeStatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format(time.UnixDate)
}

func formatNodeName(name string) string {
	if name == "" {
		return "-"
	}
	return name
}
