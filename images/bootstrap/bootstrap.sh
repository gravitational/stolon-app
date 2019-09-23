#!/bin/sh
kubectl apply -f /var/lib/gravity/resources/security.yaml
kubectl apply -f /var/lib/gravity/resources/telegraf.yaml

/usr/local/bin/stolonboot -sentinels 1

if [ $(kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
then
    kubectl scale --replicas=3 deployment stolon-sentinel
fi

kubectl create -f /var/lib/gravity/resources/pgbouncer.yaml
kubectl wait --for=condition=complete --timeout=120s job/bootstrap-auth-function

kubectl create -f /var/lib/gravity/resources/alerts.yaml
