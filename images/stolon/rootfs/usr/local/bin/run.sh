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

function chmod_keys() {
	cp -R /etc/ssl/cluster-default /etc/ssl/cluster-default-postgres
	chown -R stolon /etc/ssl/cluster-default-postgres
	chmod 0600 /etc/ssl/cluster-default-postgres/default-server-key.pem
}


# TODO: find how to map ENV vars which is populated by k8s services for discovery to vars, which utilities uses
function setup_stolonctl() {
	announce_step 'Setup stolon CLI'

	_create_pg_pass "$STOLON_POSTGRES_SERVICE_HOST" \
				   "$STOLON_POSTGRES_SERVICE_PORT" \
				   "*" \
				   "$STOLONCTL_DB_USERNAME" \
				   "$STOLON_DB_PASSWORD"

	cat > /usr/local/bin/stolonctl-cluster <<'EOF'
	export STOLONCTL_DB_HOST=$STOLON_POSTGRES_SERVICE_HOST
	export STOLONCTL_DB_PORT=$STOLON_POSTGRES_SERVICE_PORT
	stolonctl "$@"
EOF
	chmod +x /usr/local/bin/stolonctl-cluster
}

function setup_stolonrpc() {
	announce_step 'Setup stolon RPC'

	_create_pg_pass "$STOLON_POSTGRES_SERVICE_HOST" \
				   "$STOLON_POSTGRES_SERVICE_PORT" \
				   "*" \
				   "$STOLONRPC_DB_USERNAME" \
				   "$STOLON_DB_PASSWORD"

	cat > /usr/local/bin/stolonrpc-cluster <<'EOF'
	export STOLONRPC_DB_HOST=$STOLON_POSTGRES_SERVICE_HOST
	export STOLONRPC_DB_PORT=$STOLON_POSTGRES_SERVICE_PORT
	stolonrpc "$@"
EOF
	chmod +x /usr/local/bin/stolonrpc-cluster
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

	if [[ ! -e /stolon-data ]]; then
		echo "stolon data doesn't exist, data won't be persistent!"
		mkdir /stolon-data
	fi
	chown stolon:stolon /stolon-data

	export STKEEPER_LISTEN_ADDRESS=$POD_IP
	export STKEEPER_PG_LISTEN_ADDRESS=$POD_IP

	su stolon -c "stolon-keeper --data-dir /stolon-data"
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

function launch_ctl() {
	announce_step 'Launching stolon CLI'

	stolonctl-cluster "$@"
}

function launch_rpc() {
	announce_step 'Launching stolon RPC'

	stolonrpc-cluster
}

function main() {
	announce_step 'Start'

	announce_step 'Dump environment variables'
	env

	setup_cluster_ca

	announce_step 'Select which component to start'
	if [[ "${KEEPER}" == "true" ]]; then
		chmod_keys
		launch_keeper
		exit 0
	fi

	if [[ "${SENTINEL}" == "true" ]]; then
		launch_sentinel
		exit 0
	fi

	if [[ "${PROXY}" == "true" ]]; then
		chmod_keys
		launch_proxy
		exit 0
	fi

	if [[ "${CTL}" == "true" ]]; then
		setup_stolonctl
		launch_ctl "$@"
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
