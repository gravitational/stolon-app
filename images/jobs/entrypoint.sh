#!/usr/bin/env bash
# -*- mode: sh; -*-

# File: entrypoint.sh
# Time-stamp: <2017-08-31 11:35:47>
# Copyright (C) 2017 Gravitational Inc
# Description: Entrypoint for one time jobs

# set -o xtrace
set -o nounset
set -o errexit
set -o pipefail

function announce_step() {
    echo
    echo "===> $*"
    echo
}

function s3cmd_config() {
    announce_step "Configuring s3cmd"

    cat <<EOF > ~/.s3cfg
[default]
host_bucket = $PITHOS_BACKUP_BUCKET
host_base = $PITHOS_S3_ENDPOINT
access_key = $AWS_ACCESS_KEY_ID
secret_key = $AWS_SECRET_ACCESS_KEY
use_https = True
signature_v2 = True
EOF
}

if [ "$1" = "delete-backup-bucket" ]
then
    s3cmd_config
    announce_step "Deleting stolon-backup bucket"
    s3cmd --ca-cert=/etc/ssl/cluster-ca/ca.pem rb --recursive --force s3://${PITHOS_BACKUP_BUCKET}

elif [ "$1" = "create-backup-bucket" ]
then
    s3cmd_config
    announce_step "Creating stolon-backup bucket"
    s3cmd --ca-cert=/etc/ssl/cluster-ca/ca.pem mb --force s3://${PITHOS_BACKUP_BUCKET}

else
    echo "Missing argument"
fi
