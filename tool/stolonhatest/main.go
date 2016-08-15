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
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/trace"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

var (
	ErrCantParseConfig      = errors.New("Can't parse config")
	ErrDBConnAttemptsFailed = errors.New("All attempts failed")
)

var (
	DBDriverName = "postgres"
)

type Config struct {
	LogLevel            string `envconfig:"LOG_LEVEL"`
	DBUsername          string `envconfig:"DB_USERNAME"`
	DBPassword          string `envconfig:"DB_PASSWORD"`
	DBHost              string `envconfig:"STOLON_POSTGRES_SERVICE_HOST"`
	DBPort              string `envconfig:"STOLON_POSTGRES_SERVICE_PORT"`
	DBName              string `envconfig:"DB_NAME"`
	KeeperLabelSelector string `envconfig:"STOLON_KEEPER_LABEL_SELECTOR"`
	ClusterName         string `envconfig:"STOLON_CLUSTER_NAME"`
}

func GetConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, trace.Wrap(ErrCantParseConfig)
	}
	return &config, nil
}

type Database struct {
	*sql.DB
}

func GetDatabase(dsn string) (*Database, error) {
	db, err := sql.Open(DBDriverName, dsn)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return &Database{db}, nil
}

func (d *Database) Connect() error {
	var dbError error

	maxAttempts := 30
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		dbError = d.Ping()
		if dbError == nil {
			break
		}
		nextAttemptWait := time.Duration(attempt) * time.Second

		log.Errorf("Attempt %v: could not establish a connection with the database. Wait for %v.", attempt, nextAttemptWait)
		time.Sleep(nextAttemptWait)
	}

	if dbError != nil {
		return trace.Wrap(ErrDBConnAttemptsFailed)
	}
	return nil
}

func (d *Database) Close() error {
	if err := d.Close(); err != nil {
		return trace.Wrap(err, "Can't close database")
	}
	return nil
}

func setupLogging(level string) error {
	lvl := strings.ToLower(level)

	if lvl == "debug" {
		trace.SetDebug(true)
	}

	sev, err := log.ParseLevel(lvl)
	if err != nil {
		return err
	}
	log.SetLevel(sev)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	return nil
}

// Integration test for stolon in k8s.
// Logic is straightforward:
//	* Connect to DB using stolon-proxy IP and port
//	* Create DB schema
//	* Write test data
//	* Find and kill pod with stolon's keeper in master state
//	* Drop stolon's labels to prevent rescheduling of container on the same node again
//	* Execute test query and compare with test data
// If the test data were returned, replication works.
func main() {
	c, err := GetConfig()
	if err != nil {
		trace.Wrap(err)
	}

	if err = setupLogging(c.LogLevel); err != nil {
		trace.Wrap(err)
	}
	log.Infof("Start with config: %+v", c)

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		DBDriverName,
		c.DBUsername,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
	db, err := GetDatabase(dsn)
	if err != nil {
		trace.Wrap(err)
	}

	log.Debugf("Connecting to database at %s:%s", c.DBHost, c.DBPort)
	if err = db.Connect(); err != nil {
		trace.Wrap(err)
	}

	schema := `CREATE TABLE IF NOT EXISTS test (
                    id INT PRIMARY KEY NOT NULL,
                    value TEXT NOT NULL);`
	log.Debugf("Applying schema: %v", schema)
	_, err = db.Exec(schema)
	if err != nil {
		log.Errorf("Fail to apply schema: %v\n, error: %v", schema, err)
		trace.Wrap(err)
	}

	log.Info("Create test record ...")

	query := `SELECT value FROM test WHERE id=1`
	log.Debugf("Execute query: %v", query)
	var expectedTestString string
	err = db.QueryRow(query).Scan(&expectedTestString)
	if err == sql.ErrNoRows {
		insert := `INSERT INTO test VALUES (1, 'value1');`
		log.Debugf("No results found, trying to create, execute statement: %v", insert)
		_, err = db.Exec(insert)
		if err != nil {
			log.Errorf("Fail to execute statement: %v\n, error: %v", insert, err)
			trace.Wrap(err)
		}
	}
	if err != nil {
		trace.Wrap(err)
	}

	if expectedTestString == "" {
		log.Debugf("Now we should get it, Execute query: %v", query)
		err = db.QueryRow(query).Scan(&expectedTestString)
		if err != nil {
			trace.Wrap(err)
		}
	}
	log.Infof("Expected value: %v", expectedTestString)

	log.Infof("Killing PG master ...")
	log.Debugf("Finding master keeper IP in stolon cluster %v ...", c.ClusterName)
	kIP, err := GetMasterKeeperIP(c.ClusterName)
	if err != nil {
		trace.Wrap(err)
	}
	log.Debugf("Master keeper is listening on: %v", kIP)

	log.Debugf("Finding pod with IP %v from pods with %v label ...", kIP, c.KeeperLabelSelector)
	podName, err := GetPodName(c.KeeperLabelSelector, kIP)
	if err != nil {
		trace.Wrap(err)
	}

	log.Infof("Killing pod: %v", podName)
	if err = DeletePod(podName, true); err != nil {
		trace.Wrap(err)
	}

	log.Infof("Remove label %v from node %v", c.KeeperLabelSelector, kIP)
	selectorParts := strings.Split(c.KeeperLabelSelector, "=")
	if err = Label("node", kIP, selectorParts[0]+"-"); err != nil {
		trace.Wrap(err)
	}

	log.Info("Consistency test ...")
	log.Debugf("Execute query: %v", query)
	var actualTestString string
	err = db.QueryRow(query).Scan(&actualTestString)
	if err != nil {
		trace.Wrap(err)
	}
	log.Debugf("Got string: %v", actualTestString)

	if actualTestString != expectedTestString {
		log.Fatalf("Actual: %v. Expected: %v.", actualTestString, expectedTestString)
	}
	log.Info("Success!")

	log.Info("Cleanup ...")
	drop := `DROP TABLE test;`
	log.Debugf("Execute statement: %v", drop)
	_, err = db.Exec(drop)
	if err != nil {
		log.Errorf("Fail to execute statement: %v\n, error: %v", drop, err)
		trace.Wrap(err)
	}

	log.Infof("Put label %v to node %v back", c.KeeperLabelSelector, kIP)
	if err = Label("node", kIP, c.KeeperLabelSelector); err != nil {
		trace.Wrap(err)
	}

	log.Info("Done")
}
