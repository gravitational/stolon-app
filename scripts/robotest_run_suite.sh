#!/bin/bash
set -exu -o pipefail

readonly UPGRADE_FROM_DIR=${1:-$(pwd)/../upgrade_from}
readonly TOP_DIR=$(pwd)/../

declare -A UPGRADE_MAP
# gravity version -> list of OS releases to exercise on
UPGRADE_MAP[1.10.59]="redhat:7"
readonly RUN_UPGRADE=${RUN_UPGRADE:-0}

readonly OPS_APIKEY=${API_KEY:?API key for distribution Ops Center required}
readonly APP_BUILDDIR=$TOP_DIR/build
readonly ROBOTEST_SCRIPT=$(mktemp -d)/runsuite.sh

# number of environment variables are expected to be set
# see https://github.com/gravitational/robotest/blob/master/suite/README.md
export ROBOTEST_VERSION=${ROBOTEST_VERSION:-stable-gce}
export ROBOTEST_REPO=quay.io/gravitational/robotest-suite:$ROBOTEST_VERSION
export WAIT_FOR_INSTALLER=true
export INSTALLER_URL=$(pwd)/build/installer.tar
export GRAVITY_URL=$(pwd)/bin/gravity
export TELE=$TOP_DIR/bin/tele
export DEPLOY_TO=${DEPLOY_TO:-gce}
export TAG=stolon-$(git rev-parse --short HEAD)
export GCL_PROJECT_ID=${GCL_PROJECT_ID:-"kubeadm-167321"}
export GCE_REGION="northamerica-northeast1,us-west1,us-east1,us-east4,us-central1"
export GCE_PREEMPTIBLE="false"
export DOCKER_RUN_FLAGS="--user $(id -u):$(id -g)"

function build_upgrade_step {
  local usage="$FUNCNAME os release storage-driver cluster-size"
  local os=${1:?$usage}
  local release=${2:?$usage}
  local storage_driver=${3:?$usage}
  local cluster_size=${4:?$usage}
  local suite=''
  suite+=$(cat <<EOF
 upgrade3lts={${cluster_size},"os":"${os}","storage_driver":"${storage_driver}","from":"/installer_${release}.tar"}
EOF
)
  echo $suite
}

function build_upgrade_suite {
  local suite=''
  local cluster_size='"flavor":"three","nodes":3,"role":"node"'
  for release in ${!UPGRADE_MAP[@]}; do
    for os in ${UPGRADE_MAP[$release]}; do
      suite+=$(build_upgrade_step $os $release 'devicemapper' $cluster_size)
      suite+=" $(build_upgrade_step $os $release 'overlay2' $cluster_size)"
    done
  done
  echo $suite
}


function build_install_suite {
  local usage="$FUNCNAME os storage-driver cluster-size"
  local suite=''
  local test_os=${1:-'redhat:7'}
  local storage_driver=${2:-'overlay2'}
  local cluster_size=${3:-'"flavor":"three","nodes":3,"role":"node"'}
  suite+=$(cat <<EOF
 install={${cluster_size},"os":"${test_os}","storage_driver":"${storage_driver}"}
EOF
)
  echo $suite
}
set -u

function build_volume_mounts {
  for release in ${!UPGRADE_MAP[@]}; do
    echo "-v $UPGRADE_FROM_DIR/installer_${release}.tar:/installer_${release}.tar"
  done
}

export EXTRA_VOLUME_MOUNTS=$(build_volume_mounts)

suite="$(build_install_suite)"
if [ "$RUN_UPGRADE" == "1" ]; then
  suite="$suite $(build_upgrade_suite)"
fi

echo $suite

mkdir -p $UPGRADE_FROM_DIR
for release in ${!UPGRADE_MAP[@]}; do
    if [ ! -f $UPGRADE_FROM_DIR/installer_$release.tar ]; then
        tele pull $EXTRA_GRAVITY_OPTIONS stolon-app:$release --output=$UPGRADE_FROM_DIR/installer_$release.tar
    fi
done

docker pull $ROBOTEST_REPO
docker run $ROBOTEST_REPO cat /usr/bin/run_suite.sh > $ROBOTEST_SCRIPT
chmod +x $ROBOTEST_SCRIPT
$ROBOTEST_SCRIPT "$suite"
