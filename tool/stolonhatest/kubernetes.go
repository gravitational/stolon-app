// Copyright 2016 Gravitational, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

type PodMetadata struct {
	Name string `json:"name"`
}

type PodCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type PodStatus struct {
	Phase      string         `json:"phase"`
	Conditions []PodCondition `json:"conditions"`
	PodIP      string         `json:"podIP"`
}

type Pod struct {
	Status   PodStatus   `json:"status"`
	Metadata PodMetadata `json:"metadata"`
}

type PodList struct {
	Items []Pod `json:"items"`
}

func kubeCommand(args ...string) *exec.Cmd {
	return exec.Command("kubectl", args...)
}

func GetPodList(selector string) (*PodList, error) {
	cmd := kubeCommand("get", "pods", "--selector", selector, "--output", "json")
	out, err := cmd.CombinedOutput()
	log.Debugf("cmd output: %s", string(out))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var pods PodList
	err = json.Unmarshal(out, &pods)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &pods, nil
}

func GetPodName(selector, ip string) (string, error) {
	pods, err := GetPodList(selector)
	if err != nil {
		return "", trace.Wrap(err)
	}

	var podName string
	for _, pod := range pods.Items {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				if pod.Status.PodIP == ip {
					podName = podName + pod.Metadata.Name
					break
				}
			}
		}
	}
	if podName == "" {
		return "", errors.New(fmt.Sprintf("Can't find any active pod with IP %v", ip))
	}

	return podName, nil
}

func DeletePod(name string, force bool) error {
	cmd := kubeCommand("delete", "pod", name)
	if force {
		cmd = kubeCommand("delete", "pod", name, "--now=true")
	}
	out, err := cmd.CombinedOutput()
	log.Debugf("cmd output: %s", string(out))
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func Label(rType, rName, label string) error {
	cmd := kubeCommand("label", rType, rName, label)
	out, err := cmd.CombinedOutput()
	log.Debugf("cmd output: %s", string(out))
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
