#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: entrypoint.sh
# Time-stamp: <2019-05-20 13:59:04>
# Copyright (C) 2018 Sergei Antipov
# Description: Entrypoint for PostgreSQL upgarder tool

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

export ETCD_ENDPOINTS=${NODE_NAME}:2379

/usr/bin/stolonctl server
