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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/crd"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/kubernetes"
	"github.com/gravitational/stolon/common"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pgDumpCommand    = "pg_dumpall"
	pgRestoreCommand = "pg_restore"
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
	res, err := crdclient.CreateOrRead(objMeta)
	if err != nil {
		return trace.Wrap(err)
	}

	// temporary
	log.Info(res)
	res, err = crdclient.MarkStep(res, crd.StolonUpgradePhaseInit, crd.StolonUpgradeStatusCompleted)
	if err != nil {
		return trace.Wrap(err)
	}

	// Backup Postgres
	res, err = crdclient.MarkStep(res, crd.StolonUpgradePhaseBackupPostgres, crd.StolonUpgradeStatusInProgress)
	if err != nil {
		return trace.Wrap(err)
	}
	if err = backupPostgres(config); err != nil {
		return trace.Wrap(err)
	}
	res, err = crdclient.MarkStep(res, crd.StolonUpgradePhaseBackupPostgres, crd.StolonUpgradeStatusCompleted)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func backupPostgres(config Config) error {

	if err := createPgPassFile(config); err != nil {
		return trace.Wrap(err)
	}
	options := []string{fmt.Sprintf("-h%s", config.PostgresHost),
		fmt.Sprintf("-p%s", config.PostgresPort),
		fmt.Sprintf("-U%s", config.PostgresUser), "-c",
		fmt.Sprintf("-f%s", config.PostgresBackupPath)}

	_, err := exec.Command(pgDumpCommand, options...).CombinedOutput()
	if err != nil && !isExitError(err) {
		return trace.Wrap(err)
	}

	return nil
}

func isExitError(err error) bool {
	if _, ok := err.(*exec.ExitError); ok {
		return true
	}
	return false
}

func createPgPassFile(config Config) error {
	content := fmt.Sprintf("%s:%s:%s:%s:%s", config.PostgresHost, config.PostgresPort,
		"*", config.PostgresUser, config.PostgresPassword)

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return trace.BadParameter("$HOME environment variable are not set. Cannot write .pgpass file.")
	}
	return common.WriteFileAtomic(filepath.Join(homeDir, ".pgpass"), []byte(content), 0600)
}
