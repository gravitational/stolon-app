#!/usr/bin/env groovy
def propagateParamsToEnv() {
  for (param in params) {
    if (env."${param.key}" == null) {
      env."${param.key}" = param.value
    }
  }
}

properties([
  disableConcurrentBuilds(),
  parameters([
    string(name: 'TAG',
           defaultValue: 'master',
           description: 'Git tag to build'),
    string(name: 'VERSION',
           defaultValue: '',
           description: 'Override automatic versioning'),
    booleanParam(defaultValue: false,
           description: 'Run or skip robotest system wide tests.',
           name: 'RUN_ROBOTEST'),
    choice(choices: ["true", "false"].join("\n"),
           description: 'Destroy all VMs on success.',
           name: 'DESTROY_ON_SUCCESS'),
    choice(choices: ["true", "false"].join("\n"),
           description: 'Destroy all VMs on failure.',
           name: 'DESTROY_ON_FAILURE'),
    choice(choices: ["true", "false"].join("\n"),
           description: 'Abort all tests upon first failure.',
           name: 'FAIL_FAST'),
    choice(choices: ["gce"].join("\n"),
           description: 'Cloud provider to deploy to.',
           name: 'DEPLOY_TO'),
    string(name: 'PARALLEL_TESTS',
           defaultValue: '4',
           description: 'Number of parallel tests to run.'),
    string(name: 'REPEAT_TESTS',
           defaultValue: '1',
           description: 'How many times to repeat each test.'),
    string(name: 'RETRIES',
           defaultValue: '0',
           description: 'How many times to retry each failed test'),
    string(name: 'ROBOTEST_VERSION',
           defaultValue: '2.2.1',
           description: 'Robotest tag to use.'),
    string(name: 'GRAVITY_VERSION',
           defaultValue: '5.5.57',
           description: 'gravity/tele binaries version'),
    string(name: 'TELE_VERSION',
           defaultValue: '7.0.30',
           description: 'Version of tele binary to build application'),
    string(name: 'CLUSTER_SSL_APP_VERSION',
           defaultValue: '0.8.5',
           description: 'cluster-ssl-app version'),
    string(name: 'EXTRA_GRAVITY_OPTIONS',
           defaultValue: '',
           description: 'Gravity options to add when calling tele'),
    string(name: 'TELE_BUILD_EXTRA_OPTIONS',
           defaultValue: '',
           description: 'Extra options to add when calling tele build'),
    booleanParam(name: 'ADD_GRAVITY_VERSION',
                 defaultValue: false,
                 description: 'Appends "-${GRAVITY_VERSION}" to the tag to be published'),
    booleanParam(name: 'BUILD_CLUSTER_IMAGE',
                 defaultValue: true,
                 description: 'Generate a Gravity Cluster Image(Self-sufficient tarball)'),
  ]),
])

node {
  workspace {
    stage('checkout') {
      checkout([
        $class: 'GitSCM',
        branches: [[name: "${params.TAG}"]],
        doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
        extensions: [[$class: 'CloneOption', noTags: false, shallow: false]],
        submoduleCfg: [],
        userRemoteConfigs: scm.userRemoteConfigs,
      ])
    }
    stage('params') {
      echo "${params}"
      propagateParamsToEnv()
    }
    stage('clean') {
      sh "make clean"
    }

    APP_VERSION = sh(script: 'make what-version', returnStdout: true).trim()
    APP_VERSION = params.ADD_GRAVITY_VERSION ? "${APP_VERSION}-${GRAVITY_VERSION}" : APP_VERSION
    STATEDIR = "${pwd()}/state/${APP_VERSION}"
    BINARIES_DIR = "${pwd()}/bin"
    MAKE_ENV = [
      "STATEDIR=${STATEDIR}",
      "PATH+GRAVITY=${BINARIES_DIR}",
      "VERSION=${APP_VERSION}"
    ]

    stage('download gravity/tele binaries') {
      withEnv(MAKE_ENV + ["BINARIES_DIR=${BINARIES_DIR}"]) {
        sh 'make download-binaries'
      }
    }

    stage('populate state directory with gravity and cluster-ssl packages') {
      withEnv(MAKE_ENV + ["BINARIES_DIR=${BINARIES_DIR}"]) {
        sh 'make install-dependent-packages'
      }
    }

    stage('build-app') {
      if (params.BUILD_CLUSTER_IMAGE) {
        withEnv(MAKE_ENV) {
          sh 'make build-app'
        }
      } else {
        echo 'skipped build of gravity cluster image'
      }
    }

    stage('build gravity app') {
      if (params.BUILD_GRAVITY_APP) {
        withEnv(MAKE_ENV) {
          sh 'make build-gravity-app'
          archiveArtifacts "build/application.tar"
        }
      } else {
        echo 'skipped build gravity app'
      }
    }
  }
}

void workspace(Closure body) {
  timestamps {
    ws("${pwd()}-${BUILD_ID}") {
      body()
    }
  }
}
