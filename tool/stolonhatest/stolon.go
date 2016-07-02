// Copyright 2016 Gravitational, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

type KeeperState struct {
	ID              string `json:"ID"`
	Healthy         bool   `json:"Healthy"`
	ListenAddress   string `json:"ListenAddress"`
	Port            string `json:"Port"`
	PGListenAddress string `json:"PGListenAddress"`
	PGPort          string `json:"PGPort"`
}

func stolonCommand(args ...string) *exec.Cmd {
	return exec.Command("stolonctl", args...)
}

func GetMasterKeeperState() (*KeeperState, error) {
	cmd := stolonCommand("master-status", "--output", "json")
	out, err := cmd.Output()
	log.Debugf("cmd output: %s", string(out))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var ks KeeperState
	err = json.Unmarshal(out, &ks)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &ks, nil
}

func GetMasterKeeperIP() (string, error) {
	ks, err := GetMasterKeeperState()
	if err != nil {
		return "", trace.Wrap(err)
	}

	return ks.ListenAddress, nil
}
