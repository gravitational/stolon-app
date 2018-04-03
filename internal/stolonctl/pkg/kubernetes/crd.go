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

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateCRD() error {
	crd := &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaults.StolonUpgradeName,
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   defaults.StolonUpgradeGroup,
			Version: defaults.StolonUpgradeVersion,
			Scope:   defaults.StolonUpgradeScope,
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:     defaults.StolonUpgradeKind,
				Plural:   defaults.StolonUpgradePlural,
				Singular: defaults.StolonUpgradeSingular,
			},
		},
	}
	return nil
}
