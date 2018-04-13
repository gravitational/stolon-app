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

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/crd"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/kubernetes"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/utils"
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
	status, err := GetStatus(config)
	if err != nil {
		return trace.Wrap(err)
	}
	masterStatus, err := status.getMasterStatus()
	if err != nil {
		return trace.Wrap(err, "Cannot start upgrade.")
	}
	if !masterStatus.Healthy {
		return trace.Errorf("Cannot start upgrade. Keeper master %s is unhealthy.",
			masterStatus.PodName)
	}
	log.Infof("Keeper master: pod=%s, healthy=%v, host=%s", masterStatus.PodName,
		masterStatus.Healthy, masterStatus.HostIP)

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

	// TODO delete debugging info
	log.Info(res)
	if !crdclient.IsPhaseCompleted(res, crd.StolonUpgradePhaseInit) {
		res, err = crdclient.MarkPhase(res, crd.StolonUpgradePhaseInit, crd.StolonUpgradeStatusCompleted)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	// Backup Postgres
	if !crdclient.IsPhaseCompleted(res, crd.StolonUpgradePhaseBackupPostgres) {
		res, err = crdclient.MarkPhase(res, crd.StolonUpgradePhaseBackupPostgres, crd.StolonUpgradeStatusInProgress)
		if err != nil {
			return trace.Wrap(err)
		}
		if err = backupPostgres(config); err != nil {
			return trace.Wrap(err)
		}
		res, err = crdclient.MarkPhase(res, crd.StolonUpgradePhaseBackupPostgres, crd.StolonUpgradeStatusCompleted)
		if err != nil {
			return trace.Wrap(err)
		}
	}
	return nil
}

func backupPostgres(config Config) error {
	if err := createPgPassFile(config); err != nil {
		return trace.Wrap(err)
	}
	options := []string{fmt.Sprintf("-h%s", config.Postgres.Host),
		fmt.Sprintf("-p%s", config.Postgres.Port),
		fmt.Sprintf("-U%s", config.Postgres.User), "-c",
		fmt.Sprintf("-f%s", config.Postgres.BackupPath)}

	cmd := exec.Command(pgDumpCommand, options...)
	env := os.Environ()
	// Add PGPASSFILE env variable to set path to the password file
	env = append(env, fmt.Sprintf("PGPASSFILE=%s", config.Postgres.PgPassPath))
	cmd.Env = env

	output, err := utils.Run(cmd)
	if err != nil {
		return trace.Wrap(err, "Command output: %v", string(output))
	}
	log.Infof("Backed up Postgres to %s", config.Postgres.BackupPath)

	return nil
}

func createPgPassFile(config Config) error {
	content := fmt.Sprintf("%s:%s:%s:%s:%s", config.Postgres.Host, config.Postgres.Port,
		"*", config.Postgres.User, config.Postgres.Password)

	return common.WriteFileAtomic(config.Postgres.PgPassPath, []byte(content), 0600)
}
