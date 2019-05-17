/*
Copyright (C) 2019 Gravitational, Inc.

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
	"net/http"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:          "server",
	Short:        "HTTP listener",
	SilenceUsage: true,
	RunE:         server,
}

func init() {
	stolonctlCmd.AddCommand(serverCmd)
}

func server(ccmd *cobra.Command, args []string) error {
	if err := clusterConfig.CheckAndSetDefaults(); err != nil {
		return trace.Wrap(err)
	}

	log.Info("Starting endpoint.")
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		clusterStatus, err := Status()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}

		reason, isHealthy := isClusterHealthy(clusterStatus)
		if !isHealthy {
			log.Errorf("Cluster is unhealthy. Reason: %s", reason)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(reason))
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":8080", nil)
}
