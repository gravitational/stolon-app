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
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

func CrdClient(client *rest.RESTClient, scheme *runtime.Scheme, namespace string) *crdclient {
	return &crdclient{
		client:    client,
		namespace: namespace,
		plural:    defaults.StolonUpgradePlural,
		codec:     runtime.NewParameterCodec(scheme),
	}
}

type crdclient struct {
	client    *rest.RESTClient
	namespace string
	plural    string
	codec     runtime.ParameterCodec
}
