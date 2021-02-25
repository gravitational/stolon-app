#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/utils.sh

readonly RUN_UPGRADE=${RUN_UPGRADE:-0}

# UPGRADE_MAP maps gravity version -> list of linux distros to upgrade from
declare -A UPGRADE_MAP

# TODO (Sergei): enable upgrades from recommended tags
UPGRADE_MAP[1.12.11]="ubuntu:18" # this branch

# via intermediate upgrade
UPGRADE_MAP[1.10.61]="centos:7"

# customer specific scenarios, these versions have traction in the field
UPGRADE_MAP[1.12.7]="centos:7"

UPGRADE_VERSIONS=${!UPGRADE_MAP[@]}

function build_upgrade_suite {
  local suite=''
  local cluster_size='"flavor":"three","nodes":3,"role":"node"'
  for release in ${!UPGRADE_MAP[@]}; do
    local from_tarball=/$(tag_to_tarball $release)
    for os in ${UPGRADE_MAP[$release]}; do
      suite+=$(build_upgrade_step $from_tarball $os $cluster_size)
      suite+=' '
    done
  done
  echo -n $suite
}


function build_resize_suite {
  local suite=$(cat <<EOF
 resize={"to":3,"flavor":"one","nodes":1,"role":"node","state_dir":"/var/lib/telekube","os":"ubuntu:18","storage_driver":"overlay2"}
EOF
)
    echo -n $suite
}

function build_install_suite {
  local suite=''
  local test_os="redhat:7"
  local cluster_size='"flavor":"three","nodes":3,"role":"node"'
  suite+=$(cat <<EOF
 install={${cluster_size},"os":"${test_os}","storage_driver":"overlay2"}
EOF
)
  suite+=' '
  echo -n $suite
}

SUITE=$(build_install_suite)
#SUITE="$SUITE $(build_resize_suite)"
if [ "$RUN_UPGRADE" == "1" ]; then
  SUITE="$SUITE $(build_upgrade_suite)"
fi

echo "$SUITE" | tr ' ' '\n'

export SUITE UPGRADE_VERSIONS
