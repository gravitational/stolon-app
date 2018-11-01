#!/usr/bin/env bash

# create telegraf user
curl -XPOST "http://influxdb.kube-system.svc:8086/query?u=root&p=root" \
         --data-urlencode "q=CREATE USER ${INFLUXDB_TELEGRAF_USERNAME} WITH PASSWORD '${INFLUXDB_TELEGRAF_PASSWORD}'"
curl -XPOST "http://influxdb.kube-system.svc:8086/query?u=root&p=root" \
         --data-urlencode "q=GRANT ALL on k8s to ${INFLUXDB_TELEGRAF_USERNAME}"

/usr/bin/telegraf --config /etc/telegraf/telegraf.conf
