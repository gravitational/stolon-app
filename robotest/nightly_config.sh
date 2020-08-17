#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/utils.sh

# UPGRADE_MAP maps gravity version -> list of linux distros to upgrade from
declare -A UPGRADE_MAP

UPGRADE_MAP[$(recommended_upgrade_tag $(branch 5.5.x))]="redhat:7 centos:7 debian:9 ubuntu:18" # this branch
UPGRADE_MAP[5.5.0]="ubuntu:16"
UPGRADE_MAP[$(recommended_upgrade_tag $(branch 5.2.x))]="redhat:7" # compatible LTS version
UPGRADE_MAP[5.2.0]="centos:7"
UPGRADE_MAP[$(recommended_upgrade_tag $(branch 5.4.x))]="debian:9" # compatible non-LTS version
UPGRADE_MAP[5.4.0]="debian:8"

# via intermediate upgrade
UPGRADE_MAP[5.0.36]="centos:7"

# customer specific scenarios, these versions have traction in the field
UPGRADE_MAP[5.5.41]="centos:7"
UPGRADE_MAP[5.5.40]="centos:7"
UPGRADE_MAP[5.5.38]="centos:7"
UPGRADE_MAP[5.5.36]="centos:7"
UPGRADE_MAP[5.5.28]="centos:7"

UPGRADE_VERSIONS=${!UPGRADE_MAP[@]}

function build_upgrade_size_suite {
  local from_tarball=/$(tag_to_tarball $(recommended_upgrade_tag $(branch 5.5.x)))
  local os="redhat:7"
  local cluster_sizes=( \
    '"flavor":"three","nodes":3,"role":"node"' \
    '"flavor":"six","nodes":6,"role":"node"' \
    '"flavor":"one","nodes":1,"role":"node"')
  local suite=''
  for size in ${cluster_sizes[@]}; do
      suite+=$(build_upgrade_step $from_tarball $os $size)
    suite+=' '
  done
  echo -n $suite
}

function build_upgrade_to_release_under_test_suite {
  local cluster_size='"flavor":"three","nodes":3,"role":"node"'
  local suite=''
  for release in ${!UPGRADE_MAP[@]}; do
    local from_tarball=/$(tag_to_tarball $release)
    for os in ${UPGRADE_MAP[$release]}; do
      suite+=$(build_upgrade_step $from_tarball $os $cluster_size)
      suite+=' '
    done
  done
  echo -n $suite
}

function build_upgrade_suite {
  local suite=$(cat <<EOF
 $(build_upgrade_size_suite)
 $(build_upgrade_to_release_under_test_suite)
EOF
)
  echo -n $suite
}

function build_resize_suite {
  local suite=$(cat <<EOF
 resize={"to":3,"flavor":"one","nodes":1,"role":"node","state_dir":"/var/lib/telekube","os":"ubuntu:18","storage_driver":"overlay2"}
 resize={"to":6,"flavor":"three","nodes":3,"role":"node","state_dir":"/var/lib/telekube","os":"ubuntu:18","storage_driver":"overlay2"}
EOF
)
  echo -n $suite
}

function build_ops_install_suite {
  local suite=$(cat <<EOF
 install={"installer_url":"/installer/opscenter.tar","nodes":1,"flavor":"standalone","role":"node","os":"ubuntu:18","ops_advertise_addr":"example.com:443"}
EOF
)
  echo -n $suite
}

function build_install_suite {
  local suite=''
  local test_os="redhat:7 debian:9 ubuntu:16 ubuntu:18"
  local cluster_sizes=( \
    '"flavor":"three","nodes":3,"role":"node"' \
    '"flavor":"six","nodes":6,"role":"node"')
  for os in $test_os; do
    for size in ${cluster_sizes[@]}; do
      suite+=$(cat <<EOF
 install={${size},"os":"${os}","storage_driver":"overlay2"}
EOF
)
      suite+=' '
    done
  done
  suite+=$(build_ops_install_suite)
  echo -n $suite
}

SUITE="$(build_install_suite)"
SUITE+=" $(build_resize_suite)"
SUITE+=" $(build_upgrade_suite)"

# robotest wants the suite space separated, but it is easier to eyeball when newline separated
echo "$SUITE" | tr ' ' '\n'

export SUITE UPGRADE_VERSIONS
