pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'docker.io'
        PROJECT_NAME = 'crypto-bank'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Services') {
            parallel {
                stage('Build Bank Service') {
                    steps {
                        dir('bank-service') {
                            sh 'go mod download'
                            sh 'go build -o bin/bank-service ./cmd/server'
                        }
                    }
                }
                stage('Build Exchange Service') {
                    steps {
                        dir('exchange-service') {
                            sh 'go mod download'
                            sh 'go build -o bin/exchange-service ./cmd/server'
                        }
                    }
                }
                stage('Build Analytics Service') {
                    steps {
                        dir('analytics-service') {
                            sh 'go mod download'
                            sh 'go build -o bin/analytics-service ./cmd/server'
                        }
                    }
                }
                stage('Build Notification Service') {
                    steps {
                        dir('notification-service') {
                            sh 'go mod download'
                            sh 'go build -o bin/notification-service ./cmd/server'
                        }
                    }
                }
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./... -v -cover'
            }
        }

        stage('Build Docker Images') {
            steps {
                sh 'docker-compose build'
            }
        }

        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh 'docker-compose up -d'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}

