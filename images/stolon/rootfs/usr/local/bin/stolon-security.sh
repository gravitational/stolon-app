#!/usr/bin/env bash
# -*- mode: bash; -*-

# File: stolon-security.sh
#
# Copyright 2018 Gravitational, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o pipefail
source /usr/local/bin/common_lib.sh

function create_pg_passfile() {
	echo "${PG_SERVICE}:${PG_PORT:-5432}:*:${PG_USERNAME}:${PG_PASSWORD}" > ~/.pgpass
	chmod 0600 ~/.pgpass
}

function disable_connections_to_templates() {
    psql -h${PG_SERVICE} -U${PG_USERNAME} postgres -c "UPDATE pg_database SET datallowconn = false WHERE datistemplate = true"
}

function main() {
    echo "Checking environment variables..."
    env_check_err PG_USERNAME
    env_check_err PG_SERVICE
    env_check_err PG_PASSWORD

    create_pg_passfile
    disable_connections_to_templates
}

main "$@"
