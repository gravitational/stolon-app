#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: entrypoint.sh
# Time-stamp: <2018-03-21 15:27:07>
# Copyright (C) 2018 Sergei Antipov
# Description: Entrypoint for PostgreSQL upgarder tool

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

export ETCD_ENDPOINTS=${NODE_NAME}:2379

while true; do sleep 30; done;
