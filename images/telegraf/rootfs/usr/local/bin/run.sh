#!/usr/bin/env bash

# create influxdb user
curl -XPOST "http://influxdb.kube-system.svc:8086/query?u=root&p=root" \
         --data-urlencode "q=CREATE USER ${INFLUXDB_TELEGRAF_USERNAME} WITH PASSWORD '${INFLUXDB_TELEGRAF_PASSWORD}'"
curl -XPOST "http://influxdb.kube-system.svc:8086/query?u=root&p=root" \
         --data-urlencode "q=GRANT ALL on k8s to ${INFLUXDB_TELEGRAF_USERNAME}"
sed -i s/superSecurePassword/${INFLUXDB_TELEGRAF_PASSWORD}/ /etc/telegraf/telegraf.conf

# setup and start telegraf
echo "  address = \"host=stolon-postgres.default.svc user=${STOLON_USERNAME} password=$(cat $STOLON_PASSWORDFILE) sslmode=disable dbname=postgres\"" \
	 >> /etc/telegraf/telegraf.conf
/usr/bin/telegraf --config /etc/telegraf/telegraf.conf
