#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: run.sh
# Time-stamp: <2017-03-20 23:07:29>
# Copyright (C) 2017 Sergei Antipov
# Description: Entrypoint for pgweb application

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"

/usr/local/bin/pgweb_linux_amd64 --sessions --bind 0.0.0.0 --url postgres://${STOLON_USERNAME}:$(cat $STOLON_PASSWORDFILE)@stolon-postgres.default.svc:5432/postgres?sslmode=disable
