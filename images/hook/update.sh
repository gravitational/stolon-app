#!/usr/bin/env bash
# -*- mode: bash; -*-
set -o nounset
set -o errexit
set -o pipefail
set -o xtrace

export EXTRA_PARAMS=""
export PASSWORD=$(kubectl get secrets stolon -o jsonpath='{.data.password}'|base64 -d)
if [ -n $PASSWORD ]
then
    export EXTRA_PARAMS="$EXTRA_PARAMS --set superuser.password=$PASSWORD"
fi

export PG_REPL_PASSWORD=$(kubectl get secrets stolon -o jsonpath='{.data.pg_repl_password}'|base64 -d)
if [ -n $PG_REPL_PASSWORD ]
then
    export EXTRA_PARAMS="$EXTRA_PARAMS --set replication.password=$PG_REPL_PASSWORD"
fi

# check for existence of stolon helm release
if [[ $(helm list stolon | wc -l) -eq 0 ]]
then
    echo "Saving pre-helm resources in changeset: $RIG_CHANGESET"
    rig delete ds/stolon-keeper --force
    rig delete deployments/stolon-proxy --force
    rig delete deployments/stolon-rpc --force
    rig delete deployments/stolon-sentinel --force
    rig delete deployments/stolon-utils --force
    rig delete deployments/stolonctl --force
    rig delete configmaps/stolon-telegraf --force
    rig delete configmaps/stolon-telegraf-node --force
    rig delete configmaps/stolon-pgbouncer --force
    rig delete services/stolon-postgres --force
    rig delete services/stolonctl --force
    rig delete services/stolon-rpc --force
    rig delete serviceaccounts/stolon-keeper --force
    rig delete serviceaccounts/stolon-rpc --force
    rig delete serviceaccounts/stolon-sentinel --force
    rig delete serviceaccounts/stolon-utils --force
    rig delete roles/stolon-keeper --force
    rig delete roles/stolon-rpc --force
    rig delete roles/stolon-sentinel --force
    rig delete roles/stolon-utils --force
    rig delete rolebindings/stolon-keeper --force
    rig delete rolebindings/stolon-rpc --force
    rig delete rolebindings/stolon-sentinel --force
    rig delete rolebindings/stolon-utils --force
    rig delete secrets/stolon --force

    # backward-compatible(upgrades from 1.10.x) replication password
    export EXTRA_PARAMS="$EXTRA_PARAMS --set replication.password=replpassword"
fi

rig delete secrets/telegraf-influxdb-creds --force

# check whether stolon secret was created not under helm release
if [[ $(kubectl get secrets stolon -o jsonpath='{.metadata.labels.chart}' | wc -l) -eq 0 ]]
then
    rig delete secrets/stolon --force
fi

# delete jobs from previous run
for name in stolon-bootstrap-auth-function stolon-copy-telegraf-influxdb-creds \
					   stolon-create-alerts stolon-postgres-hardening
do
    kubectl delete job $name --ignore-not-found
done

if [[ $(helm list stolon | wc -l) -gt 0 ]]
then
    # scale up sentinel pods
    if [ $(kubectl get nodes -lstolon-keeper=yes --output=go-template --template="{{len .items}}") -gt 1 ]
    then
        kubectl scale deployment stolon-sentinel --replicas 3
    fi

    kubectl patch daemonset stolon-keeper --type json -p='[{"op": "remove", "path": "/spec/template/spec/nodeSelector/non-existing"}]'
fi

set +e
helm upgrade --install stolon /var/lib/gravity/resources/charts/stolon \
     --values /var/lib/gravity/resources/custom-values.yaml $EXTRA_PARAMS

set -e

timeout 5m bash -c "while ! kubectl get job stolon-postgres-hardening; do sleep 10; done"
kubectl wait --for=condition=complete --timeout=5m job/stolon-postgres-hardening
rig freeze
