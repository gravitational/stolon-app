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
	"fmt"
	"net"
	"net/http"

	"github.com/gravitational/stolon-app/internal/stolonctl/pkg/defaults"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Listener mode",
	RunE:  server,
}

func init() {
	stolonctlCmd.AddCommand(serverCmd)
}

func server(ccmd *cobra.Command, args []string) error {
	if err := clusterConfig.CheckAndSetDefaults(); err != nil {
		return trace.Wrap(err)
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		log.Infof("%s %s %s %s", req.RemoteAddr, req.Host, req.RequestURI, req.UserAgent())
		statusHandler(w, req)
	})

	listener, err := net.Listen("tcp", defaults.ListenerPort)
	if err != nil {
		return trace.Wrap(err)
	}

	errChan := make(chan error, 10)
	go func() {
		if err := http.Serve(listener, nil); err != nil {
			errChan <- trace.Wrap(err)
			return
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}

func statusHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	clusterStatus, err := Status()
	if err != nil {
		log.Infof("Error getting cluster status: %v", trace.Wrap(err))
		http.Error(w, "Error getting cluster status", http.StatusBadGateway)
	}
	if status := IsClusterHealthy(clusterStatus); status {
		fmt.Fprint(w, "OK")
	} else {
		http.Error(w, "Cluster not healthy", http.StatusBadGateway)
	}
}
