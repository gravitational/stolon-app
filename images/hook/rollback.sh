#!/usr/bin/env bash
# -*- mode: bash; -*-
set -o nounset
set -o errexit
set -o pipefail

if [[ $(helm list stolon --all | wc -l) -gt 0 ]]
then
    exit 0
fi

exit 0
