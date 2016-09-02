#!/usr/bin/env bash

# TODO: find how to map ENV vars which is populated by k8s services for discovery to vars, which utilities uses
function setup_stolonctl() {
	create_pg_pass "$STOLON_POSTGRES_SERVICE_HOST" \
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
	create_pg_pass "$STOLON_POSTGRES_SERVICE_HOST" \
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

function create_pg_pass() {
	local host=${1:-$STOLON_POSTGRES_SERVICE_HOST}
	local port=${2:-$STOLON_POSTGRES_SERVICE_PORT}
	local database=${3:-"postgres"}
	local username=${4:-"stolon"}
	local password=${5}

	echo "$host:$port:$database:$username:$password" > ~/.pgpass
	chmod 0600 ~/.pgpass
}

function launch_keeper() {
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
	export STSENTINEL_LISTEN_ADDRESS=$POD_IP
	stolon-sentinel
}

function launch_proxy() {
	export STPROXY_LISTEN_ADDRESS=$POD_IP
	stolon-proxy
}

function launch_ctl() {
	stolonctl-cluster "$@"
}

function launch_rpc() {
	stolonrpc-cluster
}

function main() {
	echo "start"
	# use hostname command to get our pod's ip until downward api are less racy (sometimes the podIP from downward api is empty)
	export POD_IP=$(hostname -i)
	env

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

	exit 1
}

main "$@"
