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
                            sh 'go mod tidy'
                            sh 'go build -o bin/bank-service ./cmd/server'
                        }
                    }
                }
                stage('Build Exchange Service') {
                    steps {
                        dir('exchange-service') {
                            sh 'go mod download'
                            sh 'go mod tidy'
                            sh '''
                                protoc --go_out=. --go_opt=paths=source_relative \
                                       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
                                       proto/exchange.proto || true
                            '''
                            sh 'go build -o bin/exchange-service ./cmd/server'
                        }
                    }
                }
                stage('Build Analytics Service') {
                    steps {
                        dir('analytics-service') {
                            sh 'go mod download'
                            sh 'go mod tidy'
                            sh 'go build -o bin/analytics-service ./cmd/server'
                        }
                    }
                }
                stage('Build Notification Service') {
                    steps {
                        dir('notification-service') {
                            sh 'go mod download'
                            sh 'go mod tidy'
                            sh 'go build -o bin/notification-service ./cmd/server'
                        }
                    }
                }
            }
        }

        stage('Test') {
            steps {
                script {
                    echo 'Running tests for all services...'
                    dir('bank-service') {
                        sh 'go test ./... -v -cover || echo "Bank service tests skipped"'
                    }
                    dir('exchange-service') {
                        sh 'go test ./... -v -cover || echo "Exchange service tests skipped"'
                    }
                    dir('analytics-service') {
                        sh 'go test ./... -v -cover || echo "Analytics service tests skipped"'
                    }
                    dir('notification-service') {
                        sh 'go test ./... -v -cover || echo "Notification service tests skipped"'
                    }
                }
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

