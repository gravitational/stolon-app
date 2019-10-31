#!/bin/sh

for resource_file in alerts keeper rpc sentinel stolontool telegraf utils security
do
	kubectl delete -f /var/lib/gravity/resources/${resource_file}.yaml
done
kubectl delete secret stolon
etcdctl rm --recursive /stolon
