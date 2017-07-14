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
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
)

const (
	RandomlyGeneratedDefault = "<randomly generated>"
	DefaultPasswordLength    = 20

	// InfluxDBServiceAddr is the address of InfluxDB service
	InfluxDBServiceAddr = "http://influxdb.kube-system.svc:8086"
	// InfluxDBAdminUser is the InfluxDB admin user name
	InfluxDBAdminUser = "root"
	// InfluxDBAdminPassword is the InfluxDB admin user password
	InfluxDBAdminPassword = "root"
	// InfluxDBDatabase is the name of the database where all metrics go
	InfluxDBDatabase = "k8s"
	// InfluxDBTelegrafUser is the InfluxDB user for Telegraf
	InfluxDBTelegrafUser = "telegraf"
	// InfluxDBTelegrafPassword is the InfluxDB password for Telegraf
	InfluxDBTelegrafPassword = "telegraf"
)

func main() {
	sentinels := flag.Int("sentinels", 2, "number of sentinels")
	rpc := flag.Int("rpc", 1, "number of RPC")
	password := flag.String("password", RandomlyGeneratedDefault, "initial database user password")

	flag.Parse()

	log.Infof("starting stolonboot")

	var err error
	if *password == RandomlyGeneratedDefault {
		*password, err = randomPassword(DefaultPasswordLength)
		if err != nil {
			log.Error(trace.DebugReport(err))
			fmt.Printf("ERROR: %v\n", err.Error())
			os.Exit(255)
		}
	}

	err = bootCluster(*sentinels, *rpc, *password)
	if err != nil {
		log.Error(trace.DebugReport(err))
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(255)
	}

	log.Infof("Creating telegraf user in InfluxDB")
	influxDBClient, err := NewInfluxDBClient()
	if err != nil {
		log.Error(trace.DebugReport(err))
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(255)
	}

	err = influxDBClient.Setup()
	if err != nil {
		log.Error(trace.DebugReport(err))
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(255)
	}

	os.Exit(0)
}

func randomPassword(length int) (string, error) {
	data := make([]byte, length)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data)[:length], nil
}
