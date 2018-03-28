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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
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

// Pods returns Stolon pods with matched label
func (c *Client) Pods(filterKey, filterValue, namespace string) ([]v1.Pod, error) {
	label, err := matchLabel(filterKey, filterValue)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	podList, err := c.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: label.String()})
	if err != nil {
		return nil, trace.Wrap(err)
	}

	if len(podList.Items) == 0 {
		return nil, trace.NotFound("pod(s) with label %s' not found", label.String())
	}

	return podList.Items, nil
}

// getClientConfig returns client config, if master not specified assumes in cluster config
func getClientConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()

}

// Label represents a Kubernetes label which is used
// as a search target for Pods
type Label struct {
	Key   string
	Value string
}

// matchLabel matches a resource with the specified label
func matchLabel(key, value string) (labels.Selector, error) {
	req, err := labels.NewRequirement(key, selection.In, []string{value})
	if err != nil {
		return nil, trace.Wrap(err, "failed to build a requirement from %v=%q", key, value)
	}
	selector := labels.NewSelector()
	selector = selector.Add(*req)
	return selector, nil
}
