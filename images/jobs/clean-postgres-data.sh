#!/usr/bin/env bash

CURRENT_VERSION=$(cat /stolon-data/postgres/PG_VERSION)
if [ $CURRENT_VERSION == $PG_VERSION_NEW ]; then exit 0; fi

rm -rf /stolon-data/postgres
