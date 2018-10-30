#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: entrypoint.sh
# Time-stamp: <2018-10-30 13:23:15>
# Copyright (C) 2018 Gravitational Inc
# Description:

set -o xtrace
set -o nounset
set -o errexit

# delete old versions of data directories
for dir in $(ls -d /stolon-data/postgres-old-* 2>/dev/null)
do
    if [ -f $dir/PG_VERSION ] && [ $(cat $dir/PG_VERSION) != $PG_VERSION_OLD ] && [ $(cat $dir/PG_VERSION) != $PG_VERSION_NEW ]
    then
        rm -rf $dir
    fi
done

sed -i -e "s#\(ssl_cert_file = '\).*#\1/home/stolon/secrets/cluster-default/default-server.pem'#" /stolon-data/postgres/postgresql.conf
sed -i -e "s#\(ssl_key_file = '\).*#\1/home/stolon/secrets/cluster-default/default-server-key.pem'#" /stolon-data/postgres/postgresql.conf
sed -i -e "s#\(ssl_ca_file = '\).*#\1/home/stolon/secrets/cluster-ca/ca.pem'#" /stolon-data/postgres/postgresql.conf
if [ $(cat /stolon-data/postgres/PG_VERSION) != $PG_VERSION_NEW ]
then
    mkdir -p /stolon-data/postgres-new /stolon-data/upgrade-state
    cd /stolon-data/upgrade-state
    if [ ! -f /stolon-data/postgres-new/postgresql.conf ]; then
        PGDATA=/stolon-data/postgres-new ${PGBINNEW}/initdb
    fi
    if [ ! -f /stolon-data/upgrade-state/delete_old_cluster.sh ]; then
        ${PGBINNEW}/pg_upgrade -d /stolon-data/postgres -D /stolon-data/postgres-new
        cp /stolon-data/postgres/pg_hba.conf /stolon-data/postgres-new/pg_hba.conf
        rsync -av /stolon-data/postgres/conf.d /stolon-data/postgres-new/
        rsync -av /stolon-data/postgres/postgresql-base.conf /stolon-data/postgres-new/
        if [ -d /stolon-data/postgres-old-$PG_VERSION_OLD ]; then rm -rf /stolon-data/postgres-old-$PG_VERSION_OLD; fi
        mv /stolon-data/postgres /stolon-data/postgres-old-$PG_VERSION_OLD
        mv /stolon-data/postgres-new /stolon-data/postgres
    fi
fi
