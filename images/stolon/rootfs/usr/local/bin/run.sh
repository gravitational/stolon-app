#!/usr/bin/env bash

set -e

export STORE_ENDPOINTS=$NODE_NAME:2379
export STSENTINEL_STORE_ENDPOINTS=$STORE_ENDPOINTS
export STKEEPER_STORE_ENDPOINTS=$STORE_ENDPOINTS
export STPROXY_STORE_ENDPOINTS=$STORE_ENDPOINTS
export STOLONCTL_STORE_ENDPOINTS=$STORE_ENDPOINTS

function die() {
    echo "ERROR: $*" >&2
    exit 1
}

function announce_step() {
    echo
    echo "===> $*"
    echo
}

function setup_cluster_ca() {
	announce_step 'Setup cluster CA'

	if [ -f /usr/local/bin/kubectl ]; then
		mkdir -p /usr/share/ca-certificates/extra
		kubectl get secret cluster-ca
		if [ $? -eq 0 ]; then
			kubectl get secret cluster-ca -o yaml|grep ca.pem|awk '{print $2}'|base64 -d > /usr/local/share/ca-certificates/cluster.crt
			update-ca-certificates
		fi
	fi
}

function setup_stolonrpc() {
	announce_step 'Setup stolon RPC'

	_create_pg_pass "$STOLON_POSTGRES_SERVICE_HOST" \
				   "$STOLON_POSTGRES_SERVICE_PORT" \
				   "*" \
				   "$STOLONRPC_DB_USERNAME" \
				   "$STOLON_DB_PASSWORD"

	cat > /home/stolon/bin/stolonrpc-cluster <<EOF
	export STOLONRPC_DB_HOST=$STOLON_POSTGRES_SERVICE_HOST
	export STOLONRPC_DB_PORT=$STOLON_POSTGRES_SERVICE_PORT
	stolonrpc "$@"
EOF
	chmod +x /home/stolon/bin/stolonrpc-cluster
}

function _create_pg_pass() {
	local host=${1:-$STOLON_POSTGRES_SERVICE_HOST}
	local port=${2:-$STOLON_POSTGRES_SERVICE_PORT}
	local database=${3:-"postgres"}
	local username=${4:-"stolon"}
	local password=${5}

	echo "$host:$port:$database:$username:$password" > ~/.pgpass
	chmod 0600 ~/.pgpass
}

function launch_keeper() {
	announce_step 'Launching stolon keeper'

	export STKEEPER_LISTEN_ADDRESS=$POD_IP
	export STKEEPER_PG_LISTEN_ADDRESS=$POD_IP

	stolon-keeper --data-dir /stolon-data
}

function launch_sentinel() {
	announce_step 'Launching stolon sentinel'

	export STSENTINEL_LISTEN_ADDRESS=$POD_IP
	stolon-sentinel
}

function launch_proxy() {
	announce_step 'Launching stolon proxy'

	export STPROXY_LISTEN_ADDRESS=$POD_IP
	stolon-proxy
}

function launch_rpc() {
	announce_step 'Launching stolon RPC'

	/home/stolon/bin/stolonrpc-cluster
}

function main() {
	announce_step 'Start'

#	setup_cluster_ca

	announce_step 'Select which component to start'
	if [[ "${KEEPER}" == "true" ]]; then
		launch_keeper
		exit 0
	fi

	if [[ "${SENTINEL}" == "true" ]]; then
		launch_sentinel
		exit 0
	fi

	if [[ "${PROXY}" == "true" ]]; then
		launch_proxy
		exit 0
	fi

	if [[ "${RPC}" == "true" ]]; then
		setup_stolonrpc
		launch_rpc
		exit 0
	fi

	die 'Nothing is selected.'
}

main "$@"
