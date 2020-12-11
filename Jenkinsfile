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
                        sh """
                        docker build \
                          --cache-from ${env.REG}/${env.NS}/${env.IMAGE}:latest \
                          --target bin \
                          --output bin/ \
                          --platform ${PLATFORM} \
                          .
                        """
                    }
                }
                stage('Publish') {
                    steps {
                        container('curl') {
                            withCredentials([usernameColonPassword(credentialsId: 'oa-jenkins', variable: 'USERPASS')]) {
                                sh "gzip bin/rdepot"
                                sh "curl -v -u $USERPASS --upload-file bin/rdepot.gz https://nexus.openanalytics.eu/repository/releases/eu/openanalytics/rdepot/rdepot-cli/${env.BRANCH_NAME}/${env.BUILD_NUMBER}/rdepot.gz"
                            }
                        }
                    }
                }
            }
        }
        stage('Build image'){
            steps {
                ecrPull "${env.REG}", "${env.NS}/${env.IMAGE}", "latest", '', 'eu-west-1'
                sh """
                docker build \
                  --cache-from ${env.REG}/${env.NS}/${env.IMAGE}:latest \
                  --target image \
                  --platform local \
                  --tag ${env.NS}/${env.IMAGE} \
                  --tag openanalytics/${env.IMAGE}:latest \
                  --tag ${env.NS}/${env.IMAGE}:${env.shortCommit} \
                  .
                """
            }
        }
    }

    post {
        success  {
            ecrPush "${env.REG}", "${env.NS}/${env.IMAGE}", "latest", '', 'eu-west-1' 
            ecrPush "${env.REG}", "${env.NS}/${env.IMAGE}", "${env.shortCommit}", '', 'eu-west-1'
            withDockerRegistry([
                    credentialsId: "openanalytics-dockerhub",
                    url: ""]) {

                sh "docker push openanalytics/${env.IMAGE}:latest"
            }

        }
    }
}

