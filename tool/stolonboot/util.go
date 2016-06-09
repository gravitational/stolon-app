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
	"fmt"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

func kubeCommand(args ...string) *exec.Cmd {
	return exec.Command("/usr/local/bin/kubectl", args...)
}

func getReplicationController(name string) (*ReplicationController, error) {
	cmd := kubeCommand("get", fmt.Sprintf("rc/%s", name), "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var rc ReplicationController
	err = json.Unmarshal(out, &rc)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &rc, nil
}

func getRCPods(name string) (*PodList, error) {
	cmd := kubeCommand("get", "pods", "-l", fmt.Sprintf("name=%s", name), "-o", "json")
	out, err := cmd.Output()
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

func scaleReplicationController(name string, replicas int, tries int) error {
	err := waitForNPods(name, 1, time.Second, tries)
	if err != nil {
		return trace.Wrap(err)
	}

	for i := 1; i < replicas; i++ {
		cmd := kubeCommand("scale", fmt.Sprintf("--replicas=%d", i+1), fmt.Sprintf("rc/%s", name))
		out, err := cmd.Output()
		if err != nil {
			return trace.Wrap(err)
		}
		log.Infof("cmd output: %s", string(out))

		err = waitForNPods(name, i+1, time.Second, tries)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	return nil
}

func waitForNPods(rcName string, desired int, delay time.Duration, tries int) error {
	for i := 0; i < tries; i++ {
		pods, err := getRCPods(rcName)
		if err != nil {
			return trace.Wrap(err)
		}

		healthy := 0
		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					healthy++
					break
				}
			}
		}

		log.Infof("looking for %d pods, have %d pods, %d healthy", desired, len(pods.Items), healthy)
		if len(pods.Items) == desired && healthy == desired {
			return nil
		}
		time.Sleep(delay)
	}

	return trace.Errorf("timed out waiting for pods")
}
