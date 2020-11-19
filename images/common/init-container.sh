#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: fix-permissions.sh
# Time-stamp: <2018-11-09 11:55:07>
# Copyright (C) 2018 Gravitational Inc.
# Description:

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

if [ -d /etc/secrets/cluster-default ]
then
    cp -R /etc/secrets/cluster-default /home/stolon/secrets/cluster-default
    chown -R stolon:stolon /home/stolon/secrets/cluster-default
	chmod 0600 /home/stolon/secrets/cluster-default/default-server-key.pem
fi

if [ -d /etc/secrets/cluster-ca ]
then
    cp -R /etc/secrets/cluster-ca /home/stolon/secrets/cluster-ca
    chown -R stolon:stolon /home/stolon/secrets/cluster-ca
	chmod 0600 /home/stolon/secrets/cluster-ca/ca.pem
fi

if [ -d /var/state ]
then
    cp -R /var/state /home/stolon/secrets/etcd
    chown -R stolon:stolon /home/stolon/secrets/etcd
fi

if [ -d /stolon-data ]; then
  if [[ ! -f /stolon-data/dummy.file ]]; then
      fallocate -l 300MB /stolon-data/dummy.file
  fi
  chown -R stolon:stolon /stolon-data
fi

if [ -n "$INIT_PGBOUNCER"]
then
	# Generate client pgbouncer certificate for connection to PostgreSQL
	openssl req -new -nodes -text -out /home/stolon/secrets/pgbouncer.csr -keyout /home/stolon/secrets/pgbouncer-key.pem -subj "/CN=pgbouncer"
	openssl x509 -req -in /home/stolon/secrets/pgbouncer.csr -text -CA /home/stolon/secrets/cluster-ca/ca.pem -CAkey /home/stolon/secrets/cluster-ca/ca-key -CAcreateserial -out /home/stolon/secrets/pgbouncer.pem

	chown stolon:stolon -R /home/stolon/secrets/pgbouncer*
	chmod 0400 /home/stolon/secrets/pgbouncer-key.pem
fi
