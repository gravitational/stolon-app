#!/usr/bin/env bash
set -o nounset

# create telegraf user
curl --silent --show-error -XPOST "http://influxdb.kube-system.svc.cluster.local:8086/query?u=root&p=root" \
     --data-urlencode "q=CREATE USER ${INFLUXDB_USERNAME} WITH PASSWORD '${INFLUXDB_PASSWORD}'"
curl --silent --show-error -XPOST "http://influxdb.kube-system.svc.cluster.local:8086/query?u=root&p=root" \
     --data-urlencode "q=GRANT ALL on k8s to ${INFLUXDB_USERNAME}"

/usr/bin/telegraf --config /etc/telegraf/telegraf.conf
