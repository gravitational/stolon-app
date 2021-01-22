#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

readonly TARGET=${1:?Usage: $0 [pr|nightly] [upgrade_from_dir]}
export UPGRADE_FROM_DIR=${2:-$(git rev-parse --show-toplevel)/upgrade_from}

readonly ROBOTEST_SCRIPT=$(mktemp -d)/runsuite.sh

# a number of environment variables are expected to be set
# see https://github.com/gravitational/robotest/blob/v2.0.0/suite/README.md
export ROBOTEST_VERSION=${ROBOTEST_VERSION:-uid-gid}
export ROBOTEST_REPO=quay.io/gravitational/robotest-suite:$ROBOTEST_VERSION
export INSTALLER_URL=$(pwd)/build/installer.tar
export GRAVITY_URL=$(pwd)/bin/gravity
export TAG=stolon-$(git rev-parse --short HEAD)
# cloud provider that test clusters will be provisioned on
# see https://github.com/gravitational/robotest/blob/v2.0.0/infra/gravity/config.go#L72
export DEPLOY_TO=${DEPLOY_TO:-gce}
export GCL_PROJECT_ID=${GCL_PROJECT_ID:-"kubeadm-167321"}
export GCE_REGION="northamerica-northeast1,us-west1,us-east1,us-east4,us-central1"
# GCE_VM tuned down from the Robotest's 7 cpu default in 09cec0e49e9d51c3603950209cec3c26dfe0e66b
# We should consider changing Robotest's default so that we can drop the override here. -- 2019-04 walt
export GCE_VM=${GCE_VM:-custom-4-8192}
# Parallelism & retry, tuned for GCE
export PARALLEL_TESTS=${PARALLEL_TESTS:-4}
export REPEAT_TESTS=${REPEAT_TESTS:-1}
export RETRIES=${RETRIES:-1}
export DOCKER_RUN_FLAGS=${DOCKER_RUN_FLAGS:-"--rm=true --user=$(id -u):$(id -g)"}

# Work around a bug in https://github.com/gravitational/robotest/blob/v2.1.0/docker/suite/run_suite.sh#L21-L30
# which mounts a volume inside a volume, resulting in docker creating the inner mountpoint
# owned by root:root if it does not already exist.
INSTALLER_BINDIR="$(dirname ${INSTALLER_URL})/bin"
mkdir -p "${INSTALLER_BINDIR}"


# set SUITE and UPGRADE_VERSIONS
case $TARGET in
  pr) source $(dirname $0)/pr_config.sh;;
  nightly) source $(dirname $0)/nightly_config.sh;;
  *) echo "Unknown target $TARGET\nUsage: $0 [pr|nightly] [upgrade_from_dir]"; exit 1;;
esac

function build_volume_mounts {
  for release in ${UPGRADE_VERSIONS[@]}; do
      local tarball=$(tag_to_tarball ${release})
      echo "-v $UPGRADE_FROM_DIR/$tarball:/$tarball"
  done
}

export EXTRA_VOLUME_MOUNTS=$(build_volume_mounts)

mkdir -p $UPGRADE_FROM_DIR
for release in ${!UPGRADE_MAP[@]}; do
  if [ ! -f $UPGRADE_FROM_DIR/$(tag_to_tarball ${release}) ]; then
      # aws s3 cp has incredibly verbose progress, disable progress for the sake of concise CI logs
      [[ -z ${CI:-} ]] || S3_FLAGS=--no-progress
      aws s3 cp ${S3_FLAGS:-} s3://builds.gravitational.io/stolon/$(tag_to_tarball ${release}) $UPGRADE_FROM_DIR/$(tag_to_tarball ${release})
  fi
done

docker pull $ROBOTEST_REPO
docker run $ROBOTEST_REPO cat /usr/bin/run_suite.sh > $ROBOTEST_SCRIPT
chmod +x $ROBOTEST_SCRIPT
$ROBOTEST_SCRIPT "$SUITE"
