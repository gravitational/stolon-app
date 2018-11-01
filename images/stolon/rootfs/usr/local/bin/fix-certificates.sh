#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: fix-permissions.sh
# Time-stamp: <2018-10-17 22:40:08>
# Copyright (C) 2018 Gravitational Inc.
# Description:

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

if [ -d /etc/secrets/cluster-default ]
then
    cp -R /etc/secrets/cluster-default /home/stolon/secrets/cluster-default
    chown -R stolon /home/stolon/secrets/cluster-default
	chmod 0600 /home/stolon/secrets/cluster-default/default-server-key.pem
fi

if [ -d /etc/secrets/cluster-ca ]
then
    cp -R /etc/secrets/cluster-ca /home/stolon/secrets/cluster-ca
    chown -R stolon /home/stolon/secrets/cluster-ca
	chmod 0600 /home/stolon/secrets/cluster-ca/ca.pem
fi

if [ -d /etc/secrets/etcd ]
then
    cp -R /etc/secrets/etcd /home/stolon/secrets/etcd
    chown -R stolon /home/stolon/secrets/etcd
fi
