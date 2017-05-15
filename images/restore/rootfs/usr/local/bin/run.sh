#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: setup-wal-e.sh
# Time-stamp: <2017-05-15 17:01:12>
# Copyright (C) 2017 Gravitational Inc
# Description: Setup WAL-e envdir

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

mkdir -p /etc/wal-e.d/env
echo ${AWS_SECRET_ACCESS_KEY} > /etc/wal-e.d/env/AWS_SECRET_ACCESS_KEY
echo ${AWS_ACCESS_KEY_ID} > /etc/wal-e.d/env/AWS_ACCESS_KEY_ID
echo "https+path://${PITHOS_S3_ENDPOINT}:443" > /etc/wal-e.d/env/WALE_S3_ENDPOINT
echo "s3://${PITHOS_BACKUP_BUCKET}/" > /etc/wal-e.d/env/WALE_S3_PREFIX
chown -R root:postgres /etc/wal-e.d

echo "*:5432:*:${STOLON_USERNAME}:${STOLON_PASSWORD}" > ~/.pgpass
chmod 0600 ~/.pgpass

export ETCDCTL_CA_FILE=/etc/etcd/secrets/root.cert
export ETCDCTL_CERT_FILE=/etc/etcd/secrets/etcd.cert
export ETCDCTL_KEY_FILE=/etc/etcd/secrets/etcd.key

while true; do sleep 30; done
