#!/usr/bin/env bash
# -*- mode: sh; -*-

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

export ETCD_ENDPOINTS=${NODE_NAME}:2379

/usr/bin/stolonctl server
