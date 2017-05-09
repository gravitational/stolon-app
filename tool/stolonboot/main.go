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
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	RandomlyGeneratedDefault = "<randomly generated>"
	ClusterCaPath            = "/etc/ssl/cluster-ca/ca.pem"
	DefaultPasswordLength    = 20
	DefaultS3Endpoint        = "pithos.default.svc"
	DefaultS3Bucket          = "stolon-backup"
	DefaultS3BucketLocation  = "G1"
)

func main() {
	sentinels := flag.Int("sentinels", 2, "number of sentinels")
	rpc := flag.Int("rpc", 1, "number of RPC")
	password := flag.String("password", RandomlyGeneratedDefault, "initial database user password")
	s3AccessKeyID := flag.String("access-key", "", "S3 access key")
	s3SecretAccessKey := flag.String("secret-key", "", "S3 secret key")
	s3Endpoint := flag.String("endpoint", DefaultS3Endpoint, "S3 endpoint address")
	s3Bucket := flag.String("bucket", DefaultS3Bucket, "S3 Bucket name")

	flag.Parse()

	log.Infof("starting stolonboot")

	var err error
	if *password == RandomlyGeneratedDefault {
		*password, err = randomPassword(DefaultPasswordLength)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	if *s3AccessKeyID != "" && *s3SecretAccessKey != "" {
		log.Infof("creating bucket in pithos for backups")
		err = createBackupBucket(&S3Config{
			AccessKeyID:     *s3AccessKeyID,
			SecretAccessKey: *s3SecretAccessKey,
			Endpoint:        *s3Endpoint,
			Bucket:          *s3Bucket,
		})
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	err = bootCluster(*sentinels, *rpc, *password)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
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
