#!/usr/bin/env bash

function setup() {
	# use hostname command to get our pod's ip until downward api are less racy (sometimes the podIP from downward api is empty)
	export POD_IP=$(hostname -i)
}

function setup_stolonctl_cluster() {
	echo "$STOLON_POSTGRES_SERVICE_HOST:$STOLON_POSTGRES_SERVICE_PORT:*:$STOLONCTL_DB_USERNAME:$(cat /etc/secrets/stolon/password)" > ~/.pgpass
	chmod 0600 ~/.pgpass

	cat > /usr/local/bin/stolonctl-cluster <<'EOF'
	export STOLONCTL_DB_HOST=$STOLON_POSTGRES_SERVICE_HOST
	export STOLONCTL_DB_PORT=$STOLON_POSTGRES_SERVICE_PORT
	stolonctl "$@"
EOF
	chmod +x /usr/local/bin/stolonctl-cluster
}

function check_data() {
	if [[ ! -e /stolon-data ]]; then
		echo "stolon data doesn't exist, data won't be persistent!"
		mkdir /stolon-data
	fi
	chown stolon:stolon /stolon-data
}

function launch_keeper() {
	check_data
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
	/usr/local/bin/stolonctl-cluster "$@"
}

function main() {
	echo "start"
	setup
	setup_stolonctl_cluster
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
		launch_ctl "$@"
		exit 0
	fi

	exit 1
}

main "$@"
