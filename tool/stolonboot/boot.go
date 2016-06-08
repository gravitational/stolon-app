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
	"encoding/base64"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

func bootCluster(sentinels int, keepers int, proxies int, password string) error {
	err := createEtcd()
	if err != nil {
		return trace.Wrap(err)
	}

	err = createSentinels(sentinels)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createSecret(password)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createKeepers(keepers)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createProxies(proxies)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func createEtcd() error {
	log.Infof("creating etcd")
	cmd := kubeCommand("create", "-f", "/resources/etcd.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}

func createSentinels(sentinels int) error {
	log.Infof("creating sentinels")
	cmd := kubeCommand("create", "-f", "/resources/sentinel.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))

	if err = scaleReplicationController("stolon-sentinel", sentinels, 30); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func createSecret(password string) error {
	log.Infof("creating secret")
	cmd := kubeCommand("create", "-f", "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return trace.Wrap(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return trace.Wrap(err)
	}

	io.WriteString(stdin, generateSecret(password))
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		log.Errorf("%v", err)
		return trace.Wrap(err)
	}

	return nil
}

func createKeepers(keepers int) error {
	log.Infof("creating initial keeper")
	cmd := kubeCommand("create", "-f", "/resources/keeper.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))

	log.Infof("scaling up keepers")
	if err := scaleReplicationController("stolon-keeper", keepers, 30); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func createSeedKeeper() error {
	log.Infof("creating initial keeper")
	cmd := kubeCommand("create", "-f", "/resources/keeper.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))
	return nil
}

func scaleUpKeepers(keepers int) error {
	log.Infof("scaling up keepers")
	if err := scaleReplicationController("stolon-keeper", keepers, 30); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func createProxies(proxies int) error {
	log.Infof("creating proxies")
	cmd := kubeCommand("create", "-f", "/resources/proxy.yml")
	out, err := cmd.Output()
	if err != nil {
		return trace.Wrap(err)
	}
	log.Infof("cmd output: %s", string(out))

	if err = scaleReplicationController("stolon-proxy", proxies, 60); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func generateSecret(password string) string {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	template := `
---
apiVersion: v1
kind: Secret
metadata:
  name: stolon
type: Opaque
data:
  password: ` + encodedPassword + `
`
	return template
}
