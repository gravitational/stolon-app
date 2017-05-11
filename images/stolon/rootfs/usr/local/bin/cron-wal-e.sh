#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: cron-wal-e.sh
# Time-stamp: <2017-05-11 21:39:37>
# Copyright (C) 2017 Sergei Antipov
# Description: Cron script to make basebackups via WAL-e

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"


master_id=$(stolonctl cluster status $CLUSTER_NAME --store-endpoints="$NODE_NAME:2379" --store-backend=${STORE_BACKEND} --store-cert=${STORE_CERT} --store-key=${STORE_KEY} --store-cacert=${STORE_CACERT} | grep master | awk '{print $1}')
keeper_id=$(cat /stolon-data/pgid)

if [[ $master_id = $keeper_id ]]
then
    PGHOST=localhost PGUSER=$STOLON_USERNAME envdir /etc/wal-e.d/env wal-e backup-push /stolon-data/postgres/ >> /var/log/cron.log 2>&1
    echo "Basebackup is complete."
else
    echo "Cronjob disabled. Not a master node."
fi
