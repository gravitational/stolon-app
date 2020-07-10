#!/bin/sh

set +e
/usr/local/bin/helm install /var/lib/gravity/resources/charts/stolon \
    --values /var/lib/gravity/resources/custom-values.yaml \
    --values /var/lib/gravity/resources/custom-build.yaml \
	--name stolon

set -e
/usr/local/bin/kubectl wait --for=condition=complete --timeout=5m job/stolon-postgres-hardening
