#!/bin/sh

for resource_file in utils rpc keeper sentinel jenkins pgweb
do
	kubectl delete -f /var/lib/gravity/resources/${resource_file}.yaml
done
kubectl delete secret stolon
kubectl label nodes -l stolon-keeper=yes stolon-keeper-
