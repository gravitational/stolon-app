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

package kubernetes

import (
	"github.com/gravitational/trace"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client is the Kubernetes API client
type Client struct {
	*kubernetes.Clientset
}

// NewClient returns a new Kubernetes API client
func NewClient(kubeConfig string) (client *Client, err error) {
	config, err := getClientConfig(kubeConfig)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &Client{
		clientset,
	}, nil
}

// Pods returns stolon pods matching the specified label
func (c *Client) Pods(filter, namespace string) ([]v1.Pod, error) {
	labelSelector, err := labels.Parse(filter)
	if err != nil {
		return nil, trace.Wrap(err, "the provided label selector %s is not valid", filter)
	}

	podList, err := c.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector.String()})
	if err != nil {
		return nil, trace.Wrap(err)
	}

	if len(podList.Items) == 0 {
		return nil, trace.NotFound("pod(s) with label selector %s not found", labelSelector.String())
	}

	return podList.Items, nil
}

// getClientConfig returns client configguration,
// if master is not specified, in-cluster configuration is assumed
func getClientConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()

}
