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
    rig delete deployments/stolon-utils --force
    kubectl create -f /var/lib/gravity/resources/delete-backup-bucket.yaml

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
    rig upsert -f /var/lib/gravity/resources/alerts.yaml --debug

    if [ $(kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
    then
        kubectl scale --replicas=3 deployment stolon-sentinel
    fi

    echo "Checking status"
    rig status $RIG_CHANGESET --retry-attempts=120 --retry-period=1s --debug
    echo "Freezing"
    rig freeze

elif [ $1 = "rollback" ]; then
    echo "Reverting changeset $RIG_CHANGESET"
    rig revert

    kubectl create -f /var/lib/gravity/resources/create-backup-bucket.yaml
    if [ $(kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
    then
        kubectl scale --replicas=3 deployment stolon-sentinel
    fi
else
    echo "Missing argument, should be either 'update' or 'rollback'"
fi
