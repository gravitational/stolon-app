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
                downloadBinaries($GRAVITY_VERSION)
            }
        }

        stage('Generate installer tarball') {
            environment {
                PATH = "\$(pwd)/bin:\$PATH"
                APP_VERSION = sh(script: 'make what-version', returnStdout: true).trim()
            }
            steps {
                lock("installer-${BRANCH}") {
                    print "Building stolon-app installer tarball version $APP_VERSION"
                    script {
                        installerTarballFileName = getInstallerTarballFileName()
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
def getInstallerTarballFileName() {
  return "stolon-app-${APP_VERSION}-installer.tar.gz"
}

def createInstallerTarball(installerTarballFileName) {
    def stateDir = "${pwd()}/state"
    sh "echo $PATH"
}
