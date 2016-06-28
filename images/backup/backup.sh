#!/bin/bash

set -e

# create .pgpass

echo "$STOLON_POSTGRES_SERVICE_HOST:$STOLON_POSTGRES_SERVICE_PORT:$STOLON_BACKUP_DB:stolon:$(cat /etc/secrets/stolon/password)" > ~/.pgpass
chmod 0600 ~/.pgpass

# pg_dump database

pg_dump --host $STOLON_POSTGRES_SERVICE_HOST --port $STOLON_POSTGRES_SERVICE_PORT -U stolon -w $STOLON_BACKUP_DB | gzip > /backups/stolon-$STOLON_BACKUP_DB-$(date +%Y%m%dT%H%M%S).sql.gz

echo "backup done"
