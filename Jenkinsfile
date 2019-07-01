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
    choice(choices: ["run", "skip"].join("\n"),
           defaultValue: 'run',
           description: 'Run or skip robotest system wide tests.',
           name: 'RUN_ROBOTEST'),
    choice(choices: ["true", "false"].join("\n"),
           defaultValue: 'true',
           description: 'Destroy all VMs on success.',
           name: 'DESTROY_ON_SUCCESS'),
    choice(choices: ["true", "false"].join("\n"),
           defaultValue: 'true',
           description: 'Destroy all VMs on failure.',
           name: 'DESTROY_ON_FAILURE'),
    choice(choices: ["true", "false"].join("\n"),
           defaultValue: 'true',
           description: 'Abort all tests upon first failure.',
           name: 'FAIL_FAST'),
    choice(choices: ["gce"].join("\n"),
           defaultValue: 'gce',
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
    string(name: 'OPS_URL',
           defaultValue: 'https://ci-ops.gravitational.io',
           description: 'Ops Center URL to download dependencies from'),
    string(name: 'GRAVITY_VERSION',
           defaultValue: '5.2.12',
           description: 'gravity/tele binaries version'),
    string(name: 'CLUSTER_SSL_APP_VERSION',
           defaultValue: '0.8.2-5.2.12',
           description: 'cluster-ssl-app version')
  ]),
])

timestamps {
  node {
    stage('checkout') {
      checkout scm
    }
    stage('params') {
      echo "${params}"
      propagateParamsToEnv()
    }
    stage('clean') {
      sh "make clean"
    }
    stage('download gravity/tele binaries') {
      sh "make download-binaries"
    }

    APP_VERSION = sh(script: 'make what-version', returnStdout: true).trim()

    stage('build-app') {
      withCredentials([
      [$class: 'StringBinding', credentialsId:'CI_OPS_API_KEY', variable: 'API_KEY'],
      ]) {
        def TELE_STATE_DIR = "${pwd()}/state/${APP_VERSION}"
        sh """
export PATH=\$(pwd)/bin:\${PATH}
rm -rf ${TELE_STATE_DIR} && mkdir -p ${TELE_STATE_DIR}
export EXTRA_GRAVITY_OPTIONS="--state-dir=${TELE_STATE_DIR}"
tele login \${EXTRA_GRAVITY_OPTIONS} -o ${OPS_URL} --token=${API_KEY}
make build-app OPS_URL=$OPS_URL"""
      }
    }
  }
  throttle(['robotest']) {
    node {
      stage('test') {
        parallel (
        robotest : {
          if (params.RUN_ROBOTEST == 'run') {
            withCredentials([
                [$class: 'FileBinding', credentialsId:'ROBOTEST_LOG_GOOGLE_APPLICATION_CREDENTIALS', variable: 'GOOGLE_APPLICATION_CREDENTIALS'],
                [$class: 'StringBinding', credentialsId:'CI_OPS_API_KEY', variable: 'API_KEY'],
                [$class: 'FileBinding', credentialsId:'OPS_SSH_KEY', variable: 'SSH_KEY'],
                [$class: 'FileBinding', credentialsId:'OPS_SSH_PUB', variable: 'SSH_PUB'],
                ]) {
                  sh """
                  make robotest-run-suite \
                    AWS_KEYPAIR=ops \
                    AWS_REGION=us-east-1 \
                    ROBOTEST_VERSION=$ROBOTEST_VERSION"""
            }
          }else {
            echo 'skipped system tests'
          }
        } )
      }
    }
  }
}
