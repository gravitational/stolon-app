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
	"fmt"

	"github.com/gravitational/stolon-app/internal/stolontool/pkg/defaults"
)

// ResourceName returns name of stolonupgrade resource
func ResourceName(config Config) string {
	return fmt.Sprintf("%s-%s", defaults.CRDName, config.Upgrade.Changeset)
}
