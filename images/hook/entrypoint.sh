#!/bin/bash
set -e

echo "Assuming changeset from the envrionment: $RIG_CHANGESET"
# note that rig does not take explicit changeset ID
# taking it from the environment variables
if [ $1 = "update" ]; then
    echo "---> Checking: $RIG_CHANGESET"
    if rig status $RIG_CHANGESET --retry-attempts=1 --retry-period=1s; then exit 0; fi

    echo "Starting update, changeset: $RIG_CHANGESET"
    rig cs delete --force -c cs/$RIG_CHANGESET

    echo "Delete old resources"
    rig delete ds/stolon-keeper --force
    rig delete deployments/stolon-proxy --force
    rig delete deployments/stolon-rpc --force
    rig delete deployments/stolon-sentinel --force
    rig delete deployments/stolon-utils --force
    rig delete deployments/stolonctl --force
    rig delete configmaps/stolon-telegraf --force
    rig delete configmaps/stolon-telegraf-node --force

    echo "Creating or updating resources"
    rig upsert -f /var/lib/gravity/resources/security.yaml --debug
    rig upsert -f /var/lib/gravity/resources/telegraf.yaml --debug
    rig upsert -f /var/lib/gravity/resources/keeper.yaml --debug
    rig upsert -f /var/lib/gravity/resources/sentinel.yaml --debug
    rig upsert -f /var/lib/gravity/resources/utils.yaml --debug
    rig upsert -f /var/lib/gravity/resources/alerts.yaml --debug
    rig upsert -f /var/lib/gravity/resources/stolonctl.yaml --debug

    if [ $(kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
    then
        sed -i 's/replicas: 1/replicas/' /var/lib/gravity/resources/sentinel.yaml
    fi

    echo "Checking status"
    rig status $RIG_CHANGESET --retry-attempts=120 --retry-period=2s --debug
    echo "Freezing"
    rig freeze

elif [ $1 = "rollback" ]; then
    echo "Reverting changeset $RIG_CHANGESET"
    rig revert

    if [ $(kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
    then
        kubectl scale --replicas=3 deployment stolon-sentinel
    fi
else
    echo "Missing argument, should be either 'update' or 'rollback'"
fi
