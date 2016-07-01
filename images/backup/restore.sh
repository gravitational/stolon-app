#!/bin/bash

set -e

# create .pgpass

echo "$STOLON_POSTGRES_SERVICE_HOST:$STOLON_POSTGRES_SERVICE_PORT:$STOLON_BACKUP_DB:stolon:$(cat /etc/secrets/stolon/password)" > ~/.pgpass
chmod 0600 ~/.pgpass

# restore database from dump

filename=/backups/$STOLON_BACKUP_FILE
echo "restoring from: $filename"

gunzip -c $filename | psql --host $STOLON_POSTGRES_SERVICE_HOST --port $STOLON_POSTGRES_SERVICE_PORT -U stolon -w $STOLON_BACKUP_DB

echo "done"
