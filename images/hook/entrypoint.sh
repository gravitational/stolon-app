#!/bin/sh
set -e

echo "Assuming changeset from the envrionment: $RIG_CHANGESET"
# note that rig does not take explicit changeset ID
# taking it from the environment variables
if [ $1 = "update" ]; then
    echo "Starting update, changeset: $RIG_CHANGESET"
    rig cs delete --force -c cs/$RIG_CHANGESET

    echo "Delete old resources"
    rig delete ds/stolon-keeper --force
    rig delete deployments/stolon-proxy --force
    rig delete deployments/stolon-rpc --force
    rig delete deployments/stolon-sentinel --force
    kubectl delete pods -l name=stolon-rpc --ignore-not-found
    kubectl delete pods -l name=stolon-sentinel --ignore-not-found
    kubectl delete pods -l name=stolon-utils --ignore-not-found

    # wait for keeper pods to go away
    while kubectl get pods --show-all|grep -q stolon-keeper
    do
    echo "waiting for keeper pods to go away..."
        sleep 5
    done

	etcdctl get /stolon/cluster/kube-stolon/clusterdata > clusterdata.json
	sed -i 's/"archive_mode":"on"/"archive_mode":"off"/g' clusterdata.json
	etcdctl set /stolon/cluster/kube-stolon/clusterdata "$(cat clusterdata.json)"

    echo "Creating or updating resources"
    rig upsert -f /var/lib/gravity/resources/keeper.yaml --debug
    rig upsert -f /var/lib/gravity/resources/rpc.yaml --debug
    rig upsert -f /var/lib/gravity/resources/sentinel.yaml --debug
    rig upsert -f /var/lib/gravity/resources/utils.yaml --debug

    echo "Checking status"
    rig status $RIG_CHANGESET --retry-attempts=120 --retry-period=1s --debug
    echo "Freezing"
    rig freeze

elif [ $1 = "rollback" ]; then
    echo "Reverting changeset $RIG_CHANGESET"
    rig revert
else
    echo "Missing argument, should be either 'update' or 'rollback'"
fi
