#!/usr/bin/env bash
set -o errexit
set -o xtrace

# check cluster health before starting an upgrade to eliminate any possible issue with that
until stolonctl status --short
do
	sleep 10
done

# scale down sentinel replicas to avoid master failover during planet upgrade
# and potentially having multi-master situation
kubectl scale --replicas=0 deployment stolon-sentinel

# delete alerts resource to fix the issue with prometheus
/opt/bin/gravity resource rm alert stolon-replication-lag

kubectl patch daemonset stolon-keeper -p '{"spec": {"template": {"spec": {"nodeSelector": {"non-existing": "true"}}}}}'
