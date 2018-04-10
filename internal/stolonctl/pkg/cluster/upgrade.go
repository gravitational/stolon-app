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
	"context"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/crd"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/kubernetes"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Upgrade(ctx context.Context, config Config) error {
	client, err := kubernetes.NewClient(config.KubeConfig)
	if err != nil {
		return trace.Wrap(err)
	}

	cfg, err := kubernetes.GetClientConfig(config.KubeConfig)
	if err != nil {
		return trace.Wrap(err)
	}

	crdclient, err := crd.NewClient(cfg, config.Namespace)
	if err != nil {
		return trace.Wrap(err)
	}

	err = crd.CreateDefinition(ctx, client, crdclient)
	if err != nil {
		return trace.Wrap(err)
	}

	objMeta := metav1.ObjectMeta{
		Name:      defaults.CRDName,
		Namespace: config.Namespace,
	}
	res, err := crdclient.Create(ctx, objMeta)
	if err != nil && !trace.IsAlreadyExists(err) {
		return trace.Wrap(err)
	}

	// temporary
	log.Info(res)
	return nil
}
