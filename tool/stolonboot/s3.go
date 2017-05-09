// Copyright 2017 Gravitational, Inc.
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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gravitational/trace"
	minio "github.com/minio/minio-go"
)

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	Bucket          string
}

func createBackupBucket(s3Config *S3Config) error {
	file, err := ioutil.ReadFile(ClusterCaPath)
	if err != nil {
		return trace.ConvertSystemError(err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(file)
	if !ok {
		return errors.New("failed to parse cluster root certificate")
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots}}

	client, err := minio.NewV2(s3Config.Endpoint, s3Config.AccessKeyID, s3Config.SecretAccessKey, true)
	if err != nil {
		return trace.Wrap(err)
	}

	time.Sleep(60 * time.Second)
	client.SetCustomTransport(tr)
	found, err := client.BucketExists(s3Config.Bucket)
	if err != nil {
		return trace.Wrap(err)
	}

	if !found {
		err = client.MakeBucket(s3Config.Bucket, DefaultS3BucketLocation)
		if err != nil {
			return trace.Wrap(err)
		}
	}
	return nil
}
