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

const (
	// StolonUpgradeName defines the name of the Custom Resource Definition
	StolonUpgradeName = "stolonupgrades.stolon.gravitational.io"
	// StolonUpgradeGroup defines Custom Resource Definition group
	StolonUpgradeGroup = "stolon.gravitational.io"
	// StolonUpgradeVersion defines Custom Resource Definition version
	StolonUpgradeVersion = "v1"
	// StolonUpgradeAPIVersion defines Custom Resource Definition API version
	StolonUpgradeAPIVersion = "stolon.gravitational.io/v1"
	// StolonUpgradeScope defines Custom Resource Definition scope
	StolonUpgradeScope = "Namespaced"
	// StolonUpgradeKind defines Custom Resource Definition kind
	StolonUpgradeKind = "StolonUpgrade"
	// StolonUpgradePlural defines resource plural form.
	// It is required to create kubernetes Custom Resource Definition
	StolonUpgradePlural = "stolonupgrades"
	// StolonUpgradeSingular defines resource singular form.
	//It is required to create kubernetes Custom Resource Definition
	StolonUpgradeSingular = "stolonupgrade"

	// Status of stolon upgrade
	StolonUpgradeStatusInProgress = "in-progress"

	// Steps of stolon upgrade
	StolonUpgradeStepInit = "init"
)
