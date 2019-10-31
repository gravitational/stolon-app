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
	"context"
	"time"

	"github.com/gravitational/stolon-app/internal/stolontool/pkg/kubernetes"
	"github.com/gravitational/stolon-app/internal/stolontool/pkg/utils"

	"github.com/gravitational/rigging"
	"github.com/gravitational/trace"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateDefinition creates Custom Resource Definition for
// StolonUpgradeResource
func CreateDefinition(ctx context.Context, kubeClient *kubernetes.Client, crdClient *Client) error {
	crd := &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: StolonUpgradeName,
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   StolonUpgradeGroup,
			Version: StolonUpgradeVersion,
			Scope:   StolonUpgradeScope,
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:     StolonUpgradeKind,
				Plural:   StolonUpgradePlural,
				Singular: StolonUpgradeSingular,
			},
		},
	}

	_, err := kubeClient.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	err = rigging.ConvertError(err)
	if err != nil && !trace.IsAlreadyExists(err) {
		return trace.Wrap(err)
	}

	// wait for the controller to init by trying to list
	// stolonupgrade resources
	return utils.Retry(ctx, 30, time.Second, func() error {
		_, err := crdClient.List()
		return err
	})
}
