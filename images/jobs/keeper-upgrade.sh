#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: entrypoint.sh
# Time-stamp: <2019-04-03 22:30:48>
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

CURRENT_VERSION=$(cat /stolon-data/postgres/PG_VERSION)
if [ $CURRENT_VERSION == $PG_VERSION_NEW ]; then exit 0; fi

# fix permissions for ssl
cp -R /etc/ssl/cluster-default /etc/ssl/cluster-default-postgres
chown -R stolon /etc/ssl/cluster-default-postgres
chmod 0600 /etc/ssl/cluster-default-postgres/default-server-key.pem

gosu stolon mkdir -p /stolon-data/postgres-new /stolon-data/upgrade-state
cd /stolon-data/upgrade-state
if [ ! -f /stolon-data/postgres-new/postgresql.conf ]; then
    PGDATA=/stolon-data/postgres-new gosu stolon ${PGBINNEW}/initdb
fi
if [ ! -f /stolon-data/upgrade-state/delete_old_cluster.sh ]; then
    sed '/idle_in_transaction_session_timeout/d' -i /stolon-data/postgres/postgresql.conf
    gosu stolon ${PGBINNEW}/pg_upgrade -d /stolon-data/postgres -D /stolon-data/postgres-new
    gosu stolon cp /stolon-data/postgres/pg_hba.conf /stolon-data/postgres-new/pg_hba.conf
    gosu stolon rsync -av /stolon-data/postgres/conf.d /stolon-data/postgres-new/
    gosu stolon rsync -av /stolon-data/postgres/postgresql-base.conf /stolon-data/postgres-new/
    if [ -d /stolon-data/postgres-old-$PG_VERSION_OLD ]; then rm -rf /stolon-data/postgres-old-$PG_VERSION_OLD; fi
    gosu stolon mv /stolon-data/postgres /stolon-data/postgres-old-$PG_VERSION_OLD
    gosu stolon mv /stolon-data/postgres-new /stolon-data/postgres
fi
