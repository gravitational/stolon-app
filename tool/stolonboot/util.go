package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

type ReplicationController struct {
	Status ReplicationControllerStatus
}

type ReplicationControllerStatus struct {
	Replicas             int
	FullyLabeledReplicas int
	ObservedGeneration   int
}

type PodCondition struct {
	Type   string
	Status string
}

type PodStatus struct {
	Phase      string
	Conditions []PodCondition
}

type Pod struct {
	Status PodStatus
}

type PodList struct {
	Items []Pod
}

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
		cmd := exec.Command("/usr/local/bin/kubectl", "scale", fmt.Sprintf("--replicas=%d", i+1), fmt.Sprintf("rc/%s", name))
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
