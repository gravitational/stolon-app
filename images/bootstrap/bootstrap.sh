#!/bin/sh

/usr/local/bin/stolonboot -sentinels 1 -rpc 1

if [ $(kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
then
    kubectl scale --replicas=3 deployment stolon-sentinel
fi

kubectl create -f /var/lib/gravity/resources/alerts.yaml
