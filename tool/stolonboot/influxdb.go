// Copyright 2016-2017 Gravitational, Inc.
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
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/roundtrip"
	"github.com/gravitational/trace"
)

var (
	// createUserQuery is the InfluxDB query to create a non-privileged user
	createUserQuery = "create user %v with password '%v'"
	// grantAllQuery is the InfluxDB query to grant write/read privileges on a database to a user
	grantAllQuery = "grant all on %q to %v"
)

// InfluxDBClient is InfluxDB API client
type InfluxDBClient struct {
	*roundtrip.Client
}

// NewInfluxDBClient creates a new client
func NewInfluxDBClient() (*InfluxDBClient, error) {
	client, err := roundtrip.NewClient(InfluxDBServiceAddr, "",
		roundtrip.BasicAuth(InfluxDBAdminUser, InfluxDBAdminPassword))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &InfluxDBClient{Client: client}, nil
}

// Setup sets up InfluxDB database
func (c *InfluxDBClient) Setup() error {
	queries := []string{
		fmt.Sprintf(createUserQuery, InfluxDBTelegrafUser, InfluxDBTelegrafPassword),
		fmt.Sprintf(grantAllQuery, InfluxDBDatabase, InfluxDBTelegrafUser),
	}
	for _, query := range queries {
		log.Infof("%v", query)

		response, err := c.PostForm(c.Endpoint("query"), url.Values{"q": []string{query}})
		if err != nil {
			return trace.Wrap(err)
		}

		log.Infof("%v %v %v", response.Code(), response.Headers(), string(response.Bytes()))
	}
	return nil
}
