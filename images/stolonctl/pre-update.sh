#!/usr/bin/env bash
set -x

# check cluster health before starting an upgrade to eliminate any possible issue with that
stolonctl status --short

# scale down sentinel replicas to avoid master failover during planet upgrade
# and potentially having multi-master situation
kubectl scale --replicas=0 deployment stolon-sentinel
kubectl patch daemonset stolon-keeper -p '{"spec": {"template": {"spec": {"nodeSelector": {"non-existing": "true"}}}}}'
