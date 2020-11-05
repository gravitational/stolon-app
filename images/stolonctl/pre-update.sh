#!/usr/bin/env bash
set -x

# check cluster health before starting an upgrade to eliminate any possible issue with that
stolonctl status --short

# scale down sentinel replicas to avoid master failover during planet upgrade
# and potentially having multi-master situation
kubectl scale --replicas=0 deployment stolon-sentinel

# delete alerts resource to fix the issue with prometheus
/usr/bin/gravity resource rm alert stolon-replication-lag

kubectl patch daemonset stolon-keeper -p '{"spec": {"template": {"spec": {"nodeSelector": {"non-existing": "true"}}}}}'
