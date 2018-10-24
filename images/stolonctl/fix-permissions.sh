#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: fix-permissions.sh
# Time-stamp: <2018-10-22 15:29:47>
# Copyright (C) 2018 Gravitational Inc
# Description:

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

mkdir -p /stolon/backup
chown stolon:stolon /stolon/backup
