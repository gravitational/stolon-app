#!/usr/bin/env bash
# -*- mode: bash; -*-
set -o nounset
set -o errexit
set -o pipefail

kubectl scale --replicas=1 deployment stolon-sentinel
kubectl delete -f /var/lib/gravity/resources/preUpdate.yaml --ignore-not-found
kubectl create -f /var/lib/gravity/resources/preUpdate.yaml
kubectl wait --for=condition=complete --timeout=120s job/stolon-app-pre-update 

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
    rig delete secrets/telegraf-influxdb-creds --force
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

set +e
helm upgrade --install stolon /var/lib/gravity/resources/charts/stolon \
     --values /var/lib/gravity/resources/custom-values.yaml \
     --set existingSecret=stolon

set -e
kubectl wait --for=condition=complete --timeout=5m job/stolon-postgres-hardening
