@Library('jenkins-ecr-libs') _
pipeline {

    agent {
        kubernetes {
            yamlFile 'agent.pod.yaml'
            defaultContainer 'dind'
        }
    }

    options {
        authorizationMatrix inheritanceStrategy: inheritingGlobal(), permissions: ['hudson.model.Item.Build:oa-infrastructure', 'hudson.model.Item.Read:oa-infrastructure']
        buildDiscarder(logRotator(numToKeepStr: '3'))
    }

    environment {
        shortCommit = sh(returnStdout: true, script: "echo ${env.GIT_COMMIT} | cut -c 1-8").trim()
        IMAGE = "rdepot-cli"
        NS = "oa-infrastructure"
        REG = "196229073436.dkr.ecr.eu-west-1.amazonaws.com"
        DOCKER_BUILDKIT = 1
        TAG = "${env.BRANCH_NAME == 'master' ? 'latest' : env.BRANCH_NAME}"
    }
    
    stages {
        stage('Binaries') {
            environment {
              PLATFORM = "linux/amd64"
              NAME = "rdepot"
            }
            stages {
                stage('Build') {
                    steps {
                        ecrPull "${env.REG}", "${env.NS}/${env.IMAGE}", "${env.TAG}", '', 'eu-west-1'
                        sh """
                        docker build \
                          --cache-from ${env.REG}/${env.NS}/${env.IMAGE}:${env.TAG} \
                          --target bin \
                          --output bin/ \
                          --platform ${PLATFORM} \
                          .
                        """
                    }
                }
                stage('Publish') {
                    when {
                        anyOf {
                            branch "master"
                            branch pattern: "\\d+\\.\\d+\\.\\d+", comparator: "REGEXP"
                        }
                    }
                    steps {
                        container('curl') {
                            withCredentials([usernameColonPassword(credentialsId: 'oa-jenkins', variable: 'USERPASS')]) {
                                sh "gzip bin/rdepot"
                                sh "curl -u $USERPASS --upload-file bin/rdepot.gz https://nexus.openanalytics.eu/repository/releases/eu/openanalytics/rdepot/rdepot-cli/${env.TAG}/rdepot.gz"
                            }
                        }
                    }
                }
            }
        }
        stage('Build image'){
            when {
                anyOf {
                    branch "master"
                    branch pattern: "\\d+\\.\\d+\\.\\d+", comparator: "REGEXP"
                }
            }
            steps {
                ecrPull "${env.REG}", "${env.NS}/${env.IMAGE}", "${env.TAG}", '', 'eu-west-1'
                sh """
                docker build \
                  --cache-from ${env.REG}/${env.NS}/${env.IMAGE}:${env.TAG} \
                  --target image \
                  --platform local \
                  --tag ${env.NS}/${env.IMAGE} \
                  --tag openanalytics/${env.IMAGE}:${env.TAG} \
                  --tag ${env.NS}/${env.IMAGE}:${env.shortCommit} \
                  --tag ${env.NS}/${env.IMAGE}:${env.TAG} \
                  .
                """
            }
            post {
                success  {
                    ecrPush "${env.REG}", "${env.NS}/${env.IMAGE}", "${env.TAG}", '', 'eu-west-1' 
                    ecrPush "${env.REG}", "${env.NS}/${env.IMAGE}", "${env.shortCommit}", '', 'eu-west-1'
                    withDockerRegistry([
                            credentialsId: "openanalytics-dockerhub",
                            url: ""]) {

                        sh "docker push openanalytics/${env.IMAGE}:${env.TAG}"
                    }

                }
            }
        }
    }

}

