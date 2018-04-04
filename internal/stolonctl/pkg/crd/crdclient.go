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

package crd

import (
	"encoding/json"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/utils"
	"github.com/gravitational/trace"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// CrdClient creates CRD client interface
func CrdClient(cfg *rest.Config, namespace string) (*crdclient, error) {
	cfg.APIPath = "/apis"
	if cfg.UserAgent == "" {
		cfg.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	cfg.ContentType = runtime.ContentTypeJSON
	cfg.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	cfg.GroupVersion = &schema.GroupVersion{Group: defaults.StolonUpgradeGroup, Version: defaults.StolonUpgradeVersion}

	clt, err := rest.RESTClientFor(cfg)
	if err != nil {
		return nil, utils.ConvertError(err)
	}

	return &crdclient{
		client:    clt,
		namespace: namespace,
		plural:    defaults.StolonUpgradePlural,
	}, nil
}

type crdclient struct {
	client    *rest.RESTClient
	namespace string
	plural    string
}

func (c *crdclient) Create(obj StolonUpgradeResource) (*StolonUpgradeResource, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var raw runtime.Unknown
	err = c.client.Post().
		Namespace(c.namespace).
		Resource(c.plural).
		Body(data).
		Do().
		Into(&raw)
	if err != nil {
		return nil, utils.ConvertError(err)
	}

	var result StolonUpgradeResource
	if err := json.Unmarshal(raw.Raw, &result); err != nil {
		return nil, trace.Wrap(err)
	}
	return &result, nil
}

func (c *crdclient) List() (*StolonUpgradeList, error) {
	var raw runtime.Unknown
	err := c.client.Get().
		Namespace(c.namespace).
		Resource(c.plural).
		Do().
		Into(&raw)
	if err != nil {
		return nil, utils.ConvertError(err)
	}

	var result StolonUpgradeList
	if err := json.Unmarshal(raw.Raw, &result); err != nil {
		return nil, trace.Wrap(err)
	}
	return &result, nil
}
