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
	"time"

	"github.com/gravitational/rigging"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/crd"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/kubernetes"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/utils"
	"github.com/gravitational/stolon/common"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pgDumpCommand             = "pg_dumpall"
	pgRestoreCommand          = "pg_restore"
	stolonKeeperDaemonset     = "stolon-keeper"
	stolonSentinelDeployment  = "stolon-sentinel"
	jobTimeout                = 600 // 10 minutes
	volumeData                = "data"
	volumeDataHostPath        = "/var/lib/data/stolon"
	volumeDataMountPath       = "/stolon-data"
	volumeClusterCA           = "cluster-ca"
	volumeClusterCAMountPath  = "/etc/ssl/cluster-ca"
	volumeDefaultSSL          = "cluster-default-ssl"
	volumeDefaultSSLMountPath = "/etc/ssl/cluster-default"
	pgMaster                  = "master"
	pgStandby                 = "standby"
)

// upgradeControl defines configuration for upgrade
type upgradeControl struct {
	kubeClient *kubernetes.Client
	crdClient  *crd.Client
	config     Config
}

// Upgrade upgrades stolon cluster
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

	upgradeControl := &upgradeControl{
		kubeClient: client,
		crdClient:  crdclient,
		config:     config,
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

	if !crdclient.IsPhaseCompleted(res, crd.StolonUpgradePhaseInit) {
		res, err = crdclient.MarkPhase(res, crd.StolonUpgradePhaseInit, crd.StolonUpgradeStatusCompleted)
		if err != nil {
			return trace.Wrap(err)
		}
	}
	// Check status of cluster
	res, err = upgradeControl.checkStatus(res, crd.StolonUpgradePhaseChecks)
	if err != nil {
		return trace.Wrap(err)
	}

	// Backup Postgres
	res, err = upgradeControl.executePhase(res, crd.StolonUpgradePhaseBackupPostgres, func() error {
		err := upgradeControl.backupPostgres()
		return err
	})
	if err != nil {
		return trace.Wrap(err)
	}

	// Delete keeper/sentinel daemonset/deployment
	res, err = upgradeControl.executePhase(res, crd.StolonUpgradePhaseDeleteDeployment, func() error {
		err := upgradeControl.deleteDeployment(stolonSentinelDeployment)
		return err
	})
	if err != nil {
		return trace.Wrap(err)
	}

	res, err = upgradeControl.executePhase(res, crd.StolonUpgradePhaseDeleteDaemonset, func() error {
		err := upgradeControl.deleteDaemonset(stolonKeeperDaemonset)
		return err
	})
	if err != nil {
		return trace.Wrap(err)
	}

	// Create job to upgrade PostgreSQL schema on keeper master
	res, err = upgradeControl.executePhase(res, crd.StolonUpgradePhaseUpgradePostgresSchema, func() error {
		var masterNodeName string
		standbyNodeNames := make(map[string]string)
		for _, pod := range res.Spec.PodsStatus {
			if pod.PodIP == res.Spec.ClusterData.KeepersState[res.Spec.ClusterData.ClusterView.Master].ListenAddress {
				masterNodeName = pod.HostIP
			} else {
				for _, keeperState := range res.Spec.ClusterData.KeepersState {
					if pod.PodIP == keeperState.ListenAddress {
						standbyNodeNames[keeperState.ID] = pod.HostIP
					}
				}
			}
		}

		log.Infof("master: %s, standbys: %v\n", masterNodeName, standbyNodeNames)
		if masterNodeName == "" {
			return trace.BadParameter("IP address of keeper master pod is empty")
		}
		err := upgradeControl.upgradePostgresSchema(ctx, pgMaster, masterNodeName, "")
		if err != nil {
			return trace.Wrap(err)
		}
		for id, standbyNodeName := range standbyNodeNames {
			err = upgradeControl.upgradePostgresSchema(ctx, pgStandby, standbyNodeName, id)
			if err != nil {
				return trace.Wrap(err)
			}
		}
		return nil

	})
	if err != nil {
		return trace.Wrap(err)
	}

	if _, err := upgradeControl.crdClient.CompleteUpgrade(res); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (u *upgradeControl) executePhase(res *crd.StolonUpgradeResource, phase string, fn func() error) (*crd.StolonUpgradeResource, error) {
	var err error
	if !u.crdClient.IsPhaseCompleted(res, phase) {
		res, err = u.crdClient.MarkPhase(res, phase, crd.StolonUpgradeStatusInProgress)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		if err := fn(); err != nil {
			return nil, trace.Wrap(err)
		}
		res, err = u.crdClient.MarkPhase(res, phase, crd.StolonUpgradeStatusCompleted)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}
	return res, nil
}

func (u *upgradeControl) checkStatus(res *crd.StolonUpgradeResource, phase string) (*crd.StolonUpgradeResource, error) {
	var err error
	if !u.crdClient.IsPhaseCompleted(res, phase) {
		res, err = u.crdClient.MarkPhase(res, phase, crd.StolonUpgradeStatusInProgress)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		status, err := GetStatus(u.config)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		masterStatus, err := status.getMasterStatus()
		if err != nil {
			return nil, trace.Wrap(err, "cannot start upgrade")
		}

		if !masterStatus.Healthy {
			return nil, trace.BadParameter("cannot start upgrade, keeper master %s is unhealthy.",
				masterStatus.PodName)
		}
		log.Infof("Keeper master: pod=%s, healthy=%v, host=%s", masterStatus.PodName,
			masterStatus.Healthy, masterStatus.HostIP)

		res, err = u.crdClient.UpdateClusterInfo(res, *status.ClusterData, status.PodsStatus)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		res, err = u.crdClient.MarkPhase(res, phase, crd.StolonUpgradeStatusCompleted)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}
	return res, nil
}

func (u *upgradeControl) backupPostgres() error {
	if err := u.createPgPassFile(); err != nil {
		return trace.Wrap(err)
	}
	options := []string{fmt.Sprintf("-h%s", u.config.Postgres.Host),
		fmt.Sprintf("-p%s", u.config.Postgres.Port),
		fmt.Sprintf("-U%s", u.config.Postgres.User), "-c",
		fmt.Sprintf("-f%s", u.config.Postgres.BackupPath)}

	cmd := exec.Command(pgDumpCommand, options...)
	env := os.Environ()
	// Add PGPASSFILE env variable to set path to the password file
	env = append(env, fmt.Sprintf("PGPASSFILE=%s", u.config.Postgres.PgPassPath))
	cmd.Env = env

	output, err := utils.Run(cmd)
	if err != nil {
		return trace.Wrap(err, "Command output: %v", string(output))
	}
	log.Infof("Backed up Postgres to %s", u.config.Postgres.BackupPath)

	return nil
}

func (u *upgradeControl) createPgPassFile() error {
	content := fmt.Sprintf("%s:%s:%s:%s:%s", u.config.Postgres.Host, u.config.Postgres.Port,
		"*", u.config.Postgres.User, u.config.Postgres.Password)

	return common.WriteFileAtomic(u.config.Postgres.PgPassPath, []byte(content), 0600)
}

// deleteDeployment deletes sentinels deployment
func (u *upgradeControl) deleteDeployment(name string) error {
	dsClient := u.kubeClient.Client.AppsV1().Deployments(u.config.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := dsClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return trace.Wrap(err, "cannot delete stolon-sentinel deployment")
	}
	return nil
}

// deleteDaemonset deletes keepers daemonset
func (u *upgradeControl) deleteDaemonset(name string) error {
	dsClient := u.kubeClient.Client.ExtensionsV1beta1().DaemonSets(u.config.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := dsClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return trace.Wrap(err, "cannot delete stolon-keeper daemonset")
	}
	return nil
}

func (u *upgradeControl) upgradePostgresSchema(ctx context.Context, pgRole, nodeName, id string) error {
	jobControl, err := rigging.NewJobControl(rigging.JobConfig{u.generateJob(pgRole, nodeName, id), u.kubeClient.Client})
	if err != nil {
		return trace.Wrap(err)
	}
	if err := jobControl.Upsert(ctx); err != nil {
		return trace.Wrap(err)
	}
	// wait 5 minutes for completed status of Job
	return utils.Retry(ctx, 300, time.Second, func() error {
		err := jobControl.Status()
		return err
	})
}

func (u *upgradeControl) generateJob(pgRole, nodeName, id string) *batchv1.Job {
	var jobName string
	if pgRole == pgMaster {
		jobName = fmt.Sprintf("stolon-keeper-schema-upgrade-%s-%s", pgRole, u.config.Changeset)
	} else {
		jobName = fmt.Sprintf("stolon-keeper-schema-upgrade-%s-%s-%s", pgRole, id, u.config.Changeset)
	}
	command := "/usr/local/bin/keeper-upgrade.sh"
	if pgRole == pgStandby {
		command = "/usr/local/bin/clean-postgres-data.sh"
	}
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: u.config.Namespace,
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: int64Ptr(jobTimeout),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: u.config.Namespace,
					Labels: map[string]string{
						"name": jobName,
					},
				},
				Spec: apiv1.PodSpec{
					Volumes: []apiv1.Volume{
						{
							Name: volumeData,
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: volumeDataHostPath,
								},
							},
						},
						{
							Name: volumeClusterCA,
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: volumeClusterCA,
								},
							},
						},
						{
							Name: volumeDefaultSSL,
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: volumeDefaultSSL,
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					NodeName:      nodeName,
					SecurityContext: &apiv1.PodSecurityContext{
						RunAsUser: int64Ptr(0),
					},
					ServiceAccountName: "stolon-keeper",
					Containers: []apiv1.Container{
						{
							Name:    "upgrade",
							Image:   fmt.Sprintf("leader.telekube.local:5000/stolon-jobs:%s", u.config.NewAppVersion),
							Command: []string{command},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      volumeData,
									MountPath: volumeDataMountPath,
								},
								{
									Name:      volumeClusterCA,
									MountPath: volumeClusterCAMountPath,
								},
								{
									Name:      volumeDefaultSSL,
									MountPath: volumeDefaultSSLMountPath,
								},
							},
						},
					},
				},
			},
		},
	}
}

func int64Ptr(i int64) *int64 { return &i }
