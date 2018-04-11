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

package utils

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

// Retry will retry function X times until period is reached
func Retry(ctx context.Context, times int, period time.Duration, fn func() error) error {
	if times < 1 {
		return nil
	}
	err := fn()
	for i := 1; i < times && err != nil; i += 1 {
		log.Infof("Attempt %v, result: %v, retry in %v.", i+1, trace.DebugReport(err), period)
		select {
		case <-ctx.Done():
			log.Infof("Context is closing, return")
			return err
		case <-time.After(period):
		}
		err = fn()
	}
	return err
}

// Run runs the command cmd and returns the output.
func Run(cmd *exec.Cmd) ([]byte, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return bytes.TrimSpace(output), err
	}
	return nil, nil
}
