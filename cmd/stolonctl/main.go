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
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/cluster"
	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"

	"github.com/fatih/color"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var (
	clusterConfig cluster.Config
	ctx           context.Context

	envs = map[string]string{
		"ETCD_CERT":          "etcd-cert-file",
		"ETCD_KEY":           "etcd-key-file",
		"ETCD_CACERT":        "etcd-ca-file",
		"ETCD_ENDPOINTS":     "etcd-endpoints",
		"POSTGRES_PASSWORD":  "postgres-password",
		"APP_VERSION":        "app-version",
		"RIG_CHANGESET":      "changeset",
		"NODE_NAME":          "nodename",
		"LAG_THRESHOLD":      "lag-threshold",
		"SENTINELS_SELECTOR": "sentinels-selector",
		"KEEPERS_SELECTOR":   "keepers-selector",
		"CLUSTER_NAME":       "cluster-name",
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
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.KubeConfig, "kubeconfig", "",
		"Kubernetes client config file")
	stolonctlCmd.PersistentFlags().StringVarP(&clusterConfig.Namespace, "namespace", "n",
		defaults.Namespace, "Kubernetes namespace for Stolon application")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.KeepersPodSelector, "keepers-selector",
		defaults.KeepersPodSelector, "Label to select keeper pods(ENV variable 'KEEPERS_SELECTOR')")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.SentinelsPodSelector, "sentinels-selector",
		defaults.SentinelsPodSelector, "Label to select sentinel pods(ENV variable 'SENTINELS_SELECTOR')")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.EtcdEndpoints, "etcd-endpoints",
		defaults.EtcdEndpoints, "Etcd server endpoints(ENV variable 'ETCD_ENDPOINTS')")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.EtcdCertFile, "etcd-cert-file", "",
		"Path to TLS certificate for connecting to etcd(ENV variable 'ETCD_CERT')")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.EtcdKeyFile, "etcd-key-file", "",
		"Path to TLS key for connecting to etcd(ENV variable 'ETCD_KEY')")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.EtcdCAFile, "etcd-ca-file", "",
		"Path to TLS CA for connecting to etcd(ENV variable 'ETCD_CACERT')")
	stolonctlCmd.PersistentFlags().StringVar(&clusterConfig.Name, "cluster-name",
		defaults.ClusterName, "Stolon cluster name")

	if err := bindFlagEnv(stolonctlCmd.PersistentFlags()); err != nil {
		log.WithError(err).Warn("Failed to bind environment variables to flags.")
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.TODO())
	go func() {
		exitSignals := make(chan os.Signal, 1)
		signal.Notify(exitSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		sig := <-exitSignals
		log.Infof("Caught signal: %v.", sig)
		cancel()
	}()

}

// bindFlagEnv binds environment variables to command flags
func bindFlagEnv(flagSet *flag.FlagSet) error {
	for env, flag := range envs {
		cmdFlag := flagSet.Lookup(flag)
		if cmdFlag != nil {
			if value := os.Getenv(env); value != "" {
				if err := cmdFlag.Value.Set(value); err != nil {
					return trace.Wrap(err)
				}

			}
		}
	}
	return nil
}

// printError prints the error message to the console
func printError(err error) {
	color.Red("[ERROR]: %v\n", trace.UserMessage(err))
}
