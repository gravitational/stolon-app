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

package main

import (
	"os"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"

	"github.com/fatih/color"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	kubeConfig      string
	keepersFilter   string
	sentinelsFilter string
	namespace       string
	etcdEndpoints   string
	etcdCertFile    string
	etcdKeyFile     string
	etcdCAFile      string
	clusterName     string

	envs = map[string]string{
		"ETCD_CERT":      "etcd-cert-file",
		"ETCD_KEY":       "etcd-key-file",
		"ETCD_CACERT":    "etcd-ca-file",
		"ETCD_ENDPOINTS": "etcd-endpoints",
	}

	stolonctlCmd = &cobra.Command{
		Use:   "",
		Short: "PostgreSQL major versions upgrade tool for Stolon cluster",
		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func main() {
	if err := stolonctlCmd.Execute(); err != nil {
		log.Error(trace.DebugReport(err))
		printError(err)
		os.Exit(255)
	}
}
func init() {
	stolonctlCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "Kubernetes client config file")
	stolonctlCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", defaults.Namespace, "Kubernetes namespace for Stolon application")
	stolonctlCmd.PersistentFlags().StringVar(&keepersFilter, "keepers-filter", defaults.KeepersPodFilter, "Label to filter keeper pods")
	stolonctlCmd.PersistentFlags().StringVar(&sentinelsFilter, "sentinels-filter", defaults.SentinelsPodFilter, "Label to filter sentinel pods")
	stolonctlCmd.PersistentFlags().StringVar(&etcdEndpoints, "etcd-endpoints", defaults.EtcdEndpoints, "Etcd server endpoints")
	stolonctlCmd.PersistentFlags().StringVar(&etcdCertFile, "etcd-cert-file", "", "Path to TLS certificate for connecting to etcd")
	stolonctlCmd.PersistentFlags().StringVar(&etcdKeyFile, "etcd-key-file", "", "Path to TLS key for connecting to etcd")
	stolonctlCmd.PersistentFlags().StringVar(&etcdCAFile, "etcd-ca-file", "", "Path to TLS CA for connecting to etcd")
	stolonctlCmd.PersistentFlags().StringVar(&clusterName, "cluster-name", defaults.ClusterName, "Stolon cluster name")

	for env, flag := range envs {
		cmdFlag := stolonctlCmd.PersistentFlags().Lookup(flag)
		if value := os.Getenv(env); value != "" {
			if cmdFlag != nil {
				cmdFlag.Value.Set(value)
			}
		}
	}
}

// printError prints the error message to the console

func printError(err error) {
	color.Red("[ERROR]: %v\n", trace.UserMessage(err))
}
