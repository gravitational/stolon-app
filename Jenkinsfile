#!/usr/bin/env groovy

// --> Constants
ANSIBLE_VERSION='2.7.4'
TERRAFORM_VERSION='0.11.11'
// <--

BRANCH = env.BRANCH_NAME

pipeline {
    agent any
    options {
        ansiColor(colorMapName: 'XTerm')
        disableConcurrentBuilds()
        timestamps()
    }

    parameters {
        string(
        name: 'GRAVITY_VERSION',
        defaultValue: '5.0.28',
        description: 'Version of gravity/tele binaries'
        ),
        string(
        name: 'CLUSTER_SSL_APP_VERSION',
        defaultValue: 'master',
        description: 'Version of cluster-ssl-app'
        )
    }

    environment {
        ANSIBLE_VERSION = "${ANSIBLE_VERSION}"
        GRAVITY_VERSION = "${GRAVITY_VERSION}"
    }
    stages {
        stage('Checkout source') {
            steps {
                checkout scm
            }
        }

        stage('Build docker image for Ansible') {
            steps {
                dir('tests/docker') {
                    sh "docker build -t ansible:${ANSIBLE_VERSION} --build-arg ANSIBLE_VERSION=${ANSIBLE_VERSION} ."
                }
            }
        }

        stage('Download gravity and tele binaries') {
            steps {
                print "Downloading tele and gravity version $GRAVITY_VERSION"
                downloadBinaries(GRAVITY_VERSION)
            }
        }

        stage('Generate installer tarball') {
            environment {
                APP_VERSION = sh(script: 'make what-version', returnStdout: true).trim()
            }
            steps {
                withEnv(['PATH+LOCAL=./bin']) {
                    print "Building stolon-app installer tarball version ${env.APP_VERSION}"
                    script {
                        installerTarballFileName = getInstallerTarballFileName(env.APP_VERSION)
                        createInstallerTarball(installerTarballFileName)
                    }
                }
            }
        }
    }
}

def downloadBinaries(version) {
    dir ('bin') {
        sh "curl -o tele https://get.gravitational.io/telekube/bin/${version}/linux/x86_64/tele -fSsL"
        sh "curl -o gravity https://get.gravitational.io/telekube/bin/${version}/linux/x86_64/gravity -fSsL"
        sh "chmod +x gravity tele"
    }
}

// returns installer tarball name based on application version
def getInstallerTarballFileName(appVersion) {
    return "stolon-app-${appVersion}-installer.tar.gz"
}

def createInstallerTarball(installerTarballFileName) {
    def stateDir = "${pwd()}/state"
    sh "echo $PATH"
}
