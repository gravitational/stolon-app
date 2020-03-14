#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: fix-permissions.sh
# Copyright (C) 2019 Gravitational Inc.
# Description:

set -o errexit

if [ -d /etc/secrets/cluster-default ]
then
    cp -R /etc/secrets/cluster-default /home/stolon/secrets/cluster-default
    chown -R stolon:stolon /home/stolon/secrets/cluster-default
    chmod 0600 /home/stolon/secrets/cluster-default/default-server-key.pem
fi

if [ -d /etc/secrets/cluster-ca ]
then
    cp -LR /etc/secrets/cluster-ca /home/stolon/secrets/cluster-ca
    chown -R stolon:stolon /home/stolon/secrets/cluster-ca
	  chmod 0600 /home/stolon/secrets/cluster-ca/ca.pem
fi

# Generate client pgbouncer certificate for connection to PostgreSQL
openssl req -new -nodes -text -out /home/stolon/secrets/pgbouncer.csr -keyout /home/stolon/secrets/pgbouncer-key.pem -subj "/CN=pgbouncer"
openssl x509 -req -in /home/stolon/secrets/pgbouncer.csr -text -CA /home/stolon/secrets/cluster-ca/ca.pem -CAkey /home/stolon/secrets/cluster-ca/ca-key -CAcreateserial -out /home/stolon/secrets/pgbouncer.pem

chown stolon:stolon -R /home/stolon/secrets/*
chmod 0400 /home/stolon/secrets/pgbouncer-key.pem
