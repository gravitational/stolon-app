#!/usr/bin/env bash
# -*- mode: bash; -*-
set -o nounset
set -o errexit
set -o pipefail

# check for existence of stolon helm release
if [[ $(helm list stolon | wc -l) -eq 0 ]]
then
    echo "Saving pre-helm resources in changeset: $RIG_CHANGESET"
    rig cs delete --force -c cs/$RIG_CHANGESET
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
    rig freeze
fi

rig delete secrets/telegraf-influxdb-creds --force

# delete jobs from previous run
for name in stolon-bootstrap-auth-function stolon-copy-telegraf-influxdb-creds \
					   stolon-create-alerts stolon-postgres-hardening
do
    kubectl delete job $name --ignore-not-found
done

export EXTRA_PARAMS=""
if [ -f /var/lib/gravity/resources/custom-build.yaml ]
then
    export EXTRA_PARAMS="--values /var/lib/gravity/resources/custom-build.yaml"
fi

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

set +e
helm upgrade --install stolon /var/lib/gravity/resources/charts/stolon \
     --values /var/lib/gravity/resources/custom-values.yaml $EXTRA_PARAMS

set -e
# scale up sentinel pods
if [ $(kubectl get nodes -lstolon-keeper=yes --output=go-template --template="{{len .items}}") -gt 1 ]
then
    kubectl scale deployment stolon-sentinel --replicas 3
fi

kubectl wait --for=condition=complete --timeout=5m job/stolon-postgres-hardening
