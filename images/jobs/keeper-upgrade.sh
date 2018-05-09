#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: entrypoint.sh
# Time-stamp: <2018-05-08 19:31:35>
# Copyright (C) 2018 Sergei Antipov
# Description:

set -o xtrace
set -o nounset
set -o errexit

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
    gosu stolon ${PGBINNEW}/pg_upgrade -d /stolon-data/postgres -D /stolon-data/postgres-new
    gosu stolon cp /stolon-data/postgres/pg_hba.conf /stolon-data/postgres-new/pg_hba.conf
    gosu stolon rsync -av /stolon-data/postgres/conf.d /stolon-data/postgres-new/
    gosu stolon rsync -av /stolon-data/postgres/postgresql-base.conf /stolon-data/postgres-new/
    gosu stolon mv /stolon-data/postgres /stolon-data/postgres-old
    gosu stolon mv /stolon-data/postgres-new /stolon-data/postgres
fi
