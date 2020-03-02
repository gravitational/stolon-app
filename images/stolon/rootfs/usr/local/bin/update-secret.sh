#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

if ! kubectl get secrets $SECRET_NAME -o json | jq '.data.pg_repl_password' -e
then
    echo "Adding key pg_repl_password into stolon secret"
    kubectl get secrets $SECRET_NAME -o json | jq ".data[\"pg_repl_password\"]=\"$PG_REPL_PASSWORD\"" | kubectl apply -f -
fi
