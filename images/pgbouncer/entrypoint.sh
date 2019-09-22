#!/usr/bin/env bash
# -*- mode: bash; -*-

# File: entrypoint.sh
# Copyright (C) 2019 Gravitational Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Description: Entrypoint file for pgbouncer container

set -o errexit
set -o pipefail
source /common_lib.sh

CONF_DIR=/home/stolon/pgbouncer
CONF_DIR_MOUNTED=/etc/pgbouncer

mkdir -p ${CONF_DIR}
cp -L ${CONF_DIR_MOUNTED}/* ${CONF_DIR}

function create_pg_passfile() {
	echo "${PG_SERVICE}:${PG_PORT:-5432}:*:${PG_USERNAME}:${PG_PASSWORD}" > ~/.pgpass
	echo "localhost:6432:*:${PG_USERNAME}:${PG_PASSWORD}" >> ~/.pgpass
	chmod 0600 ~/.pgpass
}

function create_users_file() {
    password_md5=$(echo -n "${PG_PASSWORD}${PG_USERNAME}" | md5sum | cut -d " " -f 1)
    echo "\"${PG_USERNAME}\" \"md5${password_md5}\"" > ${CONF_DIR?}/users.txt
}

function fill_config_file() {
    sed -i "s/DEFAULT_POOL_SIZE/${DEFAULT_POOL_SIZE:-20}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/MAX_CLIENT_CONN/${MAX_CLIENT_CONN:-100}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/MAX_DB_CONNECTIONS/${MAX_DB_CONNECTIONS:-0}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/MIN_POOL_SIZE/${MIN_POOL_SIZE:-0}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/POOL_MODE/${POOL_MODE:-session}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/RESERVE_POOL_SIZE/${RESERVE_POOL_SIZE:-0}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/RESERVE_POOL_TIMEOUT/${RESERVE_POOL_TIMEOUT:-5}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/QUERY_TIMEOUT/${QUERY_TIMEOUT:-0}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/IGNORE_STARTUP_PARAMETERS/${IGNORE_STARTUP_PARAMETERS:-extra_float_digits}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/PG_PORT/${PG_PORT:-5432}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/PG_SERVICE/${PG_SERVICE}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/CLIENT_TLS_SSLMODE/${CLIENT_TLS_SSLMODE:-require}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s#CLIENT_TLS_KEY_FILE#${CLIENT_TLS_KEY_FILE:-/home/stolon/secrets/cluster-default/default-server-key.pem}#g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s#CLIENT_TLS_CERT_FILE#${CLIENT_TLS_CERT_FILE:-/home/stolon/secrets/cluster-default/default-server.pem}#g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s#CLIENT_TLS_CA_FILE#${CLIENT_TLS_CA_FILE:-/home/stolon/secrets/cluster-ca/ca.pem}#g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/CLIENT_TLS_PROTOCOLS/${CLIENT_TLS_PROTOCOLS:-secure}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/SERVER_TLS_SSLMODE/${SERVER_TLS_SSLMODE:-require}/g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s#SERVER_TLS_KEY_FILE#${SERVER_TLS_KEY_FILE:-/home/stolon/secrets/cluster-default/default-server-key.pem}#g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s#SERVER_TLS_CERT_FILE#${SERVER_TLS_CERT_FILE:-/home/stolon/secrets/cluster-default/default-server.pem}#g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s#SERVER_TLS_CA_FILE#${SERVER_TLS_CA_FILE:-/home/stolon/secrets/cluster-ca/ca.pem}#g" ${CONF_DIR?}/pgbouncer.ini
    sed -i "s/SERVER_TLS_PROTOCOLS/${SERVER_TLS_PROTOCOLS:-secure}/g" ${CONF_DIR?}/pgbouncer.ini
}

echo "Checking environment variables..."
env_check_err PG_USERNAME
env_check_err PG_SERVICE
env_check_err PG_PASSWORD

echo "Checking non-default configuration parameters..."
env_check_info POOL_MODE "Setting POOL_MODE to ${POOL_MODE}"
env_check_info DEFAULT_POOL_SIZE "Setting DEFAULT_POOL_SIZE to ${DEFAULT_POOL_SIZE}"
env_check_info MAX_CLIENT_CONN "Setting MAX_CLIENT_CONN to ${MAX_CLIENT_CONN}"
env_check_info MAX_DB_CONNECTIONS "Setting MAX_DB_CONNECTIONS to ${MAX_DB_CONNECTIONS}"
env_check_info MIN_POOL_SIZE "Setting MIN_POOL_SIZE to ${MIN_POOL_SIZE}"
env_check_info RESERVE_POOL_SIZE "Setting RESERVE_POOL_SIZE to ${RESERVE_POOL_SIZE}"
env_check_info RESERVE_POOL_TIMEOUT "Setting RESERVE_POOL_TIMEOUT to ${RESERVE_POOL_TIMEOUT}"
env_check_info QUERY_TIMEOUT "Setting QUERY_TIMEOUT to ${QUERY_TIMEOUT}"
env_check_info IGNORE_STARTUP_PARAMETERS "Setting IGNORE_STARTUP_PARAMETERS to ${IGNORE_STARTUP_PARAMETERS}"
env_check_info PG_SERVICE "Setting PG_SERVICE to ${PG_SERVICE}"
env_check_info SERVER_TLS_CERT_FILE "Setting SERVER_TLS_CERT_FILE to ${SERVER_TLS_CERT_FILE}"
env_check_info SERVER_TLS_KEY_FILE "Setting SERVER_TLS_KEY_FILE to ${SERVER_TLS_KEY_FILE}"


echo "Filling configuration files..."
fill_config_file
create_users_file
create_pg_passfile

echo "Starting pgbouncer..."
/usr/sbin/pgbouncer ${CONF_DIR?}/pgbouncer.ini
