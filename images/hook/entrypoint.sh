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
    rig delete rc/stolon-proxy --force
    rig delete rc/stolon-rpc --force
    rig delete rc/stolon-sentinel --force

    # wait for keeper pods to go away
    while kubectl get pods --show-all|grep -q stolon-keeper
    do
    echo "waiting for keeper pods to go away..."
        sleep 5
    done

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
