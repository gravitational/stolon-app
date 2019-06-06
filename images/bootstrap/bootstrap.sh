#!/bin/sh
/opt/bin/kubectl apply -f /var/lib/gravity/resources/security.yaml
/opt/bin/kubectl apply -f /var/lib/gravity/resources/telegraf.yaml

/usr/local/bin/stolonboot -sentinels 1

if [ $(/opt/bin/kubectl get nodes -l stolon-keeper=yes -o name | wc -l) -ge 3 ]
then
    /opt/bin/kubectl scale --replicas=3 deployment stolon-sentinel
fi

/opt/bin/kubectl create -f /var/lib/gravity/resources/alerts.yaml
