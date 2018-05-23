#!/usr/bin/env bash

if [ -d /stolon-data/postgres-old ]
then
    rm -rf /stolon-data/postgres
    rm -rf /stolon-data/upgrade-state
    mv /stolon-data/postgres-old-$PG_VERSION_OLD /stolon-data/postgres
fi

exit 0
