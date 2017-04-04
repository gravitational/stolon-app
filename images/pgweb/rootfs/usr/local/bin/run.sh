#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: run.sh
# Time-stamp: <2017-04-04 15:22:07>
# Copyright (C) 2017 Sergei Antipov
# Description: Entrypoint for pgweb application

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"

printenv
/usr/local/bin/pgweb_linux_amd64 --sessions --bind 0.0.0.0 --host="stolon-postgres.default.svc" --port=5432 --user="${STOLON_USERNAME}" --pass="$(cat $STOLON_PASSWORDFILE)" --db="postgres" --ssl="disable"
