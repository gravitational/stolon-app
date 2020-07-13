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
    choice(choices: ["run", "skip"].join("\n"),
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
    string(name: 'ROBOTEST_VERSION',
           defaultValue: 'stable-gce',
           description: 'Robotest tag to use.'),
    booleanParam(name: 'ROBOTEST_RUN_UPGRADE',
           defaultValue: false,
           description: 'Run the upgrade suite as part of robotest'),
    string(name: 'OPS_URL',
           defaultValue: 'https://ci-ops.gravitational.io',
           description: 'Ops Center URL to download dependencies from'),
    string(name: 'OPS_CENTER_CREDENTIALS',
           defaultValue: 'CI_OPS_API_KEY',
           description: 'Jenkins\' key containing the Ops Center Credentials'),
    string(name: 'GRAVITY_VERSION',
           defaultValue: '7.0.12',
           description: 'gravity/tele binaries version'),
    string(name: 'CLUSTER_SSL_APP_VERSION',
           defaultValue: '0.8.2-5.5.21',
           description: 'cluster-ssl-app version'),
    string(name: 'INTERMEDIATE_RUNTIME_VERSION',
           defaultValue: '5.2.15',
           description: 'Version of runtime to upgrade with'),
    string(name: 'EXTRA_GRAVITY_OPTIONS',
           defaultValue: '',
           description: 'Gravity options to add when calling tele'),
    booleanParam(name: 'ADD_GRAVITY_VERSION',
                 defaultValue: false,
                 description: 'Appends "-${GRAVITY_VERSION}" to the tag to be published'),
    booleanParam(name: 'IMPORT_APP',
                 defaultValue: false,
                 description: 'Import application to ops center'
    ),
  ]),
])

node {
  workspace {
    stage('checkout') {
      checkout([
        $class: 'GitSCM',
        branches: [[name: "${params.TAG}"]],
        doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
        extensions: scm.extensions,
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
    TELE_STATE_DIR = "${pwd()}/state/${APP_VERSION}"
    BINARIES_DIR = "${pwd()}/bin"
    EXTRA_GRAVITY_OPTIONS = "--state-dir=${TELE_STATE_DIR} ${params.EXTRA_GRAVITY_OPTIONS}"
    MAKE_ENV = [
      "EXTRA_GRAVITY_OPTIONS=${EXTRA_GRAVITY_OPTIONS}",
      "PATH+GRAVITY=${BINARIES_DIR}",
      "VERSION=${APP_VERSION}"
    ]

    stage('download gravity/tele binaries for login') {
      withEnv(MAKE_ENV + ["BINARIES_DIR=${BINARIES_DIR}"]) {
        sh 'make download-binaries'
      }
    }

    stage('build-app') {
      withCredentials([
        string(credentialsId: params.OPS_CENTER_CREDENTIALS, variable: 'API_KEY'),
      ]) {
        withEnv(MAKE_ENV) {
          sh """
  rm -rf ${TELE_STATE_DIR} && mkdir -p ${TELE_STATE_DIR}
  tele logout ${EXTRA_GRAVITY_OPTIONS}
  tele login ${EXTRA_GRAVITY_OPTIONS} -o ${OPS_URL} --token=${API_KEY}
  make build-app"""
        }
      }
    }

    stage('test') {
      if (params.RUN_ROBOTEST == 'run') {
        throttle(['robotest']) {
          node {
              parallel (
              robotest : {
                withCredentials([
                  [$class: 'FileBinding', credentialsId:'ROBOTEST_LOG_GOOGLE_APPLICATION_CREDENTIALS', variable: 'GOOGLE_APPLICATION_CREDENTIALS'],
                  [$class: 'StringBinding', credentialsId: params.OPS_CENTER_CREDENTIALS, variable: 'API_KEY'],
                  [$class: 'FileBinding', credentialsId:'OPS_SSH_KEY', variable: 'SSH_KEY'],
                  [$class: 'FileBinding', credentialsId:'OPS_SSH_PUB', variable: 'SSH_PUB'],
                ]) {
                  def TELE_STATE_DIR = "${pwd()}/state/${APP_VERSION}"
                  sh """
                  export PATH=\$(pwd)/bin:\${PATH}
                  export EXTRA_GRAVITY_OPTIONS="--state-dir=${TELE_STATE_DIR}"
                  make robotest-run-suite \
                    AWS_KEYPAIR=ops \
                    AWS_REGION=us-east-1 \
                    ROBOTEST_VERSION=$ROBOTEST_VERSION \
                    RUN_UPGRADE=${params.ROBOTEST_RUN_UPGRADE ? 1 : 0}"""
                }
            } )
          }
        }
      } else {
        echo 'skipped system tests'
      }
    }

    stage('push') {
      if (params.IMPORT_APP) {
        withCredentials([
          string(credentialsId: params.OPS_CENTER_CREDENTIALS, variable: 'API_KEY'),
        ]) {
          withEnv(MAKE_ENV) {
            sh 'make push'
          }
        }
      } else {
        echo 'skipped application import'
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
