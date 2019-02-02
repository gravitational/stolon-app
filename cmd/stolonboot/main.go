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

	log "github.com/sirupsen/logrus"
)

const RandomlyGeneratedDefault = "<randomly generated>"
const DefaultPasswordLength = 20

func main() {
	sentinels := flag.Int("sentinels", 2, "number of sentinels")
	password := flag.String("password", RandomlyGeneratedDefault, "initial database user password")

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

	err = bootCluster(*sentinels, *password)
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
