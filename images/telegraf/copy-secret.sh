#!/usr/bin/env bash

## copy telegraf secret from monitoring namespace
if kubectl --namespace=monitoring get secret telegraf-influxdb-creds >/dev/null 2>&1
then
    kubectl --namespace=monitoring get secret telegraf-influxdb-creds -o yaml |\
        sed 's/namespace: monitoring/namespace: default/' | kubectl --namespace=default apply -f -
fi
