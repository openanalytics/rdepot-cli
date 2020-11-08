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
        stage('Build image'){
            steps {
                ecrPull "${env.REG}", "${env.NS}/${env.IMAGE}", "latest", '', 'eu-west-1'
                sh """
                docker build \
                  --cache-from ${env.REG}/${env.NS}/${env.IMAGE}:latest \
                  --target bin-unix \
                  --platform local \
                  --tag ${env.NS}/${env.IMAGE} \
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
        }
    }
}

